package iqiyi

type GetVideosId struct {
	Code string `json:"code"`
	Data struct {
		TotalNum int `json:"totalNum"`
		HasMore  int `json:"hasMore"`
		Sort     struct {
			ReprentativeWork struct {
				QipuId int64 `json:"qipuId"`
			} `json:"reprentativeWork"`
			Flows []struct {
				QipuId int64 `json:"qipuId"`
			} `json:"flows"`
		} `json:"sort"`
	} `json:"data"`
	Msg interface{} `json:"msg"`
}

type SmallVideoResp struct {
	Code string `json:"code"`
	Data struct {
		StoredNum int     `json:"storedNum"`
		Tvids     []int64 `json:"tvids"`
	} `json:"data"`
	Msg interface{} `json:"msg"`
}

type GetVideosInfo struct {
	Code string          `json:"code"`
	Data map[string]Item `json:"data"`
	Msg  interface{}     `json:"msg"`
}

type Item struct {
	QipuId                     int64         `json:"qipuId"`
	AlbumId                    int           `json:"albumId"`
	Title                      string        `json:"title"`
	PlayCount                  int           `json:"playCount"`
	PlayTime                   int           `json:"playTime"`
	Uid                        int           `json:"uid"`
	Thumbnail                  string        `json:"thumbnail"`
	VerticalThumbnail          string        `json:"verticalThumbnail"`
	FrgmentThumbnail           string        `json:"frgmentThumbnail"`
	VideoThumbnail             string        `json:"videoThumbnail"`
	FeatureThumbnail           interface{}   `json:"featureThumbnail"`
	CommentCoverThumbnail      string        `json:"commentCoverThumbnail"`
	TimeCreate                 int64         `json:"timeCreate"`
	FeedId                     int           `json:"feedId"`
	ChannelId                  int           `json:"channelId"`
	CType                      int           `json:"cType"`
	Pc                         int           `json:"pc"`
	Status                     int           `json:"status"`
	LikeCount                  int           `json:"likeCount"`
	VerticalVideoFirstFrameUrl string        `json:"verticalVideoFirstFrameUrl"`
	EntityType                 int           `json:"entityType"`
	Description                string        `json:"description"`
	FatherEpisodeId            int           `json:"father_episode_id"`
	IsFeature                  int           `json:"isFeature"`
	QichuanId                  interface{}   `json:"qichuanId"`
	Vid                        string        `json:"vid"`
	BusinessType               []interface{} `json:"businessType"`
	PageUrl                    string        `json:"pageUrl"`
	PlayModel                  int           `json:"playModel"`
	VideoType                  int           `json:"videoType"`
	EpisodeType                int           `json:"episodeType"`
	InteractionType            string        `json:"interaction_type"`
	IsEnabledInteraction       bool          `json:"is_enabled_interaction"`
	InteractionScriptUrl       string        `json:"interaction_script_url"`
	VedioSize                  interface{}   `json:"vedioSize"`
	Tag                        interface{}   `json:"tag"`
	StarList                   []interface{} `json:"starList"`
	SourceProvider             interface{}   `json:"sourceProvider"`
	HotScore                   int           `json:"hotScore"`
	TopicSummaryList           interface{}   `json:"topicSummaryList"`
	Nickname                   string        `json:"nickname"`
	UserIcon                   string        `json:"user_icon"`
	Vlog                       bool          `json:"vlog"`
	FansVip                    bool          `json:"fansVip"`
	DisplayCommentStatus       bool          `json:"displayCommentStatus"`
	DisplayUpDownStatus        bool          `json:"displayUpDownStatus"`
	DownloadAddr               string        `json:"downloadAddr"` //视频下载地址
}

type iqiyi struct {
	Code string `json:"code"`
	Data struct {
		VP struct {
			Du  string `json:"du"`
			Tkl []struct {
				Vs []struct {
					Bid   int    `json:"bid"`
					Scrsz string `json:"scrsz"`
					Vsize int64  `json:"vsize"`
					Fs    []struct {
						L string `json:"l"`
						B int64  `json:"b"`
					} `json:"fs"`
				} `json:"vs"`
			} `json:"tkl"`
		} `json:"vp"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type iqiyiURL struct {
	L string `json:"l"`
}
