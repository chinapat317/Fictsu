package models

import (
	"time"
)

type Status string
type Genre 	string

const (
	Completed	Status = "Completed"
	Ongoing		Status = "Ongoing"

	Fantasy Genre = "Fantasy"
	Romance Genre = "Romance"
	Horror  Genre = "Horror"
	Action  Genre = "Action"
)

type FictionModel struct {
	ID			int				`json:"id"`
	Title 		string 			`json:"title"`
	Subtitle 	string 			`json:"subtitle"`
	Author		string 			`json:"author"`
	Artist		string 			`json:"artist"`
	Status		Status 			`json:"status"`
	Synopsis	string 			`json:"synopsis"`
	Genre		Genre			`json:"genre"`
	Chapters	[]ChapterModel	`json:"chapters"`
	Created		time.Time		`json:"created"`
}
