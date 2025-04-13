package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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
	ua      = flag.String("ua", "", "UserAgent string")                     // single UA mode
	file    = flag.String("file", "", "File with JSON tuples")              // file mode
	url_    = flag.String("url", "", "URL to remote file with JSON tuples") // URI mode
	out     = flag.String("out", "out.diff.txt", "Output file")
	threads = flag.Int("threads", 1, "Number of threads")
	verbose = flag.Bool("verbose", false, "Verbose output")
)

func init() {
	flag.Parse()
	if len(*ua) == 0 && len(*file) == 0 && len(*url_) == 0 {
		log.Fatalln("param [ua|file|url] is mandatory")
	}
}

func main() {
	f, err := os.OpenFile(*out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	switch {
	case len(*ua) > 0:
		ctx := uaxpl.NewCtxWithSrcStr(*ua)
		_, _ = f.WriteString(fmt.Sprintf("origin ua: %s\n", ctx.GetUserAgent()))
		_, _ = f.WriteString(fmt.Sprintf("client:\n"))
		_, _ = f.WriteString(fmt.Sprintf(" * type: '%s'\n", ctx.GetClientType()))
		_, _ = f.WriteString(fmt.Sprintf(" * browser: '%s'\n", ctx.GetBrowser()))
		_, _ = f.WriteString(fmt.Sprintf(" * browser version: '%s'\n", ctx.GetBrowserVersion()))
		_, _ = f.WriteString(fmt.Sprintf(" * engine: '%s'\n", ctx.GetEngine()))
		_, _ = f.WriteString(fmt.Sprintf(" * engine version: '%s'\n", ctx.GetEngineVersion()))
		_, _ = f.WriteString(fmt.Sprintf("device:\n"))
		_, _ = f.WriteString(fmt.Sprintf(" * type: '%s'\n", ctx.GetDeviceType()))
		_, _ = f.WriteString(fmt.Sprintf(" * brand: '%s'\n", ctx.GetBrand()))
		_, _ = f.WriteString(fmt.Sprintf(" * model: '%s'\n", ctx.GetModel()))
		_, _ = f.WriteString(fmt.Sprintf(" * os: '%s'\n", ctx.GetOS()))
		_, _ = f.WriteString(fmt.Sprintf(" * os version: '%s'\n", ctx.GetOSVersion()))
	case len(*file) > 0:
		err = passFile(f, *file)
	case len(*url_) > 0:
		err = passURL(f, *url_)
	default:
		log.Fatalln("unknown source param")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func passRaw(w io.StringWriter, contents []byte) (err error) {
	var ds []tuple
	if err = json.Unmarshal(contents, &ds); err != nil {
		return err
	}

	var (
		fmux sync.Mutex
		wg   sync.WaitGroup
		ch   = make(chan *tuple, *threads*8)
		c    uint64
		now  = time.Now()
	)

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var buf []error
			for {
				select {
				case stage, ok := <-ch:
					if !ok {
						return
					}
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
						fmux.Lock()
						_, _ = w.WriteString(stage.UA)
						_, _ = w.WriteString("\n")
						for j := 0; j < len(buf); j++ {
							_, _ = w.WriteString("* ")
							_, _ = w.WriteString(buf[j].Error())
							_, _ = w.WriteString("\n")
						}
						_, _ = w.WriteString("\n")
						fmux.Unlock()
					}

					p := atomic.AddUint64(&c, 1)
					if *verbose && p%1000 == 0 {
						log.Printf("pass %d of %d rows\n", p, len(ds))
					}
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(ds); i++ {
			stage := &ds[i]
			ch <- stage
		}
		close(ch)
	}()
	wg.Wait()
	if *verbose {
		log.Printf("total rows %d passed in %s\n", len(ds), time.Since(now).String())
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
