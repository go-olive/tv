package tv

import (
	"encoding/base64"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-olive/tv/util"
)

func init() {
	registerSite("huya", &huya{})
}

type huya struct {
	base
}

func (this *huya) Snap(tv *Tv) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}

	options := []Option{
		this.setRoomOn(),
		this.setStreamURL(),
	}

	for _, option := range options {
		if err := option(tv); err != nil {
			return err
		}
	}

	return nil
}

func (this *huya) Name() string {
	return "虎牙"
}

func (this *huya) streamURL(roomID string) (string, error) {
	roomURL := fmt.Sprintf("https://m.huya.com/%s", roomID)
	userAgent := "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko); Chrome/75.0.3770.100 Mobile Safari/537.36 "
	req := &util.HttpRequest{
		URL:          roomURL,
		Method:       "GET",
		ResponseData: *new(string),
		ContentType:  "application/x-www-form-urlencoded",
		Header: map[string]string{
			"User-Agent": userAgent,
		},
	}
	if err := req.Send(); err != nil {
		return "", err
	}
	respBody := fmt.Sprint(req.ResponseData)
	re := regexp.MustCompile(`liveLineUrl":"([^"]+)",`)
	res := re.FindStringSubmatch(respBody)
	if len(res) > 0 { //有直播链接
		u := res[1]
		if len(u) > 0 {
			decodedRet, _ := base64.StdEncoding.DecodeString(u)
			decodedUrl := string(decodedRet)
			if strings.Contains(decodedUrl, "replay") { //重播
				return "https:" + u, nil
			} else {
				liveLineUrl := this.proc(decodedUrl)
				liveLineUrl = strings.Replace(liveLineUrl, "hls", "flv", -1)
				liveLineUrl = strings.Replace(liveLineUrl, "m3u8", "flv", -1)
				return "https:" + liveLineUrl, nil
			}
		}
	}

	return "", nil
}

func (this *huya) proc(e string) string {
	i := strings.Split(e, "?")[0]
	b := strings.Split(e, "?")[1]
	r := strings.Split(i, "/")
	re := regexp.MustCompile(".(flv|m3u8)")
	s := re.ReplaceAllString(r[len(r)-1], "")
	srcAntiCode := html.UnescapeString(b)

	c := strings.Split(srcAntiCode, "&")
	cc := c[:0]
	n := make(map[string]string)
	for _, x := range c {
		if len(x) > 0 {
			cc = append(cc, x)
			ss := strings.Split(x, "=")
			n[ss[0]] = ss[1]
		}
	}
	c = cc
	fm, _ := url.QueryUnescape(n["fm"])
	uu, _ := base64.StdEncoding.DecodeString(fm)
	u := string(uu)
	p := strings.Split(u, "_")[0]
	f := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	l := n["wsTime"]
	t := "0"
	h := p + "_" + t + "_" + s + "_" + f + "_" + l
	m := util.GetMd5Hash(h)
	url := fmt.Sprintf("%s?wsSecret=%s&wsTime=%s&u=%s&seqid=%s&txyp=%s&fs=%s&sphdcdn=%s&sphdDC=%s&sphd=%s&u=0&t=100&sv=", i, m, l, t, f, n["txyp"], n["fs"], n["sphdcdn"], n["sphdDC"], n["sphd"])
	return url
}

func (this *huya) setRoomOn() Option {
	return func(tv *Tv) error {
		webUserAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
		roomURL := fmt.Sprintf("https://www.huya.com/%s", tv.RoomID)
		req := &util.HttpRequest{
			URL:          roomURL,
			Method:       "GET",
			ResponseData: *new(string),
			ContentType:  "application/x-www-form-urlencoded",
			Header: map[string]string{
				"User-Agent": webUserAgent,
			},
		}
		if err := req.Send(); err != nil {
			return err
		}
		resp := fmt.Sprint(req.ResponseData)
		tv.roomOn = strings.Contains(resp, `"isOn":true`)

		titleRe := regexp.MustCompile(`host-title" title="([^"]+)">`)
		titleSubmatch := titleRe.FindAllStringSubmatch(resp, -1)
		titleRes := make([]string, 0)
		for _, v := range titleSubmatch {
			titleRes = append(titleRes, string(v[1]))
		}
		if len(titleRes) > 0 {
			tv.roomName = titleRes[0]
		}

		return nil
	}
}

func (this *huya) setStreamURL() Option {
	return func(tv *Tv) (err error) {
		if !tv.roomOn {
			return nil
		}
		tv.streamUrl, err = this.streamURL(tv.RoomID)
		if !strings.Contains(tv.streamUrl, "https") {
			tv.roomOn = false
		}
		return
	}
}
