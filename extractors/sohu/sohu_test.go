package sohu

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
)

func TestReddit(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "返回的MP4",
			args: test.Args{
				URL:   "https://tv.sohu.com/v/dXMvMzM4NDQ5OTczLzQzODcyMTM5MS5zaHRtbA==.html",
				Title: "王二妮，金婷婷同唱《白毛女》选段，精彩的演唱，永不忘的经典。",
			},
		},
		{
			name: "返回的MP4",
			args: test.Args{
				URL:   "https://tv.sohu.com/v/MjAxNzExMDkvbjYwMDI0NzEwMi5zaHRtbA==.html",
				Title: "拜见宫主大人第一季第1集",
			},
		}, {
			name: "短视频",
			args: test.Args{
				URL:   "https://tv.sohu.com/v/dXMvMzQxNzY3ODcwLzQ1NTQzMzY2NS5zaHRtbA==.html",
				Title: "拜见宫主大人第一季第1集",
			},
		},
		{
			name: "返回的列表",
			args: test.Args{
				URL:   "http://tv.sohu.com/s2017/dsjbjgzdr/",
				Title: "拜见宫主大人第一季第1集",
			},
			playlist: true,
		},
		{
			name: "主页视频",
			args: test.Args{
				URL:     "https://tv.sohu.com/user/media/video.do?uid=352155353",
				Title:   "卡尔：59杀4200法强小法师，点塔只需一下，W技能瞬秒对方",
				Quality: "1080p",
				Size:    468324298,
			},
			playlist: true,
		},
		{
			name: "主页视频",
			args: test.Args{
				URL:     "https://tv.sohu.com/user/352155353",
				Title:   "卡尔：59杀4200法强小法师，点塔只需一下，W技能瞬秒对方",
				Quality: "1080p",
				Size:    468324298,
			},
			playlist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.playlist {
				// playlist mode
				data, err := New().Extract(tt.args.URL, extractors.Options{
					Playlist:     true,
					ThreadNumber: 9,
				})
				test.CheckError(t, err)
				assert.NotEmpty(t, len(data))
			} else {
				data, err := New().Extract(tt.args.URL, extractors.Options{})
				test.CheckError(t, err)
				test.Check(t, tt.args, data[0])
			}
		})
	}
}
