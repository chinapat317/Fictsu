package models

import (
	"time"
)

type ChapterModel struct {
	ID			int			`json:"id"`
	Title 		string 		`json:"title"`
	FictionID	int			`json:"fictionID"`
	Content		string		`json:"content"`
	Created 	time.Time	`json:"created"`
}
