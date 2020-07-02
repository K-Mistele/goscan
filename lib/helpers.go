package lib

import (
	"log"
)

func Debug(msg string) {
	log.Println("[+]", msg)
}

func Warn(msg string) {
	log.Println("[!]", msg)
}