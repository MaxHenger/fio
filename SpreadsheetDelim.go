package fio

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
)

//The SpreadsheetDelim type represents spreadsheet-like files wherein columns
//within a row are seperated by a common delimeter and the rows themselves are
//seperated by a newline character.
type SpreadsheetDelim struct {
	buffer    int
	Filename  string
	delimeter string
	Data      [][]string
}

//NewSpreadsheetDelim will create a new instance of the SpreadsheetDelim type
//and return its pointer. The user has to specify the buffer size (which will
//grow to the largest row in the file) and the delimeter string to seperate
//column values by.
func NewSpreadsheetDelim(buffer int, delimeter string) *SpreadsheetDelim {
	return &SpreadsheetDelim{buffer, "", delimeter, nil}
}

//Load will load data from a delimeted spreadsheet into the SpreadsheetDelim
//instance. The user can specify a filename, if an empty one is specified then the
//filename from the previous Load(...) call will be used. Besides the filename
//the user can specify how many columns and/or rows to skip when reading the file
func (sd *SpreadsheetDelim) Load(filename string, skipCols, skipRows int) error {
	//check if a valid filename exists
	if len(filename) == 0 {
		if len(sd.Filename) == 0 {
			//no filename specified to load
			return Error{ErrorTypeInvalidArgument, "SpreadsheetDelim", "No filename specified to load"}
		}

		filename = sd.Filename
	} else {
		sd.Filename = filename
	}

	//open the file, then create a buffer and a buffered reader
	file, err := os.Open(filename)

	if err != nil {
		//failed to load data
		return Error{ErrorTypeLoading, "SpreadsheetDelim", "Failed to load file"}
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 0, sd.buffer)

	//start processing all data
	eof := false
	skipRowCounter := 0

	for !eof {
		//read a new line
		buffer = buffer[:0]
		eof, err = ReadBufferedLine(reader, &buffer)

		if err != nil {
			file.Close()
			return Error{ErrorTypeLoading, "SpreadsheetDelim", "Failed to read a new line"}
		}

		//check if this row should be ignored
		if skipRowCounter < skipRows {
			skipRowCounter++
			continue
		}

		//check if the line contains any data at all
		if len(buffer) == 0 {
			//no data
			continue
		}

		//add a new line of data to store the delimeter seperated values in
		newData := strings.Split(string(buffer), sd.delimeter)

		if len(newData) == 0 {
			//treat this as an error, as the buffer did contain data
			file.Close()
			return Error{ErrorTypeParsing, "SpreadsheetDelim", "Non-empty buffer did not produce delimter split values"}
		}

		if len(newData) <= skipCols {
			//after skipping the specified columns, no data remains
			continue
		}

		sd.Data = append(sd.Data, newData[skipCols:])
	}

	//all data is succesfully loaded
	err = file.Close()

	if err != nil {
		return Error{ErrorTypeLoading, "SpreadsheetDelim", "Failed to close the file after loading"}
	}

	return nil
}

//Save will save the current contents from the SpreadsheetDelim type to a file.
//A filename can be specified, if an empty filename is specified then the
//filename used for the last call to the Load(...) function is used.
func (sd *SpreadsheetDelim) Save(filename string) error {
	//check if a filename is specified
	if len(filename) == 0 {
		if len(sd.Filename) == 0 {
			//no filename to use
			return Error{ErrorTypeInvalidArgument, "SpreadsheetDelim", "No filename specified to save"}
		}

		filename = sd.Filename
	}

	//open the file to save the data to
	file, err := os.Create(filename)

	if err != nil {
		return Error{ErrorTypeSaving, "SpreadsheetDelim", "Failed to create/open file for writing"}
	}

	//loop through all rows
	for _, row := range sd.Data {
		//use buffer for writing strings instead of strings.join, as we need to
		//append a newline character at the end
		var buffer bytes.Buffer

		for i, col := range row {
			if i != 0 {
				buffer.WriteString(sd.delimeter)
			}

			buffer.WriteString(col)
		}

		//buffer contains row of data, append newline character and write it to
		//the file
		buffer.WriteByte('\n')

		_, err = file.WriteString(buffer.String())

		if err != nil { //uses the fact that if err != nil, then nWritten < len(buffer.String())
			file.Close()
			return Error{ErrorTypeSaving, "SpreadsheetDelim", "Failed to write buffered data row to file"}
		}
	}

	err = file.Close()

	if err != nil {
		return Error{ErrorTypeSaving, "SpreadsheetDelim", "Failed to close the file after saving"}
	}

	return nil
}

//Set will set a specifed value at the location of the specified row and column.
//Intermediate rows and columns will be created if the specified location does
//not yet exist
func (sd *SpreadsheetDelim) Set(row, col int, value string) error {
	//append non-existant rows
	for i := len(sd.Data); i <= row; i++ {
		sd.Data = append(sd.Data, []string{})
	}

	//append non-existant columns
	for i := len(sd.Data[row]); i <= col; i++ {
		sd.Data[row] = append(sd.Data[row], "")
	}

	//set specified value
	sd.Data[row][col] = value
	return nil
}

//Get will retrieve a value from the location of the specified row and column.
//In case the location does not exist within the current bounds of the
//spreadsheet then the function will return false.
func (sd *SpreadsheetDelim) Get(row, col int) (string, bool) {
	//check if the specified row exists
	if row >= len(sd.Data) {
		//nope
		return "", false
	}

	//checkif the specified column exists
	if col >= len(sd.Data[row]) {
		//nope
		return "", false
	}

	return sd.Data[row][col], true
}

//See Get(...), includes a conversion to int. In case the conversion fails the
//error will be non-nil
func (sd *SpreadsheetDelim) GetInt(row, col int) (int, bool, error) {
	//use Get(...)
	str, ok := sd.Get(row, col)

	if ok {
		//found variable, attempt conversion to int
		result, err := strconv.Atoi(str)
		return result, true, err
	}

	//not found
	return 0, false, nil
}

//See Get(...), includes a conversion to uint. In case the conversion fails the
//error will be non-nil
func (sd *SpreadsheetDelim) GetUint(row, col int) (uint, bool, error) {
	//use Get(...)
	str, ok := sd.Get(row, col)

	if ok {
		//found variable, attempt conversion to uint
		result, err := strconv.ParseUint(str, 10, strconv.IntSize)
		return uint(result), true, err
	}

	//not found
	return 0, false, nil
}

//See Get(...), includes a conversion to float32. In case the conversion fails the
//error will be non-nil
func (sd *SpreadsheetDelim) GetFloat32(row, col int) (float32, bool, error) {
	//use Get(...)
	str, ok := sd.Get(row, col)

	if ok {
		//found variable, attempt conversion to float32
		result, err := strconv.ParseFloat(str, 32)
		return float32(result), true, err
	}

	return 0, false, nil
}

//See Get(...), includes a conversion to float64. In case the conversion fails the
//error will be non-nil
func (sd *SpreadsheetDelim) GetFloat64(row, col int) (float64, bool, error) {
	//use Get(...)
	str, ok := sd.Get(row, col)

	if ok {
		//found variable, attempt conversion to float64
		result, err := strconv.ParseFloat(str, 64)
		return result, true, err
	}

	return 0, false, nil
}
