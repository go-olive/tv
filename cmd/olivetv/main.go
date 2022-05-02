package main

import (
	"flag"
	"fmt"

	"github.com/go-olive/tv"
)

var (
	cookie string
	url    string
	roomID string
	siteID string
)

func init() {
	flag.StringVar(&cookie, "c", "", "site cookie")
	flag.StringVar(&url, "u", "", "room url")
	flag.StringVar(&roomID, "rid", "", "room ID")
	flag.StringVar(&siteID, "sid", "", "site ID")
	flag.Parse()
}

func main() {
	switch {
	case url != "":
		t, err := tv.NewWithUrl(url, tv.SetCookie(cookie))
		if err != nil {
			println(err.Error())
			return
		}
		t.Snap()
		fmt.Println(t)

	case roomID != "" && siteID != "":
		t, err := tv.New(siteID, roomID, tv.SetCookie(cookie))
		if err != nil {
			println(err.Error())
			return
		}
		t.Snap()
		fmt.Println(t)

	default:
		fmt.Println("You need to specify [roomd id and site id] or [room url]\nType olive -h for more information.")
	}
}
