package tv

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Tv struct {
	SiteID string
	RoomID string

	*Info
}

type Option func(*Tv) error

type Info struct {
	Timestamp int64

	streamUrl string
	roomOn    bool

	roomName        string
	roomNameSet     bool
	streamerName    string
	streamerNameSet bool
}

func (i *Info) StreamUrl() (string, bool) {
	return i.streamUrl, i.roomOn
}

func (i *Info) RoomName() (string, bool) {
	return i.roomName, i.roomNameSet
}

func (i *Info) StreamerName() (string, bool) {
	return i.streamerName, i.streamerNameSet
}

type Streamer interface {
	Stream() *Tv
}

func (tv *Tv) Stream() *Tv {
	return tv
}

type RoomUrl string

func (this RoomUrl) SiteID() string {
	u, err := url.Parse(string(this))
	if err != nil {
		return ""
	}
	eTLDPO, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return ""
	}
	siteID := strings.Split(eTLDPO, ".")[0]
	return siteID
}

func (this RoomUrl) Stream() *Tv {
	site, ok := Sniff(this.SiteID())
	if !ok {
		return nil
	}
	return site.Permit(this)
}
