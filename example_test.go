package tv_test

import (
	"fmt"
	"testing"

	"github.com/go-olive/tv"
)

func ExampleTv() {
	tv, err := tv.Snap(tv.NewTv("huya", "518512"), nil)

	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(tv)
}

func ExampleRoomUrl() {
	tv, err := tv.Snap(tv.NewRoomUrl("https://www.huya.com/518512"), nil)
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
