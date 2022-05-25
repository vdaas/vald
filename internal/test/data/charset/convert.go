package charset

import (
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Utf8ToSjis(s string) string {
	b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.ShiftJIS.NewEncoder()))
	return string(b)
}

func Utf8ToEucjp(s string) string {
	b, _ := ioutil.ReadAll(transform.NewReader(strings.NewReader(s), japanese.EUCJP.NewEncoder()))
	return string(b)
}
