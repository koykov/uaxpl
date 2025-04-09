package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/koykov/uaxpl"
)

type tuple struct {
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

var (
	ua   = flag.String("ua", "", "UserAgent string")                     // single UA mode
	file = flag.String("file", "", "File with JSON tuples")              // file mode
	url_ = flag.String("url", "", "URL to remote file with JSON tuples") // URI mode
)

func init() {
	flag.Parse()
	if len(*ua) == 0 && len(*file) == 0 && len(*url_) == 0 {
		log.Fatalln("param [ua|file|url] is mandatory")
	}
}

func main() {
	var (
		buf bytes.Buffer
		err error
	)
	switch {
	case len(*ua) > 0:
		ctx := uaxpl.NewCtxWithSrcStr(*ua)
		buf.WriteString(fmt.Sprintf("origin ua: %s\n", ctx.GetUserAgent()))
		buf.WriteString(fmt.Sprintf("client:\n"))
		buf.WriteString(fmt.Sprintf(" * type: '%s'\n", ctx.GetClientType()))
		buf.WriteString(fmt.Sprintf(" * browser: '%s'\n", ctx.GetBrowser()))
		buf.WriteString(fmt.Sprintf(" * browser version: '%s'\n", ctx.GetBrowserVersion()))
		buf.WriteString(fmt.Sprintf(" * engine: '%s'\n", ctx.GetEngine()))
		buf.WriteString(fmt.Sprintf(" * engine version: '%s'\n", ctx.GetEngineVersion()))
		buf.WriteString(fmt.Sprintf("device:\n"))
		buf.WriteString(fmt.Sprintf(" * type: '%s'\n", ctx.GetDeviceType()))
		buf.WriteString(fmt.Sprintf(" * brand: '%s'\n", ctx.GetBrand()))
		buf.WriteString(fmt.Sprintf(" * model: '%s'\n", ctx.GetModel()))
		buf.WriteString(fmt.Sprintf(" * os: '%s'\n", ctx.GetOS()))
		buf.WriteString(fmt.Sprintf(" * os version: '%s'\n", ctx.GetOSVersion()))
	case len(*file) > 0:
		err = passFile(&buf, *file)
	case len(*url_) > 0:
		err = passURL(&buf, *url_)
	default:
		log.Fatalln("unknown source param")
	}
	if err != nil {
		log.Fatal(err)
	}
	println(buf.String())
}

func passRaw(w io.StringWriter, contents []byte) (err error) {
	var ds []tuple
	if err = json.Unmarshal(contents, &ds); err != nil {
		return err
	}

	var buf []error
	for i := 0; i < len(ds); i++ {
		stage := &ds[i]
		buf = buf[:0]
		ctx := uaxpl.AcquireWithSrcStr(stage.UA)
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
		uaxpl.Release(ctx)

		if len(buf) > 0 {
			_, _ = w.WriteString(stage.UA)
			_, _ = w.WriteString("\n")
			for j := 0; j < len(buf); j++ {
				_, _ = w.WriteString("* ")
				_, _ = w.WriteString(buf[j].Error())
				_, _ = w.WriteString("\n")
			}
			_, _ = w.WriteString("\n")
		}
	}
	return nil
}

func passFile(w io.StringWriter, filename string) error {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return passRaw(w, contents)
}

func passURL(w io.StringWriter, uri string) error {
	pos := strings.LastIndex(uri, "/")
	if pos == -1 {
		return fmt.Errorf("invalid uri: %s", uri)
	}
	fname := uri[pos+1:]
	fpath := "/tmp/uaxpl_" + fname
	if _, err := os.Stat(fpath); errors.Is(err, os.ErrNotExist) {
		resp, err := http.Get(uri)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer func() { _ = resp.Body.Close() }()
			var contents []byte
			if contents, err = io.ReadAll(resp.Body); err == nil {
				err = os.WriteFile(fpath, contents, 0644)
			}
		}
	}
	contents, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	return passRaw(w, contents)
}

func scmp(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}
