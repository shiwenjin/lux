package iqiyi

import (
	"encoding/json"
	"fmt"
	"github.com/iawia002/lia/array"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/parser"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
)

func init() {
	extractors.Register("iqiyi", New(SiteTypeIqiyi))
	extractors.Register("iq", New(SiteTypeIQ))
}

var header = map[string]string{
	"Accept":          "*/*",
	"Accept-Encoding": "*",
	"Accept-Language": "zh-CN,zh;q=0.9",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
	"Host":            "iqiyihao.iqiyi.com",
	"Connection":      "keep-alive",
}

// SiteType indicates the site type of iqiyi
type SiteType int

const (
	// SiteTypeIQ indicates the site is iq.com
	SiteTypeIQ SiteType = iota
	// SiteTypeIqiyi indicates the site is iqiyi.com
	SiteTypeIqiyi
	iqReferer    = "https://www.iq.com"
	iqiyiReferer = "https://www.iqiyi.com"
)

const (
	// getVideos 获取视频id列表
	getVideos = `https://iqiyihao.iqiyi.com/iqiyihao/entity/get_videos.action?agenttype=118&agentversion=10.7.5&authcookie=&dfp=%s&fuid=%s&m_device_id=cv2irlndqb0opl8fgsydst7q&page=%d&sign=%s&size=%s&timestamp=%d`
	// episodeInfo 获取所有视频详情
	videoEpisodeInfo = `https://iqiyihao.iqiyi.com/iqiyihao/episode_info.action?agenttype=118&agentversion=10.7.5&authcookie=&dfp=%s&m_device_id=cv2irlndqb0opl8fgsydst7q&qipuIds=%s&sign=%s&timestamp=%d`
	//getSmallVideos 获取小视频id列表
	getSmallVideos = `https://iqiyihao.iqiyi.com/iqiyihao/entity/get_small_videos.action?authcookie=&agenttype=119&agentversion=9.12.0&timestamp=%d&m_device_id=cv2irlndqb0opl8fgsydst7q&uid=%s&page=%d&size=%s&sign=%s`
	// smallVideoEpisodeInfo 获取所有小视频详情
	smallVideoEpisodeInfo = `https://iqiyihao.iqiyi.com/iqiyihao/episode_info.action?authcookie=&agenttype=119&agentversion=9.12.0&timestamp=%d&m_device_id=cv2irlndqb0opl8fgsydst7q&qipuIds=%s&sign=%s`
)

var (
	size = "20"
	//dfp 请求携带参数之一   如果失效 网页重新找cookie进行替换
	dfp         = "a17dd85b225feb497a8774f10823fa05300b6abd36e828d37a4ae009837d7b7176"
	currentPage = 1
	timeStamp   = time.Now().UnixMilli()
)

func getMacID() string {
	var macID string
	chars := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "n", "m", "o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	size := len(chars)
	for i := 0; i < 32; i++ {
		macID += chars[rand.Intn(size)]
	}
	return macID
}

func getVF(params string) string {
	var suffix string
	for j := 0; j < 8; j++ {
		for k := 0; k < 4; k++ {
			var v8 int
			v4 := 13 * (66*k + 27*j) % 35
			if v4 >= 10 {
				v8 = v4 + 88
			} else {
				v8 = v4 + 49
			}
			suffix += string(rune(v8)) // string(97) -> "a"
		}
	}
	params += suffix

	return utils.Md5(params)
}

func getVPS(tvid, vid, refer string) (*iqiyi, error) {
	t := time.Now().Unix() * 1000
	host := "https://cache.video.iqiyi.com"
	params := fmt.Sprintf(
		"/vps?tvid=%s&vid=%s&v=0&qypid=%s_12&src=01012001010000000000&t=%d&k_tag=1&k_uid=%s&rs=1",
		tvid, vid, tvid, t, getMacID(),
	)
	vf := getVF(params)
	apiURL := fmt.Sprintf("%s%s&vf=%s", host, params, vf)
	info, err := request.Get(apiURL, refer, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data := new(iqiyi)
	if err := json.Unmarshal([]byte(info), data); err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

type extractor struct {
	siteType SiteType
}

// New returns a iqiyi extractor.
func New(siteType SiteType) extractors.Extractor {
	return &extractor{
		siteType: siteType,
	}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	result := make([]*extractors.Data, 0)
	streams := make(map[string]*extractors.Stream)

	if option.Playlist {
		videos, err := extractPlaylist(url, false)
		if err != nil {
			return nil, err
		}
		needDownloadItems := utils.NeedDownloadList(option.Items, option.ItemStart, option.ItemEnd, len(videos))

		result = lo.Filter(videos, func(_ *extractors.Data, index int) bool {
			return array.ItemInArray(index+1, needDownloadItems)
		})

		return result, nil
	}

	if option.ShortPlaylist {
		videos, err := extractPlaylist(url, true)
		if err != nil {
			return nil, err
		}
		needDownloadItems := utils.NeedDownloadList(option.Items, option.ItemStart, option.ItemEnd, len(videos))

		result = lo.Filter(videos, func(_ *extractors.Data, index int) bool {
			return array.ItemInArray(index+1, needDownloadItems)
		})

		return result, nil
	}

	refer := iqiyiReferer
	headers := make(map[string]string)
	if e.siteType == SiteTypeIQ {
		headers = map[string]string{
			"Accept-Language": "zh-TW",
		}
		refer = iqReferer
	}
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0"
	html, err := request.Get(url, refer, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tvid := utils.MatchOneOf(
		url,
		`#curid=(.+)_`,
		`tvid=([^&]+)`,
	)
	if tvid == nil {
		tvid = utils.MatchOneOf(
			html,
			`data-player-tvid="([^"]+)"`,
			`param\['tvid'\]\s*=\s*"(.+?)"`,
			`"tvid":"(\d+)"`,
			`"tvId":(\d+)`,
		)
	}
	if tvid == nil || len(tvid) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	vid := utils.MatchOneOf(
		url,
		`#curid=.+_(.*)$`,
		`vid=([^&]+)`,
	)
	if vid == nil {
		vid = utils.MatchOneOf(
			html,
			`data-player-videoid="([^"]+)"`,
			`param\['vid'\]\s*=\s*"(.+?)"`,
			`"vid":"(\w+)"`,
		)
	}
	if vid == nil || len(vid) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	doc, err := parser.GetDoc(html)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cover := queryCoverFromDoc(html)

	var title string
	if e.siteType == SiteTypeIqiyi {
		title = strings.TrimSpace(doc.Find("h1>a").First().Text())
		var sub string
		for _, k := range []string{"span", "em"} {
			if sub != "" {
				break
			}
			sub = strings.TrimSpace(doc.Find("h1>" + k).First().Text())
		}
		title += sub
	} else {
		title = strings.TrimSpace(doc.Find("span#pageMetaTitle").First().Text())
		sub := utils.MatchOneOf(html, `"subTitle":"([^"]+)","isoDuration":`)
		if len(sub) > 1 {
			title += fmt.Sprintf(" %s", sub[1])
		}
	}
	if title == "" {
		title = doc.Find("title").Text()
	}
	videoDatas, err := getVPS(tvid[1], vid[1], refer)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if videoDatas.Code != "A00000" {
		return nil, errors.Errorf("can't play this video: %s", videoDatas.Msg)
	}

	urlPrefix := videoDatas.Data.VP.Du
	for _, video := range videoDatas.Data.VP.Tkl[0].Vs {
		urls := make([]*extractors.Part, len(video.Fs))
		for index, v := range video.Fs {
			realURLData, err := request.Get(urlPrefix+v.L, refer, nil)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			var realURL iqiyiURL
			if err = json.Unmarshal([]byte(realURLData), &realURL); err != nil {
				return nil, errors.WithStack(err)
			}
			_, ext, err := utils.GetNameAndExt(realURL.L)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			urls[index] = &extractors.Part{
				URL:  realURL.L,
				Size: v.B,
				Ext:  ext,
			}
		}
		streams[strconv.Itoa(video.Bid)] = &extractors.Stream{
			Parts:   urls,
			Size:    video.Vsize,
			Quality: video.Scrsz,
		}
	}

	siteName := "爱奇艺 iqiyi.com"
	if e.siteType == SiteTypeIQ {
		siteName = "爱奇艺 iq.com"
	}
	return []*extractors.Data{
		{
			Site:    siteName,
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
			Cover:   cover,
		},
	}, nil
}

func queryCoverFromDoc(doc string) string {
	var cover string
	var err error
	coverArr := utils.MatchOneOf(doc, `<meta[^>]+property="og:image"\s+content="([^"]+)`, `<meta[^>]+property="og:image"\s+content="([^"]+)`)
	if len(coverArr) > 0 {
		cover, _ = url.JoinPath("https://", coverArr[1])
		cover = strings.ReplaceAll(cover, ".jpg", "_1920_1080.jpg")

		_, err = request.Headers(cover, "")
		if err != nil {
			cover = strings.ReplaceAll(cover, "_1920_1080.jpg", "_1280_720.jpg")
			_, err = request.Headers(cover, "")
			if err != nil {
				cover = strings.ReplaceAll(cover, "_1280_720.jpg", "_480_270.jpg")
				_, err = request.Headers(cover, "")
				if err != nil {
					cover = utils.MatchOneOf(doc, `<meta\s+content="([^"]+)"\s+property="og:image">`)[1]
				}
			}
		}
	}
	return cover
}

// https://www.iqiyi.com/u/1412822955/videos
func extractPlaylist(uri string, shortPlaylist bool) ([]*extractors.Data, error) {
	uid := utils.MatchOneOf(uri, `/u/(\d+)/`)[1]

	var vInfo *GetVideosInfo
	var err error
	if !shortPlaylist {
		//获取视频的id列表
		ids, err := getVidList(uid)
		if err != nil {
			return nil, err
		}

		//获取视频的详情
		vInfo, err = getAllVideoInfo(uid, ids)
		if err != nil {
			return nil, err
		}
	} else {
		//获取视频的id列表
		ids, err := getSmallVidList(uid)
		if err != nil {
			return nil, err
		}

		//获取视频的详情
		vInfo, err = getAllSmallVideoInfo(uid, ids)
		if err != nil {
			return nil, err
		}
	}

	streams := make(map[string]*extractors.Stream)
	defaultStream := extractors.Stream{
		Ext: "mp4",
	}

	result := make([]*extractors.Data, 0)

	for _, video := range vInfo.Data {
		streams["default"] = &defaultStream
		result = append(result, &extractors.Data{
			Site:    "爱奇艺 iqiyi",
			URL:     video.PageUrl,
			Title:   video.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			Cover:   video.VerticalThumbnail,
		})
	}
	return result, err
}

func getAllVideoInfo(uid, qipuIds string) (*GetVideosInfo, error) {
	sign, err := getSign(uid, videoEpisodeInfo)
	if err != nil {
		return nil, err
	}
	requestUrl := fmt.Sprintf(videoEpisodeInfo, dfp, qipuIds, sign, timeStamp)
	vinfoResp := new(GetVideosInfo)
	if _, err := request.Client.R().SetHeaders(header).SetResult(vinfoResp).Get(requestUrl); err != nil {
		return nil, err
	}
	return vinfoResp, nil
}

func getVidList(uid string) (string, error) {
	sign, err := getSign(uid, getVideos)
	if err != nil {
		return "", err
	}
	requestUrl := fmt.Sprintf(getVideos, dfp, uid, currentPage, sign, size, timeStamp)
	idsResp := new(GetVideosId)
	if _, err := request.Client.R().SetHeaders(header).SetResult(idsResp).Get(requestUrl); err != nil {
		return "", err
	}
	idList := idsResp.Data.Sort.Flows
	//处理idList 拼接成字符串
	var ids []string
	for _, vid := range idList {
		id := cast.ToString(vid.QipuId)
		if id != "" {
			ids = append(ids, cast.ToString(vid.QipuId))
		}
	}
	return strings.Join(ids, ","), err
}

func getSign(uid, action string) (sign string, err error) {
	var requestUrl string
	requestUrl = action
	requestUrl = "GET" + strings.TrimLeft(requestUrl, "https://")
	requestUrl = strings.ReplaceAll(requestUrl, "&sign=%s", "")

	switch action {

	case getVideos:
		requestUrl = fmt.Sprintf(requestUrl, dfp, uid, currentPage, size, timeStamp)
		requestUrl, err = urlKeySort(requestUrl)
		if err != nil {
			return "", err
		}
		requestUrl = requestUrl + "NZrFGv72GYppTUxO"

	case videoEpisodeInfo:
		qipuIds, err := getVidList(uid)
		if err != nil {
			return "", err
		}
		requestUrl = fmt.Sprintf(requestUrl, dfp, qipuIds, timeStamp)
		requestUrl, err = urlKeySort(requestUrl)
		if err != nil {
			return "", err
		}
		requestUrl = requestUrl + "NZrFGv72GYppTUxO"

	case getSmallVideos:
		requestUrl = fmt.Sprintf(requestUrl, timeStamp, uid, currentPage, size)
		requestUrl, err = urlKeySort(requestUrl)
		if err != nil {
			return "", err
		}
		requestUrl = requestUrl + "QMK8e4agKNWEppKU"

	case smallVideoEpisodeInfo:
		qipuIds, err := getSmallVidList(uid)
		if err != nil {
			return "", err
		}
		requestUrl = fmt.Sprintf(requestUrl, timeStamp, qipuIds)
		requestUrl, err = urlKeySort(requestUrl)
		if err != nil {
			return "", err
		}
		requestUrl = requestUrl + "QMK8e4agKNWEppKU"

	default:
		return "", errors.New("没有对应链接")
	}

	return utils.Md5(requestUrl), err
}

func getSmallVidList(uid string) (string, error) {
	sign, err := getSign(uid, getSmallVideos)
	if err != nil {
		return "", err
	}
	requestUrl := fmt.Sprintf(getSmallVideos, timeStamp, uid, currentPage, size, sign)
	vidsResp := new(SmallVideoResp)
	if _, err := request.Client.R().SetHeaders(header).SetResult(vidsResp).Get(requestUrl); err != nil {
		return "", err
	}
	vidList := make([]int64, 0)
	vidList = vidsResp.Data.Tvids
	//处理vidList 拼接成字符串
	var ids string
	for _, vid := range vidList {
		ids = ids + cast.ToString(vid) + ","
	}
	ids = strings.TrimRight(ids, ",")
	return ids, nil
}

func getAllSmallVideoInfo(uid, qipuIds string) (videoInfo *GetVideosInfo, err error) {
	sign, err := getSign(uid, smallVideoEpisodeInfo)
	if err != nil {
		return nil, err
	}
	//对字符串进行转义
	qipuIds = url.QueryEscape(qipuIds)
	requestUrl := fmt.Sprintf(smallVideoEpisodeInfo, timeStamp, qipuIds, sign)
	vinfoResp := new(GetVideosInfo)
	if _, err := request.Client.R().SetHeaders(header).SetResult(vinfoResp).Get(requestUrl); err != nil {
		return nil, err
	}
	return vinfoResp, nil
}

// 下面的小写方法提供给上面导出方法使用
func urlKeySort(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", errors.New("解析url出错")
	}
	//得到map
	queryParams := u.Query()
	//存key
	keys := make([]string, 0, len(queryParams))
	for key := range queryParams {
		keys = append(keys, key)
	}
	//对key排序
	sort.Strings(keys)
	//重新组装params
	params := ""
	for _, key := range keys {
		params = params + fmt.Sprintf("%s=%s", key, queryParams[key][0]) + "&"
	}
	params = strings.TrimRight(params, "&")
	//返回url
	urlString = u.Path + "?" + params
	return urlString, err
}
