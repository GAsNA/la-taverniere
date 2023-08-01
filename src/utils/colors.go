package main

type color struct {
	name	string
	code	int
}

// In function because go doesn't allow const blobal array
func get_colors() []color {
	return []color{
		{ name: "Black", code: 0, },
		{ name: "Aqua", code: 1752220, },
		{ name: "Dark aqua", code: 1146986, },
		{ name: "Green", code: 3066993, },
		{ name: "Dark green", code: 2067276, },
		{ name: "Blue", code: 3447003, },
		{ name: "Dark blue", code: 2123412, },
		{ name: "Purple", code: 10181046, },
		{ name: "Dark purple", code: 7419530, },
		{ name: "Pink", code: 15277667, },
		{ name: "Dark pink", code: 11342935, },
		{ name: "Gold", code: 15844367, },
		{ name: "Dark gold", code: 15844367, },
		{ name: "Orange", code: 15105570, },
		{ name: "Dark orange", code: 11027200, },
		{ name: "Red", code: 15158332, },
		{ name: "Dark red", code: 10038562, },
		{ name: "Grey", code: 9807270, },
		{ name: "Dark grey", code: 9936031, },
		{ name: "Darker grey", code: 8359053, },
		{ name: "Light grey", code: 12370112, },
		{ name: "Navy", code: 3426654, },
		{ name: "Dark navy", code: 2899536, },
		{ name: "Yellow", code: 2899536, },
		{ name: "White", code: 16777215, },
	}
}

func get_color_by_name(name string) color {
	colors := get_colors()
	for i := 0; i < len(colors); i++ {
		if colors[i].name == name { return colors[i] }
	}
	return colors[0]
}
