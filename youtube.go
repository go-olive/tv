package tv

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-olive/tv/util"
)

func init() {
	registerSite("youtube", &youtube{})
}

type youtube struct {
	base
}

func (this *youtube) Name() string {
	return "油管"
}

func (this *youtube) Snap(tv *Tv) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}

	streamID, err := this.setRoomOn(tv)
	if err != nil {
		return err
	}
	return this.setStreamURL(tv, streamID)
}

func (this *youtube) setRoomOn(tv *Tv) (string, error) {
	channelURL := fmt.Sprintf("https://www.youtube.com/channel/%s", tv.RoomID)
	content, err := util.GetURLContent(channelURL)
	if err != nil {
		return "", err
	}
	tv.roomOn = strings.Contains(content, `icon":{"iconType":"LIVE"}}`)
	if !tv.roomOn {
		return "", nil
	}
	streamID, err := util.Match(`"videoRenderer":{"videoId":"([^"]+)",`, content)
	if err != nil {
		return "", err
	}
	return streamID, nil
}

func (this *youtube) setStreamURL(tv *Tv, streamID string) error {
	if !tv.roomOn {
		return nil
	}
	// youtube possibly have multiple lives in one channel,
	// curruently the program returns the first one.
	roomURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", streamID)
	tv.streamUrl = roomURL
	roomContent, err := util.GetURLContent(roomURL)
	if err != nil {
		return err
	}
	title, err := util.Match(`name="title" content="([^"]+)"`, roomContent)
	if err != nil {
		return err
	}
	tv.roomName = title
	return nil
}
