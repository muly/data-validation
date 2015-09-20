# data-validation


Overview

Deployment notes

Code organization:

Design Decisions:
- excel parser 
- yaml parser

Limitations
1) config file:
	* only <= >= != < > == are supported
	* only two conditions are supported combined using AND or OR
2) 




Scope of improvement
1) comments are added in the code to cover the documentation of the code. these comments can be optimized to be go doc friendly,so that go doc generates the proper documentation.
2) unit test functions are written to cover most of the cases, but can always be expanded to cover more possible cases
3) error handling is done to handle most of the errors, but can always be expanded to cover more.
4) right now, the YAML parsing is within the handler function that handles the validation of XLSX file. this means, the YAML parsing is done everytime client sends a XLSX file, which is redundant. The YAMLparsing can be changed to one time process and then pass the parsed YAML data to each validation request to validate  against.
5) Panic: right now all the server error are using PANIC. These can be handled/changed to use a logging mechanism.


Unit testing:



Misc Notes:
1) tmp location for storing excel files to be processed is configured using the env variable. 
	the name of the env variable is ________, 
	and this name can be changed by changing the value of the constant __________.
2) 




