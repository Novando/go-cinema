package service

type ScreenRequestDTO struct {
	MovieID string `query:"movieId" validate:"required"`
	Date    string `query:"date" validate:"required"`
}

type BookRequestDTO struct {
	ScreenID string   `json:"screenId" validate:"required"`
	Name     string   `json:"name" validate:"required"`
	Seats    []string `json:"seats" validate:"required"`
}
