package common

//BaseError implements the error interface
type BaseError struct {
	ErrorS string
}

func (be *BaseError) Error() string {
	return be.ErrorS
}

//TraceableError defines the error interface and can contain nested errors
type TraceableError struct {
	ErrorS   string
	RawError error
}

func (te *TraceableError) Error() string {
	if te.RawError != nil {
		return te.ErrorS + ". Caused by. " + te.Error()
	}
	return te.ErrorS
}

//ArgumentParsingError implements the error interface
type ArgumentParsingError struct {
	ErrorS string
}

func (pe *ArgumentParsingError) Error() string {
	return pe.ErrorS
}
