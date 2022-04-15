package tv

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-olive/tv/util"
)

func init() {
	RegisterSite("huya", &huya{})
}

type huya struct {
	Base
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
	submatch := re.FindAllStringSubmatch(respBody, -1)
	res := make([]string, 0)
	for _, v := range submatch {
		res = append(res, string(v[1]))
	}
	if len(res) < 1 {
		// 虎牙平台有直播是处于直播中的状态但获取不到直播源的情况，打开网页看直播也是同样的情况。俗称死亡回放。
		return "", errors.New("find stream url failed")
	}
	a, _ := base64.RawStdEncoding.DecodeString(res[0])
	return this.proc(string(a)), nil
}

func (this *huya) proc(in string) string {
	ib := strings.Split(in, "?")
	i, b := ib[0], ib[1]
	r := strings.Split(i, "/")
	s := strings.ReplaceAll(r[len(r)-1], ".flv", "")
	s = strings.ReplaceAll(s, ".m3u8", "")
	c := strings.SplitN(b, "&", 4)
	y := c[len(c)-1]
	n := make(map[string]string)
	for _, v := range c {
		if v == "" {
			continue
		}
		vs := strings.SplitN(v, "=", 2)
		n[vs[0]] = vs[1]
	}
	fm := url.PathEscape(n["fm"])
	ub, _ := base64.RawStdEncoding.DecodeString(fm)
	u := string(ub)
	p := strings.Split(u, "_")[0]
	f := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	l := n["wsTime"]
	t := "0"
	h := strings.Join([]string{p, t, s, f, l}, "_")
	m := md5.New()
	io.WriteString(m, h)
	url := fmt.Sprintf("%s?wsSecret=%x&wsTime=%s&u=%s&seqid=%s&%s", i, m.Sum(nil), l, t, f, y)
	url = "https:" + url
	url = strings.ReplaceAll(url, "hls", "flv")
	url = strings.ReplaceAll(url, "m3u8", "flv")
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
			tv.roomNameSet = true
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
		return
	}
}
