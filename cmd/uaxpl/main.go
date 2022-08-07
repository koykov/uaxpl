package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koykov/uaxpl"
)

var ua = flag.String("ua", "", "UserAgent string")

func init() {
	flag.Parse()
	if len(*ua) == 0 {
		log.Fatalln("param -ua is mandatory")
	}
}

func main() {
	ctx := uaxpl.NewCtxWithSrcStr(*ua)
	fmt.Printf("origin ua: %s\n", ctx.GetUserAgent())
	fmt.Printf("client:\n")
	fmt.Printf(" * type: '%s'\n", ctx.GetClientType())
	fmt.Printf(" * browser: '%s'\n", ctx.GetBrowser())
	fmt.Printf(" * browser version: '%s'\n", ctx.GetBrowserVersion())
	fmt.Printf(" * engine: '%s'\n", ctx.GetEngine())
	fmt.Printf(" * engine version: '%s'\n", ctx.GetEngineVersion())
	fmt.Printf("device:\n")
	fmt.Printf(" * type: '%s'\n", ctx.GetDeviceType())
	fmt.Printf(" * brand: '%s'\n", ctx.GetBrand())
	fmt.Printf(" * model: '%s'\n", ctx.GetModel())
	fmt.Printf(" * os: '%s'\n", ctx.GetOS())
	fmt.Printf(" * os version: '%s'\n", ctx.GetOSVersion())
}
