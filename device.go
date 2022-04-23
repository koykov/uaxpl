package uaxpl

type DeviceParser uint8

const (
	DeviceParserHbbTV DeviceParser = 1 << iota
	DeviceParserShellTV
	DeviceParserNotebook
	DeviceParserConsole
	DeviceParserCarBrowser
	DeviceParserCamera
	DeviceParserPortableMediaPlayer
	DeviceParserMobile
	DeviceParserAll = DeviceParserHbbTV | DeviceParserShellTV | DeviceParserNotebook | DeviceParserConsole |
		DeviceParserCarBrowser | DeviceParserCamera | DeviceParserPortableMediaPlayer | DeviceParserMobile
)
