package model

//Repository struct
type Repository struct {
	ID          int64  `json:"id,omitempty" db:"id"`
	Name        string `json:"name,omitempty" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
	URL         string `json:"html_url,omitempty" db:"url"`
	Language    string `json:"language,omitempty" db:"language"`
	Tags        []Tags `json:"tags"`
	tagid       int    `db:"tag_id"`
}

//Tags struct
type Tags struct {
	id   int    `db:"id"`
	Name string `json:"tag" db:"tag"`
}

//ValidateDuplicatedTag validate duplicated tags
func ValidateDuplicatedTag(repo Repository, tagValue Tags) bool {
	for _, tag := range repo.Tags {
		if tagValue.Name == tag.Name {
			return true
		}
	}
	return false
}
