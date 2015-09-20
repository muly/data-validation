package helpers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// all related to those directly used in the handle func

const SrvTmpDst = "SrvDtaValidTmp"

type Responses struct {
	Res []Response    `json:"validationresponse"`
	Err []ErrResponse `json:"fileerrors"`
}

type Response struct {
	FileName         string            `json:"fileName"`
	ValidationErrors []ValidationFlags `json:"validation"`
}

type ErrResponse struct {
	FileName string `json:"fileName"`
	Error    string `json:"error"`
}

type FileDetails struct {
	ClientFileName string
	SrvrFileName   string
}

// the function GetRules() prepares validation rules from YAML config file.
// 		this function wraps around the parseYaml() which actually does the reading and parsing of the YAML file
// 		and also wraps the parseRules() which converts configuration data from slice of map to map of struct
func GetRules() (rules Rules, err error) {

	// load the YAML configuration
	config, err := parseYaml()
	if err != nil {
		err = errors.New("ERROR parsing the config file: " + err.Error())
		return rules, err
	}

	//load the rules from YAML configuration
	rules = parseRules(config)

	return
}

// the function Upload(), from the http request, reads the multipart files, loops thru each part,
// uploads them one by one to a predefined temp folder on server (configured using the env variable)
// prefixes a randomly generated string to the file name saved on server
// returns the list of files (original file name and name and path of file on server)
func Upload(r *http.Request) (filesList []FileDetails, err error) {

	//initialization
	path := os.Getenv(SrvTmpDst) // get the server destination folder name from environment variable
	reader, _ := r.MultipartReader()

	// loop and process each part of the multipart
	for {

		// read a part
		src, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		// generate random string of alphanumerics of length same as that of the source file name an append it to the file name
		rnd := randomString(len(src.FileName()))
		SrvrFileName := path + "/" + rnd + "_" + src.FileName()

		// open a new file connection, to be used by the upload destination on the server side; use the randomly generated file name (in previous step)
		dstFile, err := os.OpenFile(SrvrFileName, os.O_WRONLY|os.O_CREATE, 0666) //Note: make sure the folder exists, otherwise this step will fail.
		if err != nil {                                                          //Note: this is possibally a 5XX error
			err = errors.New("ERROR with os.OpenFile: " + err.Error())
			return filesList, err
		}

		// copy the content of the part (of multipart) to a file on server and close the connection
		io.Copy(dstFile, src)
		dstFile.Close() //Note: do not defer this close; since the file processing is in the same handler, defer will lock the file,and will interfer when deleting the file (after validation is complete)

		// register the file being uploaded: its original name and the unique name saved as on the server
		filesList = append(filesList, FileDetails{src.FileName(), SrvrFileName})

	}

	return
}

// the function Validation() applies the validation rules to the data in XLSX files.
// loops thru the given list of excel files uploaded to server,
// 		reads each excel file into a slice of slice (2D slice)
// 		convert the 2D slice to slice of structure, for easy manipulation
// 		run the validation on this slice against the given validation rules
//		prepare the list of validation errors
//		prepare the list of file errors
//		delete the excel file after processing is complete
// after all the files are processed, return the prepared responses (validation error and file errors) back to caller
func Validation(filesList []FileDetails, rules Rules) (r Responses, err error) {

	// loop thru and process each file.
	for _, f := range filesList {

		// load the XLSX file stored on server into a 2D slice (slice of slice)
		d, err := readXLSX(f.SrvrFileName)
		if err != nil { //if error, prepare the error response and continue to next file in the multipart
			r.Err = append(r.Err,
				ErrResponse{
					FileName: f.ClientFileName,
					Error:    err.Error(),
				})
			continue
		}

		//convert the slice of slice to slice of struct
		data, err := process(d)
		if err != nil { //if error, prepare the error response and continue to next file in the multipart
			r.Err = append(r.Err,
				ErrResponse{
					FileName: f.ClientFileName,
					Error:    err.Error(),
				})
			continue
		}

		//validation the data against the given rules
		flags := doValidation(data, rules)

		//prepare the response; append response of this part to the rest
		r.Res = append(r.Res,
			Response{
				FileName:         f.ClientFileName,
				ValidationErrors: flags,
			})

		time.Sleep(2 * time.Second) // added delay for testing purpose, to get a chance to physically see the file on the drive before it is deleted. remove this line in final code

		//delete the temp file on server as processing is done
		//TODO: files that are not XLSX are not getting deleted. need to troubleshoot and fix
		err = os.Remove(f.SrvrFileName)
		if err != nil {
			err = errors.New("ERROR deleting the processed file:" + err.Error()) //Note: this is a 5XX error
			return r, err
		}
	}

	return

}
