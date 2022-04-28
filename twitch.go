package tv

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-olive/tv/util"
)

func init() {
	registerSite("twitch", &twitch{})
}

type twitch struct {
	base
}

func (this *twitch) Name() string {
	return "推趣"
}

func (this *twitch) Snap(tv *Tv) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}
	return this.set(tv)
}

func (this *twitch) set(tv *Tv) error {
	roomUrl := fmt.Sprintf("https://www.twitch.tv/%s", tv.RoomID)
	content, err := util.GetURLContent(roomUrl)
	if err != nil {
		return err
	}

	tv.roomOn = strings.Contains(content, `"isLiveBroadcast":true`)
	if !tv.roomOn {
		return nil
	}

	tv.streamUrl = roomUrl
	title, err := util.Match(`"description":"([^"]+)"`, content)
	if err != nil {
		return nil
	}

	tv.roomName = title

	return nil
}
