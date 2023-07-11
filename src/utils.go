package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getEnvVar(key string) string {
	err := godotenv.Load("../.env")
	if err != nil { log.Fatal(err) }

	return os.Getenv(key)
}

func isTheBot(idAuthorMessage string, idBot string) bool {
	if idAuthorMessage == idBot { return true }
	return false
}

func beginWithPrefix(str string) bool {
	prefix := getEnvVar("PREFIX_BOT")

	if strings.HasPrefix(str, prefix) { return true }
	return false
}

func trimPrefix(str string) string {
	prefix := getEnvVar("PREFIX_BOT")

	return strings.TrimPrefix(str, prefix)
}
