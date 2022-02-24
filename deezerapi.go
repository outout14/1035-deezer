package main

import "encoding/json"

type deezerTitle struct {
	Title        string `json:"title"`
	TitleShort   string `json:"title_short"`
	TitleVersion string `json:"title_version"`
	Link         string `json:"link"`
	Timestamp    int    `json:"timestamp"`
	Artist       struct {
		Name      string `json:"name"`
		Link      string `json:"link"`
		Tracklist string `json:"tracklist"`
		Type      string `json:"type"`
	} `json:"artist"`
	Album struct {
		Title     string `json:"title"`
		Link      string `json:"link"`
		Tracklist string `json:"tracklist"`
		Type      string `json:"type"`
	} `json:"album"`
	Type string `json:"type"`
}

type deezerHistory struct {
	Data []deezerTitle `json:"data"`
}

type deezerAccount struct {
	ID json.Number `json:"id"`
}
