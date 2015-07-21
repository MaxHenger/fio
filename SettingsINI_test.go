package fio

import (
	"os"
	"strconv"
	"testing"
)

const testSettingsINIFilename = "testSettings.ini"
const testSettingsINIHeaderName1 = "Header1"
const testSettingsINIHeaderName2 = "Header2WithQuiteALengthExtensionLikeValuesToTestTheSmallerBuffers"
const testSettingsINIHeaderNoBaseName = "A"
const testSettingsINIHeader1BaseName = "BaseName"
const testSettingsINIHeader2BaseName = "ThisIsAnAttemptAtQuiteTheLongVariableNameToTestTheSmallerBuffers"

func testSettingsINIAssign(numNo, num1, num2 int, si *SettingsINI, t *testing.T) {
	for i := 0; i < numNo; i++ {
		name := testSettingsINIHeaderNoBaseName + strconv.Itoa(i)
		value := strconv.Itoa(i * 10)
		err := si.Add("", name, value)

		if err != nil {
			t.Errorf("Failed to add '%s' = '%s' to headerless data\n", name, value)
		}
	}

	for i := 0; i < num1; i++ {
		name := testSettingsINIHeader1BaseName + strconv.Itoa(i)
		value := strconv.Itoa(i * 11)
		err := si.Add(testSettingsINIHeaderName1, name, value)

		if err != nil {
			t.Errorf("Failed to add '%s' = '%s' to header 1 data\n", name, value)
		}
	}

	for i := 0; i < num2; i++ {
		name := testSettingsINIHeader2BaseName + strconv.Itoa(i)
		value := strconv.Itoa(i * 15)
		err := si.Add(testSettingsINIHeaderName2, name, value)

		if err != nil {
			t.Errorf("Failed to add '%s' = '%s' to header 2 data\n", name, value)
		}
	}
}

func testSettingsINIValidity(numNo, num1, num2 int, shouldExist bool, si *SettingsINI, t *testing.T) {
	for i := 0; i < numNo; i++ {
		name := testSettingsINIHeaderNoBaseName + strconv.Itoa(i)
		expected := strconv.Itoa(i * 10)
		value, ok := si.Get("", name)

		if ok != shouldExist {
			t.Errorf("Unexpected ok=%t for '%s' in headerless data", ok, name)

			if ok {
				t.Errorf(" = '%s'\n", value)
			} else {
				t.Errorf("\n")
			}
		} else {
			if ok && expected != value {
				t.Errorf("'%s' != '%s' in headerless data\n", value, expected)
			}
		}
	}

	for i := 0; i < num1; i++ {
		name := testSettingsINIHeader1BaseName + strconv.Itoa(i)
		expected := strconv.Itoa(i * 11)
		value, ok := si.Get(testSettingsINIHeaderName1, name)

		if ok != shouldExist {
			t.Errorf("Unexpected ok=%t for '%s' in header 1\n", ok, name)

			if ok {
				t.Errorf(" = '%s'\n", value)
			} else {
				t.Errorf("\n")
			}
		} else {
			if ok && expected != value {
				t.Errorf("'%s' != '%s' in header 1\n", value, expected)
			}
		}
	}

	for i := 0; i < num2; i++ {
		name := testSettingsINIHeader2BaseName + strconv.Itoa(i)
		expected := strconv.Itoa(i * 15)

		value, ok := si.Get(testSettingsINIHeaderName2, name)

		if ok != shouldExist {
			t.Errorf("Unexpected ok=%t for '%s' in header 2\n", ok, name)

			if ok {
				t.Errorf(" = '%s'\n", value)
			} else {
				t.Errorf("\n")
			}
		} else {
			if ok && expected != value {
				t.Errorf("'%s' != '%s' in header 2\n", value, expected)
			}
		}
	}
}

func testSettingsINIOnce(numNo, num1, num2, buffer int, removeFile bool, t *testing.T) {
	//create settings ini file
	si := NewSettingsINI(buffer)

	//check if none of the values exists
	testSettingsINIValidity(numNo, num1, num2, false, si, t)

	//add all values
	testSettingsINIAssign(numNo, num1, num2, si, t)

	//check if all values exist now
	testSettingsINIValidity(numNo, num1, num2, true, si, t)

	//save the INI data to a temporary file
	err := si.Save(testSettingsINIFilename)

	if err != nil {
		t.Errorf("Error saving INI file to '%s'\n", testSettingsINIFilename)
	}

	if removeFile {
		defer func() { os.Remove(testSettingsINIFilename) }()
	}

	//load the file and check if all data can be loaded
	si2 := NewSettingsINI(buffer)

	err = si2.Load(testSettingsINIFilename)

	if err != nil {
		t.Errorf("Error loading INI file from '%s'\n", testSettingsINIFilename)
	}

	testSettingsINIValidity(numNo, num1, num2, true, si2, t)
}

func TestSettingsINI(t *testing.T) {
	testNo := [...]int{1, 5, 10, 100}
	test1 := [...]int{1, 5, 10, 100}
	test2 := [...]int{1, 5, 10, 100}
	testBuffer := [...]int{1, 2, 10, 1024, 4096}

	//test while removing the temporary file each time
	for _, numNo := range testNo {
		for _, num1 := range test1 {
			for _, num2 := range test2 {
				for _, buffer := range testBuffer {
					testSettingsINIOnce(numNo, num1, num2, buffer, true, t)
				}
			}
		}
	}

	//test without removing the file
	for _, numNo := range testNo {
		for _, num1 := range test1 {
			for _, num2 := range test2 {
				for _, buffer := range testBuffer {
					testSettingsINIOnce(numNo, num1, num2, buffer, false, t)
				}
			}
		}
	}

	//delete the file
	os.Remove(testSettingsINIFilename)
}
