package kuaishou

import (
	"github.com/iawia002/lux/utils"
	"github.com/pkg/errors"
	"regexp"
	"strings"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
)

func init() {
	extractors.Register("kuaishou", New())
}

type extractor struct{}

// New returns a kuaishou extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// fetch url and get the cookie that write by server
func fetchCookies(url string, headers map[string]string) (string, error) {
	res, err := request.Client.R().SetHeaders(headers).Get(url)
	if err != nil {
		return "", err
	}

	cookiesArr := make([]string, 0)
	cookies := res.Cookies()

	for _, c := range cookies {
		cookiesArr = append(cookiesArr, c.Name+"="+c.Value)
	}

	return strings.Join(cookiesArr, "; "), nil
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0",
	}

	headers["Cookie"] = option.Cookie

	html, err := request.Get(url, url, headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	titles := utils.MatchOneOf(html, `<title>([^<]+)</title>`)
	if titles == nil || len(titles) < 2 {
		return nil, errors.New("can not found title")
	}

	title := regexp.MustCompile(`\n+`).ReplaceAllString(strings.TrimSpace(titles[1]), " ")

	qualityRegMap := map[string]*regexp.Regexp{
		"sd": regexp.MustCompile(`"photoUrl":\s*"([^"]+)"`),
	}

	streams := make(map[string]*extractors.Stream, 1)
	for quality, qualityReg := range qualityRegMap {
		matcher := qualityReg.FindStringSubmatch(html)
		if len(matcher) != 2 {
			return nil, errors.WithStack(extractors.ErrURLParseFailed)
		}

		u := strings.ReplaceAll(matcher[1], `\u002F`, "/")

		size, err := request.Size(u, url)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		urlData := &extractors.Part{
			URL:  u,
			Size: size,
			Ext:  "mp4",
		}
		streams[quality] = &extractors.Stream{
			Parts:   []*extractors.Part{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	return []*extractors.Data{
		{
			Site:    "快手 kuaishou.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
		},
	}, nil
}

//func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
//	client := resty.New()
//	client.SetRedirectPolicy(resty.NoRedirectPolicy())
//	res, _ := client.R().
//		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").
//		Get(url)
//	//这里会返回err, auto redirect is disabled
//
//	locationRes, err := res.RawResponse.Location()
//	if err != nil {
//		return nil, err
//	}
//
//	// 分享的中间跳转链接不太一样, 有些是 /fw/long-video , 有些 /fw/photo
//	referUri := strings.ReplaceAll(locationRes.String(), "v.m.chenzhongtech.com/fw/long-video", "m.gifshow.com/fw/photo")
//	referUri = strings.ReplaceAll(referUri, "v.m.chenzhongtech.com/fw/photo", "m.gifshow.com/fw/photo")
//
//	videoId := strings.ReplaceAll(strings.Trim(locationRes.Path, "/"), "fw/long-video/", "")
//	videoId = strings.ReplaceAll(videoId, "fw/photo/", "")
//	videoId = strings.ReplaceAll(videoId, "short-video/", "")
//	if len(videoId) <= 0 {
//		return nil, errors.New("parse video id from share url fail")
//	}
//
//	postData := map[string]interface{}{
//		"photoId":     videoId,
//		"isLongVideo": false,
//	}
//	videoRes, err := client.R().
//		SetHeader("cookie", option.Cookie).
//		SetHeader("Origin", "https://m.gifshow.com").
//		SetHeader("referer", strings.ReplaceAll(referUri, "m.gifshow.com/fw/photo", "m.gifshow.com/fw/photo")).
//		SetHeader("content-type", "application/json").
//		SetHeader("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").
//		SetBody(postData).
//		Post("https://m.gifshow.com/rest/wd/photo/info?kpn=KUAISHOU&captchaToken=")
//
//	data := gjson.GetBytes(videoRes.Body(), "photo")
//	title := data.Get("caption").String()
//	videoUrl := data.Get("mainMvUrls.0.url").String()
//	cover := data.Get("coverUrls.0.url").String()
//
//	streams := make(map[string]*extractors.Stream, 1)
//
//	size, err := request.Size(videoUrl, strings.ReplaceAll(referUri, "m.gifshow.com/fw/photo", "m.gifshow.com/fw/photo"))
//	if err != nil {
//		return nil, err
//	}
//
//	urlData := &extractors.Part{
//		URL:  videoUrl,
//		Size: size,
//		Ext:  "mp4",
//	}
//
//	streams["sd"] = &extractors.Stream{
//		Parts:   []*extractors.Part{urlData},
//		Size:    size,
//		Quality: "sd",
//	}
//
//	return []*extractors.Data{
//		{
//			Site:    "快手 kuaishou.com",
//			Title:   title,
//			Type:    extractors.DataTypeVideo,
//			Streams: streams,
//			URL:     url,
//			Cover:   cover,
//		},
//	}, nil
//}
