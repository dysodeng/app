package wx

import (
	"testing"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

func TestMiniProgramRetriesAfterFailureAndRefreshesChangedConfig(t *testing.T) {
	previous := config.GlobalConfig
	config.GlobalConfig = &config.Config{}
	t.Cleanup(func() { config.GlobalConfig = previous })

	client, err := MiniProgram()
	if err == nil {
		t.Fatal("expected missing credentials to return an error")
	}
	if client != nil {
		t.Fatal("expected no client when credentials are missing")
	}

	config.GlobalConfig.ThirdParty.Wx.MiniProgram.AppId = "first-app"
	config.GlobalConfig.ThirdParty.Wx.MiniProgram.Secret = "first-secret"
	first, err := MiniProgram()
	if err != nil || first == nil {
		t.Fatalf("expected initialization to recover, got client=%v error=%v", first, err)
	}
	if first.Config().AppID != "first-app" || first.Config().AppSecret != "first-secret" {
		t.Fatalf("unexpected first client config: %+v", first.Config())
	}

	again, err := MiniProgram()
	if err != nil || again != first {
		t.Fatalf("expected to reuse client for unchanged config, got client=%v error=%v", again, err)
	}

	updated := &config.Config{}
	updated.ThirdParty.Wx.MiniProgram.AppId = "second-app"
	updated.ThirdParty.Wx.MiniProgram.Secret = "first-secret"
	config.GlobalConfig = updated
	second, err := MiniProgram()
	if err != nil || second == nil || second == first {
		t.Fatalf("expected a new client after config change, got client=%v error=%v", second, err)
	}
	if second.Config().AppID != "second-app" {
		t.Fatalf("expected second app ID, got %q", second.Config().AppID)
	}
}
