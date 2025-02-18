package models

type User struct {
	ID             int     `json:"id" binding:"ignore"`
	Username       string  `json:"username" binding:"required"`
	FirstName      string  `json:"first_name" binding:"required"`
	LastName       string  `json:"last_name" binding:"required"`
	ProfileInfo    int     `json:"profile_info" binding:"ignore"`
	Email          string  `json:"email" binding:"required" format:"email"`
	Validated      bool    `json:"validated" binding:"ignore"`
	Completed      bool    `json:"completed" binding:"ignore"`
	Password       string  `json:"password" binding:"required"`
	Fame_index     float64 `json:"fame_index" binding:"ignore"`
	ValidationCode []byte  `json:"validation_code" binding:"ignore"`
}

type ProfileInfo struct {
	ID                int      `json:"id" binding:"ignore"`
	Gender            string   `json:"gender" binding:"required"`
	Sexual_preference string   `json:"preferences" binding:"required" enum:"Men,Women,Both"`
	Bio               string   `json:"bio" binding:"required" max:"500"`
	Interests         []string `json:"interests" binding:"required" max:"500"`
	Age               int      `json:"age" binding:"required"`
	Latitude          float64  `json:"Latitude" binding:"ignore" `
	Longitude         float64  `json:"Longitude" binding:"ignore" `
}

type Image struct {
	ID         int    `json:"id" binding:"ignore"`
	UserID     int    `json:"user_id" binding:"ignorez"`
	Img_URI    string `json:"img" binding:"ignore"`
	Img_number int    `json:"img_number" binding:"required"`
}
