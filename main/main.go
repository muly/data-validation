package main

import (
	"encoding/json"
	//"errors"
	"fmt"
	//"io"
	"net/http"
	"os"
	//"time"

	"github.com/muly/data-validation/helpers"
)

func main() {
	fmt.Println("Excel Validation web service running....")
	fmt.Println("The staging area for loading the XLSX files on server side is configured to: ", os.Getenv(helpers.SrvTmpDst))
	fmt.Println("YAML config file is picked from the same path on server as that of the executable of this service")

	http.HandleFunc("/validate", handleUpldVldt)
	http.HandleFunc("/upload", handleUpld)
	http.ListenAndServe("localhost:8080", nil)

}

func handleUpld(w http.ResponseWriter, r *http.Request) {

	//// Upload xlsx
	filesList, err := helpers.Upload(r)
	if err != nil { //Note: this is a 5XX error
		panic(err.Error())
		return
	}

	res := []helpers.Data{}

	for _, f := range filesList {

		d, err := helpers.Load(f.SrvrFileName)
		if err != nil {
			fmt.Println("ERROR ERROR with helpers.Load()", err.Error())
			panic(err.Error())
			return
		}
		fmt.Println(d)

		data := helpers.Data{
			FileDetails: helpers.FileDetails{f.ClientFileName, f.SrvrFileName},
			Data:        d,
		}
		res = append(res, data)
	}

	//// send Response:
	jsonRes, _ := json.Marshal(res)  // prepare the combined response as json and
	fmt.Fprintln(w, string(jsonRes)) //send it back to client

}

func handleUpldVldt(w http.ResponseWriter, r *http.Request) {

	//// Get rules from YAML config
	//TODO: need to move this out of handler function so that it is parsed just once. we need not parse it everytime the excel is send for validation.
	rules, err := helpers.GetRules()
	if err != nil { //Note: this is a 5XX error
		panic(err.Error())
		return
	}

	//// Upload xlsx
	filesList, err := helpers.Upload(r)
	if err != nil { //Note: this is a 5XX error
		panic(err.Error())
		return
	}

	//// do Validation
	res, err := helpers.Validation(filesList, rules)
	if err != nil { //Note: this is a 5XX error
		panic(err.Error())
		return
	}

	//// send Response
	jsonRes, _ := json.Marshal(res)  // prepare the combined response as json and
	fmt.Fprintln(w, string(jsonRes)) //send it back to client
}
