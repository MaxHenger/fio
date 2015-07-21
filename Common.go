package fio

import (
	"bufio"
	"io"
)

//ReadBufferedLine is a function which wraps around the bufio.ReadLine method to
//return a full line of data. bufio.ReadLine might return an incomplete line
//(requiring multiple calls to bufio.ReadLine), while this function ensures that
//it will return the complete line. In case the provided buffer is not large
//enough, it will grow
func ReadBufferedLine(r *bufio.Reader, p *[]byte) (bool, error) {
	//declare variables used in next loop, initialize isPrefix as true
	var line []byte
	var isPrefix bool = true
	var err error

	for isPrefix {
		//read (partial) data
		line, isPrefix, err = r.ReadLine()

		if err != nil {
			//check if it the error is EOF
			if err == io.EOF {
				//store last data
				*p = append(*p, line...)
				return true, nil
			}

			return true, err
		}

		//append line to buffer
		*p = append(*p, line...)
	}

	return false, nil
}
