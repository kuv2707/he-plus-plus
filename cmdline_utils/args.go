package cmdlineutils

import (
	// "fmt"
	"os"
	"github.com/joho/godotenv"
)

func ReadArgs()map[string]string {
	godotenv.Load()
	args := make(map[string]string)
	
	// starts from index 2
	// we first expect the source file
	// then we expect the flags
	args["src"] = os.Getenv("SOURCE_FILE")
	return args
}
