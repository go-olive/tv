package tv

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type base struct {
	*Tv
}

func (b *base) Name() string {
	return "undefined"
}

func (b *base) Snap(tv *Tv) error {
	return fmt.Errorf("site(ID = %s) Snap Method not implemented", tv.SiteID)
}

func (b *base) Permit(roomUrl RoomUrl) *Tv {
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
