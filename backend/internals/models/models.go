package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int     `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	ProfileInfo int     `json:"profile_info"`
	Email       string  `json:"email"`
	Validated   bool    `json:"validated"`
	Completed   bool    `json:"completed"`
	Password    string  `json:"password"`
	Fame_index  float32 `json:"fame_index"`
}

type Models struct {
	DB *pgxpool.Pool
}

func (m *Models) CreateUser(ctx context.Context, u *User) error {
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

	stmt = `INSERT INTO users (username, first_name, last_name, profile_info, email, validated, completed, password, fame_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.Exec(ctx, stmt,
		u.Username,
		u.FirstName,
		u.LastName,
		profile_info_id,
		u.Email,
		u.Validated,
		u.Completed,
		hashedPassword,
		u.Fame_index)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return tx.Commit(ctx)
}

func (m *Models) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (m *Models) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
