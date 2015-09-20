package helpers

import (
	"testing"
	"time"
)

// unit testing functions for the functions in the rand.go file

// run the function 2 times (with few sec delay in between) compare the 2 strings
// if they are same, the FAIL, else PASS
// not very strong test but should be good enough
func TestRandomString(t *testing.T) {

	a := randomString(5)
	time.Sleep(2 * time.Nanosecond) // some delay is required to simulate possible real scenario, otherwise, a will be equal to b most of the times
	b := randomString(5)
	if a == b {
		t.Error("randomString() is not generating random strings")
	}
}
