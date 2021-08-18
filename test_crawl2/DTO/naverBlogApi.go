package DTO

import "time"

type NaverBlogApiStruct struct {
	Items []NaverBlogApiItem
}

type NaverBlogApiItem struct {
	CustomPostDate time.Time
	Title          string `json: "title"`
	Link           string `json: "link"`
	Description    string `json: "description"`
	Postdate       string `json: "postdate"`
}
