package domains

import (
	"context"
	"time"
)

// User represents a user in the system. Users are typically created via OAuth
// using the AuthService but users can also be created directly for testing.
type User struct {
	ID int64 `json:"id"`

	// User's preferred name and email
	Name  string `json:"name"`
	Email string `json:"email"`

	// Timestamps for user creation & last update
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// List of associated OAuth authentication objects.
	// Currently only Github is supported so there should only be a maximum of one.
	Auths []*Auth `json:"auths"`
}

// Validate returns an error if the user contains invalid fields.
// This only performs basic validation.
func (u *User) Validate() error {
	if u.Name == "" {
		return Errorf(EINVALID, "User name required.")
	}

	return nil
}

// AvatarURL returns a URL to the avatar image for the user.
// This loops over all auth providers to find the first available avatar.
// Returns blank string if no avatar URL available.
func (u *User) AvatarURL(size int) string {
	for _, auth := range u.Auths {
		if s := auth.AvatarURL(size); s != "" {
			return s
		}
	}

	return ""
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID along with their associated auth objects.
	// Returns ENOTFOUND if the user does not exist.
	FindUserByID(ctx context.Context, id int64) (*User, error)

	// Retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)

	// Creates a new user. This is only used for testing since users are typically
	// created during OAuth creation process in AuthService.CreateAuth().
	CreateUser(ctx context.Context, user *User) error

	// Updates a user object. Returns EUNAUTHORIZED if current user is not
	// the user that is being updated. Returns ENOTFOUND if user does not exist.
	UpdateUser(ctx context.Context, id int64, upd UserUpdate) (*User, error)

	// Permanently deletes a user. Returns EUNAUTHORIZED if current user is not
	// the user being deleted. Returns ENOTFOUND if user does not exist.
	DeleteUser(ctx context.Context, id int64) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID    *int64  `json:"id"`
	Email *string `json:"email"`

	// Restrict to subset of results
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
