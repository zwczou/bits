package bits

import "errors"

var ErrOverflow = errors.New("overflow")

type CheckError func(error) error

func checkErr(err error, isPanic bool) error {
	if err != nil {
		if isPanic {
			panic(err)
		}
	}
	return err
}

func Must() CheckError {
	return func(err error) error {
		return checkErr(err, true)
	}
}

func Check() CheckError {
	return func(err error) error {
		return checkErr(err, false)
	}
}
