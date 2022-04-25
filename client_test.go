package uaxpl

import "testing"

func TestClientParse(t *testing.T) {
	t.Run("single/browser", func(t *testing.T) {
		ua := "Mozilla/5.0 (Linux; Android 8.1.0; 5059D_RU Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/75.0.3770.101 YaBrowser/19.1.0.130 (lite) Mobile Safari/537.36"
		ctx := NewCtxWithSrcStr(ua)
		assertCVS(t, ctx.GetClientType(), ClientTypeBrowser)
		assertStr(t, "browser", ctx.GetBrowser(), "Yandex Browser Lite")
		assertStr(t, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130")
	})
}

func BenchmarkClientParse(b *testing.B) {
	b.Run("single/browser", func(b *testing.B) {
		b.ReportAllocs()
		ua := "Mozilla/5.0 (Linux; Android 8.1.0; 5059D_RU Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/75.0.3770.101 YaBrowser/19.1.0.130 (lite) Mobile Safari/537.36"
		for i := 0; i < b.N; i++ {
			ctx := AcquireWithSrcStr(ua)
			assertCVS(b, ctx.GetClientType(), ClientTypeBrowser)
			assertStr(b, "browser", ctx.GetBrowser(), "Yandex Browser Lite")
			assertStr(b, "browser version", ctx.GetBrowserVersionString(), "19.1.0.130")
			Release(ctx)
		}
	})
}
