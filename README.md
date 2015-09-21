# data-validation


##Overview 
- This is a REST API endpoint which accept the multipart file as input. the excel format supported is XLSX. 
- The given excel file is uploaded into a temporary folder on the server, and validates each cell in the file row by row
- For validation rules, a predefined yaml file (saved on server side) is been used. the rules are loaded from the yaml file and applied on the data
- Once the validation is done, the validation response JSON is returned back to client with the list of records that deviate from the rules
- In the same JSON response, there is a part that lists of error files which failed because of invalid format or if it is not a XLSX file


##Deployment notes
1. Assumption: that go installed and configured in the environment.
2. Dependencies: install the dependencies using the below commands  
  - go get gopkg.in/yaml.v2  
  - go get github.com/tealeg/xlsx  
3. download the code from "https://github.com/muly/data-validation" and extract it to the folder "github.com/muly/data-validation" on your src folder of the go workspace (as configured in the GOPATH env variable)
4. Temporary location for storing the uploaded excel files on the server is configured using the env variable called "SrvDtaValidTmp" 
   - update this environment variable to the path on the server where you want to keep the uploaded files 
   - for any reason if you have to use a different environment variable, it can be changed by changing the value of the constant SrvTmpDst in the helpers/handler.go file.
5. the path of the configuration YAML file to be saved is the same where the main executable of the service is placed. and the file name should be "config.yaml"

	
##Code organization
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


##Design Decisions

###excel parser 
* most popular 3rd party Go libraries dealing with excel files are 
   - github.com/tealeg/xlsx
   - github.com/mattn/go-ole
* Out of these two,  XLSX package is straightforward and easy to implement with. the disadvantage with this package though is it only supports XLSX format, and doesn't support the old XLS format. 
* on the other side, the go-ole package takes generic approach, and supports both XLSX and XLS formats, however, it is not as stright forward as XLSX package to implement with. 
* Since the required format to support was just XLSX, I thought XLSX package is the best choice.

###yaml parser
* available YAML parsers are as listed below
   - Goyaml: gopkg.in/yaml.v2 or github.com/go-yaml/yaml
   - yaml: Parser for YAML 1.2. By Ross Light. https://bitbucket.org/zombiezen/yaml
   - go-gypsy: YAML Parser. By Kyle Lemons. github.com/kylelemons/go-gypsy
   - and couple of other
* Out of these, the first one, goyaml is simple and straightforward. 
* Just like we do with JSON and XML parsing, by using the goyaml package, we just need to create a struct that matches the format of the yaml file and use the Unmarshal() method from this package.
* because it is straightforward to implement and looks to be very efficient, I thought the goyaml package is best choice.

###yaml format
* The YAML format used here is close to "Sequence of Mappings" style described in the example 2.4 of yaml 1.2 specs (http://www.yaml.org/spec/1.2/spec.html). 
* This best matches the requirement where in the validation rules for a column are grouped together, and that list repeats for each column

### storing Excel in memory vs. on drive
Initially I believed that storing in memory is the straightforward way, but the xlsx package used doesn't have a straightforward way to convert the io.reader from the multifile to the format required by the XLSX package. so I change the approach to store the uploaded XLSX file locally on the server, and delete it after it is processed.


##Scope of improvement and known issues
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
8. right now excel with at least 5 column are accepted. if there are more than 5 column, first 5 are considered. and if there are less than 5 fields, the file is ignored, and is listed under error files json response
9. once the file is processed, the temporary copy on the server is deleted. but in case if the file is not a valid XLSX file (or has less number of fields than expected (5)), then the file on the server is not deleted. this needs to be fixed.


##Testing

### Unit testing
Unit tests are written to cover most of the generic functions that doesn't involve in files.

###Environments testing:
Successfully tested the code in Windows (Windows 7 64 bit) as well as Linux (Ubuntu 14.04.3 LTS 64 bit) environments.

### Overall testing:
* Sample excel XLSX file (can be found in the "example" folder) is prepared with multiple records, each with a possible combination of valid and invalid cases. this excel is used to do the testing
* created a simple html form that uploads the provided file to server. this html file can be found in the "example" folder. I used this to upload the excel file to the end to end testing. it worked well and the html page shows the json response returned by the API call.

