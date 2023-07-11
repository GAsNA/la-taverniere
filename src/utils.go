package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func get_env_var(key string) string {
	err := godotenv.Load("../.env")
	if err != nil { log.Fatal(err) }

	return os.Getenv(key)
}

func is_the_bot(idAuthorMessage string, idBot string) bool {
	if idAuthorMessage == idBot { return true }
	return false
}

func begin_with_prefix(str string) bool {
	prefix := get_env_var("PREFIX_BOT")

	if strings.HasPrefix(str, prefix) { return true }
	return false
}

func trim_prefix(str string) string {
	prefix := get_env_var("PREFIX_BOT")

	return strings.TrimPrefix(str, prefix)
}
