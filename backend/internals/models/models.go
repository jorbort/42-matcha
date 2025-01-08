package models

type User struct {
	ID int
	username string
	first_name string
	last_name string
	profile_info int
	email string
	validated bool
	completed bool
	password string
	fame_index float32

}

type models struct {
	DB *pgxpool.Pool
}

func (m *models) CreateUser(ctx context.Context, u *User) error {
	tx , err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var profile_info_id int
	stmt := "INSERT INTO profile_info (gender , sexual_orientation, bio, interests, location,  profile_picture_one, profile_picture_two, profile_picture_three, profile_picture_four, profile_picture_five) VALUES (NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL) RETURNING id"
	err = tx.QueryRow(ctx, stmt).Scan(&profile_info_id)
	if err != nil {
		return err
	}

	stmt = "INSERT INTO users (username, first_name, last_name, profile_info, email, validated, completed, password, fame_index) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err = tx.Exec(ctx, stmt, u.username, u.first_name, u.last_name, profile_info_id, u.email, u.validated, u.completed, u.password, u.fame_index)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
