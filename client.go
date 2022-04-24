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

type engFn func(string) entry.Entry64

type cr struct {
	re int32         // regex index
	si entry.Entry64 // string index
	vi int8          // version index
	ed entry.Entry64 // default engine
	ef int32         // engine func index
	ul entry.Entry64 // url
	tp entry.Entry64 // type string
}
