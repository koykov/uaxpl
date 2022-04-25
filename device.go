package uaxpl

import "github.com/koykov/entry"

type DeviceType uint8

const (
	DeviceTypeNotebook DeviceType = 1 << iota
	DeviceTypeMobile
	DeviceTypeConsole
	DeviceTypeHbbTV
	DeviceTypeShellTV
	DeviceTypeCamera
	DeviceTypeCarBrowser
	DeviceTypePortableMediaPlayer
	DeviceTypeAll = DeviceTypeHbbTV | DeviceTypeShellTV | DeviceTypeNotebook | DeviceTypeConsole |
		DeviceTypeCarBrowser | DeviceTypeCamera | DeviceTypePortableMediaPlayer | DeviceTypeMobile
)

type dr struct {
	ne entry.Entry64 // brand name index
	re int32         // regex index
	si entry.Entry64 // string index
	sm int32         // single model index
	me entry.Entry64 // multiple models index
}

type dm struct {
	re int32         // regex index
	si entry.Entry64 // string index
	ne entry.Entry64 // model name index
}
