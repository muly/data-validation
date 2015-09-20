package helpers

// config rules string parsing functions and related types

import (
	"errors"
	"strings"
)

type rule struct {
	name      string
	pattern   string
	condition string
	length    string
}

type Rules map[string]rule

// the function parseRules() converts configuration data from slice of map to map of struct
// examples: //TODO:
// input:
// output:
func parseRules(c configRaw) (r Rules) {
	r = make(Rules)
	for _, c := range c.ValidationRules {
		r[c["name"]] = rule{name: c["name"], pattern: c["pattern"], condition: c["condition"], length: c["length"]}
	}
	return
}

// the function parseRule() splits the given conditions string into multiple conditions
//		looks for operators in the given string, to split by them
// supported Logical operators are: AND, OR (in capital letters)
// total number of conditions supported: 2 conditions combined using either of AND, OR
// Note: technically, the total number of conditions can be more than 2, with slight changes to the code
func parseRule(c string) (LogiOper string, Logic []string) {

	if strings.Index(c, "AND") != -1 { // if AND is found, seperate the multiple conditions using AND
		Logic = strings.Split(c, "AND")
		LogiOper = "AND"
	} else if strings.Index(c, "OR") != -1 { //  else if OR is found, seperate the multiple conditions using OR
		Logic = strings.Split(c, "OR")
		LogiOper = "OR"
	} else {
		Logic = []string{c}
	}

	// trim leading and trailing spaces
	for i := range Logic {
		Logic[i] = strings.TrimSpace(Logic[i])
	}

	return
}

// the function parseLogic() parses the logical condition string into operator and operand
func parseLogic(cond string) (oper string, rtOper string, err error) {

	// prepares the list of all the supported operators
	opers := []string{"<=", ">=", "!=", "<", ">", "=="}

	//loops thru all supported operators, and
	for _, o := range opers {
		if strings.Index(cond, o) != -1 { //if a match found,
			oper = o
			rtOper = strings.TrimSpace(strings.Replace(cond, o, "", -1)) //parses the expression
			break                                                        //and breaks
		}
	}

	if oper == "" { // if no operator found, raises error
		err = errors.New("Operator missing: only <= >= != < > == are supported")
	}
	return
}
