package binpack

// Errors in decoding and encoding are handled using panic and recover.
//
// A binpackError is used to distinguish errors (panics) generated in this package.
type binpackError struct {
	err error
}

//// errorf is like error_ but takes Printf-style arguments to construct an error.
//// It always prefixes the message with "binpack: ".
//func errorf(format string, args ...interface{}) {
//	error_(fmt.Errorf("binpack: "+format, args...))
//}
//
//// error wraps the argument error and uses it as the argument to panic.
//func error_(err error) {
//	panic(binpackError{err})
//}

// catchError is meant to be used as a deferred function to turn a panic(binpackError) into a
// plain error. It overwrites the error return of the function that deferred its call.
func catchError(err *error) {
	if e := recover(); e != nil {
		be, ok := e.(binpackError)
		if !ok {
			panic(e)
		}
		*err = be.err
	}
}
