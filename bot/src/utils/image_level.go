package main

import (
	"log"
	"io/ioutil"
	"strconv"
	"image"
	"net/http"
	"errors"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func load_image_from_url(URL string) (image.Image, error) {
    //Get the response bytes from the url
    response, err := http.Get(URL)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        return nil, errors.New("received non 200 response code")
    }

    img, _, err := image.Decode(response.Body)
    if err != nil {
        return nil, err
    }

    return img, nil
}

func get_image_level(name_file, username, pp_link, guild_name string, level, percent int64) string {
	img, err := gg.LoadImage("./images/template-level.png")
    if err != nil { log.Fatal(err) }

	w_img := img.Bounds().Max.X
	h_img := img.Bounds().Max.Y

	len_str_max := 24
	name := truncate_str(username, len_str_max)
	server_name := truncate_str(guild_name, len_str_max)
	str_level := "lvl." + strconv.Itoa(int(level)) + " - " + strconv.Itoa(int(percent)) + "%"
	name_output := "./" + name_file

	x_rod := 34.0
	y_rod := 120.0
	h_rod := 13.0
	w_rod_max := 332.0
	w_rod := float64(w_rod_max * float64(percent) / 100)

	x_pp := 20
	y_pp := 20

	x_str := 380.0
	y_login := 35.0
	y_server_name := 50.0
	
	x_level := float64(w_img) / 2
	y_level := float64(h_img) / 2 + 10
	
	// prepare image
	dc := gg.NewContext(w_img, h_img)
	dc.DrawImage(img, 0, 0)

	// draw profile picture
	pp, err := load_image_from_url(pp_link) //80
    if err != nil { log.Fatal(err) }
	dc.DrawImage(pp, x_pp, y_pp)

	// set font
	fontFilePath := "./fonts/TavernSBold/TavernSBold.ttf"
	fontBytes, err := ioutil.ReadFile(fontFilePath)
    if err != nil { log.Fatal(err) }
		// parse font
	f, err := truetype.Parse(fontBytes)
    if err != nil { log.Fatal(err) }
		// define new font face and set it on the context
	dc.SetFontFace(truetype.NewFace(f, &truetype.Options{Size: 17}))

	// draw string login and server name
	dc.SetRGB255(200, 200, 200)
	dc.DrawStringAnchored(name, x_str, y_login, 1, 0)
	dc.DrawStringAnchored(server_name, x_str, y_server_name, 1, 0)
	
	// draw string level
	dc.DrawStringAnchored(str_level, x_level, y_level, 0.5, 0.5)

	// draw rod
	dc.SetRGB255(78, 38, 29)
	dc.DrawRoundedRectangle(x_rod, y_rod, w_rod, h_rod, 8)

	// finish
	dc.Fill()
	dc.SavePNG(name_output)

	return name_output
}
