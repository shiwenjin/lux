package weishi

import (
	"encoding/json"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/pkg/errors"
	"net/url"
	"regexp"
	"strings"
)

type extractor struct {
}

func init() {
	e := new(extractor)
	extractors.Register("isee.weishi.qq.com", e)
	extractors.Register("m.weishi.qq.com", e)
	extractors.Register("weishi", e)
}

func (e extractor) Extract(uri string, option extractors.Options) ([]*extractors.Data, error) {
	if strings.Contains(uri, "video.weishi.qq.com") {
		headerResp, err := request.Client.R().Head(uri)
		if err != nil {
			return nil, err
		}
		uri = headerResp.RawResponse.Request.Response.Header.Get("Location")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("匹配视频id失败")
	}
	query, _ := url.ParseQuery(u.RawQuery)
	id := query.Get("id")
	if id == "" {
		return nil, errors.New("匹配视频id失败")
	}
	getResp, err := request.Client.R().SetHeader("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit").
		Get(uri)
	if err != nil {
		return nil, err
	}
	jsonData := regexp.MustCompile(`window.Vise.initState =(.*?);`).FindStringSubmatch(string(getResp.Body()))
	checkJson := strings.ReplaceAll(jsonData[1], "undefined", "null")
	var weishiResponse weishiResp
	err = json.Unmarshal([]byte(checkJson), &weishiResponse)
	if err != nil {
		return nil, err
	}
	if len(weishiResponse.FeedsList) < 0 {
		return nil, errors.New("解析微视视频失败")
	}
	video := weishiResponse.FeedsList[0]

	size, err := request.Size(video.VideoUrl, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urlData := &extractors.Part{
		URL:  video.VideoUrl,
		Size: size,
		Ext:  "mp4",
	}
	streams := map[string]*extractors.Stream{
		"default": {
			Parts: []*extractors.Part{urlData},
			Size:  size,
		},
	}

	return []*extractors.Data{
		{
			Site:    "微视  weishi.com",
			Title:   video.ShareInfo.BodyMap.Field1.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     uri,
			Cover:   video.VideoCover,
		},
	}, nil
}

func New() extractors.Extractor {
	return &extractor{}
}
