package qq

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"math/rand"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/iawia002/lia/array"
	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"

	lop "github.com/samber/lo/parallel"
)

func init() {
	extractors.Register("qq", New())
}

type qqVideoInfo struct {
	Fl struct {
		Fi []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Cname string `json:"cname"`
			Fs    int64  `json:"fs"`
			Br    int    `json:"br"`
		} `json:"fi"`
	} `json:"fl"`
	Vl struct {
		Vi []struct {
			Fn    string `json:"fn"`
			Ti    string `json:"ti"`
			Fvkey string `json:"fvkey"`
			Cl    struct {
				Fc int `json:"fc"`
				Ci []struct {
					Idx int `json:"idx"`
				} `json:"ci"`
			} `json:"cl"`
			Br int `json:"br"`
			Ul struct {
				UI []struct {
					URL string `json:"url"`
					Hls struct {
						Et    int         `json:"et"`
						Fbw   int         `json:"fbw"`
						Ftype string      `json:"ftype"`
						Hk    string      `json:"hk"`
						Hvl   interface{} `json:"hvl"`
						Pnl   struct {
							Pi []struct {
								Bw int    `json:"bw"`
								Fc int    `json:"fc"`
								Fn string `json:"fn"`
							} `json:"pi"`
						} `json:"pnl"`
						St    int    `json:"st"`
						Stype string `json:"stype"`
						Pname string `json:"pname"`
						Pt    string `json:"pt"`
					} `json:"hls"`
				} `json:"ui"`
			} `json:"ul"`
		} `json:"vi"`
	} `json:"vl"`
	Msg string `json:"msg"`
}

type qqKeyInfo struct {
	Key string `json:"key"`
}

const qqPlayerVersion string = "3.2.19.333"

func getVinfo(vid, defn, refer string) (qqVideoInfo, error) {
	html, err := request.Get(
		fmt.Sprintf(
			"http://vv.video.qq.com/getinfo?otype=json&platform=11&defnpayver=1&appver=%s&defn=%s&vid=%s",
			qqPlayerVersion, defn, vid,
		), refer, nil,
	)
	if err != nil {
		return qqVideoInfo{}, err
	}
	jsonStrings := utils.MatchOneOf(html, `QZOutputJson=(.+);$`)
	if jsonStrings == nil || len(jsonStrings) < 2 {
		return qqVideoInfo{}, errors.WithStack(extractors.ErrURLParseFailed)
	}
	jsonString := jsonStrings[1]
	var data qqVideoInfo
	if err = json.Unmarshal([]byte(jsonString), &data); err != nil {
		return qqVideoInfo{}, err
	}
	return data, nil
}

func getVideoApiKey(vid, videoUrl, seriesId, subtitleFormat, videoFormat, videoQuality string) (qqVideoInfo, error) {
	guid := getGuid()
	cKey := getCKey(vid, videoUrl, guid)
	query := map[string]string{
		"vid":           vid,
		"cid":           seriesId,
		"cKey":          cKey,
		"encryptVer":    "8.1",
		"spcaptiontype": lo.If(subtitleFormat == "vtt", "1").Else("0"),
		"sphls":         lo.If(videoFormat == "hls", "2").Else("0"),
		"dtype":         lo.If(videoFormat == "hls", "3").Else("0"),
		"defn":          videoQuality,
		"spsrt":         "2",
		"sphttps":       "1",
		"otype":         "json",
		"spwm":          "1",
		"hevclv":        "28",
		"drm":           "40",
		"spvideo":       "4",
		"spsfrhdr":      "100",
		"host":          "v.qq.com",
		"referer":       "v.qq.com",
		"ehost":         videoUrl,
		"appVer":        "3.5.57",
		"platform":      "10901",
		"guid":          guid,
		"flowid":        getGuid(),
	}

	body, err := request.Client.R().SetQueryParams(query).Get(`https://h5vv6.video.qq.com/getvinfo`)
	if err != nil {
		return qqVideoInfo{}, err
	}

	jsonStrings := utils.MatchOneOf(body.String(), `QZOutputJson=(.+);$`)
	if jsonStrings == nil || len(jsonStrings) < 2 {
		return qqVideoInfo{}, errors.WithStack(extractors.ErrURLParseFailed)
	}
	jsonString := jsonStrings[1]
	var data qqVideoInfo
	if err = json.Unmarshal([]byte(jsonString), &data); err != nil {
		return qqVideoInfo{}, err
	}
	return data, nil
}

func getCKey(videoId, url, guid string) string {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"

	var dd string
	if len(url) >= 48 {
		dd = url[:48]
	} else {
		dd = url
	}

	payload := fmt.Sprintf("%s|%d|mg3c3b04ba|%s|%s|%s|%s|%s||Mozilla|Netscape|Windows x86_64|00|", videoId, time.Now().Unix(), "3.5.57", guid, "10901", dd, ua[:48])

	plaintext := []byte(fmt.Sprintf("|%d|%s", len(payload), payload))

	block, err := aes.NewCipher([]byte{0x4f, 0x6b, 0xda, 0xa3, 0x9e, 0x2f, 0x8c, 0xb0, 0x7f, 0x5e, 0x72, 0x2d, 0x9e, 0xde, 0xf3, 0x14})
	if err != nil {
		panic(err)
	}

	iv := []byte{0x01, 0x50, 0x4a, 0xf3, 0x56, 0xe6, 0x19, 0xcf, 0x2e, 0x42, 0xbb, 0xa6, 0x8c, 0x3f, 0x70, 0xf9}

	// 对 plaintext 进行补全
	bs := block.BlockSize()
	plaintext = append(plaintext, bytes.Repeat([]byte{byte(bs - len(plaintext)%bs)}, bs-len(plaintext)%bs)...)

	ciphertext := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return strings.ToUpper(fmt.Sprintf("%x", ciphertext))
}

func getGuid() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func genStreams(vid, videoUrl string, data qqVideoInfo) (map[string]*extractors.Stream, error) {
	streams := make(map[string]*extractors.Stream)

	var qualities []string

	if len(data.Fl.Fi) == 0 {
		qualities = []string{"shd", "fhd"}
	} else {
		for _, item := range data.Fl.Fi {
			qualities = append(qualities, item.Name)
		}
	}

	apiResponses := make([]*qqVideoInfo, 0)
	apiResponses = append(apiResponses, &data)

	for _, q := range qualities {
		if !array.ItemInArray(q, []string{"ld", "sd", "hd"}) {
			stream, err := getVideoApiKey(vid, videoUrl, "", "vtt", "hls", q)
			if err != nil {
				fmt.Println(err)
			}
			apiResponses = append(apiResponses, &stream)
		}
	}

	for _, apiResponse := range apiResponses {
		videoResponse := apiResponse.Vl.Vi[0]
		for _, videoFormat := range videoResponse.Ul.UI {
			var parts []*extractors.Part
			var totalSize int64
			if videoFormat.Hls.Pt != "" || determineExt(videoFormat.URL) == "m3u8" {
				tsUrls, err := utils.M3u8URLs(videoFormat.URL + videoFormat.Hls.Pt)
				if err != nil {
					return nil, err
				}

				parts = lop.Map(tsUrls, func(tsUrl string, _ int) *extractors.Part {
					size, _ := request.Size(tsUrl, "")
					return &extractors.Part{
						URL:  tsUrl,
						Ext:  "ts",
						Size: size,
					}
				})
			}

			totalSize = lo.SumBy(parts, func(part *extractors.Part) int64 {
				return part.Size
			})

			identifier := videoResponse.Br
			var formatResponse string
			for _, item := range apiResponse.Fl.Fi {
				if item.Br == identifier {
					formatResponse = item.Name
				}
			}

			streams[formatResponse] = &extractors.Stream{
				Parts:   parts,
				ID:      formatResponse,
				Quality: formatResponse,
				Size:    totalSize,
			}
		}
	}
	return streams, nil
}

var KnownExtensions = map[string]bool{
	"3gp": true, "avi": true, "flv": true, "m4v": true, "mkv": true, "mov": true,
	"mp4": true, "mpeg": true, "mpg": true, "webm": true, "wmv": true,
}

func determineExt(urlStr string) string {
	if urlStr == "" || !strings.Contains(urlStr, ".") {
		return ""
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	path := u.Path
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	if ext != "" {
		return ext
	}

	guess := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if match, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, guess); match {
		return guess
	} else if KnownExtensions[strings.ToLower(guess)] {
		return strings.ToLower(guess)
	}

	return ""
}

type extractor struct{}

// New returns a qq extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(uri string, option extractors.Options) ([]*extractors.Data, error) {
	//vids := utils.MatchOneOf(uri, `vid=(\w+)`, `/(\w+)\.html`)
	//if vids == nil || len(vids) < 2 {
	//	return nil, errors.WithStack(extractors.ErrURLParseFailed)
	//}
	//vid := vids[1]

	if option.Playlist {
		return extractPlaylist(uri)
	}

	var xx []*extractors.Data
	//if len(vid) != 11 {
	u, err := request.Get(uri, uri, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	vids := utils.MatchOneOf(
		u, `vid=(\w+)`, `vid:\s*["'](\w+)`, `vid\s*=\s*["']\s*(\w+)`,
	)
	if vids == nil || len(vids) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}
	vid := vids[1]

	more := utils.MatchOneOf(u, `window.__pinia=(.*?)</script>`)
	moreJson := strings.ReplaceAll(more[1], "undefined", `""`)
	if len(more) > 0 {
		seriesJson := gjson.Parse(moreJson).Get("episodeMain.listData.0.list.0")
		seriesJson.ForEach(func(key, value gjson.Result) bool {
			addr, _ := url.JoinPath(`https://v.qq.com/x/cover`, value.Get("cid").String(), value.Get("vid").String())

			pic := lo.If(value.Get("pic").String() != "", value.Get("pic").String()).Else(value.Get("picVertial").String())

			xx = append(xx, &extractors.Data{
				Title: value.Get("playTitle").String(),
				URL:   addr,
				Site:  "腾讯视频 v.qq.com",
				Type:  extractors.DataTypeVideo,
				Cover: pic,
			})
			return true
		})
	}
	//}

	data, err := getVideoApiKey(vid, uri, "", "srt", "hls", "hd")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// API request error
	if data.Msg != "" {
		return nil, errors.New(data.Msg)
	}
	streams, err := genStreams(vid, uri, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pic := gjson.Parse(moreJson).Get("epiosode.currentVideoInfo.pic").String()

	return []*extractors.Data{
		{
			Site:    "腾讯视频 v.qq.com",
			Title:   data.Vl.Vi[0].Ti,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     uri,
			Series:  xx,
			Cover:   pic,
		},
	}, nil
}

func extractPlaylist(uri string) ([]*extractors.Data, error) {
	vcuid := utils.MatchOneOf(uri, `vcuid=(\w+)`)

	if vcuid == nil || len(vcuid) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	uri = fmt.Sprintf(`https://pbaccess.video.qq.com/trpc.creator_center.header_page.personal_page/GetUserVideoList?vcuid=%s&page_size=30&list_type=1&page=1`, vcuid[1])
	body, err := request.Client.R().Get(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*extractors.Data, 0)
	gjson.ParseBytes(body.Body()).Get("data.data.list").ForEach(func(key, value gjson.Result) bool {
		addr, _ := url.JoinPath("https://v.qq.com/x/page", value.Get("vid").String()+".html")
		result = append(result, &extractors.Data{
			Title: value.Get("title").String(),
			URL:   addr,
			Site:  "腾讯视频 v.qq.com",
			Type:  extractors.DataTypeVideo,
			Cover: value.Get("image").String(),
		})
		return true
	})

	return result, nil
}
