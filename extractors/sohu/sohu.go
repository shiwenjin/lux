package sohu

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/iawia002/lia/array"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
	"time"
)

func init() {
	extractors.Register("sohu", New())
}

type extractor struct{}

func New() extractors.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(urlAddr string, option extractors.Options) ([]*extractors.Data, error) {
	var (
		err  error
		size int64
	)
	result := make([]*extractors.Data, 0)
	streams := make(map[string]*extractors.Stream)

	if strings.Contains(urlAddr, "my.tv.sohu.com/us") {
		headerResp, err := request.Client.R().Head(urlAddr)
		if err != nil {
			return nil, err
		}
		urlAddr = headerResp.RawResponse.Request.Response.Header.Get("Location")
	}

	var htmlMeta *tvSohuComHtmlMeta
	if option.Playlist {
		var videos []*Video
		if strings.Contains(urlAddr, "tv.sohu.com/user") {
			uid := extractUserID(urlAddr)
			if uid == "" {
				return nil, errors.New("uid empty")
			}
			videos, err = extractPlaylistOfZMT(uid)
		} else {
			htmlMeta, err = getHTMLMeta(urlAddr)
			if err != nil {
				return nil, err
			}
			videos, err = extractPlaylist(htmlMeta.Aid)
		}

		if err != nil {
			return nil, err
		}
		needDownloadItems := utils.NeedDownloadList(option.Items, option.ItemStart, option.ItemEnd, len(videos))

		defaultStream := extractors.Stream{
			ID:      "",
			Quality: "",
			Parts:   nil,
			Size:    0,
			Ext:     "mp4",
			NeedMux: false,
		}

		for index, video := range videos {
			if !array.ItemInArray(index+1, needDownloadItems) {
				continue
			}
			streams["default"] = &defaultStream
			result = append(result, &extractors.Data{
				Site:    "搜狐 sohu",
				URL:     video.URL,
				Title:   video.Title,
				Type:    extractors.DataTypeVideo,
				Streams: streams,
				Cover:   video.Cover,
			})
		}

		return result, nil
	}

	var title, cover string
	if htmlMeta.Aid != "" && htmlMeta.Vid != "" {
		//_, err = checkPermission(htmlMeta.Aid, htmlMeta.Vid, htmlMeta.TVid)
		//if err != nil {
		//	return nil, errors.WithStack(err)
		//}

		videoMeta, err := getVideoClips(urlAddr, htmlMeta.Vid, "")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if len(videoMeta.Data.Su) != 0 {
			urls := make([]*extractors.Part, 0, len(videoMeta.Data.Su))
			title = videoMeta.Data.TvName
			cover = videoMeta.Data.CoverImg

			//> 1，走m3u8合并
			for i, su := range videoMeta.Data.Su {
				if strings.Index(su, "http") == 0 {
					urls = append(urls, &extractors.Part{
						URL:  su,
						Size: videoMeta.Data.ClipsBytes[i],
						Ext:  "mp4",
					})
				} else {
					playUrl := fmt.Sprintf("https://%s/ip?new=%s&num=1&key=%s&ch=%s&pt=1&pg=2&prod=h5n&uid=%d", videoMeta.Allot, su, videoMeta.Data.Ck[i], videoMeta.Data.Ch, videoMeta.Syst)
					realPlayUrl, err := request.Client.R().Get(playUrl)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					temp := &extractors.Part{
						URL:  gjson.ParseBytes(realPlayUrl.Body()).Get("servers.0.url").String(),
						Size: videoMeta.Data.ClipsBytes[i],
						Ext:  "mp4",
					}
					urls = append(urls, temp)

					//只有一条，直接返回
					if len(videoMeta.Data.Su) == 1 {
						urls = append(urls, &extractors.Part{
							URL:  gjson.ParseBytes(realPlayUrl.Body()).Get("servers.url").String(),
							Size: videoMeta.Data.ClipsBytes[i],
							Ext:  "mp4",
						})
					}
				}
			}

			streams["default"] = &extractors.Stream{
				Parts: urls,
				Size:  videoMeta.Data.TotalBytes,
			}
		} else {
			result, err := getVideoNew(htmlMeta.Vid)
			if err != nil {
				return nil, err
			}
			title = result.Title
			cover = result.CoverImg

			var mp4PlayUrlVo Mp4PlayUrlVo
			_, err = request.Client.R().SetResult(&mp4PlayUrlVo).Get(result.Mp4PlayUrl)
			if err != nil {
				return nil, err
			}

			var realAddr string
			if mp4PlayUrlVo.Servers == nil {
				realAddr = result.Mp4PlayUrl
			} else {
				realAddr = mp4PlayUrlVo.Servers[0].Url
			}

			urlData := &extractors.Part{
				URL:  realAddr,
				Size: result.TotalBytes,
				Ext:  "mp4",
			}
			streams["default"] = &extractors.Stream{
				Parts: []*extractors.Part{urlData},
				Size:  size,
			}
		}
	}

	return []*extractors.Data{
		{
			Site:    "搜狐  sohu.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     urlAddr,
			Cover:   cover,
		},
	}, nil
}

func extractUserID(url string) string {
	// 使用正则表达式提取数字部分
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(url, -1)
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}
	return ""
}

// 获取自媒体主页视频
func extractPlaylistOfZMT(uid string) ([]*Video, error) {
	uri := fmt.Sprintf(`https://my.tv.sohu.com/user/wm/ta/v.do?uid=%s&pg=1&size=50`, uid)
	body, err := request.Client.R().
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetHeader("Transfer-Encoding", "chunked").
		SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36").Get(uri)

	if err != nil {
		return nil, err
	}

	var videos []*Video
	gjson.ParseBytes(body.Body()).Get("data.list").ForEach(func(key, value gjson.Result) bool {
		videos = append(videos, &Video{
			Title: value.Get("title").String(),
			URL:   value.Get("url").String(),
			Cover: value.Get("cover640").String(),
		})
		return true
	})
	return videos, nil
}

// 获取主页视频
func extractPlaylist(aid string) ([]*Video, error) {
	uri := fmt.Sprintf(`https://pl.hd.sohu.com/videolist?playlistid=%s&pageRule=3&pagesize=100&pagenum=1`, aid)
	body, err := request.Client.R().Get(uri)
	if err != nil {
		return nil, err
	}

	var bodyUtf8 []byte
	bodyUtf8, err = utils.GbkToUtf8(body.Body())

	var videos []*Video
	gjson.ParseBytes(bodyUtf8).Get("videos").ForEach(func(key, value gjson.Result) bool {
		videos = append(videos, &Video{
			Title: value.Get("name").String(),
			URL:   value.Get("pageUrl").String(),
			Cover: value.Get("largePicUrl").String(),
		})
		return true
	})
	return videos, nil
}

func getVideoNew(vid string) (result *VideoNewVo, err error) {
	uri := "https://my.tv.sohu.com/play/videonew.do"
	body, err := request.Client.R().SetQueryParams(map[string]string{
		"vid":     vid,
		"ver":     "2",
		"ssl":     "1",
		"referer": uri,
		"t":       cast.ToString(time.Now().UnixMilli()),
	}).Get(uri)

	data := gjson.ParseBytes(body.Body()).Get("data")
	result = &VideoNewVo{
		Title:      data.Get("tvName").String(),
		Duration:   data.Get("totalDuration").Int(),
		CoverImg:   data.Get("coverImg").String(),
		Mp4PlayUrl: data.Get("mp4PlayUrl|0").String(),
		TotalBytes: data.Get("totalBytes").Int(),
	}
	return
}

func checkPermission(aid, vid, tVid string) (string, error) {
	uri := "https://api.store.sohu.com/video/pc/checkpermission?aid=" + aid + "&vid=" + vid + "&tvid=" + tVid + "&_=" + cast.ToString(time.Now().UnixMilli())
	body, err := request.Client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36").
		SetHeader("Referer", uri).
		SetHeader("Cookie", "user_isOpenedVip=1").
		Get(uri)

	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func getVideoClips(originURL, vid, mkey string) (*tvSohuComGetVideoClipsResp, error) {
	uri := fmt.Sprintf("https://hot.vrs.sohu.com/vrs_flash.action?vid=%s&ver=1&ssl=1&pflag=pch5&mkey=%s", vid, mkey)
	var resp tvSohuComGetVideoClipsResp
	_, err := resty.New().R().
		SetHeader("Referer", fmt.Sprintf("%s?user_isOpenedVip=2500", originURL)).SetResult(&resp).
		ForceContentType("application/json").Get(uri)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func getVideoURL(key string) (string, error) {
	uri := fmt.Sprintf("https://data.vod.itc.cn/ip?new=%s&num=1&ch=tv&pt=1&pg=2&prod=h5n", key)
	resp := new(tvSohuComGetVideoURLResp)

	body, err := request.Client.R().Get(uri)
	if err != nil {
		return "", err
	}

	fmt.Println(body.String())

	return resp.Servers[0].URL, nil
}

func getHTMLMeta(uri string) (*tvSohuComHtmlMeta, error) {
	html, err := request.Get(uri, "", nil)
	if err != nil {
		return nil, err
	}

	aid := ""
	vid := ""
	tVid := ""
	if strings.Contains(html, "playlistId") {
		aid = regexp.MustCompile(`playlistId\s*=\s*['"]([^'"]+)`).FindStringSubmatch(string(html))[1]
	}
	if aid == "" && strings.Contains(html, `"aid"`) {
		aid = regexp.MustCompile(`(?U)"aid"\s.*?value=['"]([^'"]+)`).FindStringSubmatch(string(html))[1]
	}
	if strings.Contains(html, " vid") {
		vid = regexp.MustCompile(`\s+vid\s*=\s*['"]([^'"]+)`).FindStringSubmatch(string(html))[1]
	}
	if vid == "" && strings.Contains(html, `"vid"`) {
		vid = regexp.MustCompile(`(?U)"vid"\s.*?value=['"]([^'"]+)`).FindStringSubmatch(string(html))[1]
	}
	if strings.Contains(html, " tvid") {
		tVid = regexp.MustCompile(`tvid\s*=\s*['"]([^'"]*)`).FindStringSubmatch(string(html))[1]
	}
	if tVid == "" && strings.Contains(html, `"tvid"`) {
		tVid = regexp.MustCompile(`(?U)"tvid"\s.*?value=['"]([^'"]*)`).FindStringSubmatch(string(html))[1]
	}

	return &tvSohuComHtmlMeta{
		Vid:  vid,
		Aid:  aid,
		TVid: tVid,
	}, nil
}

type Video struct {
	Title string
	URL   string
	Cover string
}

type tvSohuComHtmlMeta struct {
	Vid  string
	Aid  string
	TVid string
}

type tvSohuComGetVideoClipsResp struct {
	URL   string `json:"url"`
	Tvid  int    `json:"tvid"`
	Syst  int64  `json:"syst"`
	Allot string `json:"allot"`
	Data  struct {
		TvName        string    `json:"tvName"`
		SubName       string    `json:"subName"`
		Ch            string    `json:"ch"`
		Fps           int       `json:"fps"`
		IPLimit       int       `json:"ipLimit"`
		Width         int       `json:"width"`
		ClipsURL      []string  `json:"clipsURL"`
		Version       int       `json:"version"`
		ClipsBytes    []int64   `json:"clipsBytes"`
		Num           int       `json:"num"`
		CoverImg      string    `json:"coverImg"`
		Height        int       `json:"height"`
		TotalDuration float64   `json:"totalDuration"`
		TotalBytes    int64     `json:"totalBytes"`
		ClipsDuration []float64 `json:"clipsDuration"`
		Orifee        int       `json:"orifee"`
		Ck            []string  `json:"ck"`
		Hc            []string  `json:"hc"`
		Su            []string  `json:"su"`
	} `json:"data"`
}

type tvSohuComGetVideoURLResp struct {
	Servers []struct {
		Nid   int    `json:"nid"`
		Isp2P int    `json:"isp2p"`
		URL   string `json:"url"`
	} `json:"servers"`
}

type VideoNewVo struct {
	Title      string
	Duration   int64
	CoverImg   string
	Mp4PlayUrl string
	TotalBytes int64
}

type Mp4PlayUrlVo struct {
	Servers []struct {
		Nid   int    `json:"nid"`
		Isp2P int    `json:"isp2p"`
		Url   string `json:"url"`
	} `json:"servers"`
}
