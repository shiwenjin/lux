package youku

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
	"github.com/iawia002/lux/utils"
)

func init() {
	extractors.Register("youku", New())
}

type errorData struct {
	Note string `json:"note"`
	Code int    `json:"code"`
}

type segs struct {
	Size int64  `json:"size"`
	URL  string `json:"cdn_url"`
}

type stream struct {
	Size      int64  `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Segs      []segs `json:"segs"`
	Type      string `json:"stream_type"`
	AudioLang string `json:"audio_lang"`
}

type youkuVideo struct {
	Title string `json:"title"`
	Logo  string `json:"logo"`
}

type youkuShow struct {
	Title string `json:"title"`
}

type data struct {
	Error  errorData  `json:"error"`
	Stream []stream   `json:"stream"`
	Video  youkuVideo `json:"video"`
	Show   youkuShow  `json:"show"`
}

type youkuData struct {
	Data data `json:"data"`
}

const youkuReferer = "https://v.youku.com"

func getAudioLang(lang string) string {
	var youkuAudioLang = map[string]string{
		"guoyu": "国语",
		"ja":    "日语",
		"yue":   "粤语",
	}
	translate, ok := youkuAudioLang[lang]
	if !ok {
		return lang
	}
	return translate
}

func youkuUpsNew(uri string) (*youkuData, error) {
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}

	params, err := GetJoinParams(uri, cookie, cast.ToString(time.Now().UnixMilli()))
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	j, _ := json.Marshal(params)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}

	body, err := request.Client.R().SetQueryParams(m).SetHeaders(map[string]string{
		"Accept":     "*/*",
		"Host":       "acs.youku.com",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
		"cookie":     cookie,
		"Referer":    "https://v.youku.com/",
	}).Get(`https://acs.youku.com/h5/mtop.youku.play.ups.appinfo.get/1.1/`)

	realJson := gjson.ParseBytes(body.Body()).Get("data").Raw

	data := youkuData{}
	if err = json.Unmarshal([]byte(realJson), &data); err != nil {
		return nil, errors.WithStack(err)
	}
	if data.Data.Error == (errorData{}) {
		return &data, nil
	}

	return &data, nil
}

// https://g.alicdn.com/player/ykplayer/0.5.61/youku-player.min.js
// {"0505":"interior","050F":"interior","0501":"interior","0502":"interior","0503":"interior","0510":"adshow","0512":"BDskin","0590":"BDskin"}

// var ccodes = []string{"0510", "0502", "0507", "0508", "0512", "0513", "0514", "0503", "0590"}

func youkuUps(vid string, option extractors.Options) (*youkuData, error) {
	var (
		url   string
		utid  string
		utids []string
		data  youkuData
	)
	if strings.Contains(option.Cookie, "cna") {
		utids = utils.MatchOneOf(option.Cookie, `cna=(.+?);`, `cna\s+(.+?)\s`, `cna\s+(.+?)$`)
	} else {
		headers, err := request.Headers("http://log.mmstat.com/eg.js", youkuReferer)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		setCookie := headers.Get("Set-Cookie")
		utids = utils.MatchOneOf(setCookie, `cna=(.+?);`)
	}
	if utids == nil || len(utids) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}
	utid = utids[1]

	// https://g.alicdn.com/player/ykplayer/0.5.61/youku-player.min.js
	// grep -oE '"[0-9a-zA-Z+/=]{256}"' youku-player.min.js
	if option.YoukuPassword != "" {
		url = fmt.Sprintf("%s&password=%s", url, option.YoukuPassword)
	}

	// data must be emptied before reassignment, otherwise it will contain the previous value(the 'error' data)
	body, err := request.Client.R().ForceContentType("application/json").SetHeader("Referer", `https://ups.youku.com/ups/get.json`).
		SetHeader("Cookie", "__ysuid="+getYSUID()).
		SetHeader("xreferrer", "http://www.youku.com").
		SetQueryParams(map[string]string{
			"vid":       vid,
			"ccode":     "0524",
			"client_ip": "192.168.1.1",
			"utid":      utid,
			"client_ts": cast.ToString(time.Now().Unix()),
		}).
		SetResult(&data).Get(`https://ups.youku.com/ups/get.json`)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if body.StatusCode() == 200 && data.Data.Error == (errorData{}) {
		return &data, nil
	}
	return &data, nil
}

func getYSUID() string {
	t := time.Now().Unix()
	rand.Seed(t)
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%d%s", t, string(b))
}

func getBytes(val int32) []byte {
	var buff bytes.Buffer
	binary.Write(&buff, binary.BigEndian, val) // nolint
	return buff.Bytes()
}

func hashCode(s string) int32 {
	var result int32
	for _, c := range s {
		result = result*0x1f + c
	}
	return result
}

func hmacSha1(key []byte, msg []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(msg) // nolint
	return mac.Sum(nil)
}

func generateUtdid() string {
	timestamp := int32(time.Now().Unix())
	var buffer bytes.Buffer
	buffer.Write(getBytes(timestamp - 60*60*8))
	buffer.Write(getBytes(rand.Int31()))
	buffer.WriteByte(0x03)
	buffer.WriteByte(0x00)
	imei := fmt.Sprintf("%d", rand.Int31())
	buffer.Write(getBytes(hashCode(imei)))
	data := hmacSha1([]byte("d6fc3a4a06adbde89223bvefedc24fecde188aaa9161"), buffer.Bytes())
	buffer.Write(getBytes(hashCode(base64.StdEncoding.EncodeToString(data))))
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func genData(youkuData data) map[string]*extractors.Stream {
	var (
		streamString string
		quality      string
	)
	streams := make(map[string]*extractors.Stream, len(youkuData.Stream))
	for _, stream := range youkuData.Stream {
		if stream.AudioLang == "default" {
			streamString = stream.Type
			quality = fmt.Sprintf(
				"%s %dx%d", stream.Type, stream.Width, stream.Height,
			)
		} else {
			streamString = fmt.Sprintf("%s-%s", stream.Type, stream.AudioLang)
			quality = fmt.Sprintf(
				"%s %dx%d %s", stream.Type, stream.Width, stream.Height,
				getAudioLang(stream.AudioLang),
			)
		}

		ext := strings.Split(
			strings.Split(stream.Segs[0].URL, "?")[0],
			".",
		)
		urls := make([]*extractors.Part, len(stream.Segs))
		for index, data := range stream.Segs {
			urls[index] = &extractors.Part{
				URL:  data.URL,
				Size: data.Size,
				Ext:  ext[len(ext)-1],
			}
		}
		streams[streamString] = &extractors.Stream{
			Parts:   urls,
			Size:    stream.Size,
			Quality: quality,
		}
	}
	return streams
}

type extractor struct{}

// New returns a youku extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string, option extractors.Options) ([]*extractors.Data, error) {
	if option.Playlist {
		return extractPlaylist(url)
	}

	vids := utils.MatchOneOf(
		url, `id_(.+?)\.html`, `id_(.+)`,
	)
	if vids == nil || len(vids) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}
	vid := vids[1]

	var youkuData *youkuData
	var err error

	youkuData, err = youkuUps(vid, option)
	if len(youkuData.Data.Stream) == 0 {
		youkuData, err = youkuUpsNew(url)
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}
	if youkuData.Data.Error.Code != 0 {
		return nil, errors.New(youkuData.Data.Error.Note)
	}
	streams := genData(youkuData.Data)
	var title string
	if youkuData.Data.Show.Title == "" || strings.Contains(
		youkuData.Data.Video.Title, youkuData.Data.Show.Title,
	) {
		title = youkuData.Data.Video.Title
	} else {
		title = fmt.Sprintf("%s %s", youkuData.Data.Show.Title, youkuData.Data.Video.Title)
	}

	other, _ := getMoreVideo(url, option)

	return []*extractors.Data{
		{
			Site:    "优酷 youku.com",
			Title:   title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     url,
			Series:  other,
			Cover:   youkuData.Data.Video.Logo,
		},
	}, nil
}

func getMoreVideo(uri string, option extractors.Options) ([]*extractors.Data, error) {
	var cookie = `cna=ubO6HDVBcWQCAcpkM3GpM3Xc; isI18n=false; __ysuid=1683186180965tIe; __ayft=1683186180966; __aysid=1683186180966qy1; __ayscnt=1; P_F=1; P_ck_ctl=F3EC895F26EAA42A8EAD6D4F50670748; P_gck=NA%7CzhFsKk9rCnptf7Ievlk2Gw%3D%3D%7CNA%7C1683187844398; disrd=78710; xlly_s=1; P_sck=1178WbKg85JyEcFqQZiyJskhiWTG1bSMS8glc0QrZ0wzMRhiJzN3LOXO45FcFgUOZ6Pek%2FMsOUwWHyKZD8O%2BKbH57Y4thvYjpyoKXgoDKRFwixblvw5J06ADxC1b%2B3VzsmujTrdbdIhhvQy778kv7w%3D%3D; P_pck_rm=6ffGReQN270ff57471f68fZBmzjTanCZo%2FoioEcK9he3kX7bIRuLhx49yONc929EuPw%2BMHGezZbvXGVmOY%2BpcDADoq1ZrbGGXFphQdepdS8IG5fZ1bhCcYJDjDlYsq175%2F%2FrMo87oGu95F6gWAY9G%2BU2LD4Oo1ttpYyBBA%3D%3D_V2; yseid=16832729875016nXi40; yseidcount=1; ycid=0; juid=01gvle6srh2rl6; seid=01gvle6srihk5; referhost=; ypvid=1683272994273RPDtMV; ysestep=2; yseidtimeout=1683280194274; ystep=2; seidtimeout=1683274794279; _m_h5_tk=4f458120a8704bd762035fb03014d9b9_1683354628643; _m_h5_tk_enc=a0246fdb7ce30e67fde6584f2719a9a7; x5sec=7b22617365727665722d686579693b32223a226637326139656232356237663930636263663463396565653437393536353830434a335831364947454b3334392b5743754a5774465443436766444f2b502f2f2f2f384251414d3d227d; __arycid=dd-3-00; __arcms=dd-3-00; __arpvid=16833528269900SbayW-1683352827044; __aypstp=93; __ayspstp=93; __ayvstp=518; __aysvstp=518; tfstk=cSucIngS286CKZTtFxaX8T_ldakQl_8O7J0ZP4RR9B9wDVs_8S5mHgSUj7dunQYlf; l=fBaof8xHNivf2lrSBO5anurza779aIRfGsPzaNbMiIEGB6SUGFIIhqdQ2qZvxYt5WhQNEsUeR3R4Lnf6B7Yzey4tmxv9-ewGtXmrndLnwpzHU; isg=BOfnBoYAT6voZss1yXrH_Hsfdh2xbLtOqtzvpLlU0nbHqANqwTnznxcgzqg2QJPG`
	if option.Cookie != "" {
		cookie = option.Cookie
	}

	resp, err := request.Client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36").
		SetHeader("Cookie", cookie).
		SetHeader("xreferrer", "http://www.youku.com").
		Get(uri)
	if err != nil {
		return nil, err
	}

	more := utils.MatchOneOf(resp.String(), "window.__INITIAL_DATA__ =(.*?);")
	if len(more) == 0 {
		return nil, errors.New("没有更多视频")
	}

	result := make([]*extractors.Data, 0)
	//10013 Web播放页选集组件
	xuanji := gjson.Parse(more[1]).Get(`data.data.nodes.0.nodes.#(type="10013")`)

	xuanji.Get("nodes").ForEach(func(key, value gjson.Result) bool {
		result = append(result, &extractors.Data{
			Title: value.Get("data.title").String(),
			URL:   fmt.Sprintf(`https://v.youku.com/v_show/id_%s.html`, value.Get("data.action.value")),
			Site:  "优酷 youku.com",
			Type:  extractors.DataTypeVideo,
			Cover: value.Get("data.img").String(),
		})
		return true
	})

	return result, nil
}

// extractPlaylist 解析主页视频列表
func extractPlaylist(homeUrl string) ([]*extractors.Data, error) {
	reg := regexp.MustCompile(`uid=(\w+)`)
	match := reg.FindStringSubmatch(homeUrl)

	if len(match) < 2 {
		return nil, errors.WithStack(extractors.ErrURLParseFailed)
	}

	uri := fmt.Sprintf(`https://www.youku.com/profile/profile-data?type=video&pageNo=1&uid=%s`, match[1])
	resp, err := request.Client.R().SetHeader("referer", fmt.Sprintf("https://www.youku.com/profile/index?uid=%s", match[1])).Get(uri)
	if err != nil {
		return nil, err
	}

	videoJson := gjson.ParseBytes(resp.Body()).Get(`data.componentList.#(type="video")`)
	if videoJson.Exists() {
		result := make([]*extractors.Data, 0)
		videoJson.Get("moduleList").ForEach(func(key, value gjson.Result) bool {
			value = value.Get("data")
			videoUrl, _ := url.JoinPath("https://", value.Get("videoLink").String())
			result = append(result, &extractors.Data{
				Site:  "优酷 youku.com",
				Title: value.Get("title").String(),
				Type:  extractors.DataTypeVideo,
				URL:   videoUrl,
				Cover: value.Get("imgUrl").String(),
			})
			return true
		})
		return result, err
	} else {
		return nil, errors.New("视频列表为空")
	}
}

func GetStealParams(initTime string) (string, error) {
	result := &StealParams{
		Ccode:    "0502",
		ClientIp: "192.168.1.1",
		Utid:     "ubO6HDVBcWQCAcpkM3GpM3Xc",
		ClientTs: initTime,
		Version:  "2.1.63",
		Ckey:     "DIl58SLFxFNndSV1GFNnMQVYkx1PP5tKe1siZu/86PR1u/Wh1Ptd+WOZsHHWxysSfAOhNJpdVWsdVJNsfJ8Sxd8WKVvNfAS8aS8fAOzYARzPyPc3JvtnPHjTdKfESTdnuTW6ZPvk2pNDh4uFzotgdMEFkzQ5wZVXl2Pf1/Y6hLK0OnCNxBj3+nb0v72gZ6b0td+WOZsHHWxysSo/0y9D2K42SaB8Y/+aD2K42SaB8Y/+ahU+WOZsHcrxysooUeND",
	}

	j, _ := json.Marshal(result)
	return string(j), nil
}

func GetBizParams(uri string) (string, error) {
	vids := utils.MatchOneOf(
		uri, `id_(.+?)\.html`, `id_(.+)`,
	)
	if vids == nil || len(vids) < 2 {
		return "", errors.WithStack(extractors.ErrURLParseFailed)
	}
	vid := vids[1]

	showId, err := GetCurrentShowId(uri)
	if err != nil {
		return "", err
	}

	result := &BizParams{
		Vid:           vid,
		PlayAbility:   "16782592", // 写死在js里的
		CurrentShowid: showId,
		PreferClarity: "4",             // 貌似是清晰度
		Extag:         "EXT-X-PRIVINF", // 写死在js里的
		MasterM3U8:    "1",
		MediaType:     "standard,subtitle",
		AppVer:        "2.1.63",
		DrmType:       "19",
		KeyIndex:      "web01",
	}

	j, _ := json.Marshal(result)
	return string(j), nil
}

func GetAdParams() string {
	result := &AdParams{
		Vs:        "1.0",
		Pver:      "2.1.63",
		Sver:      "2.0",
		Site:      1,
		Aw:        "w",
		Fu:        0,
		D:         "0",
		Bt:        "pc",
		Os:        "mac",
		Osv:       "",
		Dq:        "auto",
		Atm:       "",
		Partnerid: "null",
		Wintype:   "interior",
		Isvert:    0,
		Vip:       0,
		P:         1,
		Rst:       "mp4",
		Needbf:    2,
		Avs:       "1.0"}

	j, _ := json.Marshal(result)
	return string(j)
}

func getData(uri, initTime string) (*Data, error) {
	steal, err := GetStealParams(initTime)
	if err != nil {
		return nil, err
	}

	biz, err := GetBizParams(uri)
	if err != nil {
		return nil, err
	}
	return &Data{StealParams: steal, BizParams: biz, AdParams: GetAdParams()}, nil
}

func GetJoinParams(uri, cookie, initTime string) (*JoinParams, error) {
	data, err := getData(uri, initTime)
	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	token := regexp.MustCompile("_m_h5_tk=(.+?)_").FindStringSubmatch(cookie)

	aa := string(dataStr)
	ss := []byte(token[1] + "&" + initTime + "&24679788&" + aa)

	m5 := md5.New()
	m5.Write(ss)
	sign := fmt.Sprintf("%x", m5.Sum(nil))

	return &JoinParams{
		Jsv:            "2.5.8",
		AppKey:         "24679788",
		T:              initTime,
		Sign:           sign,
		Api:            "mtop.youku.play.ups.appinfo.get",
		V:              "1.1",
		Timeout:        "20000",
		YKPid:          "20160317PLF000211",
		YKLoginRequest: "true",
		AntiFlood:      "true",
		AntiCreep:      "true",
		Type:           "json",
		DataType:       "json",
		Data:           aa}, nil
}

func GetCurrentShowId(uri string) (string, error) {
	body, err := request.Client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36").
		Get(uri)
	if err != nil {
		return "", err
	}
	return utils.MatchOneOf(body.String(), `id_(.*?).html`)[1], nil
}

func getCookie() (string, error) {
	body, err := request.Client.R().SetQueryParams(map[string]string{
		"appKey": "24679788",
		"api":    "mtop.youku.play.ups.appinfo.get",
		"t":      cast.ToString(time.Now().UnixMilli()),
	}).SetHeaders(map[string]string{
		"Referer":    "https://v.youku.com/",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
	}).Get(`https://acs.youku.com/h5/mtop.youku.play.ups.appinfo.get/1.1/`)
	if err != nil {
		return "", err
	}

	var result = "cna=ubO6HDVBcWQCAcpkM3GpM3Xc;"
	for _, item := range body.Cookies() {
		result += fmt.Sprintf("%s=%s;", item.Name, item.Value)
	}

	return result, nil
}
