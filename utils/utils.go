package utils

import (
	"os"
)

func DoNothing(args ...any) {}

func ReadFileContent(path string) []byte {
	filecontent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return filecontent
}
