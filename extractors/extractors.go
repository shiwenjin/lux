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
		} else if u.Host == "video.weishi.qq.com" {
			domain = "weishi"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor, ok := extractorMap[domain]
	if !ok {
		domain = GetTopLevelDomain(u)
	}
	extractor = extractorMap[domain]
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

func GetTopLevelDomain(u string) string {
	url, err := url.Parse(u)
	if err != nil {
		panic(err.Error())
	}

	host := url.Hostname()
	parts := strings.Split(host, ".")

	// Special case for IP addresses
	if len(parts) == 1 {
		return parts[0]
	}

	// Check whether the last two elements are top-level and second-level domain names
	sld := parts[len(parts)-2]

	// Check whether the second-level domain is a common suffix (e.g. co, com, gov, edu, etc.)
	// If so, return the third-level domain as the top-level domain
	commonSuffixes := []string{"com", "org", "net", "int", "edu", "gov", "mil", "arpa"}
	if contains(commonSuffixes, sld) {
		return parts[len(parts)-3]
	}

	// Otherwise, the second-level domain is the top-level domain
	return sld
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
