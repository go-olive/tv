package tv

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	EmptyRoomName     = ""
	EmptyStreamerName = ""
)

var (
	_ ITv = (*Tv)(nil)

	ErrNotSupported = errors.New("streamer not supported")
)

type ITv interface {
	Snap() error
	StreamUrl() (string, bool)
	RoomName() (string, bool)
	StreamerName() (string, bool)
}

type Tv struct {
	SiteID string
	RoomID string

	cookie string

	*Info
}

func New(siteID, roomID string, opts ...Option) (*Tv, error) {
	_, valid := Sniff(siteID)
	if !valid {
		return nil, ErrNotSupported
	}

	t := &Tv{
		SiteID: siteID,
		RoomID: roomID,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t, nil
}

func NewWithUrl(roomUrl string, opts ...Option) (*Tv, error) {
	u := RoomUrl(roomUrl)
	t, err := u.Stream()
	if err != nil {
		err = fmt.Errorf("%+v (err msg = %s)", ErrNotSupported, err.Error())
		return nil, err
	}

	for _, opt := range opts {
		opt(t)
	}
	return t, nil
}

type Option func(*Tv) error

func SetCookie(cookie string) Option {
	return func(t *Tv) error {
		t.cookie = cookie
		return nil
	}
}

type Info struct {
	Timestamp int64

	streamUrl    string
	roomOn       bool
	roomName     string
	streamerName string
}

// Snap takes the latest snapshot of the streamer info that could be retrieved individually.
func (tv *Tv) Snap() error {
	if tv == nil {
		return errors.New("tv is nil")
	}
	site, ok := Sniff(tv.SiteID)
	if !ok {
		return fmt.Errorf("site(ID = %s) not supported", tv.SiteID)
	}
	return site.Snap(tv)
}

func (tv *Tv) SiteName() string {
	if tv == nil {
		return ""
	}
	site, ok := Sniff(tv.SiteID)
	if !ok {
		return ""
	}
	return site.Name()
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
	return tv.roomName, tv.roomName != EmptyRoomName
}

func (tv *Tv) StreamerName() (string, bool) {
	if tv == nil || tv.Info == nil {
		return "", false
	}
	return tv.streamerName, tv.streamerName != EmptyStreamerName
}

func (tv *Tv) String() string {
	sb := &strings.Builder{}
	sb.WriteString("Powered by go-olive/tv\n")
	sb.WriteString(format("SiteID", tv.SiteID))
	sb.WriteString(format("SiteName", tv.SiteName()))
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

func (this RoomUrl) Stream() (*Tv, error) {
	site, ok := Sniff(this.SiteID())
	if !ok {
		return nil, nil
	}
	return site.Permit(this)
}
