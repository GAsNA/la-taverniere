package main

import (
	"log"
	"database/sql"
	"strings"
)

func db_create_table(db *sql.DB, name string, columns []string) {
	query := "CREATE TABLE IF NOT EXISTS " + name + "("

	for i := 0; i < len(columns); i++ {
		query += columns[i]
		if i < len(columns) - 1 { query += "," }
	}

	query += ")"

	_, err := db.Exec(query)
	if err != nil { log.Fatal(err) }

	log.Println("Table created.")
}

func db_insert_in_table(db *sql.DB, name string, columns []string, values []string, not_exists ...string) {

	query := "INSERT INTO " + name + "("

	for i := 0; i < len(columns); i++ {
		query += columns[i]
		if i < len(columns) - 1 { query += ", " }
	}

	query += ") SELECT "

	for i := 0; i < len(values); i++ {
		query += "'" + values[i] + "'"
		if i < len(values) - 1 { query += ", " }
	}

	if len(not_exists) > 0 {
		column_not_exists := not_exists[0]
		value_not_exists := not_exists[1]

		query += " WHERE NOT EXISTS (SELECT " + column_not_exists + " FROM " + name + " WHERE " + column_not_exists + " = '" + value_not_exists + "');"
	}

	log.Println(query)

	_, err := db.Exec(query)
	if err != nil { log.Fatal(err) }

	log.Println("Inserted into table.")
}

func db_select_in_table(db *sql.DB, name string, columns []string) *sql.Rows {
	columns_str := strings.Join(columns, ",")
	
	query := "SELECT " + columns_str + " from " + name
	rows, err := db.Query(query)
	if err != nil { log.Fatal(err) }

	return rows
}
