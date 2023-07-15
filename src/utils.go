package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func get_env_var(key string) string {
	err := godotenv.Load(".env")
	if err != nil { log.Fatal(err) }

	return os.Getenv(key)
}

func is_the_bot(idAuthorMessage string, idBot string) bool {
	if idAuthorMessage == idBot { return true }
	return false
}

func is_role_admin(id_role string) bool {
	admin_role_ids_env := get_env_var("ADMIN_ROLE_IDS")
	admin_role_ids := strings.Split(admin_role_ids_env, ",")

	for i := 0; i < len(admin_role_ids); i++ {
		if admin_role_ids[i] == id_role { return true }
	}

	return false
}

func is_good_format_date(date string) bool {
	if len(date) != 10 {
		return false
	}

	if date[2] != '/' || date[5] != '/' {
		return false
	}

	for i := 0; i < len(date); i++ {
		if i == 2 || i == 5 {
			continue
		}
		if date[i] < '0' || date[i] > '9' {
			return false
		}
	}

	return true
}
