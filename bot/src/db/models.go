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

type action struct {
	bun.BaseModel `bun:"table:action"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Name			string	`bun:"name,notnull"`
}

type channel_for_action struct {
	bun.BaseModel `bun:"table:channel_for_action"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Channel_ID		string	`bun:"channel_id,notnull"`
	Action_ID		int64	`bun:"action_id,notnull"`
	Guild_ID		string	`bun:"guild_id,notnull"`
}

type role_admin struct {
	bun.BaseModel `bun:"table:role_admin"`

	ID				int64	`bun:"id,pk,autoincrement,type:SERIAL"`
	Role_ID			string	`bun:"role_id,notnull"`
	Guild_ID		string	`bun:"guild_id,notnull"`
}
