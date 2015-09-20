package helpers

// data validation related functions and types

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ValidationFlags struct { // structure to hold the result of a validation with the original data record and flags indicating which column is invalid
	ServerInfo
	IdInValid         bool
	ServerNameInValid bool
	HostNameInValid   bool
	IPInValid         bool
	DateInValid       bool
}

// the function doValidation() accepts the data and validation rules and
// loops thru each record in the data,
// 		field by field, runs the rule against that field,
// 		and flags it the field doesn't comply with the validation rule
// 		if there is atleast one field inthat record that fails the validation,
//			that record is added to the response along with the validation result flags
//		appends this response with the response of the other records
// once all records are processed, sends the response back to the caller
// examples: //TODO:
// input:
// output:
func doValidation(data []ServerInfo, r Rules) (flags []ValidationFlags) {

	for _, d := range data {
		var IdInValid, ServerNameInValid, HostNameInValid, IPInValid, DateInValid bool

		//// validate Id: by pattern and then condition
		if validByPattern(d.Id, r["Id"].pattern) {
			v, err := validByCondition(d.Id, r["Id"].condition, false)
			if err != nil {
				//TODO: need to handle the error
			}
			IdInValid = !v
		} else {
			IdInValid = true
		}

		//// validate Server Name: by pattern and then by length
		if validByPattern(d.ServerName, r["ServerName"].pattern) {
			v, err := validByCondition(d.ServerName, r["ServerName"].length, true)
			if err != nil {
				//TODO: need to handle the error
			}
			ServerNameInValid = !v
		} else {
			ServerNameInValid = true
		}

		//// validate Host Name: by pattern and then by length
		if validByPattern(d.HostName, r["HostName"].pattern) {
			v, err := validByCondition(d.HostName, r["HostName"].length, true)
			if err != nil {
				//TODO: need to handle the error
			}
			HostNameInValid = !v
		} else {
			HostNameInValid = true
		}

		//// validate IP address: if the pattern is "IP"
		IPInValid = !validIP(d.IP, r["IP"].pattern)

		//// validate Date:
		DateInValid = !validDate(d.Date, r["Date"].pattern)

		// if there is any invalid field, store the records along with the invalid flags
		if IdInValid || ServerNameInValid || HostNameInValid || IPInValid || DateInValid {
			flags = append(flags,
				ValidationFlags{
					ServerInfo:        d,
					IdInValid:         IdInValid,
					ServerNameInValid: ServerNameInValid,
					HostNameInValid:   HostNameInValid,
					IPInValid:         IPInValid,
					DateInValid:       DateInValid,
				})
		}

	}

	return
}

// the function validDate() validates the given date string against the given pattern
//		returns false if validation fails, else returns true
func validDate(val string, patt string) bool {

	// defining the date parts as per Go's Reference Time value
	y, mm, dd, hh, mi, ss := "2006", "01", "02", "15", "04", "05"

	// convert the Go's reference time value into the format of the given pattern
	if patt == strings.Replace(patt, "YYYY", y, 1) {
		patt = strings.Replace(patt, "YY", y[2:], 1)
	} else {
		patt = strings.Replace(patt, "YYYY", y, 1)
	}
	patt = strings.Replace(patt, "MM", mm, 1)
	patt = strings.Replace(patt, "DD", dd, 1)
	patt = strings.Replace(patt, "HH", hh, 1) //TODO: need to add support for 12hr format. right now it is supporting only 24hr format
	patt = strings.Replace(patt, "MI", mi, 1)
	patt = strings.Replace(patt, "SS", ss, 1)

	// compare the given date value with the given pattern (in Go's reference time value), no error means given date is in given format
	_, err := time.Parse(patt, val)
	if err == nil {
		return true
	}
	return false

}

// the function validIP() validates if the given string is an valid IP address are not.
// works for IPv4, IPv6, IPv6 collapsed formats
func validIP(val string, patt string) (result bool) {

	if IP := net.ParseIP(val); IP != nil {
		switch patt {
		case "IP":
			result = true
		case "IPv4":
			if IP.To4() != nil {
				result = true
			}
		case "IPv6":
			if IP.To4() == nil {
				result = true
			}
		}
	}
	return
}

//the function validByPattern() validates the given string against the given pattern.
//		the pattern is expected to be a simple regex indicating what every character in the string should be
func validByPattern(val string, patt string) (result bool) {
	return regexp.MustCompile(`^` + patt + `+$`).MatchString(val)
}

// the function validByCondition() validates the value of the given string against the given condition.
//		if the isLength parameter is false, then the comparition is applicable for numerical values.
// 		if the islength parameter is true, the the condition is applied on the lenght of the given string, instead of its value
// this function supports more than one condition in the condition parameter,
// 		in which case, it loops thru each condition
// 		within the loop the helper function validValue() is called which actually evaluates the single condition
// once the loop is code, this function then applies the logical operation in the condition to give the final validation result
func validByCondition(ltOper string, logi string, isLength bool) (result bool, err error) {

	var validations []bool
	LogiOper, Logic := parseRule(logi)

	if isLength {
		ltOper = strconv.Itoa(len(ltOper))
	}

	for _, cond := range Logic {
		o, rtOper, err := parseLogic(cond)
		if err != nil {
			return result, err
		}

		v, err := validValue(ltOper, o, rtOper)
		if err != nil {
			return result, err
		}
		validations = append(validations, v)
	}

	switch LogiOper { // calculate the combined effect of all the validations considering the given logical operator. since the number of conditions is expected to be not more than 2, the below code hardcodes it to 2
	case "AND":
		result = validations[0] && validations[1]
	case "OR":
		result = validations[0] || validations[1]
	case "":
		result = validations[0]
	default:
		return result, errors.New("Unsupported Logical Operator: only AND, OR are supported")
	}

	return result, nil
}

// function validValue() takes the operator and operands (left and right), and evaluates them and returns if the expression is true or false
// this function is called within the validByCondition() function, which is called with in the loop for each of the logical conditions
func validValue(ltOper, oper, rtOper string) (v bool, err error) {

	lt, err := strconv.Atoi(ltOper)
	if err != nil {
		return
	}
	rt, err := strconv.Atoi(rtOper)
	if err != nil {
		return
	}
	switch oper {
	case "<=":
		v = lt <= rt
	case ">=":
		v = lt >= rt
	case "!=":
		v = lt != rt
	case "<":
		v = lt < rt
	case ">":
		v = lt > rt
	case "==":
		v = lt == rt
	default:
		return v, errors.New("Unsupported Conditional Operator: only <= >= != < > == are supported")
	}

	return
}
