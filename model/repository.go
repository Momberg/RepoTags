package model

import (
	"strconv"
)

//Repository struct
type Repository struct {
	ID                int64  `json:"id,omitempty" db:"id"`
	Name              string `json:"name,omitempty" db:"name"`
	Description       string `json:"description,omitempty" db:"description"`
	URL               string `json:"html_url,omitempty" db:"url"`
	Language          string `json:"language,omitempty" db:"language"`
	Tagrecommendation string `json:"tagrecommendation,omitempty"`
	Tags              []Tags `json:"tags,omitempty"`
	tagid             int    `db:"tag_id"`
}

//Tags struct
type Tags struct {
	ID   int64  `json:"-" db:"id"`
	Name string `json:"name" db:"name"`
}

//ValidateDuplicatedTag validate duplicated tags
func ValidateDuplicatedTag(name string, tags []Tags, id string) bool {
	for _, tag := range tags {
		if name == tag.Name && strconv.FormatInt(tag.ID, 10) == id {
			return true
		}
	}
	return false
}
