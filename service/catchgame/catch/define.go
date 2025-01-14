package catch

import (
	"strconv"
	"strings"
)

type catchNum string

func (c catchNum) IsAll() bool {
	return strings.EqualFold(string(c), "all")
}

func (c catchNum) GetNum() int64 {
	num, err := strconv.ParseInt(string(c), 10, 64)
	if err == nil && num > 0 {
		return num
	}
	return 1
}
