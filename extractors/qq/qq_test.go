package qq

import (
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:     "https://page.om.qq.com/page/OrJrlquA4iPzVETMWh-7ZpqQ0",
				Title:   "世界杯第一期：100秒速成！“伪球迷”世界杯生存指南",
				Size:    23759683,
				Quality: "蓝光;(1080P)",
			},
		}, {
			name: "normal test",
			args: test.Args{
				URL:     "https://v.qq.com/x/page/n0687peq62x.html",
				Title:   "世界杯第一期：100秒速成！“伪球迷”世界杯生存指南",
				Size:    23759683,
				Quality: "蓝光;(1080P)",
			},
		},
		{
			name: "movie and vid test",
			args: test.Args{
				URL:     "https://v.qq.com/x/cover/e5qmd3z5jr0uigk.html",
				Title:   "赌侠（粤语版）",
				Size:    1046910811,
				Quality: "超清;(720P)",
			},
		},
		{
			name: "短视频",
			args: test.Args{
				URL:     "https://v.qq.com/x/page/t0046y8r0bs.html",
				Title:   "跟郭采洁逛AMI影展 分享“家”的意义",
				Size:    14112979,
				Quality: "超清;(720P)",
			},
		}, {
			name: "长视频列表",
			args: test.Args{
				URL:     "https://v.qq.com/x/cover/mzc00200fhhxx8d/v0046des2cr.html",
				Title:   "漫长的季节",
				Size:    14112979,
				Quality: "超清;(720P)",
			},
		}, {
			name: "短视频列表",
			args: test.Args{
				URL:     "https://v.qq.com/biu/creator/home?vcuid=9000009005",
				Title:   "腾讯时尚",
				Size:    14112979,
				Quality: "超清;(720P)",
			},
			playlist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.args.URL, extractors.Options{Playlist: tt.playlist})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
	}
}

func TestName(t *testing.T) {
	getVideoApiKey("e0765r4mwcr", `https://v.qq.com/x/cover/2aya3ibdmft6vdw/e0765r4mwcr.html`, "", "srt", "hls", "hd")

}

func TestM3u8(t *testing.T) {
	t.Log(getGuid(32))
}
