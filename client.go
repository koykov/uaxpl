package uaxpl

import (
	"github.com/koykov/entry"
)

type ClientParser uint8

const (
	ClientParserFeedReader ClientParser = 1 << iota
	ClientParserMobileApp
	ClientParserMediaPlayer
	ClientParserPIM
	ClientParserBrowser
	ClientParserLibrary
	ClientParserAll = ClientParserFeedReader | ClientParserMobileApp | ClientParserMediaPlayer | ClientParserPIM |
		ClientParserBrowser | ClientParserLibrary
)
