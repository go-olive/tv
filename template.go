package tv

import (
	"time"
)

func init() {
	registerSite("tmpl", &tmpl{})
}

type tmpl struct {
	base
}

func (this *tmpl) Name() string {
	return "tmpl"
}

func (this *tmpl) Snap(tv *Tv) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}
	return this.set(tv)
}

func (this *tmpl) set(tv *Tv) error {
	tv.roomName = "tmpl room name"
	tv.streamerName = "tmpl streamer name"
	tv.roomOn = true
	tv.streamUrl = "tmpl stream url"
	return nil
}
