package uaxpl

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type osDS struct {
	UA string `json:"user_agent"`
	OS struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		Version   string `json:"version"`
		Platform  string `json:"platform,omitempty"`
		Family    string `json:"family,omitempty"`
	} `json:"os"`
	Headers map[string]string `json:"headers,omitempty"`
}

func TestOSParse(t *testing.T) {
	testDS := func(filename string) error {
		var ds []osDS
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
			if o := ctx.GetOS(); o != stage.OS.Name {
				d.add("os", o, stage.OS.Name)
			}
			if v := ctx.GetOSVersionString(); v != stage.OS.Version {
				d.add("version", v, stage.OS.Version)
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
	t.Run("os", func(t *testing.T) {
		if err := testDS("testdata/os.json"); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkOSParse(b *testing.B) {
	benchDS := func(b *testing.B, filename string) error {
		var ds []deviceDS
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
			_ = ctx.GetOS()
			_ = ctx.GetOSVersionString()
			Release(ctx)
		}

		return nil
	}
	b.Run("os", func(b *testing.B) {
		if err := benchDS(b, "testdata/os.json"); err != nil {
			b.Error(err)
		}
	})
}
