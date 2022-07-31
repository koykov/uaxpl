package uaxpl

import "github.com/koykov/entry"

type DeviceType uint16

const (
	DeviceTypeNotebook DeviceType = 1 << iota
	DeviceTypeMobile
	DeviceTypePhablet
	DeviceTypeTablet
	DeviceTypeConsole
	DeviceTypeTV
	DeviceTypeShellTV
	DeviceTypeCamera
	DeviceTypeCarBrowser
	DeviceTypePortableMediaPlayer
	DeviceTypeAll = DeviceTypeTV | DeviceTypeShellTV | DeviceTypeNotebook | DeviceTypeConsole |
		DeviceTypeCarBrowser | DeviceTypeCamera | DeviceTypePortableMediaPlayer | DeviceTypeMobile | DeviceTypePhablet |
		DeviceTypeTablet
)

type deviceTuple struct {
	type64   entry.Entry64 // device type index
	brand64  entry.Entry64 // brand name index
	matchRI  int32         // regex index
	match64  entry.Entry64 // string index
	modelSI  int32         // single model index
	models64 entry.Entry64 // multiple models index
}

type modelTuple struct {
	type64  entry.Entry64 // model type index
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
	case DeviceTypeTablet:
		return "tablet"
	case DeviceTypePhablet:
		return "phablet"
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
