package helpers

import (
	"errors"
	"testing"
)

// unit testing functions for the functions in the validation.go file

// examples:

func TestValidDate(t *testing.T) {

	type output struct {
		result bool
	}
	type input struct {
		val  string
		patt string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"09-11-2015 19:38", "MM-DD-YYYY HH:SS"}, output{true}},
		{input{"2015-09-11 19:38", "MM-DD-YYYY HH:SS"}, output{false}},
		{input{"09-11-2015 19:38", "MM-DD-YY HH:SS"}, output{false}},
		{input{"09-11-15 19:38", "MM-DD-YY HH:SS"}, output{true}},
		{input{"09-11-2015 19:38", "MM-DD-YYYY HH:MI"}, output{true}},
		{input{"09-11-2015 19:38:38", "MM-DD-YYYY HH:MI:SS"}, output{true}},
		{input{"09-11-2015 19:38", "MM-DD-YYYY HH:MI:SS"}, output{false}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.result = validDate(d.input.val, d.input.patt)

		if got.result != want.result {
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		}
	}

}

//examples:

func TestValidIP(t *testing.T) {

	type output struct {
		result bool
	}
	type input struct {
		val  string
		patt string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"1.2.3.4", "IP"}, output{true}},
		{input{"1.2.3.4", "IP"}, output{true}},
		{input{"FE80:0000:0000:0000:0202:B3FF:FE1E:8329", "IP"}, output{true}},
		{input{"FE80::0202:B3FF:FE1E:8329", "IP"}, output{true}},
		{input{"1/2.3.4", "IP"}, output{false}},
		{input{"FE80..0202.B3FF.FE1E.8329", "IP"}, output{false}},
		{input{"FE80.0000.0000.0000.0202.B3FF.FE1E.8329", "IP"}, output{false}},
		{input{"FE80:0000:0000:0000:0202:B3FF:FE1E:8329", "IPv4"}, output{false}},
		{input{"FE80:0000:0000:0000:0202:B3FF:FE1E:8329", "IPv6"}, output{true}},
		{input{"1.2.3.4", "IPv4"}, output{true}},
		{input{"1.2.3.4", "IPv6"}, output{false}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.result = validIP(d.input.val, d.input.patt)

		if got.result != want.result {
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		}
	}

}

func TestValidByPattern(t *testing.T) {

	type output struct {
		result bool
	}
	type input struct {
		val  string
		patt string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"abcxyz", "[0-9]"}, output{false}},
		{input{"abcxyz", "[A-Z]"}, output{false}},
		{input{"abcxyz", "[a-z]"}, output{true}},
		{input{"abcxyz", "[a-zA-Z]"}, output{true}},
		{input{"aBcXyZ", "[a-zA-Z]"}, output{true}},
		{input{"abcxyz1", "[a-zA-Z]"}, output{false}},
		{input{"123456", "[a-zA-Z]"}, output{false}},
		{input{"123456", "[0-9]"}, output{true}},
		{input{"abcxyz123", "[a-zA-Z0-9]"}, output{true}},
		{input{"abcx_yz", "[a-zA-Z0-9]"}, output{false}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.result = validByPattern(d.input.val, d.input.patt)

		if got.result != want.result {
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		}
	}

}

//// examples: validation of value

//// examples: validation of length

//func validByCondition(ltOper string, logi string, isLength bool) (result bool, err error) {
func TestValidByCondition(t *testing.T) {

	type output struct {
		result bool
		err    error
	}
	type input struct {
		ltOper   string
		logi     string
		isLength bool
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"15", ">10 AND <100", false}, output{true, nil}},
		{input{"10", ">10 OR <100", false}, output{true, nil}},
		{input{"10", ">10 AND <100", false}, output{false, nil}},
		{input{"15", ">10 ORR <100", false}, output{err: errors.New("some error")}},
		{input{"15YYY", ">10 OR <100", false}, output{err: errors.New("some error")}},
		{input{"15", ">10YYY OR <100", false}, output{err: errors.New("some error")}},
		{input{"15", ">10 OR <100YY", false}, output{err: errors.New("some error")}},
		{input{"10", "== 10", false}, output{true, nil}},
		{input{"abcdefghi", ">5 AND <10", true}, output{true, nil}},
		{input{"abcdefghij", ">5 AND <10", true}, output{false, nil}},
		{input{"abcdefghijk", ">5 AND <10", true}, output{false, nil}},
		{input{"abcd", ">5 AND <10", true}, output{false, nil}},
		{input{"abcde", ">5 AND <10", true}, output{false, nil}},
		{input{"abcdef", ">5 AND <10", true}, output{true, nil}},
		{input{"abcd", "<5", true}, output{true, nil}},
		{input{"123456", ">5 AND <10", true}, output{true, nil}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.result, got.err = validByCondition(d.input.ltOper, d.input.logi, d.input.isLength)

		if (got.err != nil && want.err == nil) || //if got error but not want or
			(got.err == nil && want.err != nil) { // if want error but not got
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		} else if want.result != got.result { // else if no errors compare the other fields
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		}

	}

}

// examples:

//func validValue(ltOper, oper, rtOper string) (v bool, err error) {
func TestValidValue(t *testing.T) {

	type output struct {
		result bool
		err    error
	}
	type input struct {
		ltOper, oper, rtOper string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"100", "<=", "999"}, output{true, nil}},
		{input{"999", "<=", "999"}, output{true, nil}},
		{input{"1000", ">=", "999"}, output{true, nil}},
		{input{"1000", ">=", "1000"}, output{true, nil}},
		{input{"100", "!=", "999"}, output{true, nil}},
		{input{"999", "!=", "999"}, output{false, nil}},
		{input{"1000", ">", "999"}, output{true, nil}},
		{input{"100", "<", "999"}, output{true, nil}},
		{input{"1000", ">", "1000"}, output{false, nil}},
		{input{"100", "<", "100"}, output{false, nil}},
		{input{"999", "==", "999"}, output{true, nil}},
		{input{"10", "<=", "99"}, output{true, nil}},
		{input{"100AAA", "<=", "99"}, output{err: errors.New("some error")}},
		{input{"100", "<=", "99AAA"}, output{err: errors.New("some error")}},
		{input{"100", "!", "99"}, output{err: errors.New("some error")}},
		//{input{}, output{}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.result, got.err = validValue(d.input.ltOper, d.input.oper, d.input.rtOper)

		if (got.err != nil && want.err == nil) || //if got error but not want or
			(got.err == nil && want.err != nil) { // if want error but not got
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		} else if want.result != got.result { // else if no errors compare the other fields
			t.Error("for input:", d.input, " wanted: ", want, ", but got: ", got)
		}

	}

}
