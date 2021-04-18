package errx

// PanicError is a wrapper used when recover() is called at completion of a goroutine
// and the recovered item needs to be passed on as an error to another goroutine.
// The code that receives the error can then rethrow the item by calling panic with the
// wrapped item as argument.
type PanicError struct {
	item interface{}
}
