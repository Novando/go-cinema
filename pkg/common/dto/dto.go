package dto

type StdResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StdService struct {
	Code int
	Err  error
	Data interface{}
}
