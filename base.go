package tv

import (
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Base struct{}

func (b *Base) Name() string {
	return "undefined"
}

func (b *Base) Snap(*Tv) (*Info, error) {
	return nil, errors.New("not implemented")
}

func (b *Base) Permit(roomUrl RoomUrl) *Tv {
	u, err := url.Parse(string(roomUrl))
	if err != nil {
		return nil
	}
	eTLDPO, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		return nil
	}
	siteID := strings.Split(eTLDPO, ".")[0]
	base := strings.TrimPrefix(u.Path, "/")
	roomID := strings.Split(base, "/")[0]
	return &Tv{
		SiteID: siteID,
		RoomID: roomID,
	}
}
