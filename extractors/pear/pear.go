package pear

import (
	"fmt"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

func init() {
	extractors.Register("pearvideo", New())
}

type extractor struct{}

func (e extractor) Extract(uri string, option extractors.Options) ([]*extractors.Data, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.New("匹配视频id失败")
	}
	split := strings.Split(u.Path, "_")
	if len(split) <= 0 {
		return nil, errors.New("匹配视频id失败")
	}
	contId := split[1]
	api := fmt.Sprintf("https://www.pearvideo.com/videoStatus.jsp?contId=%s&mrd=0.4469249512219813", contId)

	var result pearResp
	_, err = request.Client.R().SetHeaders(map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "UTF-8,*;q=0.5",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "en-US,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
	}).SetResult(&result).ForceContentType("application/json").SetHeader("referer", uri).Get(api)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	srcUrl := result.VideoInfo.Videos.SrcUrl
	str := fmt.Sprintf("cont-%s", contId)
	replace := strings.ReplaceAll(srcUrl, result.SystemTime, str)

	size, err := request.Size(replace, "")
	if err != nil {
		return nil, err
	}

	body, err := request.Client.R().Get(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	titles := utils.MatchOneOf(body.String(), `data-title="(.+?)"`)
	if titles == nil || len(titles) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	title := titles[1]

	urlData := &extractors.Part{
		URL:  replace,
		Size: size,
		Ext:  "mp4",
	}

	streams := make(map[string]*extractors.Stream, 1)
	streams["sd"] = &extractors.Stream{
		Parts:   []*extractors.Part{urlData},
		Size:    size,
		Quality: "sd",
	}

	return []*extractors.Data{
		{
			Site:    "梨视频 pearvideo.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     uri,
			Cover:   result.VideoInfo.VideoImage,
		},
	}, nil

}

// New returns a netease extractor.
func New() extractors.Extractor {
	return &extractor{}
}

type pearResp struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	ReqId      string `json:"reqId"`
	SystemTime string `json:"systemTime"`
	VideoInfo  struct {
		PlaySta    string `json:"playSta"`
		VideoImage string `json:"video_image"`
		Videos     struct {
			HdUrl    string `json:"hdUrl"`
			HdflvUrl string `json:"hdflvUrl"`
			SdUrl    string `json:"sdUrl"`
			SdflvUrl string `json:"sdflvUrl"`
			SrcUrl   string `json:"srcUrl"`
		} `json:"videos"`
	} `json:"videoInfo"`
}
