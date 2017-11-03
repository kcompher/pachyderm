package require

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"
)

// Matches checks that a string matches a regular-expression.
func Matches(tb testing.TB, expectedMatch string, actual string, msgAndArgs ...interface{}) {
	tb.Helper()
	r, err := regexp.Compile(expectedMatch)
	if err != nil {
		fatal(tb, msgAndArgs, "Match string provided (%v) is invalid", expectedMatch)
	}
	if !r.MatchString(actual) {
		fatal(tb, msgAndArgs, "Actual string (%v) does not match pattern (%v)", actual, expectedMatch)
	}
}

// Equal checks equality of two values.
func Equal(tb testing.TB, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(expected, actual) {
		fatal(
			tb,
			msgAndArgs,
			"Not equal: %#v (expected)\n"+
				"        != %#v (actual)", expected, actual)
	}
}

// NotEqual checks inequality of two values.
func NotEqual(tb testing.TB, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	if reflect.DeepEqual(expected, actual) {
		fatal(
			tb,
			msgAndArgs,
			"Equal: %#v (expected)\n"+
				"    == %#v (actual)", expected, actual)
	}
}

// EqualOneOf checks if a value is equal to one of the elements of a slice.
func EqualOneOf(tb testing.TB, expecteds []interface{}, actual interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	equal := false
	for _, expected := range expecteds {
		if reflect.DeepEqual(expected, actual) {
			equal = true
			break
		}
	}
	if !equal {
		fatal(
			tb,
			msgAndArgs,
			"Not equal 1 of: %#v (expecteds)\n"+
				"        != %#v (actual)", expecteds, actual)
	}
}

// oneOfEquals is a helper function for OneOfEquals and NoneEquals, that simply
// returns a bool indicating whether 'expected' is in the slice 'actuals'.
func oneOfEquals(expected interface{}, actuals interface{}) (bool, error) {
	e := reflect.ValueOf(expected)
	as := reflect.ValueOf(actuals)
	if as.Kind() != reflect.Slice {
		return false, fmt.Errorf("\"actuals\" must a be a slice, but instead was %s", as.Type().String())
	}
	if e.Type() != as.Type().Elem() {
		return false, nil
	}
	for i := 0; i < as.Len(); i++ {
		if reflect.DeepEqual(e.Interface(), as.Index(i).Interface()) {
			return true, nil
		}
	}
	return false, nil
}

// OneOfEquals checks whether one element of a slice equals a value.
func OneOfEquals(tb testing.TB, expected interface{}, actuals interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	equal, err := oneOfEquals(expected, actuals)
	if err != nil {
		fatal(tb, msgAndArgs, err.Error())
	}
	if !equal {
		fatal(tb, msgAndArgs,
			"Not equal : %#v (expected)\n"+
				" one of  != %#v (actuals)", expected, actuals)
	}
}

// NoneEquals checks one element of a slice equals a value.
func NoneEquals(tb testing.TB, expected interface{}, actuals interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	equal, err := oneOfEquals(expected, actuals)
	if err != nil {
		fatal(tb, msgAndArgs, err.Error())
	}
	if equal {
		fatal(tb, msgAndArgs,
			"Equal : %#v (expected)\n one of == %#v (actuals)", expected, actuals)
	}
}

// NoError checks for no error.
func NoError(tb testing.TB, err error, msgAndArgs ...interface{}) {
	tb.Helper()
	if err != nil {
		fatal(tb, msgAndArgs, "No error is expected but got %s", err.Error())
	}
}

// NoErrorWithinT checks that 'f' finishes within time 't' and does not emit an
// error
func NoErrorWithinT(tb testing.TB, t time.Duration, f func() error, msgAndArgs ...interface{}) {
	tb.Helper()
	errCh := make(chan error)
	go func() {
		// This goro will leak if the timeout is exceeded, but it's okay because the
		// test is failing anyway
		errCh <- f()
	}()
	select {
	case err := <-errCh:
		if err != nil {
			fatal(tb, msgAndArgs, "No error is expected but got %s", err.Error())
		}
	case <-time.After(t):
		fatal(tb, msgAndArgs, "operation did not finish within %s", t.String())
	}
}

// YesError checks for an error.
func YesError(tb testing.TB, err error, msgAndArgs ...interface{}) {
	tb.Helper()
	if err == nil {
		fatal(tb, msgAndArgs, "Error is expected but got %v", err)
	}
}

// NotNil checks a value is non-nil.
func NotNil(tb testing.TB, object interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	success := true

	if object == nil {
		success = false
	} else {
		value := reflect.ValueOf(object)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			success = false
		}
	}

	if !success {
		fatal(tb, msgAndArgs, "Expected value not to be nil.")
	}
}

// Nil checks a value is nil.
func Nil(tb testing.TB, object interface{}, msgAndArgs ...interface{}) {
	tb.Helper()
	if object == nil {
		return
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return
	}

	fatal(tb, msgAndArgs, "Expected value to be nil.")
}

// True checks a value is true.
func True(tb testing.TB, value bool, msgAndArgs ...interface{}) {
	tb.Helper()
	if !value {
		fatal(tb, msgAndArgs, "Should be true.")
	}
}

// False checks a value is false.
func False(tb testing.TB, value bool, msgAndArgs ...interface{}) {
	tb.Helper()
	if value {
		fatal(tb, msgAndArgs, "Should be false.")
	}
}

func logMessage(tb testing.TB, msgAndArgs []interface{}) {
	tb.Helper()
	if len(msgAndArgs) == 1 {
		tb.Logf(msgAndArgs[0].(string))
	}
	if len(msgAndArgs) > 1 {
		tb.Logf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
}

func fatal(tb testing.TB, userMsgAndArgs []interface{}, msgFmt string, msgArgs ...interface{}) {
	tb.Helper()
	logMessage(tb, userMsgAndArgs)
	tb.Fatalf(msgFmt, msgArgs...)
}
