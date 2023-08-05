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

	Guild_ID	string	`bun:"guild_id,pk,unique"`
}

type handler_reaction_role struct {
	bun.BaseModel `bun:"table:handler_reaction_role"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Msg_Link		string	`bun:"msg_link,notnull"`
	Msg_ID			string	`bun:"msg_id,notnull"`
	Reaction		string	`bun:"reaction,notnull"`
	Reaction_ID		string	`bun:"reaction_id"`
	Reaction_Name	string	`bun:"reaction_name,notnull"`
	Role_ID			string	`bun:"role_id,notnull"`
	Guild_ID		string	`bun:"guild_id,notnull"`
}

type nb_msg struct {
	bun.BaseModel `bun:"table:nb_msg"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	User_ID			string	`bun:"user_id,notnull"`
	Guild_ID		string	`bun:"guild_id,notnull"`
	Nb_Msg			int64	`bun:"nb_msg,notnull,default:1"`
}
