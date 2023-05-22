package iqiyi

import (
	"testing"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name          string
		args          test.Args
		playlist      bool
		shortPlaylist bool
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:     "https://www.iqiyi.com/v_1llt29653xc.html",
				Title:   "新一轮降水将至 冷空气影响中东部地区",
				Size:    2952228,
				Quality: "896x504",
			},
		},
		{
			name: "title test 1",
			args: test.Args{
				URL:     "http://www.iqiyi.com/v_19rqy2z83w.html",
				Title:   "收了创意视频2018 :58天环球飞行记",
				Size:    76186786,
				Quality: "1920x1080",
			},
		},
		{
			name: "curid test 1",
			args: test.Args{
				URL:     "https://www.iqiyi.com/v_19rro0jdls.html#curid=350289100_6e6601aae889d0b1004586a52027c321",
				Title:   "Shawn Mendes - Never Be Alone",
				Size:    79921894,
				Quality: "1920x800",
			},
		},
		{
			name: "主页视频",
			args: test.Args{
				URL:     "https://www.iqiyi.com/u/520143092318950/videos",
				Title:   "Shawn Mendes - Never Be Alone",
				Size:    79921894,
				Quality: "1920x800",
			},
			playlist: true,
		},
		{
			name: "主页视频",
			args: test.Args{
				URL:     "https://www.iqiyi.com/u/520143092318950/videos",
				Title:   "Shawn Mendes - Never Be Alone",
				Size:    79921894,
				Quality: "1920x800",
			},
		},
		{
			name: "主页短视频",
			args: test.Args{
				URL:     "https://www.iqiyi.com/u/520143092318950/videos",
				Title:   "Shawn Mendes - Never Be Alone",
				Size:    79921894,
				Quality: "1920x800",
			},
			shortPlaylist: true,
		},
		{
			name: "高清视频",
			args: test.Args{
				URL:     "https://www.iqiyi.com/v_14aoajaoai0.html",
				Title:   "Shawn Mendes - Never Be Alone",
				Size:    79921894,
				Quality: "1920x800",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New(SiteTypeIqiyi).Extract(tt.args.URL, extractors.Options{
				Playlist:      tt.playlist,
				ShortPlaylist: tt.shortPlaylist,
			})
			test.CheckError(t, err)
			test.Check(t, tt.args, data[0])
		})
	}
}
