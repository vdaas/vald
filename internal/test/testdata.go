package test

import (
	"os"
	"path/filepath"
	"strings"
)

func GetTestdataPath(filename string) string {
	b := datasetDir()
	basepath := filepath.Dir(b)
	fp, _ := filepath.Abs(basepath + "/vald/internal/test/data/" + filename)
	return fp
}

func datasetDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for cur := filepath.Dir(wd); cur != "/"; cur = filepath.Dir(cur) {
		if strings.HasSuffix(cur, "vald") {
			return cur
		}
	}
	return ""
}
