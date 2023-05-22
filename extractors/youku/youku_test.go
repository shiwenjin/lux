package youku

import (
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "长视频",
			args: test.Args{
				URL: "https://v.youku.com/v_show/id_XNDEzNTc5NTY4OA==.html?s=d88efd308ea811e69e06&spm=a2hje.13141534.1_3.d_1_1&scm=20140719.apircmd.240015.video_XNDEzNTc5NTY4OA==",
			},
		},
		{
			name: "短视频",
			args: test.Args{
				URL: "https://v.youku.com/v_show/id_XNTk2MjM4ODIyMA==.html",
			},
		}, {
			name: "主页视频",
			args: test.Args{
				URL: "https://www.youku.com/profile/index?uid=UMTIwMDA0MDIwMTI",
			},
			playlist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.args.URL, extractors.Options{
				Playlist: tt.playlist,
			})

			assert.NotEmpty(t, data)
			assert.NoError(t, err)
		})
	}
}
