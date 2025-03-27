package models

import (
	"time"
)

type Status string

const (
	Completed Status = "Completed"
	Ongoing   Status = "Ongoing"
)

type FictionModel struct {
	ID               int       `json:"id"`
	Contributor_Name string    `json:"contributor_name"`
	Cover            string    `json:"cover"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Author           string    `json:"author"`
	Artist           string    `json:"artist"`
	Status           Status    `json:"status"`
	Synopsis         string    `json:"synopsis"`
	CreatedTime      time.Time `json:"created_time"`
}

type GenreModel struct {
	ID         int    `json:"id"`
	Genre_Name string `json:"genre_name"`
}
