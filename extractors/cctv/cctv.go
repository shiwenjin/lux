package cctv

import (
	"fmt"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
	"github.com/tidwall/gjson"
	"time"
)

func init() {
	extractors.Register("cctv", New())
}

type extractor struct{}

func New() extractors.Extractor {
	return &extractor{}
}

var qualityString = map[string]string{"lowChapters": "low", "chapters": "normal", "chapters2": "hd1", "chapters3": "hd2", "chapters4": "hd3"}

func (e *extractor) Extract(urlAddr string, option extractors.Options) ([]*extractors.Data, error) {
	var (
		err error
	)

	body, err := request.Client.R().Get(urlAddr)
	if err != nil {
		return nil, err
	}

	vid, err := getVid(body.String())
	fmt.Println(vid)

	body, err = request.Client.SetTimeout(60 * time.Second).R().Get(fmt.Sprintf(`https://vdn.apps.cntv.cn/api/getHttpVideoInfo.do?pid=%s`, vid))
	if err != nil {
		return nil, err
	}

	var title string
	jsonObj := gjson.ParseBytes(body.Body())
	title = jsonObj.Get("title").String()

	streams := make(map[string]*extractors.Stream)
	for s, v := range jsonObj.Get("video").Map() {
		if hd, ok := qualityString[s]; ok {

			var totalSize int64
			parts := make([]*extractors.Part, 0, 10)
			partsArr := v.Array()
			for _, p := range partsArr {
				size, _ := request.Size(p.Get("url").String(), "")
				totalSize += size
				parts = append(parts, &extractors.Part{
					URL:  p.Get("url").String(),
					Size: size,
					Ext:  "mp4"})
			}
			streams[hd] = &extractors.Stream{
				Quality: hd,
				Parts:   parts,
				Ext:     "mp4",
				Size:    totalSize,
			}
		}
	}

	return []*extractors.Data{
		{
			Site:    "央视频 cctv.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     urlAddr,
		},
	}, nil
}

func getVid(html string) (string, error) {
	vIds := utils.MatchOneOf(html, `guid\s*=\s*['"]([^'"]+)['"]`)
	return vIds[1], nil
}
