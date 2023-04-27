package cctv

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
		//{
		//	name: "返回的MP4",
		//	args: test.Args{
		//		URL:   "https://my.tv.sohu.com/us/338449973/438721391.shtml",
		//		Title: "王二妮，金婷婷同唱《白毛女》选段，精彩的演唱，永不忘的经典。",
		//	},
		//},
		//{
		//	name: "返回的MP4",
		//	args: test.Args{
		//		URL:   "https://tv.sohu.com/v/MjAxNzExMDkvbjYwMDI0NzEwMi5zaHRtbA==.html",
		//		Title: "拜见宫主大人第一季第1集",
		//	},
		//},
		{
			name: "返回的列表",
			args: test.Args{
				URL:   "https://tv.cctv.com/2021/01/15/VIDE9wBaPD0WracfCxc6ORcM210115.shtml",
				Title: "拜见宫主大人第一季第1集",
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
