package main

import (
	"github.com/uptrace/bun"
)

type test_users struct {
	bun.BaseModel `bun:"table:test_users"`

	ID			int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Name		string	`bun:"name,notnull"`
	Roll_number	int		`bun:roll_number,notnull`
}

type guild struct {
	bun.BaseModel `bun:"table:guild"`

	ID			int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Guild_ID	string	`bun:"guild_id,notnull"`
}
