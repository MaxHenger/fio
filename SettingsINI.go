package fio

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
)

//SettingsINIHeader represents a header within the SettingsINI struct. All .ini
//variables are stored in a map. The key is the variable name, the corresponding
//value is the variable value.
type SettingsINIHeader struct {
	Values map[string]string
}

//SettingsINI represents the contents of a .ini file. It implemented the
//Settinger interface defined within the FileIO framework. New instances of this
//type should be created using the NewSettingsINI(...) function.
type SettingsINI struct {
	buffer   int
	Filename string
	Headers  map[string]*SettingsINIHeader
}

//NewSettingsINI creates a new SettingsINI instance and returns the pointer. The
//user of this function is required to set a buffer size used while loading the
//file. During loading this buffer will grow to the largest line encountered, an
//initially adequate buffer reduces the times it will have to be resized.
func NewSettingsINI(buffer int) *SettingsINI {
	return &SettingsINI{buffer, "", make(map[string]*SettingsINIHeader)}
}

//Load is capable of loading a file styled like a .ini file. Headers should be
//defined using the '[HeaderName]' syntax, variables as 'Name = Value'.
func (si *SettingsINI) Load(filename string) error {
	//check argument for errors
	if len(filename) == 0 {
		if len(si.Filename) == 0 {
			//no valid filename specified
			return Error{ErrorTypeInvalidArgument, "SettingsINI", "Internal and argument filename are empty"}
		}

		filename = si.Filename
	} else {
		si.Filename = filename
	}

	//open file and create buffer and reader
	file, err := os.Open(filename)

	if err != nil {
		return Error{ErrorTypeLoading, "SettingsINI", "Failed to open the file"}
	}

	buffer := make([]byte, 0, si.buffer)
	reader := bufio.NewReader(file)

	//create the map
	si.Headers = make(map[string]*SettingsINIHeader)

	//process file contents
	var currentHeader *SettingsINIHeader = nil
	eof := false

	for !eof {
		//read a new line
		buffer = buffer[:0]
		eof, err = ReadBufferedLine(reader, &buffer)

		if err != nil {
			file.Close()
			return Error{ErrorTypeLoading, "SettingsINI", "Failed to read new line"}
		}

		line := strings.TrimSpace(string(buffer))

		//if line is empty continue with the next line
		if len(line) == 0 {
			continue
		}

		if len(line) >= 2 {
			//if the line contains a comment, continue
			if line[0] == '/' && line[1] == '/' {
				//this is a comment
				continue
			}

			if line[0] == '[' {
				if line[len(line)-1] == ']' {
					//dealing with a header, retrieve the header name and check its validity
					headerName := line[1 : len(line)-1]

					if len(headerName) == 0 {
						file.Close()
						return Error{ErrorTypeParsing, "SettingsINI", "No header name specified between brackets"}
					}

					//check if the header doesn't already exist
					_, ok := si.Headers[headerName]

					if ok {
						file.Close()
						return Error{ErrorTypeParsing, "SettingsINI", "Header name is specified twice"}
					}

					//create and append the new header, then continue processing next line
					currentHeader = &SettingsINIHeader{make(map[string]string)}
					si.Headers[headerName] = currentHeader
				} else {
					//invalid INI syntax: an opening bracket '[', but no matching closing bracket
					file.Close()
					return Error{ErrorTypeParsing, "SettingsINI", "Invalid header syntax encountered"}
				}
			} else {
				//this line is not empty, a comment or a header, so it must
				//contain a line with a value and a name
				equal := strings.IndexByte(line, '=')

				if equal == -1 {
					file.Close()
					return Error{ErrorTypeParsing, "SettingsINI", "Expected to find an equal-character"}
				}

				name := strings.TrimSpace(line[:equal])
				value := strings.TrimSpace(line[equal+1:])

				if len(name) == 0 {
					file.Close()
					return Error{ErrorTypeParsing, "SettingsINI", "Value pair encountered without name"}
				}

				if len(value) == 0 {
					file.Close()
					return Error{ErrorTypeParsing, "SettingsINI", "Value pair encountered without value"}
				}

				//check if there is a header to put this value pair under
				if currentHeader == nil {
					//nope, create default header
					currentHeader = &SettingsINIHeader{make(map[string]string)}
					si.Headers[""] = currentHeader
				}

				//add value to current header
				currentHeader.Values[name] = value
			}
		}
	}

	//everything is loaded, close the file
	err = file.Close()

	if err != nil {
		return Error{ErrorTypeLoading, "SettingsINI", "Failed to close the file after reading"}
	}

	return nil
}

//Save will store the current SettingsINI type contents to a file.
func (si *SettingsINI) Save(filename string) error {
	//make sure a valid filename exists
	if len(filename) == 0 {
		if len(si.Filename) == 0 {
			return Error{ErrorTypeInvalidArgument, "SettingsINI", "Internal and argument filenames are empty"}
		}

		filename = si.Filename
	}

	//all data is in the settings file now, save to file
	file, err := os.Create(filename)

	if err != nil {
		return Error{ErrorTypeSaving, "SettingsINI", "Failed to open file for writing"}
	}

	//Note: If there is a 'headerless' header (that is: a header without a name),
	//then it should be written first. If not, then the next time the file is read
	//the headerless header will be interpreted as belonging to the previous header
	header, ok := si.Headers[""]

	if ok {
		//write the headerless data
		for valueKey, value := range header.Values {
			var valueBuffer bytes.Buffer

			valueBuffer.WriteString(valueKey)
			valueBuffer.WriteString(" = ")
			valueBuffer.WriteString(value)
			valueBuffer.WriteByte('\n')

			_, err = file.WriteString(valueBuffer.String())

			if err != nil { //uses the fact that if err != nil, then nWritten < len(buffer.String())
				file.Close()
				return Error{ErrorTypeSaving, "SettingsINI", "Failed to write headerless value pair string to file"}
			}
		}
	}

	//loop through all headers
	for headerKey, header := range si.Headers {
		//Note: Only write headered data here
		if len(headerKey) != 0 {
			var headerBuffer bytes.Buffer

			//write the header
			headerBuffer.WriteByte('[')
			headerBuffer.WriteString(headerKey)
			headerBuffer.WriteString("]\n")

			_, err = file.WriteString(headerBuffer.String())

			if err != nil { //uses the fact that if err != nil, then nWritten < len(buffer.String())
				file.Close()
				return Error{ErrorTypeSaving, "SettingsINI", "Failed to write header string to file"}
			}

			//write all the values
			for valueKey, value := range header.Values {
				var valueBuffer bytes.Buffer

				//write the value
				valueBuffer.WriteString(valueKey)
				valueBuffer.WriteString(" = ")
				valueBuffer.WriteString(value)
				valueBuffer.WriteByte('\n')

				_, err = file.WriteString(valueBuffer.String())

				if err != nil { //uses the fact that if err != nil, then nWritten < len(buffer.String())
					file.Close()
					return Error{ErrorTypeSaving, "SettingsINI", "Failed to write headered value pair string to file"}
				}
			}
		}
	}

	err = file.Close()

	if err != nil {
		//failed to close the file
		return Error{ErrorTypeSaving, "SettingsINI", "Failed to close the file after saving"}
	}

	return nil
}

//HeaderExists returns true if the specified header name is stored in the
//SettingsINI type, false otherwise.
func (si *SettingsINI) HeaderExists(header string) (ok bool) {
	_, ok = si.Headers[header]
	return ok
}

//ValueExists returns true if the specified variable exists within the specified
//header, false otherwise
func (si *SettingsINI) ValueExists(header, name string) (ok bool) {
	h, ok := si.Headers[header]

	if !ok {
		return
	}

	_, ok = h.Values[name]
	return ok
}

//Add will store the specified variable value in the specified header using the
//variable name as key. Non-existant headers will be created. If a variable
//already exists with a similar name in the specified header this function will
//return an error, nil otherwise.
func (si *SettingsINI) Add(header, name, value string) error {
	//check if the header exists
	h, ok := si.Headers[header]

	if !ok {
		//header doesn't exist, create it
		h = &SettingsINIHeader{make(map[string]string)}
		si.Headers[header] = h
	}

	//check if the value already exists
	_, ok = h.Values[name]

	if ok {
		//value already exists, cannot add the value
		return Error{ErrorTypeExists, "SettingsINI", "Value pair already exists in the specified header"}
	}

	//set the new value
	h.Values[name] = value
	return nil
}

//Set will set the variable, indicated by name, existing within the specified
//header to the new specified value. The function will return an error if the
//specified header doesn't exist or if the specified variable does not exist
//in the specified header.
func (si *SettingsINI) Set(header, name, value string) error {
	//attempt to find the header
	h, ok := si.Headers[header]

	if !ok {
		//did not find header
		return Error{ErrorTypeNotFound, "SettingsINI", "Could not find header while setting value"}
	}

	_, ok = h.Values[name]

	if !ok {
		//did not find value pair
		return Error{ErrorTypeNotFound, "SettingsINI", "Could not find value pair while setting value"}
	}

	//set value
	h.Values[name] = value
	return nil
}

//Get will return the specified variable's value if it exists in the specified
//header. If either the variable or the header does not exist then the function's
//boolean return value will be false.
func (si *SettingsINI) Get(header, name string) (result string, ok bool) {
	//attempt to find the header
	h, ok := si.Headers[header]

	if !ok {
		//did not find the header
		return "", false
	}

	//return value and value indicating if it exists
	result, ok = h.Values[name]
	return
}

//See Get(...), includes a conversion to int. In case the conversion fails the
//error will be non-nil
func (si *SettingsINI) GetInt(header, name string) (int, bool, error) {
	//use Get(...)
	str, ok := si.Get(header, name)

	if ok {
		//found variable, attempt conversion to int
		result, err := strconv.Atoi(str)
		return result, true, err
	}

	//did not find variable
	return 0, false, nil
}

//See Get(...), includes a conversion to uint. In case the conversion fails the
//error will be non-nil
func (si *SettingsINI) GetUint(header, name string) (uint, bool, error) {
	//use Get(...)
	str, ok := si.Get(header, name)

	if ok {
		//found variable, attempt conversion to uint
		result, err := strconv.ParseUint(str, 10, strconv.IntSize)
		return uint(result), true, err
	}

	//did not find variable
	return 0, false, nil
}

//See Get(...), includes a conversion to float32. In case the conversion fails the
//error will be non-nil
func (si *SettingsINI) GetFloat32(header, name string) (float32, bool, error) {
	//use Get(...)
	str, ok := si.Get(header, name)

	if ok {
		//found variable, attempt conversion to float32
		result, err := strconv.ParseFloat(str, 32)
		return float32(result), true, err
	}

	//did not find variable
	return 0, false, nil
}

//See Get(...), includes a conversion to float64. In case the conversion fails the
//error will be non-nil
func (si *SettingsINI) GetFloat64(header, name string) (float64, bool, error) {
	//use Get(...)
	str, ok := si.Get(header, name)

	if ok {
		//found variable, attempt conversion to float64
		result, err := strconv.ParseFloat(str, 64)
		return result, true, err
	}

	//did not find variable
	return 0, false, nil
}
