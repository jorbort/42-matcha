package models

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twpayne/go-geos"
	_ "github.com/twpayne/pgx-geos"
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	DB *pgxpool.Pool
}

func (m *Models) CreateTables(ctx context.Context) error {
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS profile_info(
            id bigserial PRIMARY KEY,
            gender varchar(50),
            sexual_orientation varchar(50),
            age int,
            bio varchar(255),
            interests text[],
            location geometry(Point, 4326)
        );`,

		`CREATE TABLE IF NOT EXISTS users (
            id bigserial PRIMARY KEY,
            username varchar(255) UNIQUE not null,
            first_name varchar(255) not null,
            last_name varchar(255) not null,
            profile_info int UNIQUE REFERENCES profile_info(id) on delete cascade,
            email varchar(255) UNIQUE not null,
            validated boolean not null default false,
            completed boolean not null default false,
            password varchar(255) not null,
            fame_index float not null,
            validation_code bytea not null
        );`,

		`CREATE TABLE IF NOT EXISTS user_images(
            id BIGSERIAL PRIMARY KEY,
            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            image_number INT NOT NULL CHECK (image_number BETWEEN 1 AND 5),
            image_url varchar(255) NOT NULL,
            UNIQUE(user_id, image_number)
        );`,

		`CREATE TABLE IF NOT EXISTS conversations(
        	id SERIAL PRIMARY KEY,
         user1_id integer REFERENCES users(id),
         user2_id integer REFERENCES users(id),
         created_at TIMESTAMP with TIME ZONE DEFAULT CURRENT_TIMESTAMP,
         unique(user1_id, user2_id)
        );`,

		`CREATE TABLE IF NOT EXISTS messages(
        	id SERIAL PRIMARY KEY,
         	conversation_id integer REFERENCES conversations(id),
          	sender_id integer REFERENCES users(id),
            receiver_id integer REFERENCES users(id),
            message TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for _, query := range createTableQueries {
		_, err := m.DB.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}

func (m *Models) InsertUser(ctx context.Context, u *User) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var profile_info_id int
	stmt := "INSERT INTO profile_info (gender , sexual_orientation, bio, interests, location) VALUES (NULL, NULL, NULL, NULL, NULL) RETURNING id"
	err = tx.QueryRow(ctx, stmt).Scan(&profile_info_id)
	if err != nil {
		return err
	}
	hashedPassword, err := m.HashPassword([]byte(u.Password))

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
		string(hashedPassword),
		u.Fame_index,
		u.ValidationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	stmt = `
		INSERT  INTO user_images (user_id, image_number, image_url) VALUES
		($1, 1, ''),
		($1, 2, ''),
		($1, 3, ''),
		($1, 4, ''),
		($1, 5, '')
		`
	_, err = tx.Exec(ctx, stmt, profile_info_id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return tx.Commit(ctx)
}

func (m *Models) HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return bytes, err
}

func (m *Models) VerifyPassword(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func (m *Models) UserValidation(ctx context.Context, code string) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `UPDATE users SET validated = true WHERE validation_code = $1`
	_, err = tx.Exec(ctx, stmt, code)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *Models) GetUserByUsername(ctx context.Context, username string) (*User, error) {

	stmt := `SELECT id, username, first_name, last_name, profile_info, email, validated, completed, password, fame_index, validation_code FROM users WHERE username = $1`
	row := m.DB.QueryRow(ctx, stmt, username)

	u := &User{}
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.ProfileInfo, &u.Email, &u.Validated, &u.Completed, &u.Password, &u.Fame_index, &u.ValidationCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", username)
		}
		return nil, err
	}
	return u, nil
}

func (m *Models) InsertProfileInfo(ctx context.Context, p *ProfileInfo) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	geoctx := geos.NewContext()
	coord := []float64{p.Longitude, p.Latitude}
	point := geoctx.NewPoint(coord)

	wkb := point.ToWKB()
	stmt := `UPDATE profile_info SET gender = $1, sexual_orientation = $2, bio = $3, interests = $4, age = $5 , location = ST_SetSRID(ST_GeomFromWKB($6), 4326) WHERE id = $7`
	_, err = tx.Exec(ctx, stmt, p.Gender, p.Sexual_preference, p.Bio, p.Interests, p.Age, wkb, p.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *Models) UpdateUserCompleted(ctx context.Context, id int) error {
	stmt := `UPDATE users SET completed = true WHERE id = $1`
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Models) UpdateUser(ctx context.Context, validationCode []byte, email string) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `UPDATE users SET validation_code = $1 WHERE email = $2`
	result, err := tx.Exec(ctx, stmt, validationCode, email)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("email not found")
	}
	return tx.Commit(ctx)
}

func (m *Models) UpdatePassword(ctx context.Context, code, newPassword string) error {
	log.Println("code: ", code)
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	stmt := `UPDATE users SET password = $1 WHERE validation_code = $2`
	hashedPassword, err := m.HashPassword([]byte(newPassword))
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, stmt, string(hashedPassword), code)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		log.Println("code not found")
		return fmt.Errorf("code not found")
	}
	return tx.Commit(ctx)
}
