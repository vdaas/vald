package asset

import (
	"os"
	"path/filepath"
	"strings"
)

func GetTestdataPath(filename string) string {
	b := datasetDir()
	basepath := filepath.Dir(b)
	fp, _ := filepath.Abs(basepath + "/vald/internal/asset/testdata/" + filename)
	return fp
}

func datasetDir() string {
	cur, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if strings.HasSuffix(cur, "vald") {
			return cur
		} else {
			cur = filepath.Dir(cur)
		}
	}
	return ""
}
