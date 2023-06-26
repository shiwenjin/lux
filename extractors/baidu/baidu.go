package baidu

import (
	"encoding/json"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func init() {
	e := new(extractor)
	extractors.Register("m.baidu.com", e)
	extractors.Register("my.mbd.baidu.com", e)
	extractors.Register("baidu", e)
}

type extractor struct {
}

func (e extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	getResp, err := request.Client.R().
		SetHeader("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").
		Get(url)
	if err != nil {
		return nil, err
	}
	jsonData := regexp.MustCompile(`window.jsonData = (.*?);`).FindStringSubmatch(string(getResp.Body()))
	checkJson := strings.ReplaceAll(jsonData[1], "undefined", "null")
	var result baiduResp
	err = json.Unmarshal([]byte(checkJson), &result)
	videoUrl := result.Data.VideoInfo.PlayUrl
	if err != nil {
		//序列化会因为字母大小写报类型错误 但是可以拿到值 如果是这种错误 可以放行 (不这样做 就需要更改结构体 很麻烦)
		if videoUrl == "" {
			return nil, errors.New("百度解析失败")
		}
	}

	size, err := request.Size(videoUrl, "")
	if err != nil {
		return nil, err
	}

	urlData := &extractors.Part{
		URL:  videoUrl,
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
			Site:    "百度视频 mbd.baidu.com",
			Title:   result.Data.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
			Cover:   result.Data.VideoInfo.PosterImage,
		},
	}, nil
}

func New() extractors.Extractor {
	return &extractor{}
}

type baiduResp struct {
	Errno string `json:"errno"`
	Data  struct {
		ShareInvoke  int    `json:"shareInvoke"`
		Id           string `json:"id"`
		Nid          string `json:"nid"`
		Layout       string `json:"layout"`
		IsPayColumn  int    `json:"is_pay_column"`
		Title        string `json:"title"`
		Status       string `json:"status"`
		ViewCounts   int    `json:"view_counts"`
		VideoWidth   string `json:"videoWidth"`
		VideoHeight  string `json:"videoHeight"`
		ResourceType string `json:"resourceType"`
		VideoInfo    struct {
			PosterImage  string        `json:"posterImage"`
			Title        string        `json:"title"`
			Vid          string        `json:"vid"`
			Duration     string        `json:"duration"`
			PlayUrl      string        `json:"play_url"`
			Height       string        `json:"height"`
			Width        string        `json:"width"`
			From         string        `json:"from"`
			PageUrl      string        `json:"pageUrl"`
			IsMicrovideo interface{}   `json:"is_microvideo"`
			VideoCropPos []interface{} `json:"video_crop_pos"`
			ResourceType string        `json:"resourceType"`
			ClarityText  string        `json:"clarityText"`
			FreeDuration string        `json:"freeDuration"`
		} `json:"videoInfo"`
		SpecialColumn string `json:"specialColumn"`
		TitleZone     struct {
			Title string `json:"title"`
		} `json:"titleZone"`
		Author  []interface{} `json:"author"`
		Comment struct {
			Count    int    `json:"count"`
			ThreadId string `json:"threadId"`
		} `json:"comment"`
		Praise struct {
			Count string `json:"count"`
			Liked string `json:"liked"`
		} `json:"praise"`
		ExtLog struct {
		} `json:"extLog"`
		Sharer struct {
		} `json:"sharer"`
		Mount struct {
		} `json:"mount"`
		Mpd            []interface{} `json:"mpd"`
		ShowGuideTips  int           `json:"showGuideTips"`
		SecShareRelate int           `json:"secShareRelate"`
		NewStyle       int           `json:"newStyle"`
		IsThird        string        `json:"isThird"`
		PaidVideoTips  []interface{} `json:"paidVideoTips"`
		ExtRequest     struct {
			Frsrcid     string `json:"frsrcid"`
			Lid         string `json:"lid"`
			Word        string `json:"word"`
			Oword       string `json:"oword"`
			HaokanVideo string `json:"haokan_video"`
			Height      string `json:"height"`
			Title       string `json:"title"`
			Loc         string `json:"loc"`
			LogLoc      string `json:"log_loc"`
		} `json:"extRequest"`
	} `json:"data"`
	Timestamp string `json:"timestamp"`
}
