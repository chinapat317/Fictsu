package models

import (
	"time"
)

type ChapterModel struct {
	Fiction_ID	int			`json:"fiction_id"`
	ID			int			`json:"id"`
	Title 		string 		`json:"title"`
	Content		string		`json:"content"`
	Created 	time.Time	`json:"created"`
}
