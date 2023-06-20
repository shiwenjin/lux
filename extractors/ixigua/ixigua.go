package ixigua

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/request"
)

func init() {
	extractors.Register("ixigua", New())
	extractors.Register("toutiao", New())
}

type extractor struct{}

type Video struct {
	Title     string `json:"title"`
	Qualities []struct {
		Quality string `json:"quality"`
		Size    int64  `json:"size"`
		URL     string `json:"url"`
		Ext     string `json:"ext"`
	} `json:"qualities"`
}

// New returns a ixigua extractor.
func New() extractors.Extractor {
	return &extractor{}
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(uri string, option extractors.Options) ([]*extractors.Data, error) {
	if option.Playlist {
		return extractPlaylist(uri)
	}

	//headers := map[string]string{
	//	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:98.0) Gecko/20100101 Firefox/98.0",
	//	"Content-Type": "application/json",
	//}

	// ixigua 有三种格式的 URL
	// 格式一 https://www.ixigua.com/7053389963487871502
	// 格式二 https://v.ixigua.com/RedcbWM/
	// 格式三 https://m.toutiao.com/is/dtj1pND/
	// 格式二会跳转到格式一
	// 格式三会跳转到 https://www.toutiao.com/a7053389963487871502

	var finalURL string
	if strings.HasPrefix(uri, "https://www.ixigua.com/") {
		finalURL = uri
	}

	if strings.HasPrefix(uri, "https://v.ixigua.com/") || strings.HasPrefix(uri, "https://m.toutiao.com/") {
		resp, err := http.Get(uri)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer resp.Body.Close() // nolint
		// follow redirects, https://stackoverflow.com/a/16785343
		finalURL = resp.Request.URL.String()
	}

	if strings.Contains(finalURL, "https://www.toutiao.com/a") {
		finalURL = strings.ReplaceAll(finalURL, "https://www.toutiao.com/a", "https://www.ixigua.com/")
	}

	if strings.Contains(finalURL, "https://www.toutiao.com/video") {
		finalURL = strings.ReplaceAll(finalURL, "https://www.toutiao.com/video", "https://www.ixigua.com")
	}

	r := regexp.MustCompile(`(ixigua.com/)(\w+)?`)
	id := r.FindSubmatch([]byte(finalURL))[2]
	url2 := fmt.Sprintf("https://toutiao.com/video/%s", id)

	video, err := Parse(url2)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	streams := make(map[string]*extractors.Stream)

	videoList := video.Data.InitialVideo.VideoPlayInfo.VideoList
	// 音视频分开
	if len(videoList) == 0 {
		part := extractors.Part{}

		audio := video.Data.InitialVideo.VideoPlayInfo.DynamicVideo.DynamicAudioList[0]
		part.URL = audio.MainUrl
		part.Size = audio.AudioMeta.Size
		part.Ext = "m4a"

		for _, quality := range video.Data.InitialVideo.VideoPlayInfo.DynamicVideo.DynamicVideoList {
			streams[quality.VideoMeta.Definition] = &extractors.Stream{
				Size:    quality.VideoMeta.Size,
				Quality: quality.VideoMeta.Definition,
				Parts: []*extractors.Part{
					{
						URL:  quality.MainUrl,
						Size: quality.VideoMeta.Size,
						Ext:  "mp4",
					}, &part,
				},
				NeedMux: true,
			}
		}
	} else {
		for _, quality := range videoList {
			streams[quality.VideoMeta.Definition] = &extractors.Stream{
				Size:    quality.VideoMeta.Size,
				Quality: quality.VideoMeta.Definition,
				Parts: []*extractors.Part{
					{
						URL:  quality.MainUrl,
						Size: quality.VideoMeta.Size,
						Ext:  quality.VideoMeta.Vtype,
					},
				},
			}
		}
	}

	cover := "https:" + video.Data.InitialVideo.CoverUrl
	return []*extractors.Data{
		{
			Site:    "西瓜视频 ixigua.com",
			Title:   video.Data.InitialVideo.Title,
			Type:    extractors.DataTypeVideo,
			Streams: streams,
			URL:     uri,
			Cover:   cover,
		},
	}, nil
}

var homeReg = regexp.MustCompile(`^https:\/\/www\.ixigua\.com\/home\/(\d+)\/`)

func extractPlaylist(u string) ([]*extractors.Data, error) {
	match := homeReg.FindStringSubmatch(u)
	if len(match) != 2 {
		return nil, extractors.ErrURLParseFailed
	}
	vo, err := getNewVideoListByXiGua(match[1])
	if err != nil {
		return nil, err
	}

	data := make([]*extractors.Data, 0)
	for _, s := range vo.Data.VideoList {
		data = append(data, &extractors.Data{
			URL:   fmt.Sprintf(`https://www.ixigua.com/%s`, s.GroupId),
			Site:  "西瓜视频 ixigua.com",
			Title: s.Title,
			Type:  extractors.DataTypeVideo,
			Cover: s.VideoDetailInfo.DetailVideoLargeImage.Url,
		})
	}

	return data, nil
}

func getNewVideoListByXiGua(userId string) (*AuthorXiGuaVideoVo, error) {
	v := new(AuthorXiGuaVideoVo)
	path := fmt.Sprintf(`https://www.ixigua.com/api/videov2/author/new_video_list?to_user_id=%s`, userId)

	//获取 请求参数 _signature值
	//signVo, signErr := a.getSign(path)
	//if signErr != nil {
	//	return nil, errors.Trace(signErr)
	//}

	//请求 url
	pathUrl := fmt.Sprintf(path+"&_signature=%s", "")

	//请求头Header  需要 rerfer 参数
	refererValue := fmt.Sprintf(`https://www.ixigua.com/home/%s/video/?preActiveKey=pseries&list_entrance=userdetail`, userId)

	//发起请求 存入 v
	_, err := request.Client.R().SetHeader("Referer", refererValue).SetResult(v).Get(pathUrl)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func base64Decode(t string) string {
	d, _ := base64.StdEncoding.DecodeString(t)
	return string(d)
}

func Parse(videoUrl string) (*VideoHomeVo, error) {
	cookie := "MONITOR_WEB_ID=7143612877764527652; __ac_nonce=06354165e0031397cdbb2; __ac_signature=_02B4Z6wo00f01fMszzAAAIDAeGdU0ds5es3zDMuAAB-hLkcyw3fUcsBLTtRrE0F5G49ooKwZz6ndN47fnhZx9zSVC6fgul9Gm0gDXVB77WhH564acTz3U67bXgn2Ve2-vAbsBDgsLMUWo2ive3; _tea_utm_cache_1768={%22utm_source%22:%22copy_link%22%2C%22utm_medium%22:%22android%22%2C%22utm_campaign%22:%22client_share%22}; _tea_utm_cache_1300={%22utm_source%22:%22copy_link%22%2C%22utm_medium%22:%22android%22%2C%22utm_campaign%22:%22client_share%22}; _tea_utm_cache_2285={%22utm_source%22:%22copy_link%22%2C%22utm_medium%22:%22android%22%2C%22utm_campaign%22:%22client_share%22}; ixigua-a-s=1; tt_scid=1RzrmUQ51j8q5QQYOPq-5V6RvaJewfqGuGfAF6294SbxsJGjStTZsFiOEL.YpAp.3e2c; ttwid=1%7Ccu9m9yb45Ydazbt1ZywzV5oW-kcMjhLOL6wl3BCCLfw%7C1666455550%7C46d6198ce31b93fe21f5e4ead3928e9316c6a87ab6691c696d9eaf0e33ef665a;"

	body, err := request.Client.R().SetHeader("Cookie", cookie).Get(videoUrl)
	if err != nil {
		return nil, err
	}

	res := bytes.NewReader(body.Body())

	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jsonData := doc.Find("#RENDER_DATA").Text()
	realJsonData, err := url.QueryUnescape(jsonData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var video VideoHomeVo
	err = json.Unmarshal([]byte(realJsonData), &video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

type VideoHomeVo struct {
	Data struct {
		ItemId           string `json:"itemId"`
		GroupId          string `json:"groupId"`
		BizId            string `json:"bizId"`
		VideoType        string `json:"videoType"`
		Pathname         string `json:"pathname"`
		ImmerseEnterFrom string `json:"immerseEnterFrom"`
		InitialVideo     struct {
			GroupId       string `json:"group_id"`
			VideoType     string `json:"videoType"`
			GroupId1      string `json:"groupId"`
			ItemId        string `json:"itemId"`
			GroupSource   int    `json:"groupSource"`
			Title         string `json:"title"`
			DetailUrl     string `json:"detailUrl"`
			CoverUrl      string `json:"coverUrl"`
			Poster        string `json:"poster"`
			PublishTime   int64  `json:"publishTime"`
			VideoPlayInfo struct {
				Status        int     `json:"status"`
				Message       string  `json:"message"`
				Version       int     `json:"version"`
				VideoId       string  `json:"video_id"`
				EnableSsl     bool    `json:"enable_ssl"`
				VideoDuration float64 `json:"video_duration"`
				MediaType     string  `json:"media_type"`
				UrlExpire     int     `json:"url_expire"`
				BigThumbs     []struct {
					ImgNum   int      `json:"img_num"`
					ImgUri   string   `json:"img_uri"`
					ImgUrl   string   `json:"img_url"`
					ImgXSize int      `json:"img_x_size"`
					ImgYSize int      `json:"img_y_size"`
					ImgXLen  int      `json:"img_x_len"`
					ImgYLen  int      `json:"img_y_len"`
					Duration float64  `json:"duration"`
					Interval int      `json:"interval"`
					Fext     string   `json:"fext"`
					ImgUrls  []string `json:"img_urls"`
				} `json:"big_thumbs"`
				FallbackApi struct {
					FallbackApi string `json:"fallback_api"`
				} `json:"fallback_api"`
				DynamicVideo DynamicVideoResp `json:"dynamic_video"`
				VideoList    []struct {
					MainUrl   string `json:"main_url"`
					BackupUrl string `json:"backup_url"`
					VideoMeta struct {
						Definition    string `json:"definition"`
						Quality       string `json:"quality"`
						Vtype         string `json:"vtype"`
						Vwidth        int    `json:"vwidth"`
						Vheight       int    `json:"vheight"`
						Bitrate       int    `json:"bitrate"`
						CodecType     string `json:"codec_type"`
						Size          int64  `json:"size"`
						FileId        string `json:"file_id"`
						Fps           int    `json:"fps"`
						FileHash      string `json:"file_hash"`
						RealBitrate   int    `json:"real_bitrate"`
						AudioChannels string `json:"audio_channels"`
						AudioLayout   string `json:"audio_layout"`
					} `json:"video_meta"`
					P2PInfo struct {
						P2PVerifyUrl string `json:"p2p_verify_url"`
					} `json:"p2p_info"`
					CheckInfo struct {
						CheckInfo string `json:"check_info"`
					} `json:"check_info"`
					Volume struct {
						Loudness float64 `json:"loudness"`
						Peak     float64 `json:"peak"`
					} `json:"volume"`
					QualityType int    `json:"quality_type"`
					PktOffset   string `json:"pkt_offset"`
				} `json:"video_list"`
				EnableAdaptive bool `json:"enable_adaptive"`
				Volume         struct {
					Loudness float64 `json:"loudness"`
					Peak     float64 `json:"peak"`
				} `json:"volume"`
				SubtitleLangs []int `json:"subtitle_langs"`
				SubtitleInfos []struct {
					SubId      int    `json:"sub_id"`
					LanguageId int    `json:"language_id"`
					Format     string `json:"format"`
					Version    string `json:"version"`
					Size       int    `json:"size"`
				} `json:"subtitle_infos"`
				HasEmbeddedSubtitle bool `json:"has_embedded_subtitle"`
			} `json:"videoPlayInfo"`
			UserInfo struct {
				Name         string `json:"name"`
				UserId       string `json:"userId"`
				UserAuthInfo struct {
					AuthInfo  string `json:"auth_info"`
					AuthType  string `json:"auth_type"`
					OtherAuth struct {
						Interest string `json:"interest"`
					} `json:"other_auth"`
				} `json:"userAuthInfo"`
				AvatarUrl   string `json:"avatarUrl"`
				IsFollowing bool   `json:"isFollowing"`
			} `json:"userInfo"`
			ProfileUrl   string `json:"profileUrl"`
			PlayCount    int    `json:"playCount"`
			DiggCount    int    `json:"diggCount"`
			UserDigg     int    `json:"userDigg"`
			CommentCount int    `json:"commentCount"`
			RepinCount   int    `json:"repinCount"`
			UserRepin    int    `json:"userRepin"`
			Duration     int    `json:"duration"`
			VideoId      string `json:"videoId"`
			MediaId      int64  `json:"mediaId"`
			LogPb        struct {
				AuthorId    string `json:"author_id"`
				BizId       string `json:"biz_id"`
				GroupId     string `json:"group_id"`
				GroupSource string `json:"group_source"`
				ImprId      string `json:"impr_id"`
				IsFollowing string `json:"is_following"`
			} `json:"log_pb"`
		} `json:"initialVideo"`
		SeoTDK struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Keywords    string `json:"keywords"`
		} `json:"seoTDK"`
		LogId    string `json:"logId"`
		Identity struct {
			WebId       string `json:"web_id"`
			UserIsLogin bool   `json:"user_is_login"`
		} `json:"identity"`
		AbtestInfo struct {
			RspType     int    `json:"rsp_type"`
			VersionName string `json:"version_name"`
			Parameters  struct {
				Filter struct {
					DebugEnablePcSmallVideo bool `json:"debug_enable_pc_small_video"`
					EnablePcSmallVideo      bool `json:"enable_pc_small_video"`
				} `json:"filter"`
				HomeNavConf struct {
					DcdOut int `json:"dcd_out"`
				} `json:"home_nav_conf"`
				Optimus struct {
					RuleRankRules string `json:"rule_rank_rules"`
				} `json:"optimus"`
				PageUpgrade struct {
					NewProfile        bool `json:"new_profile"`
					VideoDoubleColumn bool `json:"video_double_column"`
				} `json:"page_upgrade"`
				RandomSource struct {
					FoldedHeight int  `json:"folded_height"`
					FoldComment  bool `json:"fold_comment"`
				} `json:"random_source"`
				Recall struct {
					DebugFilterReasonList   []interface{} `json:"debug_filter_reason_list"`
					DebugRecallReasonList   []int         `json:"debug_recall_reason_list"`
					EnableDebugFilterReason bool          `json:"enable_debug_filter_reason"`
					EnableDebugRecallReason bool          `json:"enable_debug_recall_reason"`
					FriendFeed              struct {
						ControlNums                     int  `json:"control_nums"`
						Count                           int  `json:"count"`
						Enable                          bool `json:"enable"`
						EnableCppFriendFeed             bool `json:"enable_cpp_friend_feed"`
						EnableFriendFeedContainerFilter bool `json:"enable_friend_feed_container_filter"`
						EnableGroupStatusFilter         bool `json:"enable_group_status_filter"`
						Params                          struct {
						} `json:"params"`
						TopK               int   `json:"top_k"`
						UseNewFriendRecall bool  `json:"use_new_friend_recall"`
						ValidGroupSource   []int `json:"valid_group_source"`
						Weight             int   `json:"weight"`
					} `json:"friend_feed"`
				} `json:"recall"`
				Seraph struct {
					RuleRankRulesPcSmall       string `json:"rule_rank_rules_pc_small"`
					RuleRankRulesPcSmallWindow string `json:"rule_rank_rules_pc_small_window"`
				} `json:"seraph"`
				SmallSort struct {
					EnableFilterZhanwai       bool `json:"enable_filter_zhanwai"`
					FilterGenrePlog           bool `json:"filter_genre_plog"`
					FilterGenreSmall          bool `json:"filter_genre_small"`
					SkipSmallAppVersionFilter bool `json:"skip_small_app_version_filter"`
					SkipSmallGenreFilter      bool `json:"skip_small_genre_filter"`
				} `json:"small_sort"`
				Sort struct {
					AllowedTicai                  []string `json:"allowed_ticai"`
					EnableOptimusGenPcMvCard      bool     `json:"enable_optimus_gen_pc_mv_card"`
					EnableOptimusGenPcSvCard      bool     `json:"enable_optimus_gen_pc_sv_card"`
					EnablePcSkipAppSmallvideoCard bool     `json:"enable_pc_skip_app_smallvideo_card"`
				} `json:"sort"`
				UgcSort struct {
					ExporeSmallvideo bool `json:"expore_smallvideo"`
				} `json:"ugc_sort"`
				VideoChannel struct {
					UseFeed int `json:"use_feed"`
					Rank    int `json:"rank"`
				} `json:"video_channel"`
			} `json:"parameters"`
			EnvFlag int `json:"env_flag"`
		} `json:"abtestInfo"`
		LocalCityInfo struct {
			Name      string `json:"name"`
			Code      string `json:"code"`
			ChannelId int64  `json:"channelId"`
		} `json:"localCityInfo"`
		IsGreyTheme bool        `json:"isGreyTheme"`
		SearchBot   interface{} `json:"searchBot"`
	} `json:"data"`
}

type DynamicVideoResp struct {
	DynamicType      string `json:"dynamic_type"`
	DynamicVideoList []struct {
		MainUrl   string `json:"main_url"`
		BackupUrl string `json:"backup_url"`
		VideoMeta struct {
			Definition  string `json:"definition"`
			Quality     string `json:"quality"`
			Vtype       string `json:"vtype"`
			Vwidth      int    `json:"vwidth"`
			Vheight     int    `json:"vheight"`
			Bitrate     int    `json:"bitrate"`
			CodecType   string `json:"codec_type"`
			Size        int64  `json:"size"`
			FileId      string `json:"file_id"`
			Fps         int    `json:"fps"`
			FileHash    string `json:"file_hash"`
			RealBitrate int    `json:"real_bitrate"`
		} `json:"video_meta"`
		BaseRangeInfo struct {
			InitRange  string `json:"init_range"`
			IndexRange string `json:"index_range"`
		} `json:"base_range_info"`
		CheckInfo struct {
			CheckInfo string `json:"check_info"`
		} `json:"check_info"`
		QualityType int `json:"quality_type"`
	} `json:"dynamic_video_list"`
	DynamicAudioList []struct {
		MainUrl   string `json:"main_url"`
		BackupUrl string `json:"backup_url"`
		AudioMeta struct {
			Quality     string `json:"quality"`
			Atype       string `json:"atype"`
			Bitrate     int    `json:"bitrate"`
			CodecType   string `json:"codec_type"`
			Size        int64  `json:"size"`
			FileId      string `json:"file_id"`
			FileHash    string `json:"file_hash"`
			RealBitrate int    `json:"real_bitrate"`
		} `json:"audio_meta"`
		BaseRangeInfo struct {
			InitRange  string `json:"init_range"`
			IndexRange string `json:"index_range"`
		} `json:"base_range_info"`
		CheckInfo struct {
			CheckInfo string `json:"check_info"`
		} `json:"check_info"`
	} `json:"dynamic_audio_list"`
}
