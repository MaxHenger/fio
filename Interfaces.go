/*
The fio (FileIO) package implements various common file formats in different
types and provides general interfaces to these types. The currently provided
interfaces are:

- FileLoader: Provides a single Load(...) function
- FileSaver: Provides a single Save(...) function
- Settinger: Provides methods to add/set/get variables besides implementing load/save methods
- Spreadsheeter: Provides set/get methods and implements load/save methods

The currently implemented file types are:

- SettingsINI: Implements the Settinger interface for .ini-like files
- SpreadsheetDelim: Implements the Spreadsheeter interface for .csv-like files
*/
package fio

//The FileLoader interface defines a single 'Load(string) error' function
type FileLoader interface {
	Load(file string) error
}

//The FileSaver interface defines a single 'Save(string) error' function
type FileSaver interface {
	Save(file string) error
}

//The Settinger interface provides an interface for general settings files,
//mainly based on the manner in which .ini files are commonly defined. A combination
//of a variable name and value are stored in a file which can (but do not have
//to) be stored under a header name. Setting the value of a variable should not
//be allowed if the value doesn't exists and adding a value should not be allowed
//if the value exists. In the case that a value is added to a non-existant
//header than this header should be created in the process.
type Settinger interface {
	FileLoader
	FileSaver
	HeaderExists(header string) bool
	ValueExists(header, name, value string) bool
	Add(header, name, value string) error
	Set(header, name, value string) error
	Get(header, name string) (string, bool)

	//various methods derived from Get(...)
	GetInt(header, name string) (int, bool, error)
	GetUint(header, name string) (uint, bool, error)
	GetFloat32(header, name string) (float32, bool, error)
	GetFloat64(header, name string) (float64, bool, error)
}

//The Spreadsheeter interface provides an interface for a general spreadsheet
//file. The user should be able to Set and Get values at any given column and
//row. In the case the column/row does not exist yet then all columns/rows
//between the closest-existant column/row and the specified column/row will be
//created.
type Spreadsheeter interface {
	FileLoader
	FileSaver
	Set(row, col int, value string) error
	Get(row, col int) (string, bool)

	//various methods derived from Get(...)
	GetInt(row, col int) (int, bool, error)
	GetUint(row, col int) (uint, bool, error)
	GetFloat32(row, col int) (float32, bool, error)
	GetFloat64(row, col int) (float64, bool, error)
}
