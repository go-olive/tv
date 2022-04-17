package tv

import (
	"errors"
	"fmt"
	"sync"
)

// map[string]site
var sites sync.Map

type Site interface {
	Name() string
	Snap(*Tv) error
	Permit(RoomUrl) *Tv
}

func RegisterSite(siteID string, site Site) {
	if _, dup := sites.LoadOrStore(siteID, site); dup {
		panic("site already registered")
	}
}

func Sniff(siteID string) (Site, bool) {
	s, ok := sites.Load(siteID)
	if !ok {
		return nil, ok
	}
	return s.(Site), ok
}

func Snap(streamer Streamer, parms *Parms) (*Tv, error) {
	tv := streamer.Stream()
	if tv == nil {
		return nil, errors.New("streamer not valid")
	}
	if parms != nil {
		tv.Parms = parms
	}
	site, ok := Sniff(tv.SiteID)
	if !ok {
		return nil, fmt.Errorf("site(ID = %s) not supported", tv.SiteID)
	}
	if err := site.Snap(tv); err != nil {
		return nil, err
	}
	return tv, nil
}
