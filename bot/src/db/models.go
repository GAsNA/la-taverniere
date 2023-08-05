package main

import (
	"github.com/uptrace/bun"
)

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

type level struct {
	bun.BaseModel `bun:"table:level"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	User_ID			string	`bun:"user_id,notnull"`
	Guild_ID		string	`bun:"guild_id,notnull"`
	Nb_Msg			int64	`bun:"nb_msg,notnull,default:1"`
	Level			int64	`bun:"level,notnull,default:0"`
}
