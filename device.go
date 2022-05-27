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

type deviceTuple struct {
	brand64  entry.Entry64 // brand name index
	matchRI  int32         // regex index
	match64  entry.Entry64 // string index
	modelSI  int32         // single model index
	models64 entry.Entry64 // multiple models index
}

type modelTuple struct {
	matchRI int32         // regex index
	match64 entry.Entry64 // string index
	model64 entry.Entry64 // model name index
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
