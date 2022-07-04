package cache

import (
	"github.com/xinzf/kit/container/kcfg"
	"testing"
)

func init() {
	kcfg.New("/Users/xiangzhi/Work/gohome/src/saas.xunray.com/lowcode/config/config.yaml")
}

func TestSearch(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		{
			name: "search",
			args: args{pattern: "*admin:auth:1*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.args.pattern)

			if err := got.Error(); err != nil {
				t.Error(err)
				return
			}

			t.Log("strs", got.Strings())
		})
	}
}
