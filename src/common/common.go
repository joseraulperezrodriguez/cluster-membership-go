package common

//BasicError implements the error interface
type BasicError struct {
	ErrorS string
}

//ArgumentParsingError implements the error interface
type ArgumentParsingError struct {
	ErrorS string
}

func (be *BasicError) Error() string {
	return be.ErrorS
}

func (pe *ArgumentParsingError) Error() string {
	return pe.ErrorS
}
