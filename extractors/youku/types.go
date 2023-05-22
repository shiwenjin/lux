package youku

type JoinParams struct {
	Jsv            string `json:"jsv"`
	AppKey         string `json:"appKey"`
	T              string `json:"t"`
	Sign           string `json:"sign"`
	Api            string `json:"api"`
	V              string `json:"v"`
	Timeout        string `json:"timeout"`
	YKPid          string `json:"YKPid"`
	YKLoginRequest string `json:"YKLoginRequest"`
	AntiFlood      string `json:"AntiFlood"`
	AntiCreep      string `json:"AntiCreep"`
	Type           string `json:"type"`
	DataType       string `json:"dataType"`
	Callback       string `json:"callback"`
	Data           string `json:"data"`
}

type StealParams struct {
	Ccode    string `json:"ccode"`
	ClientIp string `json:"client_ip"`
	Utid     string `json:"utid"`
	ClientTs string `json:"client_ts"`
	Version  string `json:"version"`
	Ckey     string `json:"ckey"`
}

type Data struct {
	StealParams string `json:"steal_params,omitempty"`
	BizParams   string `json:"biz_params,omitempty"`
	AdParams    string `json:"ad_params,omitempty"`
}

type AdParams struct {
	Vs        string `json:"vs"`
	Pver      string `json:"pver"`
	Sver      string `json:"sver"`
	Site      int    `json:"site"`
	Aw        string `json:"aw"`
	Fu        int    `json:"fu"`
	D         string `json:"d"`
	Bt        string `json:"bt"`
	Os        string `json:"os"`
	Osv       string `json:"osv"`
	Dq        string `json:"dq"`
	Atm       string `json:"atm"`
	Partnerid string `json:"partnerid"`
	Wintype   string `json:"wintype"`
	Isvert    int    `json:"isvert"`
	Vip       int    `json:"vip"`
	P         int    `json:"p"`
	Rst       string `json:"rst"`
	Needbf    int    `json:"needbf"`
	Avs       string `json:"avs"`
}

type BizParams struct {
	Vid           string `json:"vid"`
	PlayAbility   string `json:"play_ability"`
	CurrentShowid string `json:"current_showid"`
	PreferClarity string `json:"preferClarity"`
	Extag         string `json:"extag"`
	MasterM3U8    string `json:"master_m3u8"`
	MediaType     string `json:"media_type"`
	AppVer        string `json:"app_ver"`
	DrmType       string `json:"drm_type"`
	KeyIndex      string `json:"key_index"`
}
