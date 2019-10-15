package model

//Repository struct
type Repository struct {
	ID          int64  `json:"id,omitempty" db:"id"`
	Name        string `json:"name,omitempty" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
	URL         string `json:"html_url,omitempty" db:"url"`
	Language    string `json:"language,omitempty" db:"language"`
	Tags        []Tags `json:"tags,omitempty"`
	tagid       int    `db:"tag_id"`
}

//Tags struct
type Tags struct {
	id   int    `db:"id"`
	Name string `json:"name" db:"name"`
}

//ValidateDuplicatedTag validate duplicated tags
func ValidateDuplicatedTag(name string, tags []Tags) bool {
	for _, tag := range tags {
		if name == tag.Name {
			return true
		}
	}
	return false
}
