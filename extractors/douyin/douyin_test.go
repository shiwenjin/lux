package douyin

import (
	"testing"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name string
		args test.Args
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://www.douyin.com/video/6967223681286278436?previous_page=main_page&tab_name=home",
				Title: "是爱情，让父子相认#陈翔六点半  #关于爱情",
			},
		},
		{
			name: "image test",
			args: test.Args{
				URL:   "https://v.douyin.com/LvCYKvV",
				Title: "黑发限定#开春必备",
			},
		}, {
			name: "image test",
			args: test.Args{
				URL:   "7.64 Ljp:/ 妈妈发现儿子会跷二郎腿，想分享一下喜悦，结果是爸爸穿反了两条腿# 爸爸带娃 # 搞笑 # 萌娃  https://v.douyin.com/U1Wjjf5/ 复制此链接，打开Dou音搜索，直接观看视频！",
				Title: "黑发限定#开春必备",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.args.URL, extractors.Options{})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
	}
}

func TestName(t *testing.T) {
	extractPlaylist("https://www.douyin.com/user/MS4wLjABAAAAhudOqr4jfJeVfF283LrAm73kgX-g2RtGIOV99KhaWsc")
}
