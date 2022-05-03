# OliveTV

[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/go-olive/tv?tab=doc)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/go-olive/tv/goreleaser?style=for-the-badge)](https://github.com/go-olive/tv/actions/workflows/release.yml)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/go-olive/tv)

OliveTV is a CLI utility which gets stream url along with other streamer details.

## Installation

* build from source

    `go install github.com/go-olive/tv/cmd/olivetv@latest`

* download from [**releases**](https://github.com/go-olive/tv/releases)

## Quickstart

After installing, simply use:

```sh
olivetv -u https://www.huya.com/518512
```

or

```sh
olivetv -sid huya -rid 518512
```

> Some platforms might need a cookie, use -c to set one.
>
> eg.  `olivetv -u https://live.douyin.com/xxx -c cookie`

| site     | cookie example                                               |
| -------- | ------------------------------------------------------------ |
| douyin   | `"__ac_nonce=06245c89100e7ab2dd536; __ac_signature=_02B4Z6wo00f01LjBMSAAAIDBwA.aJ.c4z1C44TWAAEx696;"` |
| kuaishou | `"clientid=3; client_key=65890b29; kpn=GAME_ZONE; userId=2832074080; ksliveShowClipTip=true; didv=1651558610000; did=web_e0b9484b9e54432e8cb60673641734bf; Hm_lvt_86a27b7db2c5c0ae37fee4a8a35033ee=1651558613; userId=2832074080; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAU23I2W2S8FeExWzD5xSXt3HGrVhP-jGfjXZ6yGgpqNtjhntnw-_7XcGqM7zdI5uwQvUoMWLPYtrJDNSipu3HXoejsaj0VibuLJ2EBneayxBBH2STgnQZ7fzooB-8r5HvskmLZGYWHutuFnw1tg2y5xQcQQaFihc2Le6EnONj4Lerle6a63TTOFYOa8pRRcVLprdmq0YMOvqIYZ5CnEq5lkaEhrHsWfESUHgv806qk-5eqStgCIgoFeVe2MYPP_0i8KkGLcbbx3nYN1rqS7JUfxNZXQ3cgwoBTAB; kuaishou.live.web_ph=9180f279025aaec3e4c5bd0740c4dc100d3a; kuaishou.live.bfb1s=9b8f70844293bed778aade6e0a8f9942"` |

## API Guide

This API is what powers the cli but is also available to developers that wish to make use of the data OliveTV can retrieve in their own application.

### Extracting streams

```go
package main

import (
	"github.com/go-olive/tv"
)

func main() {
	t, err := tv.New("huya", "518512")
	if err != nil {
		return
	}
	if err := t.Snap(); err != nil {
		return
	}
	if url, liveOn := t.StreamUrl(); liveOn {
		println("stream url: ", url)
	}
}

```

## Contributing

All contributions are welcome. Feel free to open a new thread on the issue tracker or submit a new pull request.

For developer, check out [template file](template.go) if you want to add a new site.

## Credits

This project is inspired by [real-url](https://github.com/wbt5/real-url) and [streamlink](https://github.com/streamlink/streamlink).
