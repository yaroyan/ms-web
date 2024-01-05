package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/yaroyan/ms/authn/domain/model/authentication"
	"golang.org/x/crypto/bcrypt"
)

const ConnectionTimeout = 3 * time.Second

type UserRepository struct {
	Connection *sql.DB
}

// FindAll returns a slice of all users, sorted by last name
func (u *UserRepository) FindAll() ([]*authentication.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users order by last_name`

	rows, err := u.Connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*authentication.User

	for rows.Next() {
		var user authentication.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// FindByEmail returns one user by email
func (u *UserRepository) FindByEmail(email string) (*authentication.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user authentication.User
	var pw authentication.Password
	var em authentication.Email
	row := u.Connection.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&em.Email,
		&user.FirstName,
		&user.LastName,
		&pw.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	user.Password = pw
	user.Email = em

	if err != nil {
		log.Println(user)
		return nil, err
	}

	return &user, nil
}

// FindOne returns one user by id
func (u *UserRepository) FindOne(id int) (*authentication.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user authentication.User
	row := u.Connection.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *UserRepository) Update(user authentication.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := u.Connection.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *UserRepository) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := u.Connection.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *UserRepository) Insert(user authentication.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := u.Connection.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *UserRepository) ResetPassword(password string, user authentication.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = u.Connection.ExecContext(ctx, stmt, hashedPassword, user.ID)
	if err != nil {
		return err
	}

	return nil
}
