# UserAgent eXPLain library.

Go port of [device-detector](https://github.com/matomo-org/device-detector) library.
Designed to use in high-load.

> Library doesn't provide 100% match with original PHP library due to different
> Regexp engine (Go regexp vs PHP PCRE). Difference is ~ 5%. 

### Usage

```go
ctx := uaxpl.Acquire()
defer uaxpl.Release(ctx)
ctx.SetUserAgentStr("Mozilla/5.0 (Linux; U; Android 9; RMX1941 Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.134 Mobile Safari/537.36 RealmeBrowser/35.5.0.8")
fmt.Println("client type:", ctx.GetClientType())               // browser
fmt.Println("browser:", ctx.GetBrowser())                      // Realme Browser
fmt.Println("browser version:", ctx.GetBrowserVersionString()) // 35.5.0.8
fmt.Println("browser version full:", ctx.GetBrowserVersion())  // 35.5.0.8
fmt.Println("engine:", ctx.GetEngine())                        // WebKit
fmt.Println("engine version:", ctx.GetEngineVersionString())   // 537.36
fmt.Println("engine version full:", ctx.GetEngineVersion())    // 537.36.0.0
fmt.Println("device type:", ctx.GetDeviceType())               // smartphone
fmt.Println("brand:", ctx.GetBrand())                          // Realme
fmt.Println("model:", ctx.GetModel())                          // C2
fmt.Println("OS:", ctx.GetOS())                                // Android
fmt.Println("OS version:", ctx.GetOSVersionString())           // 9
fmt.Println("OS version full:", ctx.GetOSVersion())            // 9.0.0.0
```

### CLI Installation

`uaxpl` has two cli commands: [uaxpl](cmd/uaxpl) and [uaxplc](cmd/uaxplc). First on is a simple cli tool to parse UA in
terminals.

`uaxplc` uses to recompile internal repositories from [device-detector](https://github.com/matomo-org/device-detector)'s
[YAML](https://github.com/matomo-org/device-detector/tree/master/regexes) files.

To install the tool run
```bash
go install github.com/koykov/uaxpl/cmd/uaxplc
```
As result, you must have binary `$GOPATH/bin/uaxplc`.

Then recompile the repos
```bash
go get github.com/koykov/uaxpl
cd $GOPATH/src/github.com/koykov/uaxpl
go generate
```

Run tests to make sure repos were compiled successfully.
