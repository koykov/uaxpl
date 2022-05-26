package uaxpl

import (
	"github.com/koykov/entry"
)

type ClientType uint8

const (
	ClientTypeBrowser ClientType = 1 << iota
	ClientTypeMobileApp
	ClientTypeLibrary
	ClientTypeFeedReader
	ClientTypeMediaPlayer
	ClientTypePIM
	ClientTypeAll = ClientTypeBrowser | ClientTypeMobileApp | ClientTypeLibrary | ClientTypeFeedReader |
		ClientTypeMediaPlayer | ClientTypePIM
)

type engineFn func(string) entry.Entry64

type clientTuple struct {
	matchRI   int32         // matching regex index
	match64   entry.Entry64 // static string matching
	browser64 entry.Entry64 // browser name
	browserVI int8          // version index (from regex matches)
	engine64  entry.Entry64 // default engine
	engineFI  int32         // engine func index
	url64     entry.Entry64 // url
	type64    entry.Entry64 // type string
}

func (t ClientType) String() string {
	switch t {
	case ClientTypeBrowser:
		return "browser"
	case ClientTypeMobileApp:
		return "mobile app"
	case ClientTypeLibrary:
		return "library"
	case ClientTypeFeedReader:
		return "feed reader"
	case ClientTypeMediaPlayer:
		return "mediaplayer"
	case ClientTypePIM:
		return "pim"
	default:
		return Unknown
	}
}
