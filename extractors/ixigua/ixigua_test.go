package ixigua

import (
	"testing"

	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "test 1",
			args: test.Args{
				URL:     "https://www.ixigua.com/7040164572094792195",
				Title:   "漫威斥巨资拍的《永恒族》，刚上架就被多国禁播，究竟拍了什么？",
				Quality: "1080p",
				Size:    313091514,
			},
		},
		{
			name: "test 2",
			args: test.Args{
				URL:     "https://v.ixigua.com/RedcbWM/",
				Title:   "为长生不老，竟然连小鲛人都杀@中视频伙伴计划官号",
				Quality: "1080p",
				Size:    64980732,
			},
		},
		{
			name: "test 3",
			args: test.Args{
				URL:     "https://m.toutiao.com/is/dtj1pND/",
				Title:   "卡尔：59杀4200法强小法师，点塔只需一下，W技能瞬秒对方",
				Quality: "1080p",
				Size:    468324298,
			},
		}, {
			name: "test 3",
			args: test.Args{
				URL:     "https://www.toutiao.com/video/7226902781586571816",
				Title:   "非洲小伙娶9个老婆，生活像皇帝，非洲“韦小宝”生活大揭秘！",
				Quality: "1080p",
				Size:    468324298,
			},
		}, {
			name: "主页视频",
			args: test.Args{
				URL:     "https://www.ixigua.com/home/53959751401/?list_entrance=homepage&video_card_type=shortvideo",
				Title:   "卡尔：59杀4200法强小法师，点塔只需一下，W技能瞬秒对方",
				Quality: "1080p",
				Size:    468324298,
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
	t.Log(extractors.GetTopLevelDomain("https://www.bilibili.com/"))
	t.Log(extractors.GetTopLevelDomain("https://haokan.baidu.com/"))
	t.Log(extractors.GetTopLevelDomain("https://tv.cctv.com/"))
	t.Log(extractors.GetTopLevelDomain("https://v.cctv.com"))
}
