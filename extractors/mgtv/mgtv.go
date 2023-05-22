package mgtv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
)

func init() {
	extractors.Register("mgtv", New())
}

type mgtvVideoStream struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Def  string `json:"def"`
}

type mgtvVideoInfo struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type mgtvVideoData struct {
	Stream       []mgtvVideoStream `json:"stream"`
	StreamDomain []string          `json:"stream_domain"`
	Info         mgtvVideoInfo     `json:"info"`
}

type mgtv struct {
	Data mgtvVideoData `json:"data"`
}

type mgtvVideoAddr struct {
	Info string `json:"info"`
}

type mgtvURLInfo struct {
	URL  string
	Size int64
}

type mgtvPm2Data struct {
	Data struct {
		Atc struct {
			Pm2 string `json:"pm2"`
		} `json:"atc"`
		Info mgtvVideoInfo `json:"info"`
	} `json:"data"`
}

func mgtvM3u8(url string) ([]mgtvURLInfo, int64, error) {
	var data []mgtvURLInfo
	var temp mgtvURLInfo
	var size, totalSize int64
	urls, err := utils.M3u8URLs(url)
	if err != nil {
		return nil, 0, err
	}
	m3u8String, err := request.Get(url, url, nil)
	if err != nil {
		return nil, 0, err
	}
	sizes := utils.MatchAll(m3u8String, `#EXT-MGTV-File-SIZE:(\d+)`)
	// sizes: [[#EXT-MGTV-File-SIZE:1893724, 1893724]]
	for index, u := range urls {
		size, err = strconv.ParseInt(sizes[index][1], 10, 64)
		if err != nil {
			return nil, 0, err
		}
		totalSize += size
		temp = mgtvURLInfo{
			URL:  u,
			Size: size,
		}
		data = append(data, temp)
	}
	return data, totalSize, nil
}

func encodeTk2(str string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(str))
	r1 := regexp.MustCompile(`/\+/g`)
	r2 := regexp.MustCompile(`///g`)
	r3 := regexp.MustCompile(`/=/g`)
	r1.ReplaceAllString(encodeString, "_")
	r2.ReplaceAllString(encodeString, "~")
	r3.ReplaceAllString(encodeString, "-")
	encodeString = utils.Reverse(encodeString)
	return encodeString
}

type extractor struct{}

// New returns a mgtv extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	if option.Playlist {
		return extractPlaylist(url)
	}

	html, err := request.Get(url, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	vid := utils.MatchOneOf(
		url,
		`https?://www.mgtv.com/(?:b|l)/\d+/(\d+).html`,
		`https?://www.mgtv.com/hz/bdpz/\d+/(\d+).html`,
	)
	if vid == nil {
		vid = utils.MatchOneOf(html, `vid: (\d+),`)
	}
	if vid == nil || len(vid) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	// API extract from https://js.mgtv.com/imgotv-miniv6/global/page/play-tv.js
	// getSource and getPlayInfo function
	// Chrome Network JS panel
	headers := map[string]string{
		"Cookie": "PM_CHKID=1",
	}
	clit := fmt.Sprintf("clit=%d", time.Now().Unix()/1000)
	pm2DataString, err := request.Get(
		fmt.Sprintf(
			"https://pcweb.api.mgtv.com/player/video?video_id=%s&tk2=%s",
			vid[1],
			encodeTk2(fmt.Sprintf(
				"did=f11dee65-4e0d-4d25-bfce-719ad9dc991d|pno=1030|ver=5.5.1|%s", clit,
			)),
		),
		url,
		headers,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var pm2 mgtvPm2Data
	if err = json.Unmarshal([]byte(pm2DataString), &pm2); err != nil {
		return nil, errors.WithStack(err)
	}

	dataString, err := request.Get(
		fmt.Sprintf(
			"https://pcweb.api.mgtv.com/player/getSource?video_id=%s&tk2=%s&pm2=%s",
			vid[1], encodeTk2(clit), pm2.Data.Atc.Pm2,
		),
		url,
		headers,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var mgtvData mgtv
	if err = json.Unmarshal([]byte(dataString), &mgtvData); err != nil {
		return nil, errors.WithStack(err)
	}

	title := strings.TrimSpace(
		pm2.Data.Info.Title + " " + pm2.Data.Info.Desc,
	)
	mgtvStreams := mgtvData.Data.Stream
	var addr mgtvVideoAddr
	streams := make(map[string]*extractors.Stream)
	for _, stream := range mgtvStreams {
		if stream.URL == "" {
			continue
		}
		// real download address
		addr = mgtvVideoAddr{}
		addrInfo, err := request.GetByte(mgtvData.Data.StreamDomain[0]+stream.URL, url, headers)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err = json.Unmarshal(addrInfo, &addr); err != nil {
			return nil, errors.WithStack(err)
		}

		m3u8URLs, totalSize, err := mgtvM3u8(addr.Info)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		urls := make([]*extractors.Part, len(m3u8URLs))
		for index, u := range m3u8URLs {
			urls[index] = &extractors.Part{
				URL:  u.URL,
				Size: u.Size,
				Ext:  "ts",
			}
		}
		streams[stream.Def] = &extractors.Stream{
			Parts:   urls,
			Size:    totalSize,
			Quality: stream.Name,
		}
	}

	return []*extractors.Data{
		{
			Site:    "芒果TV mgtv.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}

func extractPlaylist(uri string) ([]*extractors.Data, error) {
	uid := utils.MatchOneOf(uri, `/b/(\d+)/`)[1]
	listUrl := fmt.Sprintf(`https://pcweb.api.mgtv.com/episode/list?collection_id=%s&size=100&callback=&page=1`, uid)

	body, err := request.Client.R().Get(listUrl)
	if err != nil {
		return nil, err
	}

	bodyJson := gjson.ParseBytes(body.Body())
	if bodyJson.Get("code").Int() == 200 {
		result := make([]*extractors.Data, 0)
		bodyJson.Get("data.list").ForEach(func(key, value gjson.Result) bool {
			realUrl, _ := url.JoinPath("https://www.mgtv.com", value.Get("url").String())
			result = append(result, &extractors.Data{
				Site:    "芒果TV mgtv.com",
				URL:     realUrl,
				Title:   value.Get("t3").String(),
				Type:    extractors.DataTypeVideo,
				Streams: make(map[string]*extractors.Stream),
			})
			return true
		})
		return result, nil
	} else {
		return nil, errors.New(bodyJson.Get("msg").String())
	}
}
