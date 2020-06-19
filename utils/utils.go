package utils

import (
	"encoding/json"
	"log"
	"strings"
)
const(
	printDebug = true
	printError = true
	spliter = "\x00"
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

func ByteToStringArray(b []byte) []string{
	rawStr := string(b[:])
	return strings.Split(rawStr,spliter)
}

func StringArrayToByte (s []string) []byte{
	stringByte := strings.Join(s,spliter)
	return []byte(stringByte)
}

func ByteToKeyValueMap(b []byte) map[string] string{
	var resultMap map[string]string
	if err := json.Unmarshal(b, &resultMap); err != nil{
		log.Fatal("Unmarshal error: " + err.Error())
	}
	return resultMap
}

func KeyValueMapToByte(m map[string]string) []byte{
	b, err := json.Marshal(m)
	if err != nil{
		log.Fatal("Marshal error: " + err.Error())
	}
	return b
}