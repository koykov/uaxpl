package uaxpl

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

type deviceDS struct {
	UA     string `json:"user_agent"`
	Device struct {
		Type  string `json:"type"`
		Brand string `json:"brand"`
		Model string `json:"model"`
	} `json:"device"`
}

func TestDeviceParse(t *testing.T) {
	testDS := func(filename string, deviceType DeviceType) error {
		var ds []deviceDS
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
			if dt := ctx.GetDeviceType(); dt != deviceType {
				d.add("device type", dt.String(), deviceType.String())
			}
			if b := ctx.GetBrand(); b != stage.Device.Brand {
				d.add("brand", b, stage.Device.Brand)
			}
			if m := ctx.GetModel(); m != stage.Device.Model {
				d.add("model", m, stage.Device.Model)
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

	t.Run("camera", func(t *testing.T) {
		if err := testDS("testdata/device/camera.json", DeviceTypeCamera); err != nil {
			t.Error(err)
		}
	})
	t.Run("car_browser", func(t *testing.T) {
		if err := testDS("testdata/device/car_browser.json", DeviceTypeCarBrowser); err != nil {
			t.Error(err)
		}
	})
	t.Run("console", func(t *testing.T) {
		if err := testDS("testdata/device/console.json", DeviceTypeConsole); err != nil {
			t.Error(err)
		}
	})
	t.Run("notebook", func(t *testing.T) {
		if err := testDS("testdata/device/notebook.json", DeviceTypeNotebook); err != nil {
			t.Error(err)
		}
	})
}

func BenchmarkDeviceParse(b *testing.B) {
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
			_ = ctx.GetDeviceType()
			_ = ctx.GetBrand()
			_ = ctx.GetModel()
			Release(ctx)
		}

		return nil
	}
	b.Run("camera", func(b *testing.B) {
		if err := benchDS(b, "testdata/device/camera.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("car_browser", func(b *testing.B) {
		if err := benchDS(b, "testdata/device/car_browser.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("console", func(b *testing.B) {
		if err := benchDS(b, "testdata/device/console.json"); err != nil {
			b.Error(err)
		}
	})
	b.Run("notebook", func(b *testing.B) {
		if err := benchDS(b, "testdata/device/notebook.json"); err != nil {
			b.Error(err)
		}
	})
}
