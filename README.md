# data-validation

**This document is still DRAFT**

##Overview

##Deployment notes

##Code organization:

##Design Decisions:
###excel parser 
###yaml parser

##Limitations
1) config file:
	* only <= >= != < > == are supported
	* only two conditions are supported combined using AND or OR
2) 




##Scope of improvement
1) comments are added in the code to cover the documentation of the code. these comments can be optimized to be go doc friendly,so that go doc generates the proper documentation.
2) unit test functions are written to cover most of the cases, but can always be expanded to cover more possible cases
3) error handling is done to handle most of the errors, but can always be expanded to cover more.
4) right now, the YAML parsing is within the handler function that handles the validation of XLSX file. this means, the YAML parsing is done everytime client sends a XLSX file, which is redundant. The YAMLparsing can be changed to one time process and then pass the parsed YAML data to each validation request to validate  against.
5) Panic: right now all the server error are using PANIC. These can be handled/changed to use a logging mechanism.


##Unit testing:
Unit tests are written to cover most of the generic functions that doesn't involve in files.
Code coverage is 45% (as checked on 2015/09/20) 

##Testing Environments:
Sucessfully tested the code in Windows (Windows 7 64 bit) aswell as Linux (Ubuntu 14.04.3 LTS 64 bit) environments.

##Misc Notes:
1) Temperory location for storing excel files on the server (which are uploaded) is configured using the env variable. 
	the name of the env variable is SrvDtaValidTmp is used in the code, 
	hoewver, this environment variable can be changed by changing the value of the constant SrvTmpDst in the helpers/handler.go .
2) 




