package helpers

import (
	"errors"
	"testing"
)

// unit testing functions for the functions in the parse.go file

// pending test cases:
// input: "> 100 abcd <999"
// output: error  //TODO: this case is not handled yet in parseRule()

//func parseRule(c string) (LogiOper string, Logic []string) {
func TestParseRule(t *testing.T) {

	type output struct {
		LogiOper string
		Logic    []string
	}
	type input struct {
		c string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"> 100 AND <999"}, output{"AND", []string{"> 100", "<999"}}},
		{input{"> 100 OR <999"}, output{"OR", []string{"> 100", "<999"}}},
		{input{"> 100 OR <999"}, output{"OR", []string{"> 100", "<999"}}},
		{input{"> 100"}, output{"", []string{"> 100"}}},
		{input{""}, output{"", []string{"", ""}}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}
		got.LogiOper, got.Logic = parseRule(d.input.c)
		if want.LogiOper != got.LogiOper || !sliceMatch(want.Logic, got.Logic) {
			t.Error("wanted: ", want, ", but got: ", got)
		}

	}

}

// test function to test parseLogic()
// pending test cases:
// input: "> < 100"
// output: _, _, error //multiple operators are not expected in the condition. need to raise error in this case. Note: this case is not yet handled
func TestParseLogic(t *testing.T) {

	type output struct {
		oper   string
		rtOper string
		err    error
	}
	type input struct {
		cond string
	}
	type unittest struct {
		input
		want output
	}

	unittests := []unittest{
		{input{"> 100"}, output{">", "100", nil}},
		{input{"! 100"}, output{err: errors.New("some error")}},
	}

	for _, d := range unittests {
		want := d.want
		got := output{}

		got.oper, got.rtOper, got.err = parseLogic(d.input.cond)

		if (got.err != nil && want.err == nil) || //if got error but not want or
			(got.err == nil && want.err != nil) { // if want error but not got
			t.Error("wanted: ", want, ", but got: ", got)
		} else if want.oper != got.oper || want.rtOper != got.rtOper { // else if no errors compare the other fields
			t.Error("wanted: ", want, ", but got: ", got)
		}
	}

}

// the function sliceMatch is a helper function to compare two slices and return if a boolean value indicating a successful match or not
func sliceMatch(want []string, got []string) bool {
	w := map[string]int{}
	g := map[string]int{}

	for _, j := range want {
		w[j] = 1
	}
	for _, j := range got {
		g[j] = 1
	}

	for s := range w { // comparing wanted against got:
		if w[s] != g[s] {
			return false
		}
	}

	for s := range g { // comparing got against wanted:
		if w[s] != g[s] {
			return false
		}
	}
	return true
}
