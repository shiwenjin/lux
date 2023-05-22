package ixigua

type AuthorXiGuaVideoVo struct {
	Code int `json:"code"`
	Data struct {
		VideoList []struct {
			GroupId        string `json:"group_id"`
			Gid            string `json:"gid"`
			ItemId         string `json:"item_id"`
			AggrType       int    `json:"aggr_type"`
			Title          string `json:"title"`
			PublishTime    int    `json:"publish_time"`
			HasVideo       bool   `json:"has_video"`
			ArticleType    int    `json:"article_type"`
			ArticleSubType int    `json:"article_sub_type"`
			ArticleUrl     string `json:"article_url"`
			DisplayUrl     string `json:"display_url"`
			Abstract       string `json:"abstract"`
			IsOriginal     bool   `json:"is_original"`
			VideoExclusive bool   `json:"video_exclusive"`
			DetailSchema   string `json:"detail_schema"`
			CellType       string `json:"cell_type"`
			GroupFlags     string `json:"group_flags"`
			VideoStyle     int    `json:"video_style"`
			Composition    int    `json:"composition"`
			Source         string `json:"source"`
			MediaName      string `json:"media_name"`
			MediaInfo      struct {
				AvatarUrl    string `json:"avatar_url"`
				Name         string `json:"name"`
				UserVerified bool   `json:"user_verified"`
				MediaId      string `json:"media_id"`
				Subscribed   int    `json:"subscribed"`
				Subcribed    int    `json:"subcribed"`
			} `json:"media_info"`
			IsSubscribe bool `json:"is_subscribe"`
			UserInfo    struct {
				AvatarUrl         string `json:"avatar_url"`
				Name              string `json:"name"`
				Description       string `json:"description"`
				UserId            string `json:"user_id"`
				SecUserId         string `json:"sec_user_id"`
				UserVerified      bool   `json:"user_verified"`
				VerifiedContent   string `json:"verified_content"`
				Follow            bool   `json:"follow"`
				FollowerCount     string `json:"follower_count"`
				IsLiving          bool   `json:"is_living"`
				AuthorDesc        string `json:"author_desc"`
				VideoTotalCount   string `json:"video_total_count"`
				IsDiscipulus      bool   `json:"is_discipulus"`
				IsBlocking        bool   `json:"is_blocking"`
				IsBlocked         bool   `json:"is_blocked"`
				FollowersCountStr string `json:"followers_count_str"`
				IsFollower        bool   `json:"is_follower"`
			} `json:"user_info"`
			UserVerified    int    `json:"user_verified"`
			VideoDuration   string `json:"video_duration"`
			HistoryDuration string `json:"history_duration"`
			VideoId         string `json:"video_id"`
			VideoDetailInfo struct {
				GroupFlags            string `json:"group_flags"`
				DetailVideoLargeImage struct {
					Url     string `json:"url"`
					Width   int    `json:"width"`
					UrlList []struct {
						Url string `json:"url"`
					} `json:"url_list"`
					Uri    string `json:"uri"`
					Height int    `json:"height"`
				} `json:"detail_video_large_image"`
				VideoId             string        `json:"video_id"`
				DirectPlay          int           `json:"direct_play"`
				ShowPgcSubscribe    int           `json:"show_pgc_subscribe"`
				VideoWatchCount     int           `json:"video_watch_count"`
				VideoType           int           `json:"video_type"`
				VideoPreloadingFlag int           `json:"video_preloading_flag"`
				VideoUrl            []interface{} `json:"video_url"`
				LastPlayDuration    string        `json:"last_play_duration"`
				UseLastDuration     bool          `json:"use_last_duration"`
			} `json:"video_detail_info"`
			VideoProportion        float64 `json:"video_proportion"`
			VideoProportionArticle float64 `json:"video_proportion_article"`
			ShowPortrait           bool    `json:"show_portrait"`
			ShowPortraitArticle    bool    `json:"show_portrait_article"`
			LargeImageList         []struct {
				Url     string `json:"url"`
				Width   int    `json:"width"`
				UrlList []struct {
					Url string `json:"url"`
				} `json:"url_list"`
				Uri    string `json:"uri"`
				Height int    `json:"height"`
			} `json:"large_image_list"`
			MiddleImage struct {
				Url     string `json:"url"`
				Width   int    `json:"width"`
				UrlList []struct {
					Url string `json:"url"`
				} `json:"url_list"`
				Uri    string `json:"uri"`
				Height int    `json:"height"`
			} `json:"middle_image"`
			FirstFrameImage struct {
				Url     string `json:"url"`
				Width   int    `json:"width"`
				UrlList []struct {
					Url string `json:"url"`
				} `json:"url_list"`
				Uri    string `json:"uri"`
				Height int    `json:"height"`
			} `json:"first_frame_image"`
			PlayAuthToken   string `json:"play_auth_token"`
			PlayBizToken    string `json:"play_biz_token"`
			PreviewUrl      string `json:"preview_url"`
			ShareUrl        string `json:"share_url"`
			VideoUserLike   int    `json:"video_user_like"`
			BanDanmaku      int    `json:"ban_danmaku"`
			BanDanmakuSend  int    `json:"ban_danmaku_send"`
			DefaultDanmaku  int    `json:"default_danmaku"`
			BanComment      int    `json:"ban_comment"`
			BanDownload     int    `json:"ban_download"`
			CanCommentLevel int    `json:"can_comment_level"`
			DanmakuCount    string `json:"danmaku_count"`
			VideoLikeCount  int    `json:"video_like_count"`
			DiggCount       int    `json:"digg_count"`
			BuryCount       int    `json:"bury_count"`
			CommentCount    int    `json:"comment_count"`
			ImpressionCount int    `json:"impression_count"`
			RepinCount      int    `json:"repin_count"`
			Tag             string `json:"tag"`
			ShareCount      int    `json:"share_count"`
			LogPb           struct {
				CategoryName string `json:"category_name"`
				ImprId       string `json:"impr_id"`
				EnterFrom    string `json:"enter_from"`
				IsFollowing  string `json:"is_following"`
				AuthorId     string `json:"author_id"`
				GroupId      string `json:"group_id"`
				GroupSource  string `json:"group_source"`
				ImprType     string `json:"impr_type"`
				AlbumType    string `json:"album_type,omitempty"`
				AlbumId      string `json:"album_id,omitempty"`
			} `json:"log_pb"`
			VerifyStatus int      `json:"verify_status"`
			VerifyReason string   `json:"verify_reason"`
			Categories   []string `json:"categories"`
			PreAdParams  string   `json:"pre_ad_params"`
			GroupIdStr   string   `json:"group_id_str"`
			NearId       string   `json:"near_id"`
			NearId2      string   `json:"near_id2"`
			NearId3      string   `json:"near_id3"`
			IsTop        bool     `json:"is_top"`
			CanTop       bool     `json:"can_top"`
			BehotTime    string   `json:"behot_time"`
			Cursor       string   `json:"cursor"`
			GroupSource  string   `json:"group_source"`
			StickerList  []struct {
				StartTime       string `json:"start_time"`
				Duration        string `json:"duration"`
				MarginLeft      string `json:"margin_left"`
				MarginTop       string `json:"margin_top"`
				StickType       int    `json:"stick_type"`
				StickerId       string `json:"sticker_id"`
				HeightRatio     string `json:"height_ratio"`
				StickerEffectId string `json:"sticker_effect_id"`
				Extra           struct {
					FollowSource string `json:"follow_source,omitempty"`
					Source       string `json:"source"`
				} `json:"extra,omitempty"`
				FollowInfo struct {
					UserId    string `json:"user_id"`
					AvatarUrl string `json:"avatar_url"`
					IsFollow  bool   `json:"is_follow"`
				} `json:"follow_info,omitempty"`
				DiggInfo struct {
					IsDigged       bool `json:"is_digged"`
					IsSupperDigged bool `json:"is_supper_digged"`
				} `json:"digg_info,omitempty"`
				DanmakuInfo struct {
					Title   string `json:"title"`
					Content string `json:"content"`
				} `json:"danmaku_info,omitempty"`
				VoteInfo struct {
					Title                string `json:"title"`
					Type                 int    `json:"type"`
					MaxSelectedOptionCnt string `json:"max_selected_option_cnt"`
					Options              []struct {
						OptionId int    `json:"option_id"`
						Type     int    `json:"type"`
						Text     string `json:"text"`
						ImageUri string `json:"image_uri"`
						ImageUrl string `json:"image_url"`
						Count    string `json:"count"`
						IsChosen bool   `json:"is_chosen"`
					} `json:"options"`
					VoteId         string `json:"vote_id"`
					IsVoted        bool   `json:"is_voted"`
					TotalCount     string `json:"total_count"`
					TotalUserCount string `json:"total_user_count"`
				} `json:"vote_info,omitempty"`
			} `json:"sticker_list"`
			SuperDiggControl struct {
				AnimeKey     string `json:"anime_key"`
				AnimeHeadUrl string `json:"anime_head_url"`
				AnimeBodyUrl string `json:"anime_body_url"`
				AnimeLokiId  string `json:"anime_loki_id"`
				Audio        struct {
					Uri       string `json:"uri"`
					Url       string `json:"url"`
					AudioType string `json:"audio_type"`
					AudioName string `json:"audio_name"`
				} `json:"audio"`
			} `json:"super_digg_control"`
			BanDanmakuReason string `json:"ban_danmaku_reason"`
			TimerPublishText string `json:"timer_publish_text"`
			XgVideoRichText  struct {
			} `json:"xg_video_rich_text"`
			XiRelated             bool   `json:"xi_related"`
			BanDownloadReason     string `json:"ban_download_reason"`
			CanJumpDetail         bool   `json:"can_jump_detail"`
			CanNotJumpReason      string `json:"can_not_jump_reason"`
			HideType              int    `json:"hide_type"`
			BanAudioComment       bool   `json:"ban_audio_comment"`
			BanAudioCommentReason string `json:"ban_audio_comment_reason"`
			IpInfo                struct {
				Address string `json:"address"`
			} `json:"ip_info"`
			BanShare       int `json:"ban_share"`
			Offset         int `json:"offset"`
			Rank           int `json:"rank"`
			HomoLvideoInfo struct {
				AlbumId                string `json:"album_id"`
				AlbumGroupId           string `json:"album_group_id"`
				EpisodeId              string `json:"episode_id"`
				Duration               int    `json:"duration"`
				Title                  string `json:"title"`
				SubTitle               string `json:"sub_title"`
				PlayButtonText         string `json:"play_button_text"`
				LandscapeGoDetailTitle string `json:"landscape_go_detail_title"`
				LandscapeGoDetailHint  string `json:"landscape_go_detail_hint"`
				Cover                  struct {
					Url     string `json:"url"`
					Width   int    `json:"width"`
					UrlList []struct {
						Url string `json:"url"`
					} `json:"url_list"`
					Uri    string `json:"uri"`
					Height int    `json:"height"`
				} `json:"cover"`
				VerticalCover struct {
					Url    string `json:"url"`
					Width  int    `json:"width"`
					Uri    string `json:"uri"`
					Height int    `json:"height"`
				} `json:"vertical_cover"`
				AlbumType      int    `json:"album_type"`
				BubbleStyle    int    `json:"bubble_style"`
				SectionControl string `json:"section_control"`
				ActionUrl      string `json:"action_url"`
				SliceInfoList  []struct {
					ShortSliceStartTime       string `json:"short_slice_start_time"`
					ShortSliceEndTime         string `json:"short_slice_end_time"`
					LongMatchedSliceStartTime string `json:"long_matched_slice_start_time"`
					LongMatchedSliceEndTime   string `json:"long_matched_slice_end_time"`
				} `json:"slice_info_list"`
				BeltStyle               int    `json:"belt_style"`
				TvHomePartnerAction     string `json:"tv_home_partner_action"`
				FavoriteStatus          bool   `json:"favorite_status"`
				FavoriteStyle           int    `json:"favorite_style"`
				CanSubscribe            bool   `json:"can_subscribe"`
				HasSubscribed           bool   `json:"has_subscribed"`
				HasCopyRight            bool   `json:"has_copy_right"`
				IfDiverse               bool   `json:"if_diverse"`
				CompassInfoSchema       string `json:"compass_info_schema"`
				SubscribeHint           string `json:"subscribe_hint"`
				SubscribeOnlineTimeHint string `json:"subscribe_online_time_hint"`
				ToastHint               string `json:"toast_hint"`
				SubscribeButtonSchema   string `json:"subscribe_button_schema"`
				GuideText               string `json:"guide_text"`
			} `json:"homo_lvideo_info,omitempty"`
			PseriesRank string `json:"pseries_rank,omitempty"`
			DanmakuMask int    `json:"danmaku_mask,omitempty"`
		} `json:"videoList"`
	} `json:"data"`
}
