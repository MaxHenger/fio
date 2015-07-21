package fio

//ErrorType is the typedefinition for the various error types that can be
//encountered while using the FileIO framework. Any error can be converted to
//a Error type using a type assertion, after which the error type can be
//inspected for more information about the returned error
type ErrorType byte

//The various error consts to use in conjunction with the ErrorType type
const (
	ErrorTypeParsing         ErrorType = iota //error while parsing the file
	ErrorTypeInvalidArgument                  //an invalid argument was specified to the function
	ErrorTypeNotFound                         //the requested element was not found
	ErrorTypeExists                           //the requested element already exists
	ErrorTypeSaving                           //failed while saving/writing
	ErrorTypeLoading                          //failed while loading/reading
	ErrorTypeTotal                            //unused, indicates the total number of error types
)

//The strings used to describe the ErrorType value when it is printed
const (
	stringErrorTypeParsing         = "Parsing Error"
	stringErrorTypeInvalidArgument = "Invalid Argument"
	stringErrorTypeNotFound        = "Not Found"
	stringErrorTypeExists          = "Already Exists"
	stringErrorTypeSaving          = "Saving Error"
	stringErrorTypeLoading         = "Loading Error"
)

func (et ErrorType) String() string {
	switch et {
	case ErrorTypeParsing:
		return stringErrorTypeParsing
	case ErrorTypeInvalidArgument:
		return stringErrorTypeInvalidArgument
	case ErrorTypeNotFound:
		return stringErrorTypeNotFound
	case ErrorTypeExists:
		return stringErrorTypeExists
	case ErrorTypeSaving:
		return stringErrorTypeSaving
	case ErrorTypeLoading:
		return stringErrorTypeLoading
	}

	return "UNKNOWN"
}

//Error is the struct type that is used throughout the FileIO framework to
//return erros. It includes a type, source and message. The error type is
//indicative of the origin of the error.
type Error struct {
	t       ErrorType
	Source  string
	Message string
}

func (e Error) Error() string {
	return e.t.String() + "[" + e.Source + "]:" + e.Message
}
