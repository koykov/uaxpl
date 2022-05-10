package uaxpl

import "github.com/koykov/entry"

type DeviceType uint8

const (
	DeviceTypeNotebook DeviceType = 1 << iota
	DeviceTypeMobile
	DeviceTypeConsole
	DeviceTypeTV
	DeviceTypeShellTV
	DeviceTypeCamera
	DeviceTypeCarBrowser
	DeviceTypePortableMediaPlayer
	DeviceTypeAll = DeviceTypeTV | DeviceTypeShellTV | DeviceTypeNotebook | DeviceTypeConsole |
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

func (t DeviceType) String() string {
	switch t {
	case DeviceTypeNotebook:
		return "desktop"
	case DeviceTypeMobile:
		return "smartphone"
	case DeviceTypeConsole:
		return "console"
	case DeviceTypeTV:
		return "tv"
	case DeviceTypeShellTV:
		return "tv"
	case DeviceTypeCamera:
		return "camera"
	case DeviceTypeCarBrowser:
		return "car browser"
	case DeviceTypePortableMediaPlayer:
		return "portable media player"
	default:
		return Unknown
	}
}
