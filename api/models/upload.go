package models

type Url struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}

type MultipleFileUploadResponse struct {
	Url []*Url `json:"url"`
}
