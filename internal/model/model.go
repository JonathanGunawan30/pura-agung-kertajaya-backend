package model

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Paging *PageMetadata `json:"paging,omitempty"`
	Errors string        `json:"errors,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type ImageVariants struct {
	Blur   string `json:"blur"`
	Avatar string `json:"avatar"`
	Xs     string `json:"xs"`
	Sm     string `json:"sm"`
	Md     string `json:"md"`
	Lg     string `json:"lg"`
	Xl     string `json:"xl"`
	TwoXl  string `json:"2xl"`
	Fhd    string `json:"fhd"`
}
