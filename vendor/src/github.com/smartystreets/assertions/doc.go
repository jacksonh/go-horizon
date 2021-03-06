// Package assertions contains the implementations for all assertions which
// are referenced in goconvey's `convey` package
// (github.com/smartystreets/goconvey/convey) for use with the So(...) method.
// They can also be used in traditional Go test functions and even in
// applicaitons.
package assertions

// By default we use a no-op serializer. The actual Serializer provides a JSON
// representation of failure results on selected assertions so the goconvey
// web UI can display a convenient diff.
var serializer Serializer = new(noopSerializer)

// GoConveyMode provides control over JSON serialization of failures. When
// using the assertions in this package from the convey package JSON results
// are very helpful and can be rendered in a DIFF view. In that case, this function
// will be called with a true value to enable the JSON serialization. By default,
// the assertions in this package will not serializer a JSON result, making
// standalone ussage more convenient.
func GoConveyMode(yes bool) {
	if yes {
		serializer = newSerializer()
	} else {
		serializer = new(noopSerializer)
	}
}

type testingT interface {
	Error(args ...interface{})
}

type Assertion struct{ t testingT }

func New(t testingT) *Assertion {
	return &Assertion{t}
}

func (this *Assertion) So(actual interface{}, assert assertion, expected ...interface{}) bool {
	if result := so(actual, assert, expected...); len(result) > 0 {
		this.t.Error("\n" + result)
		return false
	} else {
		return true
	}
}

// So is a convenience function (as opposed to an inconvenience function?)
// for running assertions on arbitrary arguments in any context, be it for testing or even
// application logging. It allows you to perform assertion-like behavior (and get nicely
// formatted messages detailing discrepancies) but without the program blowing up or panicking.
// All that is required is to import this package and call `So` with one of the assertions
// exported by this package as the second parameter.
// The first return parameter is a boolean indicating if the assertion was true. The second
// return parameter is the well-formatted message showing why an assertion was incorrect, or
// blank if the assertion was correct.
//
// Example:
//
//   if ok, message := So(x, ShouldBeGreaterThan, y); !ok {
//        log.Println(message)
//   }
//
func So(actual interface{}, assert assertion, expected ...interface{}) (bool, string) {
	if result := so(actual, assert, expected...); len(result) == 0 {
		return true, result
	} else {
		return false, result
	}
}

// so is like So, except that it only returns the string message, which is blank if the
// assertion passed. Used to facilitate testing.
func so(actual interface{}, assert func(interface{}, ...interface{}) string, expected ...interface{}) string {
	return assert(actual, expected...)
}

// assertion is an alias for a function with a signature that the So()
// function can handle. Any future or custom assertions should conform to this
// method signature. The return value should be an empty string if the assertion
// passes and a well-formed failure message if not.
type assertion func(actual interface{}, expected ...interface{}) string

////////////////////////////////////////////////////////////////////////////
