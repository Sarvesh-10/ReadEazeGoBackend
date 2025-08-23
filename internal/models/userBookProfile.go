package models

import "encoding/json"

type UserBookProfile struct {
	UserID            int     `json:"user_id"`
	BookID            int     `json:"book_id"`
	BookName          string  `json:"book_name"`
	Mode              string  `json:"mode"`
	TotalPages        int     `json:"total_pages"`
	CurrentPage       int     `json:"current_page"`
	ReadingPercentage float64 `json:"reading_percentage"`
}

func (u UserBookProfile) ToJSON() ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return data, nil
}
