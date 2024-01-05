package response

import (
	"time"

	"github.com/yaroyan/ms/authn/domain/model/authentication"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewResponse(u *authentication.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
