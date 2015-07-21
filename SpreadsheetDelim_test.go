package fio

import (
	"os"
	"strconv"
	"testing"
)

const testFilenameSpreadsheetDelim = "test_file.csv"

func testSpreadsheetDelimAssign(rows, cols int, ss *SpreadsheetDelim, t *testing.T) {
	//add values in columns and rows 0 to 5
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			toSet := strconv.Itoa(row * (col + 10))
			err := ss.Set(row, col, toSet)

			if err != nil {
				t.Errorf("Failed assigning '%s' [r:%d, c:%d]", toSet, row, col)
			}
		}
	}
}

func testSpreadsheetDelimValues(rows, cols int, ss *SpreadsheetDelim, t *testing.T) {
	//retrieve all values and check their validity
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			expected := strconv.Itoa(row * (col + 10))
			retrieved, err := ss.Get(row, col)

			if err != nil {
				t.Errorf("Failed to retrieve [r:%d, c:%d], error: %s\n", row, col, err.Error())
			}

			if expected != retrieved {
				t.Errorf("'%s' [r:%d, c:%d] != '%s'\n", retrieved, row, col, expected)
			}
		}
	}

	//check nonexistant columns
	for col := 0; col < cols; col++ {
		retrieved, err := ss.Get(rows, col)

		if err == nil {
			t.Errorf("Expected '%s' [r:%d, c:%d] to be unretrievable\n", retrieved, rows, col)
		}
	}

	//check nonexistant rows
	for row := rows; row < rows+10; row++ {
		retrieved, err := ss.Get(row, 0)

		if err == nil {
			t.Errorf("Expected '%s' [r:%d, c:%d] to be unretrievable\n", retrieved, row, 0)
		}
	}
}

func testSpreadsheetOnce(rows, cols int, buffer int, removeFile bool, t *testing.T) {
	//create a new delimeted spreadsheet
	ss := NewSpreadsheetDelim(buffer, ",")

	testSpreadsheetDelimAssign(rows, cols, ss, t)
	testSpreadsheetDelimValues(rows, cols, ss, t)

	//save temporarily
	err := ss.Save(testFilenameSpreadsheetDelim)

	if err != nil {
		t.Errorf("Failed to write to file, error: %s\n", err.Error())
	}

	if removeFile {
		defer func() { os.Remove(testFilenameSpreadsheetDelim) }()
	}

	//attempt to reloader
	ss2 := NewSpreadsheetDelim(buffer, ",")
	err = ss2.Load(testFilenameSpreadsheetDelim, 0, 0)

	if err != nil {
		t.Errorf("Failed to reload the file, error: %s\n", err.Error())
	}

	testSpreadsheetDelimValues(rows, cols, ss2, t)
}

func TestSpreadsheetDelim(t *testing.T) {
	testRows := [...]int{1, 5, 10, 50, 100}
	testCols := [...]int{1, 5, 10, 50, 100}
	testBuffer := [...]int{1, 16, 128, 1024, 4096}

	for _, buffer := range testBuffer {
		for _, rows := range testRows {
			for _, cols := range testCols {
				//fmt.Printf("Buffer = %d (rows = %d, cols = %d) No Removal\n", buffer, rows, cols)
				testSpreadsheetOnce(rows, cols, buffer, false, t)
			}
		}
	}

	for _, buffer := range testBuffer {
		for _, rows := range testRows {
			for _, cols := range testCols {
				//fmt.Printf("Buffer = %d (rows = %d, cols = %d) With Removal\n", buffer, rows, cols)
				testSpreadsheetOnce(rows, cols, buffer, true, t)
			}
		}
	}

	os.Remove(testFilenameSpreadsheetDelim)
}
