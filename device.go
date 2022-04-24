package uaxpl

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
