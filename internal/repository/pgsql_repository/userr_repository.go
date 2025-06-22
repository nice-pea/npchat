package pgsqlRepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	sqlxRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/sqlx_repo"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

type UserrRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *UserrRepository) List(filter userr.Filter) ([]userr.User, error) {
	// Запросить пользователей
	var users []dbUser
	if err := r.DB().Select(&users, `
		SELECT u.*
		FROM users u
			LEFT JOIN oauth_users o ON u.id = o.user_id
		WHERE ($1 = '' OR $1 = u.id)
			AND ($2 = '' OR $2 = u.login)
			AND ($3 = '' OR $3 = u.password)
			AND ($4 = '' OR $4 = o.provider)
			AND ($5 = '' OR $5 = o.user_id)
		GROUP BY u.id
	`, filter.ID, filter.BasicAuthLogin, filter.BasicAuthPassword,
		filter.OAuthProvider, filter.OAuthUserID); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Если пользователей нет, сразу вернуть пустой список
	if len(users) == 0 {
		return nil, nil
	}

	// Собрать ID найденных пользователей
	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	// Найти oauth пользователей для пользователей
	var oauthUsers []dbOAuthUser
	if err := r.DB().Select(&oauthUsers, `
		SELECT *
		FROM oauth_users
		WHERE user_id = ANY($1)
	`, userIDs); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Создать карту, где ключ это ID пользователя, а значение это список его oauth пользователей
	oauthUsersMap := make(map[string][]dbOAuthUser, len(users))
	for _, u := range oauthUsers {
		oauthUsersMap[u.UserID] = append(oauthUsersMap[u.UserID], u)
	}

	return toDomainUsers(users, oauthUsersMap), nil
}

func (r *UserrRepository) Upsert(user userr.User) error {
	if user.ID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if r.IsTx() {
		return r.upsert(user)
	} else {
		return r.InTransaction(func(txRepo userr.Repository) error {
			return txRepo.Upsert(user)
		})
	}
}

func (r *UserrRepository) upsert(user userr.User) error {
	if _, err := r.DB().NamedExec(`
		INSERT INTO users(id, name, nick, login, password) 
		VALUES (:is, :name, :nick, :login, :password)
		ON CONFLICT DO UPDATE SET
			name = excluded.name,
			nick = excluded.nick,
			login = excluded.login,
			password = excluded.password
	`, toDBUser(user)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	if _, err := r.DB().NamedExec(`
		DELETE
		FROM oauth_users
		WHERE user_id =  :user_id;
		INSERT INTO oauth_users(id, user_id, provider, email, name, picture, access_token, token_type, refresh_token, expiry) 
		VALUES (:id, :user_id, :provider, :email, :name, :picture, :access_token, :token_type, :refresh_token, :expiry)
	`, toDBOAuthUsers(user)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	return nil
}

func (r *UserrRepository) InTransaction(fn func(txRepo userr.Repository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&UserrRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbUser struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Nick     string `db:"nick"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func toDBUser(user userr.User) dbUser {
	return dbUser{
		ID:       user.ID.String(),
		Name:     user.Name,
		Nick:     user.Nick,
		Login:    user.BasicAuth.Login,
		Password: user.BasicAuth.Password,
	}
}

func toDomainUser(user dbUser, oauthUsers []dbOAuthUser) userr.User {
	return userr.User{
		ID:            uuid.MustParse(user.ID),
		Name:          user.Name,
		Nick:          user.Nick,
		OpenAuthUsers: toDomainOAuthUsers(oauthUsers),
		BasicAuth: userr.BasicAuth{
			Login:    user.Login,
			Password: user.Password,
		},
	}
}

func toDomainUsers(users []dbUser, oauthUsers map[string][]dbOAuthUser) []userr.User {
	domainUsers := make([]userr.User, len(users))
	for i, u := range users {
		domainUsers[i] = toDomainUser(u, oauthUsers[u.ID])
	}

	return domainUsers
}

type dbOAuthUser struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	Provider     string    `db:"provider"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Picture      string    `db:"picture"`
	AccessToken  string    `db:"access_token"`
	TokenType    string    `db:"token_type"`
	RefreshToken string    `db:"refresh_token"`
	Expiry       time.Time `db:"expiry"`
}

func toDBOAuthUsers(user userr.User) []dbOAuthUser {
	dbUsers := make([]dbOAuthUser, len(user.OpenAuthUsers))
	for i, oauthUser := range user.OpenAuthUsers {
		dbUsers[i] = dbOAuthUser{
			ID:           oauthUser.ID,
			UserID:       user.ID.String(),
			Provider:     oauthUser.Provider,
			Email:        oauthUser.Email,
			Name:         oauthUser.Name,
			Picture:      oauthUser.Picture,
			AccessToken:  oauthUser.Token.AccessToken,
			TokenType:    oauthUser.Token.TokenType,
			RefreshToken: oauthUser.Token.RefreshToken,
			Expiry:       oauthUser.Token.Expiry,
		}
	}

	return dbUsers
}

func toDomainOAuthUsers(users []dbOAuthUser) []userr.OpenAuthUser {
	domainUsers := make([]userr.OpenAuthUser, len(users))
	for i, u := range users {
		domainUsers[i] = userr.OpenAuthUser{
			ID:       u.ID,
			Provider: u.Provider,
			Email:    u.Email,
			Name:     u.Name,
			Picture:  u.Picture,
			Token: userr.OpenAuthToken{
				AccessToken:  u.AccessToken,
				TokenType:    u.TokenType,
				RefreshToken: u.RefreshToken,
				Expiry:       u.Expiry,
			},
		}
	}

	return domainUsers
}
