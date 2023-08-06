package main

type action_db struct {
	id		int64
	name	string
}

var actions_db = []action_db {
					{ id: 0, name: "Youtube Video Announcements" },
					{ id: 0, name: "Youtube Live Announcements" },
					{ id: 0, name: "Logs" },
					{ id: 0, name: "Levels" },
					{ id: 0, name: "Blacklist Logs" },
				}

func get_action_db_by_name(name string) action_db {
	for i := 0; i < len(actions_db); i++ {
		if actions_db[i].name == name {
			return actions_db[i]
		}
	}

	return actions_db[0]
}
