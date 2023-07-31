package entity

type Request struct {
	ID             int64
	ChatID         int64
	UserID         int64
	ImagePathName  string
	ResImgPathName string
	Description    string
}

type DetectionResult struct {
	ID          int64
	Description string
	Image       string
}

type Description struct {
	Name string `json:"name"`
}

type User struct {
	Login    string
	Password string
}
