package weishi

type weishiResp struct {
	CollectionInfo interface{} `json:"collectionInfo"`
	FeedsList      []struct {
		Index    int    `json:"index"`
		Id       string `json:"id"`
		PosterId string `json:"posterId"`
		Poster   struct {
			Id         string `json:"id"`
			Type       int    `json:"type"`
			Uid        string `json:"uid"`
			Createtime int    `json:"createtime"`
			Nick       string `json:"nick"`
			Avatar     string `json:"avatar"`
		} `json:"poster"`
		Video struct {
			FileId   string `json:"fileId"`
			Duration int    `json:"duration"`
		} `json:"video"`
		Images []struct {
			Url          string `json:"url"`
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			Type         int    `json:"type"`
			SpriteWidth  int    `json:"sprite_width"`
			SpriteHeight int    `json:"sprite_height"`
			SpriteSpan   int    `json:"sprite_span"`
			Priority     int    `json:"priority"`
			PhotoRgb     string `json:"photo_rgb"`
			Format       string `json:"format"`
		} `json:"images"`
		DingCount       int    `json:"dingCount"`
		TotalCommentNum int64  `json:"totalCommentNum"`
		MaterialDesc    string `json:"materialDesc"`
		IsDing          int    `json:"isDing"`
		PlayNum         int64  `json:"playNum"`
		VideoUrl        string `json:"videoUrl"`
		VideoCover      string `json:"videoCover"`
		Reserve         struct {
			Field1 struct {
				SchemeMultiStyle       string `json:"schemeMultiStyle"`
				SchemeDynamicStyleIcon string `json:"schemeDynamicStyleIcon"`
				SchemeDynamicStyleText string `json:"schemeDynamicStyleText"`
				IconDynamicStyle       string `json:"iconDynamicStyle"`
				TextDynamicStyle       string `json:"textDynamicStyle"`
				ExposureMaterialName   string `json:"exposureMaterialName"`
				MaterialId             string `json:"materialId"`
				H5ShareTitleContent    string `json:"h5ShareTitleContent"`
				H5ShareButtonContent   string `json:"h5ShareButtonContent"`
				CommonLabelScheme      string `json:"commonLabelScheme"`
				MaterialCategory       string `json:"materialCategory"`
				IosMinVersionStr       string `json:"iosMinVersionStr"`
				AndroidMinVersionStr   string `json:"androidMinVersionStr"`
			} `json:"51"`
		} `json:"reserve"`
		ShareInfo struct {
			BodyMap struct {
				Field1 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"0"`
				Field2 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"1"`
				Field3 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"2"`
				Field4 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"3"`
				Field5 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"4"`
				Field6 struct {
					Title    string `json:"title"`
					Desc     string `json:"desc"`
					ImageUrl string `json:"image_url"`
					Url      string `json:"url"`
				} `json:"5"`
			} `json:"bodyMap"`
			JumpUrl string `json:"jumpUrl"`
		} `json:"shareInfo"`
		TopicId string `json:"topicId"`
		Topic   struct {
			Id             string        `json:"id"`
			Name           string        `json:"name"`
			ThumbUrl1      string        `json:"thumbUrl1"`
			ThumbUrl2      string        `json:"thumbUrl2"`
			ThumbUrl3      string        `json:"thumbUrl3"`
			Detail         string        `json:"detail"`
			Createtime     int           `json:"createtime"`
			FeedlistTimeId string        `json:"feedlist_time_id"`
			FeedlistHotId  string        `json:"feedlist_hot_id"`
			MaterialIds    []interface{} `json:"material_ids"`
			Mask           int           `json:"mask"`
			Type           int           `json:"type"`
			Reserve        struct {
			} `json:"reserve"`
			ViewNum    int   `json:"view_num"`
			StartTime  int   `json:"start_time"`
			EndTime    int   `json:"end_time"`
			AppVersion int   `json:"appVersion"`
			WorkNum    int   `json:"workNum"`
			LikeNum    int64 `json:"likeNum"`
			Person     struct {
				Id                string `json:"id"`
				Type              int    `json:"type"`
				Uid               string `json:"uid"`
				Createtime        int    `json:"createtime"`
				Nick              string `json:"nick"`
				Avatar            string `json:"avatar"`
				Sex               int    `json:"sex"`
				FeedlistTimeId    string `json:"feedlist_time_id"`
				FeedlistHotId     string `json:"feedlist_hot_id"`
				RelatedFeedlistId string `json:"related_feedlist_id"`
				FollowerlistId    string `json:"followerlist_id"`
				InteresterlistId  string `json:"interesterlist_id"`
				ChatlistId        string `json:"chatlist_id"`
				RichFlag          int    `json:"rich_flag"`
				Age               int    `json:"age"`
				Address           string `json:"address"`
				Wealth            struct {
					FlowerNum int `json:"flower_num"`
					Score     int `json:"score"`
				} `json:"wealth"`
				Background        string `json:"background"`
				Status            string `json:"status"`
				FollowStatus      int    `json:"followStatus"`
				ChartScore        int    `json:"chartScore"`
				ChartRank         int    `json:"chartRank"`
				FeedGoldNum       int    `json:"feedGoldNum"`
				AvatarUpdatetime  int    `json:"avatar_updatetime"`
				DescFromOperator  string `json:"desc_from_operator"`
				SyncContent       int    `json:"sync_content"`
				FeedlistPraiseId  string `json:"feedlist_praise_id"`
				Settingmask       int    `json:"settingmask"`
				Originalavatar    string `json:"originalavatar"`
				BlockTime         string `json:"block_time"`
				Grade             int    `json:"grade"`
				Medal             int    `json:"medal"`
				BlockReason       string `json:"block_reason"`
				Qq                int    `json:"qq"`
				RecommendReason   string `json:"recommendReason"`
				LastUpdateFeedNum int    `json:"lastUpdateFeedNum"`
				Updateinfo        struct {
					Flag int    `json:"flag"`
					Tip  string `json:"tip"`
					Num  int    `json:"num"`
				} `json:"updateinfo"`
				NickUpdatetime     int    `json:"nick_updatetime"`
				LastDownloadAvatar string `json:"lastDownloadAvatar"`
				RealName           string `json:"realName"`
				PinyinFirst        string `json:"pinyin_first"`
				CertifDesc         string `json:"certif_desc"`
				PrivateInfo        struct {
					PhoneNum string `json:"phone_num"`
					Name     string `json:"name"`
					IdNum    string `json:"id_num"`
				} `json:"privateInfo"`
				ExternInfo struct {
					MpEx struct {
					} `json:"mpEx"`
					BindAcct  []interface{} `json:"bind_acct"`
					BgPicUrl  string        `json:"bgPicUrl"`
					LevelInfo struct {
						Level           int `json:"level"`
						Score           int `json:"score"`
						PrevUpgradeTime int `json:"prev_upgrade_time"`
					} `json:"level_info"`
					WeishiId            string `json:"weishiId"`
					WeishiidModifyCount string `json:"weishiid_modify_count"`
					WatermarkType       int    `json:"watermark_type"`
					RealNick            string `json:"real_nick"`
					CmtLevel            struct {
						Level           int `json:"level"`
						Cmtscore        int `json:"cmtscore"`
						Dingscore       int `json:"dingscore"`
						PrevUpgradeTime int `json:"prev_upgrade_time"`
					} `json:"cmt_level"`
					FlexibilityFlag int `json:"flexibility_flag"`
					LiveStatus      int `json:"live_status"`
					NowLiveRoomId   int `json:"now_live_room_id"`
					MedalInfo       struct {
						TotalScore int           `json:"total_score"`
						MedalList  []interface{} `json:"medal_list"`
					} `json:"medal_info"`
					H5HasLogin   int    `json:"h5_has_login"`
					RelationType int    `json:"relation_type"`
					Feedid       string `json:"feedid"`
					CoverType    int    `json:"cover_type"`
					IndustryType int    `json:"industry_type"`
					InditeType   int    `json:"indite_type"`
				} `json:"extern_info"`
				CertifData struct {
					CertifIcon    string `json:"certif_icon"`
					CertifJumpurl string `json:"certif_jumpurl"`
				} `json:"certifData"`
				IsShowPOI    int `json:"isShowPOI"`
				IsShowGender int `json:"isShowGender"`
				FormatAddr   struct {
					Country  string `json:"country"`
					Province string `json:"province"`
					City     string `json:"city"`
				} `json:"formatAddr"`
				AuthorizeTime int `json:"authorize_time"`
				ActivityInfo  struct {
					InvitePersonid string `json:"invitePersonid"`
				} `json:"activityInfo"`
				SpecialIdentity struct {
				} `json:"special_identity"`
				TmpMark      int `json:"tmpMark"`
				PmtMark      int `json:"pmtMark"`
				IndustryInfo struct {
					PrimaryIndustry struct {
						IndustryId   int    `json:"industry_id"`
						IndustryDesc string `json:"industry_desc"`
					} `json:"primary_industry"`
					SecondaryIndustry struct {
						IndustryId   int    `json:"industry_id"`
						IndustryDesc string `json:"industry_desc"`
					} `json:"secondary_industry"`
				} `json:"industry_info"`
				Homeland struct {
					Country  string `json:"country"`
					Province string `json:"province"`
					City     string `json:"city"`
				} `json:"homeland"`
			} `json:"person"`
			FeedId            string `json:"feed_id"`
			PendantMaterialId string `json:"pendant_material_id"`
			MusicMaterialId   string `json:"music_material_id"`
			MusicInfo         struct {
				Id              string        `json:"id"`
				Name            string        `json:"name"`
				Desc            string        `json:"desc"`
				Type            string        `json:"type"`
				ThumbUrl        string        `json:"thumbUrl"`
				Version         int           `json:"version"`
				MiniSptVersion  int           `json:"miniSptVersion"`
				PackageUrl      string        `json:"packageUrl"`
				FeedlistTimeId  string        `json:"feedlist_time_id"`
				FeedlistHotId   string        `json:"feedlist_hot_id"`
				TopicIds        []interface{} `json:"topic_ids"`
				Mask            int           `json:"mask"`
				ShortName       string        `json:"shortName"`
				RichFlag        int           `json:"rich_flag"`
				EffectId        string        `json:"effectId"`
				Rgbcolor        string        `json:"rgbcolor"`
				IsCollected     int           `json:"isCollected"`
				BubbleStartTime int           `json:"bubbleStartTime"`
				BubbleEndTime   int           `json:"bubbleEndTime"`
				CollectTime     int           `json:"collectTime"`
				SdkInfo         struct {
					IsSdk            int           `json:"isSdk"`
					SdkMinVersion    int           `json:"sdkMinVersion"`
					SdkMaxVersion    int           `json:"sdkMaxVersion"`
					SdkMinSptVersion int           `json:"sdkMinSptVersion"`
					SdkAndroidGrays  []interface{} `json:"sdkAndroidGrays"`
					SdkMinVersionStr string        `json:"sdkMinVersionStr"`
					SdkMaxVersionStr string        `json:"sdkMaxVersionStr"`
				} `json:"sdkInfo"`
				BigThumbUrl string        `json:"bigThumbUrl"`
				Priority    int           `json:"priority"`
				MusicIDs    []interface{} `json:"musicIDs"`
				Platform    string        `json:"platform"`
				Reserve     struct {
				} `json:"reserve"`
				Category                string        `json:"category"`
				ShootingTips            string        `json:"shooting_tips"`
				VecSubcategory          []interface{} `json:"vec_subcategory"`
				DemoVideoList           []interface{} `json:"demo_video_list"`
				RecommendTemplateTags   []interface{} `json:"recommendTemplateTags"`
				InspirationButtonText   string        `json:"inspirationButtonText"`
				InspirationButtonSchema string        `json:"inspirationButtonSchema"`
				RelationMaterialId      string        `json:"relationMaterialId"`
				MoreMaterialLink        string        `json:"moreMaterialLink"`
				StartTime               int           `json:"startTime"`
				EndTime                 int           `json:"endTime"`
				RandomPackageUrls       struct {
				} `json:"randomPackageUrls"`
				HideType   int `json:"hideType"`
				FollowInfo struct {
					IsFollowShotShown         int    `json:"isFollowShotShown"`
					SchemeType                int    `json:"schemeType"`
					Scheme                    string `json:"scheme"`
					IconUrl                   string `json:"iconUrl"`
					SharePageMark             string `json:"sharePageMark"`
					ShareButtonMark           string `json:"shareButtonMark"`
					Name                      string `json:"name"`
					FollowNewPagIconUrl       string `json:"followNewPagIconUrl"`
					FollowNewDefaultIconUrl   string `json:"followNewDefaultIconUrl"`
					FollowNewIconScheme       string `json:"followNewIconScheme"`
					FollowNewNameScheme       string `json:"followNewNameScheme"`
					FollowNewMultiLabelScheme string `json:"followNewMultiLabelScheme"`
					CategoryFollowShotInfo    struct {
						IsFollowShotShown int `json:"isFollowShotShown"`
						Priority          int `json:"priority"`
						FollowShotStyle   int `json:"followShotStyle"`
					} `json:"categoryFollowShotInfo"`
				} `json:"followInfo"`
				State       int    `json:"state"`
				Deleted     int    `json:"deleted"`
				Packages    string `json:"packages"`
				PackageUrls struct {
				} `json:"packageUrls"`
				MaterialPackageUrls struct {
				} `json:"materialPackageUrls"`
				RoundThumbUrl    string `json:"roundThumbUrl"`
				BigRoundThumbUrl string `json:"bigRoundThumbUrl"`
				Title            string `json:"title"`
				CarID            string `json:"carID"`
				ComposedInfo     struct {
					VideoDuration       int           `json:"videoDuration"`
					IncludeMIDList      []interface{} `json:"includeMIDList"`
					AbilityList         []interface{} `json:"abilityList"`
					SlotInfo            string        `json:"slotInfo"`
					ThumbResolution     string        `json:"thumbResolution"`
					SlotInfoList        []interface{} `json:"slotInfoList"`
					IncludeMaterialInfo struct {
					} `json:"includeMaterialInfo"`
				} `json:"composedInfo"`
				AuthorID             string `json:"authorID"`
				UseCount             int    `json:"useCount"`
				CardID               string `json:"cardID"`
				CompressedPackageUrl string `json:"compressedPackageUrl"`
			} `json:"music_info"`
			PendantMaterialIdIos string `json:"pendant_material_id_ios"`
			MediaMaterialUrl     string `json:"media_material_url"`
			BubbleStartTime      int    `json:"bubble_start_time"`
			BubbleEndTime        int    `json:"bubble_end_time"`
			BubbleCopywrite      string `json:"bubble_copywrite"`
			Rgbcolor             int    `json:"rgbcolor"`
			Lplaynum             int    `json:"lplaynum"`
			QqMusicInfo          struct {
				AlbumInfo struct {
					UiId    int    `json:"uiId"`
					StrMid  string `json:"strMid"`
					StrName string `json:"strName"`
					StrPic  string `json:"strPic"`
				} `json:"albumInfo"`
				SingerInfo struct {
					UiId         int           `json:"uiId"`
					StrMid       string        `json:"strMid"`
					StrName      string        `json:"strName"`
					StrPic       string        `json:"strPic"`
					StrSchema    string        `json:"strSchema"`
					OtherSingers []interface{} `json:"otherSingers"`
				} `json:"singerInfo"`
				SongInfo struct {
					UiId               int    `json:"uiId"`
					StrMid             string `json:"strMid"`
					StrName            string `json:"strName"`
					StrGenre           string `json:"strGenre"`
					IIsOnly            int    `json:"iIsOnly"`
					StrLanguage        string `json:"strLanguage"`
					IPlayable          int    `json:"iPlayable"`
					ITrySize           int    `json:"iTrySize"`
					ITryBegin          int    `json:"iTryBegin"`
					ITryEnd            int    `json:"iTryEnd"`
					IPlayTime          int    `json:"iPlayTime"`
					StrH5Url           string `json:"strH5Url"`
					StrPlayUrl         string `json:"strPlayUrl"`
					StrPlayUrlStandard string `json:"strPlayUrlStandard"`
					StrPlayUrlHq       string `json:"strPlayUrlHq"`
					StrPlayUrlSq       string `json:"strPlayUrlSq"`
					ISize              int    `json:"iSize"`
					ISizeStandard      int    `json:"iSizeStandard"`
					ISizeHq            int    `json:"iSizeHq"`
					ISizeSq            int    `json:"iSizeSq"`
					Copyright          int    `json:"copyright"`
					ISource            int    `json:"iSource"`
					PublicTime         string `json:"publicTime"`
					SongTitle          string `json:"songTitle"`
					SongDescription    string `json:"songDescription"`
					State              int    `json:"state"`
					Deleted            int    `json:"deleted"`
					StartTime          string `json:"startTime"`
					EndTime            string `json:"endTime"`
				} `json:"songInfo"`
				LyricInfo struct {
					UiSongId      int    `json:"uiSongId"`
					StrSongMid    string `json:"strSongMid"`
					StrFormat     string `json:"strFormat"`
					StrLyric      string `json:"strLyric"`
					StrMatchLyric string `json:"strMatchLyric"`
				} `json:"lyricInfo"`
				ConfInfo struct {
					IType               int    `json:"iType"`
					IStartPos           int    `json:"iStartPos"`
					StrLabel            string `json:"strLabel"`
					IsCollected         int    `json:"isCollected"`
					CollectTime         int    `json:"collectTime"`
					Exclusive           int    `json:"exclusive"`
					FollowFeed          string `json:"followFeed"`
					UseCount            int    `json:"useCount"`
					TogetherFeed        string `json:"togetherFeed"`
					TogetherType        int    `json:"togetherType"`
					FeedUseType         int    `json:"feedUseType"`
					DefaultFeedPosition int    `json:"defaultFeedPosition"`
					DefaultTogetherFeed int    `json:"defaultTogetherFeed"`
					BubbleStartTime     int    `json:"bubbleStartTime"`
					BubbleEndTime       int    `json:"bubbleEndTime"`
					SongLabels          struct {
					} `json:"songLabels"`
					SongLabelCategory struct {
					} `json:"songLabelCategory"`
					IsStuckPoint      bool   `json:"isStuckPoint"`
					StuckPointJsonUrl string `json:"stuckPointJsonUrl"`
					TrackBeatInfo     struct {
						TrackBeatFinished int `json:"trackBeatFinished"`
						Vocal             struct {
							JsonURL      string `json:"jsonURL"`
							AudioFileURL string `json:"audioFileURL"`
						} `json:"vocal"`
						Drums struct {
							JsonURL      string `json:"jsonURL"`
							AudioFileURL string `json:"audioFileURL"`
						} `json:"drums"`
						Accompaniment struct {
							JsonURL      string `json:"jsonURL"`
							AudioFileURL string `json:"audioFileURL"`
						} `json:"accompaniment"`
						Bass struct {
							JsonURL      string `json:"jsonURL"`
							AudioFileURL string `json:"audioFileURL"`
						} `json:"bass"`
					} `json:"trackBeatInfo"`
					ExtraInfo string `json:"extraInfo"`
				} `json:"confInfo"`
				SubtitleInfo struct {
					UiSongId      int    `json:"uiSongId"`
					StrSongMid    string `json:"strSongMid"`
					StrFormat     string `json:"strFormat"`
					StrLyric      string `json:"strLyric"`
					StrMatchLyric string `json:"strMatchLyric"`
				} `json:"subtitleInfo"`
				Foreignlyric struct {
					UiSongId      int    `json:"uiSongId"`
					StrSongMid    string `json:"strSongMid"`
					StrFormat     string `json:"strFormat"`
					StrLyric      string `json:"strLyric"`
					StrMatchLyric string `json:"strMatchLyric"`
				} `json:"foreignlyric"`
				RecommendInfo struct {
					TraceStr      string `json:"traceStr"`
					AnalyseResult string `json:"analyse_result"`
					RecomReason   string `json:"recom_reason"`
				} `json:"recommendInfo"`
				UnplayableInfo struct {
					UnplayableCode int    `json:"unplayableCode"`
					UnplayableMsg  string `json:"unplayableMsg"`
				} `json:"unplayableInfo"`
				LabelInfo       []interface{} `json:"labelInfo"`
				MusicType       int           `json:"musicType"`
				MusicSrcType    int           `json:"musicSrcType"`
				CacheUpdateTime int           `json:"cacheUpdateTime"`
				State           int           `json:"state"`
			} `json:"qqMusicInfo"`
			TopicType           int    `json:"topicType"`
			TopicSource         int    `json:"topicSource"`
			Creater             string `json:"creater"`
			LastOperator        string `json:"lastOperator"`
			SecurityAuditstate  int    `json:"securityAuditstate"`
			SecurityAuditReason string `json:"securityAuditReason"`
			ManualAuditstate    int    `json:"manualAuditstate"`
			ManualAuditReason   string `json:"manualAuditReason"`
			Status              int    `json:"status"`
			UpdateTime          int    `json:"updateTime"`
			NewAppVersion       string `json:"newAppVersion"`
			TopicMusicName      string `json:"TopicMusicName"`
			PendantMaterialCate string `json:"pendant_material_cate"`
			Schema              string `json:"schema"`
			SchemaType          int    `json:"schemaType"`
			ExternalLink        struct {
				LinkUrl  string `json:"link_url"`
				LinkName string `json:"link_name"`
			} `json:"external_link"`
			InteractiveNews struct {
				InteractiveDetails []interface{} `json:"interactive_details"`
			} `json:"interactive_news"`
			ActivityInfo struct {
				Name          string `json:"name"`
				Label         string `json:"label"`
				SubName       string `json:"sub_name"`
				BtnTxt        string `json:"btn_txt"`
				ShowStartTime int    `json:"show_start_time"`
				ShowEndTime   int    `json:"show_end_time"`
				StartTime     int    `json:"start_time"`
				EndTime       int    `json:"end_time"`
				RuleInfo      struct {
					RuleDetails []interface{} `json:"rule_details"`
				} `json:"rule_info"`
				NeedShow          int    `json:"need_show"`
				ResourceBackColor string `json:"resource_back_color"`
				Status            int    `json:"status"`
			} `json:"activity_info"`
			PublishInfo struct {
				BtnStyle          int    `json:"btn_style"`
				BtnText           string `json:"btn_text"`
				BtnPic            string `json:"btn_pic"`
				BlueCollarPublish struct {
					DefaultCallModel    int           `json:"default_call_model"`
					DefaultCamera       int           `json:"default_camera"`
					BlueCollarMaterials []interface{} `json:"blue_collar_materials"`
					PendantId           string        `json:"pendant_id"`
					MusicId             string        `json:"music_id"`
					SongListId          string        `json:"song_list_id"`
					TeleprompterDesc    string        `json:"teleprompter_desc"`
				} `json:"blue_collar_publish"`
			} `json:"publish_info"`
			LatestPublishTime int `json:"latest_publish_time"`
			CollarType        int `json:"collar_type"`
			BlueCollar        struct {
				BlueCollarTags          []interface{} `json:"blue_collar_tags"`
				BlueCollarProfessionTag string        `json:"blue_collar_profession_tag"`
			} `json:"blue_collar"`
			UserGroupId string `json:"user_group_id"`
		} `json:"topic"`
		FeedDesc string `json:"feedDesc"`
		GeoInfo  struct {
			Country   string `json:"country"`
			Province  string `json:"province"`
			City      string `json:"city"`
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
			Altitude  int    `json:"altitude"`
			District  string `json:"district"`
			Name      string `json:"name"`
			Distance  int    `json:"distance"`
			PolyGeoID string `json:"polyGeoID"`
		} `json:"geoInfo"`
		MusicId             string `json:"musicId"`
		FeedDescWithat      string `json:"feedDescWithat"`
		FeedRecommendReason string `json:"feedRecommendReason"`
		ExternInfo          struct {
			MpEx struct {
			} `json:"mpEx"`
		} `json:"externInfo"`
		CollectionId string `json:"collectionId"`
		RichDing     struct {
			H5RichDingDisplay struct {
				Url   string `json:"url"`
				Text  string `json:"text"`
				Count int    `json:"count"`
			} `json:"h5RichDingDisplay"`
		} `json:"richDing"`
		VideoConfig struct {
		} `json:"videoConfig"`
		VideoSpecUrls struct {
			Field1 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"0"`
			Field2 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"1"`
			Field3 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"2"`
			Field4 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"5"`
			Field5 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"6"`
			Field6 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"8"`
			Field7 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"11"`
			Field8 struct {
				Url           string `json:"url"`
				Size          int    `json:"size"`
				Hardorsoft    int    `json:"hardorsoft"`
				RecommendSpec int    `json:"recommendSpec"`
				HaveWatermark int    `json:"haveWatermark"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				VideoCoding   int    `json:"videoCoding"`
				VideoQuality  int    `json:"videoQuality"`
				ExternInfo    struct {
				} `json:"externInfo"`
				StaticCover struct {
					Url          string `json:"url"`
					Width        int    `json:"width"`
					Height       int    `json:"height"`
					Type         int    `json:"type"`
					SpriteWidth  int    `json:"sprite_width"`
					SpriteHeight int    `json:"sprite_height"`
					SpriteSpan   int    `json:"sprite_span"`
					Priority     int    `json:"priority"`
					PhotoRgb     string `json:"photo_rgb"`
					Format       string `json:"format"`
				} `json:"staticCover"`
				Fps int `json:"fps"`
			} `json:"999"`
		} `json:"videoSpecUrls"`
	} `json:"feedsList"`
	IsCollection        bool        `json:"isCollection"`
	ShowCollection      bool        `json:"showCollection"`
	ActiveIndex         int         `json:"activeIndex"`
	ShareUi             interface{} `json:"shareUi"`
	ActivePlayerPlaying bool        `json:"activePlayerPlaying"`
	IsNewUI             bool        `json:"isNewUI"`
	ZzConfig            interface{} `json:"zzConfig"`
	IsInstall           interface{} `json:"isInstall"`
	GrowthTestData      struct {
		CollectionPullNew bool `json:"collectionPullNew"`
		CollectionPullOld bool `json:"collectionPullOld"`
		CollectionStyle   int  `json:"collectionStyle"`
		DataSource        int  `json:"dataSource"`
		Ext               struct {
		} `json:"ext"`
		JumpMiniPlayer int    `json:"jumpMiniPlayer"`
		JumpNewPage    bool   `json:"jumpNewPage"`
		ModuleDataSrc  int    `json:"moduleDataSrc"`
		ModulePullNew  bool   `json:"modulePullNew"`
		ModulePullOld  bool   `json:"modulePullOld"`
		Msg            string `json:"msg"`
		PlayBtnPullNew bool   `json:"playBtnPullNew"`
		PlayBtnPullOld bool   `json:"playBtnPullOld"`
		Ret            int    `json:"ret"`
		ShowModule     bool   `json:"showModule"`
		StrategyID     string `json:"strategyID"`
	} `json:"growthTestData"`
	Banner struct {
		Person interface{} `json:"person"`
	} `json:"banner"`
	ErrorMsg      string `json:"errorMsg"`
	IsGrowthPopup string `json:"isGrowthPopup"`
}
