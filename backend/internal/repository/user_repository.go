package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, can_create_group, created_at, updated_at)
		VALUES ($1, $2, $3, COALESCE($4, FALSE), NOW(), NOW())
		RETURNING id, can_create_group, created_at, updated_at
	`
	err := database.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.CanCreateGroup,
	).Scan(&user.ID, &user.CanCreateGroup, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	var avatar sql.NullString
	var canCreateGroup sql.NullBool
	// Use case-insensitive search
	query := `
		SELECT id, username, email, password, avatar, can_create_group, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&avatar,
		&canCreateGroup,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if canCreateGroup.Valid {
		user.CanCreateGroup = canCreateGroup.Bool
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	var avatar sql.NullString
	var canCreateGroup sql.NullBool
	query := `
		SELECT id, username, email, password, avatar, can_create_group, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&avatar,
		&canCreateGroup,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if canCreateGroup.Valid {
		user.CanCreateGroup = canCreateGroup.Bool
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByIDs(ids []int) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	// Build query with IN clause for better compatibility
	query := `
		SELECT id, username, email, avatar, can_create_group, created_at, updated_at
		FROM users
		WHERE id = ANY($1::int[])
	`
	rows, err := database.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		var avatar sql.NullString
		var canCreateGroup sql.NullBool
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&avatar,
			&canCreateGroup,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if avatar.Valid {
			user.Avatar = avatar.String
		}
		if canCreateGroup.Valid {
			user.CanCreateGroup = canCreateGroup.Bool
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SearchUsers(query string, excludeUserID int, limit int) ([]*model.User, error) {
	searchQuery := `
		SELECT id, username, email, avatar, can_create_group, created_at, updated_at
		FROM users
		WHERE id != $1
		AND (LOWER(username) LIKE LOWER($2) OR LOWER(email) LIKE LOWER($2))
		ORDER BY 
			CASE 
				WHEN LOWER(username) = LOWER($3) THEN 1
				WHEN LOWER(username) LIKE LOWER($4) THEN 2
				WHEN LOWER(email) LIKE LOWER($2) THEN 3
				ELSE 4
			END,
			username
		LIMIT $5
	`
	
	searchPattern := "%" + query + "%"
	exactPattern := query
	startsPattern := query + "%"
	
	rows, err := database.DB.Query(searchQuery, excludeUserID, searchPattern, exactPattern, startsPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		var avatar sql.NullString
		var canCreateGroup sql.NullBool
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&avatar,
			&canCreateGroup,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if avatar.Valid {
			user.Avatar = avatar.String
		}
		if canCreateGroup.Valid {
			user.CanCreateGroup = canCreateGroup.Bool
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) UpdateAvatar(userID int, avatar string) error {
	query := `UPDATE users SET avatar = $1, updated_at = NOW() WHERE id = $2`
	_, err := database.DB.Exec(query, avatar, userID)
	return err
}

func (r *UserRepository) CheckCanCreateGroup(userID int) (bool, error) {
	var canCreateGroup bool
	query := `SELECT can_create_group FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, userID).Scan(&canCreateGroup)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return canCreateGroup, nil
}