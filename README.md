# data-validation

**This document is still DRAFT**

##Overview: 
- this is a REST endpoint/API which accept the multipart file as input. the excel format supported is XLSX. 
- the given excel file is uploaded into a temporary folder on the server, and validates each cell in the file row by row
- for validation rules, a predefined yaml file (saved on server side) is been used. the rules are loaded from the yaml file and applied on the data
- once the validation is done, the validation response JSON is returned back to client with the list of records that deviate from the rules
- in the same JSON response, there is a part that lists of error files which failed because of invalid format or if it is not a XLSX file

##Deployment notes:
1. Assumption: that go installed and configured in the environment.
2. Dependencies: install the dependencies using the below commands  
  - go get gopkg.in/yaml.v2  
  - go get github.com/tealeg/xlsx  
3. download the code from "https://github.com/muly/data-validation" and extract to the folder "github.com/muly/data-validation" on your src folder of the go workspace (as configured in the GOPATH env variable
4. Temporary location for storing excel files on the server (which are uploaded) is configured using the env variable. 
   - the name of the env variable is SrvDtaValidTmp as used in the code, 
   - however, this environment variable can be changed by changing the value of the constant SrvTmpDst in the helpers/handler.go file.
5. the path of the configuration file to be saved is the same where the main executable of the service is placed. and the file name of the config file should be "config.yaml"

	
	
	
##Code organization:
there are three folders (1) "main" (2) "helpers" (3) "example"
- main/main.go: this file has the main() function for this API, which has the code related to http handler endpoint (/validate)
- main/config.yaml: this is the YAML config file that has the validation rules for each column in the excel file
- helpers/: this folder contains the rest of code separated into multiple files based on the functionality. 
- helpers/handler.go: has the functions and types which are used in the handler function. these cover the basic functionality to upload a XLSX file, parse the YAML file, and run validations. rest of the code in other files are called with in one of these functions here
- helpers/xlsx.go: has the code related to parsing the XLSX file and loading the data into variable and transforming it into the desired format
- helpers/yaml.go: has the code related to parsing the YAML file, and loading the data into variable and transforming it into the desired format
- helpers/parse.go: has the functions and types related to parsing the validation rules (which are loaded from the YAML config file) further down into small pieces in such a way that they can be used to easily applied in the validation steps.
- helpers/validation.go: has the functions and types related to actually doing the validation if each cell.
- helpers/rand.go: has a function that generates a random string (using the seed as current time at nanosecond grain), which is used to prefix the uploaded file on the server. this is to ensure that there is no ambiguity when two files with the same name are received by the server (may be by different clients).
- example/: for examples excel files and client HTML form


##Design Decisions:
###excel parser 
available 3rd party Go libraries..................
the format of excel supported is XLSX.
the 3rd party library from "github.com/tealeg/xlsx" is used here.
###yaml parser
there are 2 libraries most talked about...................................
the 3rd party library from "gopkg.in/yaml.v2" is used here.
the yaml format used is close to "Sequence of Mappings" as described in the example 2.4 of yaml 1.2 specs (http://www.yaml.org/spec/1.2/spec.html)







##Scope of improvement and known issues:
1. comments are added in the code to cover the documentation of the code. these comments can be optimized to be go doc friendly, so that go doc generates the proper documentation.
2. unit test functions are written to cover most of the cases, but can always be expanded to cover more possible cases. right now the code coverage is 46.2% (as checked on 2015/09/20)
3. error handling is done to handle most of the errors, but can always be expanded to cover more.
4. right now, the YAML parsing is within the handler function that handles the validation of XLSX file. this means, the YAML parsing is done every time client sends a XLSX file, which is redundant. The YAMLparsing can be changed to one time process and then pass the parsed YAML data to each validation request to validate  against.
5. Panic: right now all the server error are using PANIC. These can be handled/changed to use a logging mechanism.
6. the json response is complex, and can be simplified to make it easy for the consumers to use.
7. validation rules supported:  
   - only <= >= != < > == comparative operators are supported  
   - only two conditions are supported combined either using AND or OR.  
   - more complex expression are not supported  
8. right now excel with atleast 5 column are accepted. if there are more than 5 column, first 5 are considered. and if there are less than 5 fields, the file is ignored, and is listed under error files json response
9. once the file is processed, the temporary copy on the server is deleted. but in case if the file is not a valid XLSX file (or has less number of fields than expected (5)), then the file on the server is not deleted. this needs to be fixed.



##testing:
### Unit testing
Unit tests are written to cover most of the generic functions that doesn't involve in files.

###Environments testing:
Successfully tested the code in Windows (Windows 7 64 bit) aswell as Linux (Ubuntu 14.04.3 LTS 64 bit) environments.

### Overall testing:
sample excel XLSX file (can be found in the "example" folder) is prepared with multiple records, each with a possible combination of valid and invalid cases. this excel is used to do the testing





