package helpers

// excel file related functions and types

import (
	"errors"
	//"io"
	//"net/http"
	//"os"

	"github.com/tealeg/xlsx"
)

type ServerInfo struct { //Note: when structure of the excel file changes update this structure accordingly
	Id         string
	ServerName string
	HostName   string
	IP         string
	Date       string
}

const (
	firstRowHeader bool = true // indicates if the excel file being loaded has a header row or not.
)

// the function ReadXLSX() read the given excel file (file name parameter also has the local path)
// and returns the excel data (of all sheets) into a tabular 2D slice
// and also returns an error if any
// examples:
// input: create an excel file with one sheet and add a record with two columns. fill in "1", "Hello", "World" in the A1, B1, C1 cells. and pass the excel file name (along with path)
// output: [[Hello World]]
// input: to the same file created above, add another row and fill in "2", "Hello", "Web" in A2, B2, C2 cells. and pass the excel file name (along with path)
// output: [[Hello World][Hello Web]]
// input: to the same file creatd above, add another sheet and and fill in "3", "Hello", "Go" in A1, B1, C1 cells. and pass the excel file name (along with path)
// output: [[Hello World][Hello Web][Hello Go]]
// input: to the same file creatd above, add another row (in the second sheet) and fill in "4", "Hello", "Hello", "world" in A2, B2, C2, D2 cells. and pass the excel file name (along with path)
// output: [[Hello World][Hello Web][Hello Go]]
// input: pass the incorrect path of excel file, where the the file doesn't exist
// output: the second return value, err, should have a non nil value
// input: pass the path of an existing file but of a xls file (instead of xlsx file)
// output: the second return value, err, should have a non nil value
// input: pass the path of an existing file but of a different file (let's say a PDF file)
// output: the second return value, err, should have a non nil value
func readXLSX(excelFileName string) (data [][]string, err error) {

	// open the excel file
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil { // return error if file is not accessible
		return data, err
	}

	// read...
	for _, sheet := range xlFile.Sheets { // for each sheet in the excel file //TODO: not yet tested for multiple sheets in the excel file
		for _, row := range sheet.Rows { // for each row in the sheet
			r := []string{}
			for _, cell := range row.Cells { //for each cell in the row
				r = append(r, cell.String()) // add the cell value into the row slice
			}
			data = append(data, r) // add the row processed into the 2D table slice
		}
	}

	return data, err
}

// the function process() converts the 2D slice (slice of slice) to a slice of struct
// this process ignores the first row if indicated as header (by the global constant 'firstRowHeader')
// examples:
// input: [[Hello World Two Three Four]]
// output: [{Hello World Two Three Four}], nil
// input: [[Hello World Two Three Four][Hello Web Two Three Four]]
// output: [{Hello World Two Three Four}{Hello Web Two Three Four}], nil
// input: [[Hello World]]
// output: _, error //because the length if input is less than expected. Note: the struct in destination is
func process(data [][]string) (s []ServerInfo, err error) {
	//TODO: need to convert the parameter passing by value to reference, inorder to save on memory (as pass by value creates a copy which adds to the memory), which helps especially if the excel file is huge

	//check the number of fields in the excel file. atleast 5 columns are expected.
	if len(data[0]) < 5 { // 5 being the count of fields in the ServerInfo struct 	//TODO: need to explore ways to dynamically check the length of the struct, but looks like there is no direct way
		return s, errors.New("Source data row does not have sufficient fields")
	}

	// for each row in the 2D slice,
	for i, r := range data {
		if firstRowHeader && i == 0 { // ignore the first row if it is indicated as header
			continue
		}

		// convert the row slice to struct and append it to the slice of struct
		s = append(s, ServerInfo{Id: r[0], ServerName: r[1], HostName: r[2], IP: r[3], Date: r[4]}) ////Note: when structure of the excel file changes update this line of code accordingly

	}

	return
}
