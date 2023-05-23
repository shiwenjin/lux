package extractors

import (
	"mvdan.cc/xurls/v2"
	"net/url"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/utils"
)

var lock sync.RWMutex
var extractorMap = make(map[string]Extractor)

// Register registers an Extractor.
func Register(domain string, e Extractor) {
	lock.Lock()
	extractorMap[domain] = e
	lock.Unlock()
}

var rxRelaxed = xurls.Relaxed()

// Extract is the main function to extract the data.
func Extract(u string, option Options) ([]*Data, error) {
	u = rxRelaxed.FindString(u)
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else if u.Host == "tv.cctv.com" {
			domain = "cctv"
		} else if u.Host == "my.tv.sohu.com" {
			domain = "sohu"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
