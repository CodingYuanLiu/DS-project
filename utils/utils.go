package utils

import (
	"log"
)
const(
	printDebug = true
	printError = true
)
func Debug(format string, args ...interface{}){
	if printDebug{
		log.Printf("[DEBUG] " + format, args...)
	}
}

func Error(format string, args ...interface{}){
	if printError{
		log.Printf("[ERROR] " + format, args...)
	}
}