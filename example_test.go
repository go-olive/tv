package tv_test

import (
	"fmt"
	"testing"

	"github.com/go-olive/tv"
)

func ExampleTv() {
	tv, err := tv.Snap(&tv.Tv{
		SiteID: "huya",
		RoomID: "518512",
	})

	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(tv)
}

func ExampleRoomUrl() {
	tv, err := tv.Snap(tv.RoomUrl("https://www.huya.com/518512"))
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(tv)
}

func TestExampleTv(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleTv()
}

func TestExampleRoomUrl(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleRoomUrl()
}
