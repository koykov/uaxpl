package uaxpl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

type customDS struct {
	UA     string `json:"user_agent"`
	Client struct {
		Type          string `json:"type,omitempty"`
		Name          string `json:"name,omitempty"`
		Version       string `json:"version,omitempty"`
		Engine        string `json:"engine,omitempty"`
		EngineVersion string `json:"engine_version,omitempty"`
		Family        string `json:"family,omitempty"`
	} `json:"client,omitempty"`
	Device struct {
		Type      string `json:"type,omitempty"`
		Brand     string `json:"brand,omitempty"`
		Model     string `json:"model,omitempty"`
		OS        string `json:"os,omitempty"`
		OSVersion string `json:"os_version,omitempty"`
	} `json:"device,omitempty"`
}

func TestCustomParse(t *testing.T) {
	scmp := func(a, b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}

	testDS := func(filename string) error {
		var ds []customDS
		contents, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(contents, &ds); err != nil {
			return err
		}

		var buf []error
		for i := 0; i < len(ds); i++ {
			stage := &ds[i]
			buf = buf[:0]
			ctx := AcquireWithSrcStr(stage.UA)
			if ct := ctx.GetClientType(); len(stage.Client.Type) > 0 && !scmp(ct.String(), stage.Client.Type) {
				buf = append(buf, fmt.Errorf("client type mismatch: need '%s' got '%s'", stage.Client.Type, ct))
			}
			if b := ctx.GetBrowser(); len(stage.Client.Name) > 0 && !scmp(b, stage.Client.Name) {
				buf = append(buf, fmt.Errorf("browser mismatch: need '%s' got '%s'", stage.Client.Name, b))
			}
			if bv := ctx.GetBrowserVersionString(); len(stage.Client.Version) > 0 && !scmp(bv, stage.Client.Version) {
				buf = append(buf, fmt.Errorf("browser version mismatch: need '%s' got '%s'", stage.Client.Version, bv))
			}
			if e := ctx.GetEngine(); len(stage.Client.Engine) > 0 && !scmp(e, stage.Client.Engine) {
				buf = append(buf, fmt.Errorf("engine mismatch: need '%s' got '%s'", stage.Client.Engine, e))
			}
			if ev := ctx.GetEngineVersionString(); len(stage.Client.EngineVersion) > 0 && !scmp(ev, stage.Client.EngineVersion) {
				buf = append(buf, fmt.Errorf("engine version mismatch: need '%s' got '%s'", stage.Client.EngineVersion, ev))
			}

			if dt := ctx.GetDeviceType(); len(stage.Device.Type) > 0 && !scmp(dt.String(), stage.Device.Type) {
				buf = append(buf, fmt.Errorf("device type mismatch: need '%s' got '%s'", stage.Device.Type, dt))
			}
			if b := ctx.GetBrand(); len(stage.Device.Brand) > 0 && !scmp(b, stage.Device.Brand) {
				buf = append(buf, fmt.Errorf("brand mismatch: need '%s' got '%s'", stage.Device.Brand, b))
			}
			if m := ctx.GetModel(); len(stage.Device.Model) > 0 && !scmp(m, stage.Device.Model) {
				buf = append(buf, fmt.Errorf("model mismatch: need '%s' got '%s'", stage.Device.Model, m))
			}
			if o := ctx.GetOS(); len(stage.Device.OS) > 0 && !scmp(o, stage.Device.OS) {
				buf = append(buf, fmt.Errorf("OS mismatch: need '%s' got '%s'", stage.Device.OS, o))
			}
			if ov := ctx.GetOSVersionString(); len(stage.Device.OSVersion) > 0 && !scmp(ov, stage.Device.OSVersion) {
				buf = append(buf, fmt.Errorf("OS version mismatch: need '%s' got '%s'", stage.Device.OSVersion, ov))
			}
			Release(ctx)

			if len(buf) > 0 {
				t.Log(stage.UA)
				for j := 0; j < len(buf); j++ {
					t.Log("*", buf[j])
				}
			}
		}
		return nil
	}

	t.Run("single0", func(t *testing.T) {
		_ = testDS("testdata/custom/single0.json")
	})
	t.Run("single1", func(t *testing.T) {
		_ = testDS("testdata/custom/single1.json")
	})
	t.Run("single2", func(t *testing.T) {
		_ = testDS("testdata/custom/single2.json")
	})
	t.Run("single3", func(t *testing.T) {
		_ = testDS("testdata/custom/single3.json")
	})
	t.Run("single4", func(t *testing.T) {
		_ = testDS("testdata/custom/single4.json")
	})
	t.Run("ds0", func(t *testing.T) {
		_ = testDS("testdata/custom/ds0.json")
	})
	t.Run("ds1", func(t *testing.T) {
		var ds []string
		contents, err := os.ReadFile("testdata/custom/ds1.json")
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(contents, &ds); err != nil {
			t.Fatal(err)
		}
		ctx := Acquire()
		defer Release(ctx)
		for i := 0; i < len(ds); i++ {
			ctx.SetUserAgentStr(ds[i])
			_ = ctx.GetClientType()
			_ = ctx.GetDeviceType()
		}
	})
}
