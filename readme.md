# UserAgent parser.

Go port of [device-detector](https://github.com/matomo-org/device-detector) library.
Designed to use in high-load.

> Library doesn't provide 100% match with original PHP library due to different
> Regexp engine (Go regexp vs PHP PCRE). Difference is ~ 5%. 

### Usage

```go
ctx := uaxpl.Acquire()
ctx.SetUserAgentStr("Mozilla/5.0 (Linux; U; Android 9; RMX1941 Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.134 Mobile Safari/537.36 RealmeBrowser/35.5.0.8")
fmt.Println("client type:", ctx.GetClientType())               // browser
fmt.Println("browser:", ctx.GetBrowser())                      // Realme Browser
fmt.Println("browser version:", ctx.GetBrowserVersionString()) // 35.5.0.8
fmt.Println("browser version full:", ctx.GetBrowserVersion())  // 35.5.0.8
fmt.Println("engine:", ctx.GetEngine())                        // WebKit
fmt.Println("engine version:", ctx.GetEngineVersionString())   // 537.36
fmt.Println("engine version full:", ctx.GetEngineVersion())    // 537.36.0.0
fmt.Println("device type:", ctx.GetDeviceType())               // smartphone
fmt.Println("brand:", ctx.GetBrand())                          // Motorola
fmt.Println("model:", ctx.GetModel())                          // DROID 9
fmt.Println("OS:", ctx.GetOS())                                // Android
fmt.Println("OS version:", ctx.GetOSVersionString())           // 9
fmt.Println("OS version full:", ctx.GetOSVersion())            // 9.0.0.0
Release(ctx)
```
