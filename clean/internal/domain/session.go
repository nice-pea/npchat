package domain

type Session struct {
	ID     string
	UserID string
	Token  string
	Status int
}

const (
	SessionStatusNew      = 1
	SessionStatusPending  = 2
	SessionStatusVerified = 3
	SessionStatusExpired  = 4
	SessionStatusRevoked  = 5
	SessionStatusFailed   = 6
)

type SessionsRepository interface {
	Save(session Session) error
	List(filter SessionFilter) ([]Session, error)
	Delete(id string) error
}

type SessionFilter struct {
	Token string
}
