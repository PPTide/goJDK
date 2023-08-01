package validate

import (
	"errors"
	"github.com/PPTide/gojdk/parse"
)

func Validate(file parse.ClassFile) (res interface{}, err error) {
	if file.Magic != 0xCAFEBABE {
		return res, errors.New("magic doesn't match")
	}
	return
}
