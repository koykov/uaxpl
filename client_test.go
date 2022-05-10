package uaxpl

import (
	"strconv"
	"testing"
)

type browserDS struct {
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
		assertCVS(t, ctx.GetClientType(), ClientTypeBrowser, false)
		assertStr(t, "browser", ctx.GetBrowser(), "Yandex Browser Lite", false)
		assertVer(t, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130", false)
		assertStr(t, "engine", ctx.GetEngine(), "Blink", false)
		assertVer(t, "engine version", ctx.GetEngineVersionString(), "75.0.3770.101", false)
	})

	ds, err := testLoadBrowserDS("testdata/browser.json")
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(ds); i++ {
		stage := &ds[i]
		t.Run("dataset/browser"+strconv.Itoa(i), func(t *testing.T) {
			ok := true
			ctx := AcquireWithSrcStr(stage.UA)
			ctx.SetRequestedWith(testGetRWH(stage))
			ok = ok && assertCVS(t, ctx.GetClientType(), ClientTypeBrowser, false)
			ok = ok && assertStr(t, "browser", ctx.GetBrowser(), stage.Client.Name, false)
			ok = ok && assertVer(t, "browser version", ctx.GetBrowserVersionString(), stage.Client.Version, false)
			ok = ok && assertStr(t, "engine", ctx.GetEngine(), stage.Client.Engine, false)
			ok = ok && assertVer(t, "engine", ctx.GetEngineVersionString(), stage.Client.EngineVersion, false)
			if !ok {
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
			assertCVS(b, ctx.GetClientType(), ClientTypeBrowser, true)
			assertStr(b, "browser", ctx.GetBrowser(), "Yandex Browser Lite", true)
			assertVer(b, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130", true)
			assertStr(b, "engine", ctx.GetEngine(), "Blink", false)
			assertVer(b, "engine version", ctx.GetEngineVersionString(), "75.0.3770.101", false)
			Release(ctx)
		}
	})
	b.Run("dataset/browser", func(b *testing.B) {
		ds, err := testLoadBrowserDS("testdata/browser.json")
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			stage := &ds[i%len(ds)]
			ctx := AcquireWithSrcStr(stage.UA)
			ctx.SetRequestedWith(testGetRWH(stage))
			assertCVS(b, ctx.GetClientType(), ClientTypeBrowser, true)
			assertStr(b, "browser", ctx.GetBrowser(), stage.Client.Name, true)
			assertVer(b, "browser version", ctx.GetBrowserVersionString(), stage.Client.Version, true)
			assertStr(b, "engine", ctx.GetEngine(), stage.Client.Engine, true)
			assertVer(b, "engine version", ctx.GetEngineVersionString(), stage.Client.EngineVersion, true)
			Release(ctx)
		}
	})
}
