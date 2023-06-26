package baidu

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
				URL:     "https://my.mbd.baidu.com/r/MVOgEBgMa4?f=cp&u=0d2e1f73bb319b1c",
				Title:   "看到这个字，你第一时间想起谁",
				Size:    2548137,
				Quality: "sd",
			},
		}, {
			name: "normal test",
			args: test.Args{
				URL:     "https://mbd.baidu.com/newspage/data/videolanding?nid=sv_10396018751292201200&sourceFrom=pc_feedlist",
				Title:   "看到这个字，你第一时间想起谁",
				Size:    2548137,
				Quality: "sd",
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
