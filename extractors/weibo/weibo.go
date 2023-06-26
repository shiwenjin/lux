package weibo

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	netURL "net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
)

func init() {
	extractors.Register("weibo", New())
}

type playInfo struct {
	Title string            `json:"title"`
	URLs  map[string]string `json:"urls"`
}

type playData struct {
	PlayInfo playInfo `json:"Component_Play_Playinfo"`
}

type weiboData struct {
	Code string   `json:"code"`
	Data playData `json:"data"`
	Msg  string   `json:"msg"`
}

func getXSRFToken() (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := "https://weibo.com/ajax/getversion"
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close() // nolint

	token := utils.MatchOneOf(res.Header.Get("Set-Cookie"), `XSRF-TOKEN=(.+?);`)[1]
	return token, nil
}

func downloadWeiboVideo(url string) ([]*extractors.Data, error) {
	urldata, err := netURL.Parse(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	api := fmt.Sprintf(
		"https://video.h5.weibo.cn/s/video/object?object_id=%s&mid=%s",
		strings.Split(urldata.Path, "/")[1], strings.Split(urldata.Path, "/")[2],
	)
	jsonString, err := request.Get(api, "", nil)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	rawSummary := utils.MatchOneOf(jsonString, `"summary":"(.+?)",`)[1]
	summary, err := strconv.Unquote(strings.Replace(strconv.Quote(rawSummary), `\\u`, `\u`, -1))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rawhdURL := utils.MatchOneOf(jsonString, `"hd_url":"([^"]+)",`)[1]
	unescapedhdURL, err := strconv.Unquote(strings.Replace(strconv.Quote(rawhdURL), `\\u`, `\u`, -1))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	realhdURL := strings.ReplaceAll(unescapedhdURL, `\/`, `/`)
	hdsize, err := request.Size(realhdURL, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	streams := make(map[string]*extractors.Stream, 2)
	streams["hd"] = &extractors.Stream{
		Parts: []*extractors.Part{
			{
				URL:  realhdURL,
				Size: hdsize,
				Ext:  "mp4",
			},
		},
		Size:    hdsize,
		Quality: "hd",
	}
	rawURL := utils.MatchOneOf(jsonString, `"url":"([^"]+)",`)[1]
	unescapedURL, err := strconv.Unquote(strings.Replace(strconv.Quote(rawURL), `\\u`, `\u`, -1))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	realURL := strings.ReplaceAll(unescapedURL, `\/`, `/`)
	size, err := request.Size(realURL, "")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	streams["sd"] = &extractors.Stream{
		Parts: []*extractors.Part{
			{
				URL:  realhdURL,
				Size: size,
				Ext:  "mp4",
			},
		},
		Size:    size,
		Quality: "sd",
	}
	return []*extractors.Data{
		{
			Site:    "微博 weibo.com",
			Title:   summary,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}

func downloadWeiboTV(url string) ([]*extractors.Data, error) {
	APIEndpoint := "https://weibo.com/tv/api/component?page="
	urldata, err := netURL.Parse(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	APIURL := APIEndpoint + netURL.QueryEscape(urldata.Path)
	token, err := getXSRFToken()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	headers := map[string]string{
		"Cookie":       "SUB=_2AkMpogLYf8NxqwJRmP0XxG7kbo10ww_EieKf_vMDJRMxHRl-yj_nqm4NtRB6AiIsKFFGRY4-UuGD5B1-Kf9glz3sp7Ii; XSRF-TOKEN=" + token,
		"Referer":      utils.MatchOneOf(url, `^([^?]+)`)[1],
		"content-type": `application/x-www-form-urlencoded`,
		"x-xsrf-token": token,
	}
	oid := utils.MatchOneOf(url, `tv/show/([^?]+)`)[1]
	postData := "data=" + netURL.QueryEscape("{\"Component_Play_Playinfo\":{\"oid\":\""+oid+"\"}}")
	payload := strings.NewReader(postData)
	res, err := request.Request(http.MethodPost, APIURL, payload, headers)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close() // nolint
	var dataReader io.ReadCloser
	if res.Header.Get("Content-Encoding") == "gzip" {
		dataReader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		dataReader = res.Body
	}
	var data weiboData
	if err = json.NewDecoder(dataReader).Decode(&data); err != nil {
		return nil, errors.WithStack(err)
	}

	if data.Data.PlayInfo.URLs == nil {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}
	realURLs := map[string]string{}
	for k, v := range data.Data.PlayInfo.URLs {
		if strings.HasPrefix(v, "http") {
			continue
		}
		realURLs[k] = "https:" + v
	}

	streams := make(map[string]*extractors.Stream, len(realURLs))
	for q, u := range realURLs {
		size, err := request.Size(u, "")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		streams[q] = &extractors.Stream{
			Parts: []*extractors.Part{
				{
					URL:  u,
					Size: size,
					Ext:  "mp4",
				},
			},
			Size:    size,
			Quality: q,
		}
	}
	return []*extractors.Data{
		{
			Site:    "微博 weibo.com",
			Title:   data.Data.PlayInfo.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}

type extractor struct{}

// New returns a weibo extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	urlInfo, err := netURL.Parse(url)
	if err != nil {
		return nil, errors.New("parse share url fail")
	}
	var videoId string
	if strings.Contains(url, "show?fid=") {
		if len(urlInfo.Query()["fid"]) <= 0 {
			return nil, errors.New("can not parse video id from share url")
		}
		videoId = urlInfo.Query()["fid"][0]
	} else {
		videoId = strings.ReplaceAll(urlInfo.Path, "/tv/show/", "")
	}

	info, err := parseVideoID(videoId)
	if err != nil {
		return nil, err
	}

	resp, err := request.Client.R().Get(info.VideoUrl)

	size := resp.Size()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urlData := &extractors.Part{
		URL:  info.VideoUrl,
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
			Site:    "微博 weibo.com",
			Title:   info.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
			Cover:   info.CoverUrl,
		},
	}, nil
}

func parseVideoID(videoId string) (*VideoParseInfo, error) {
	reqUrl := fmt.Sprintf("https://h5.video.weibo.com/api/component?page=/show/%s", videoId)
	client := request.Client
	videoRes, err := client.R().
		SetHeader("cookie", "login_sid_t=6b652c77c1a4bc50cb9d06b24923210d; cross_origin_proto=SSL; WBStorage=2ceabba76d81138d|undefined; _s_tentry=passport.weibo.com; Apache=7330066378690.048.1625663522444; SINAGLOBAL=7330066378690.048.1625663522444; ULV=1625663522450:1:1:1:7330066378690.048.1625663522444:; TC-V-WEIBO-G0=35846f552801987f8c1e8f7cec0e2230; SUB=_2AkMXuScYf8NxqwJRmf8RzmnhaoxwzwDEieKh5dbDJRMxHRl-yT9jqhALtRB6PDkJ9w8OaqJAbsgjdEWtIcilcZxHG7rw; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9W5Qx3Mf.RCfFAKC3smW0px0; XSRF-TOKEN=JQSK02Ijtm4Fri-YIRu0-vNj").
		SetHeader("referer", "https://h5.video.weibo.com/show/"+videoId).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetHeader("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").
		SetBody([]byte(`data={"Component_Play_Playinfo":{"oid":"` + videoId + `"}}`)).
		Post(reqUrl)
	if err != nil {
		return nil, err
	}
	data := gjson.GetBytes(videoRes.Body(), "data.Component_Play_Playinfo")
	var videoUrl string
	data.Get("urls").ForEach(func(key, value gjson.Result) bool {
		if len(videoUrl) == 0 {
			// 第一条码率最高
			videoUrl = "https:" + value.String()
		}
		return true
	})
	parseInfo := &VideoParseInfo{
		Title:    data.Get("title").String(),
		VideoUrl: videoUrl,
		CoverUrl: "https:" + data.Get("cover_image").String(),
	}
	parseInfo.Author.Name = data.Get("author").String()
	parseInfo.Author.Avatar = "https:" + data.Get("avatar").String()

	return parseInfo, nil
}

type VideoParseInfo struct {
	Author struct {
		Uid    string `json:"uid"`    // 作者id
		Name   string `json:"name"`   // 作者名称
		Avatar string `json:"avatar"` // 作者头像
	} `json:"author"`
	Title    string `json:"title"`     // 描述
	VideoUrl string `json:"video_url"` // 视频播放地址
	MusicUrl string `json:"music_url"` // 音乐播放地址
	CoverUrl string `json:"cover_url"` // 视频封面地址
}
