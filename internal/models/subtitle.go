package models

type LineFormat struct {
	Prefix string
	Suffix string
}

type SubtitleBlock struct {
	ID         string
	Time       string
	Text       string
	Formatting []LineFormat
}
