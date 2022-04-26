package uaxpl

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

type dsBrowser struct {
	UA     string `json:"user_agent"`
	Client struct {
		Type          string `json:"type"`
		Name          string `json:"name"`
		Version       string `json:"version"`
		Engine        string `json:"engine"`
		EngineVersion string `json:"engine_version"`
		Family        string `json:"family"`
	} `json:"client"`
	Headers map[string]string
}

func TestClientParse(t *testing.T) {
	t.Run("single/browser", func(t *testing.T) {
		ua := "Mozilla/5.0 (Linux; Android 8.1.0; 5059D_RU Patch/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/75.0.3770.101 YaBrowser/19.1.0.130 (lite) Mobile Safari/537.36"
		ctx := NewCtxWithSrcStr(ua)
		assertCVS(t, ctx.GetClientType(), ClientTypeBrowser)
		assertStr(t, "browser", ctx.GetBrowser(), "Yandex Browser Lite")
		assertStr(t, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130")
	})

	var ds []dsBrowser
	contents, err := os.ReadFile("dataset/browser.json")
	if err != nil {
		t.Error(err)
		return
	}
	if err = json.Unmarshal(contents, &ds); err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(ds); i++ {
		stage := &ds[i]
		t.Run("dataset/browser"+strconv.Itoa(i), func(t *testing.T) {
			ctx := AcquireWithSrcStr(stage.UA)
			if rwh, ok := stage.Headers["http-x-requested-with"]; ok {
				ctx.SetRequestedWith(rwh)
			}
			assertCVS(t, ctx.GetClientType(), ClientTypeBrowser)
			if !assertStr(t, "browser", ctx.GetBrowser(), stage.Client.Name) {
				t.Log("->", stage.UA)
			}
			Release(ctx)
		})
	}
}

func BenchmarkClientParse(b *testing.B) {
	b.Run("single/browser", func(b *testing.B) {
		b.ReportAllocs()
		ua := "Mozilla/5.0 (Linux; Android 8.1.0; 5059D_RU Patch/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/75.0.3770.101 YaBrowser/19.1.0.130 (lite) Mobile Safari/537.36"
		for i := 0; i < b.N; i++ {
			ctx := AcquireWithSrcStr(ua)
			assertCVS(b, ctx.GetClientType(), ClientTypeBrowser)
			assertStr(b, "browser", ctx.GetBrowser(), "Yandex Browser Lite")
			assertStr(b, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130")
			Release(ctx)
		}
	})
}
