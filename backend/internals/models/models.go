package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int     `json:"id" binding:"ignore"`
	Username    string  `json:"username" binding:"required"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	ProfileInfo int     `json:"profile_info" binding:"ignore"`
	Email       string  `json:"email" binding:"required" format:"email"`
	Validated   bool    `json:"validated" binding:"ignore"`
	Completed   bool    `json:"completed" binding:"ignore"`
	Password    string  `json:"password" binding:"required"`
	Fame_index  float64 `json:"fame_index" binding:"ignore"`
	ValidationCode string `json:"validation_code" binding:"ignore"`
}

type Models struct {
	DB *pgxpool.Pool
}

func (m *Models) InsertUser(ctx context.Context, u *User) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}	
	defer tx.Rollback(ctx)

	var profile_info_id int
	stmt := "INSERT INTO profile_info (gender , sexual_orientation, bio, interests, location,  profile_picture_one, profile_picture_two, profile_picture_three, profile_picture_four, profile_picture_five) VALUES (NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) RETURNING id"
	err = tx.QueryRow(ctx, stmt).Scan(&profile_info_id)
	if err != nil {
		return err
	}
	hashedPassword, err := m.HashPassword(u.Password)

	stmt = `INSERT INTO users (username, first_name, last_name, profile_info, email, validated, completed, password, fame_index , validation_code)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, stmt,
		u.Username,
		u.FirstName,
		u.LastName,
		profile_info_id,
		u.Email,
		u.Validated,
		u.Completed,
		hashedPassword,
		u.Fame_index,
		u.ValidationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (m *Models) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func (m *Models) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (m *Models) userValidation(ctx context.Context, code string) (map[int]bool ,error) {
	var userInfo map[int]bool
	var id int
	var completed bool
	stmt := `UPDATE users SET validated = true WHERE validation_code = $1 RETURNING id , completed`
	err := m.DB.QueryRow(ctx, stmt, code).Scan(&id, &completed)
	if err != nil {
		return nil, err
	}
	userInfo[id] = completed
	return userInfo ,nil
}