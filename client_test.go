package uaxpl

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type clientDS struct {
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
	testDS := func(filename string, clientType ClientType) error {
		var ds []clientDS
		contents, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(contents, &ds); err != nil {
			return err
		}

		var dl []*diff
		for i := 0; i < len(ds); i++ {
			stage := &ds[i]
			d := newDiff(stage.UA)
			ctx := AcquireWithSrcStr(stage.UA)
			ctx.SetRequestedWith(testGetRWH(stage))
			if ct := ctx.GetClientType(); ct != clientType {
				d.add("client type", ct.String(), clientType.String())
			}
			if b := ctx.GetBrowser(); b != stage.Client.Name {
				d.add("browser", b, stage.Client.Name)
			}
			if bv := ctx.GetBrowserVersionString(); len(stage.Client.Version) > 0 && bv != stage.Client.Version {
				d.add("browser version", bv, stage.Client.Version)
			}
			if e := ctx.GetEngine(); e != stage.Client.Engine {
				d.add("engine", e, stage.Client.Engine)
			}
			if ev := ctx.GetEngineVersionString(); len(stage.Client.EngineVersion) > 0 && ev != stage.Client.EngineVersion {
				d.add("engine version", ev, stage.Client.EngineVersion)
			}
			if d.len() > 0 {
				dl = append(dl, d)
			}
			Release(ctx)
		}

		if len(dl) > 0 {
			var buf bytes.Buffer
			for i := 0; i < len(dl); i++ {
				dl[i].write(&buf)
				_ = buf.WriteByte('\n')
			}
			t.Log(buf.String())
		}

		return nil
	}
	t.Run("browser", func(t *testing.T) {
		if err := testDS("testdata/client/browser.json", ClientTypeBrowser); err != nil {
			t.Error(err)
		}
	})
	t.Run("feed_reader", func(t *testing.T) {
		if err := testDS("testdata/client/feed_reader.json", ClientTypeFeedReader); err != nil {
			t.Error(err)
		}
	})
	t.Run("library", func(t *testing.T) {
		if err := testDS("testdata/client/library.json", ClientTypeLibrary); err != nil {
			t.Error(err)
		}
	})
	t.Run("mediaplayer", func(t *testing.T) {
		if err := testDS("testdata/client/mediaplayer.json", ClientTypeMediaPlayer); err != nil {
			t.Error(err)
		}
	})
	t.Run("mobile_app", func(t *testing.T) {
		if err := testDS("testdata/client/mobile_app.json", ClientTypeMobileApp); err != nil {
			t.Error(err)
		}
	})
	t.Run("pim", func(t *testing.T) {
		if err := testDS("testdata/client/pim.json", ClientTypePIM); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkClientParse(b *testing.B) {
	benchDS := func(b *testing.B, filename string) error {
		var ds []clientDS
		contents, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(contents, &ds); err != nil {
			return err
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			stage := &ds[i%len(ds)]
			ctx := AcquireWithSrcStr(stage.UA)
			ctx.SetRequestedWith(testGetRWH(stage))
			_ = ctx.GetClientType()
			_ = ctx.GetBrowser()
			_ = ctx.GetBrowserVersionString()
			_ = ctx.GetEngine()
			_ = ctx.GetEngineVersionString()
			Release(ctx)
		}

		return nil
	}
	b.Run("browser", func(b *testing.B) {
		if err := benchDS(b, "testdata/client/browser.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("library", func(b *testing.B) {
		if err := benchDS(b, "testdata/client/library.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("mediaplayer", func(b *testing.B) {
		if err := benchDS(b, "testdata/client/mediaplayer.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("mobile_app", func(b *testing.B) {
		if err := benchDS(b, "testdata/client/mobile_app.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("pim", func(b *testing.B) {
		if err := benchDS(b, "testdata/client/pim.json"); err != nil {
			b.Error(err)
		}
	})
}
