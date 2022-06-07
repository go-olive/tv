package tv

import (
	"log"
	"time"

	"github.com/Davincible/gotiktoklive"
)

func init() {
	registerSite("tiktok", &tiktok{})
}

type tiktok struct {
	base
}

func (this *tiktok) Name() string {
	return "tiktok"
}

func (this *tiktok) Snap(tv *Tv) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}
	return this.set(tv)
}

func (this *tiktok) set(tv *Tv) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("tiktok panic: ", err)
		}
	}()

	tiktok := gotiktoklive.NewTikTok()
	info, err := tiktok.GetRoomInfo(tv.RoomID)
	if err != nil {
		return err
	}

	candi := []string{
		info.StreamURL.FlvPullURL.FullHd1,
		info.StreamURL.FlvPullURL.Hd1,
		info.StreamURL.FlvPullURL.Sd1,
		info.StreamURL.FlvPullURL.Sd2,
	}
	var streamUrl string
	for _, v := range candi {
		if v != "" {
			streamUrl = v
			break
		}
	}

	if streamUrl != "" {
		tv.roomName = info.Owner.Nickname + " is LIVE now"
		tv.streamerName = info.Owner.Nickname
		tv.roomOn = true
		tv.streamUrl = streamUrl
	}

	return nil
}
