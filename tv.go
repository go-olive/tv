package tv

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var _ ITv = (*Tv)(nil)

type ITv interface {
	Refresh()
	StreamUrl() (string, bool)
	RoomName() (string, bool)
	StreamerName() (string, bool)
}

type Tv struct {
	SiteID string
	RoomID string
	*Parms
	*Info
}

func NewTv(siteID, roomID string) *Tv {
	tv := &Tv{
		SiteID: siteID,
		RoomID: roomID,
	}
	return tv
}

type Parms struct {
	Cookie string
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

func (tv *Tv) Refresh() {
	site, ok := Sniff(tv.SiteID)
	if !ok {
		return
	}
	site.Snap(tv)
}

func (tv *Tv) StreamUrl() (string, bool) {
	if tv == nil || tv.Info == nil {
		return "", false
	}
	return tv.streamUrl, tv.roomOn
}

func (tv *Tv) RoomName() (string, bool) {
	if tv == nil || tv.Info == nil {
		return "", false
	}
	return tv.roomName, tv.roomNameSet
}

func (tv *Tv) StreamerName() (string, bool) {
	if tv == nil || tv.Info == nil {
		return "", false
	}
	return tv.streamerName, tv.streamerNameSet
}

func (tv *Tv) String() string {
	sb := &strings.Builder{}
	sb.WriteString("Powered by go-olive/tv\n")
	sb.WriteString(format("SiteID", tv.SiteID))
	sb.WriteString(format("RoomID", tv.RoomID))
	if roomName, ok := tv.RoomName(); ok {
		sb.WriteString(format("RoomName", roomName))
	}
	if streamerName, ok := tv.StreamerName(); ok {
		sb.WriteString(format("Streamer", streamerName))
	}
	if streamUrl, ok := tv.StreamUrl(); ok {
		sb.WriteString(format("StreamUrl", streamUrl))
	}
	return sb.String()
}

func format(k, v string) string {
	return fmt.Sprintf("  %-12s%-s\n", k, v)
}

type Streamer interface {
	Stream() *Tv
}

func (tv *Tv) Stream() *Tv {
	return tv
}

type RoomUrl string

func NewRoomUrl(roomUrl string) RoomUrl {
	return RoomUrl(roomUrl)
}

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
