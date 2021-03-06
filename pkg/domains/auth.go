package domains

import (
	"context"
	"fmt"
	"time"
)

// Authentication providers.
// Currently we only support GitHub, but any OAuth provider can be supported.
const (
	AuthSourceGitHub = "github"
)

// Auth represents a set of OAuth credentials. These are linked to a User so a
// single user could authenticate through multiple providers.
//
// The authentication system links users by email address, however, some GitHub
// users don't provide their email publicly, so we may not be able to link them
// by email address.
type Auth struct {
	ID int64 `json:"id"`

	// User can have one or more methods of authentication.
	// However, only one per source is allowed per user.
	UserID int64 `json:"userID"`
	User   *User `json:"user"`

	// The authentication source & the source provider's user ID.
	// Source can only be "github" currently.
	Source   string `json:"source"`
	SourceID string `json:"sourceID"`

	// OAuth fields returned from the authentication provider.
	// GitHub does not use refresh tokens but the field exists for future providers.
	AccessToken  string     `json:"-"`
	RefreshToken string     `json:"-"`
	Expiry       *time.Time `json:"-"`

	// Timestamps of creation & last update
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate returns an error if any fields are invalid on the Auth object.
// This can be called by the SQLProvider implementation to do some basic checks.
func (a *Auth) Validate() error {
	if a.UserID == 0 {
		return Errorf(EINVALID, "User required.")
	} else if a.Source == "" {
		return Errorf(EINVALID, "Source required.")
	} else if a.SourceID == "" {
		return Errorf(EINVALID, "Source ID required.")
	} else if a.AccessToken == "" {
		return Errorf(EINVALID, "Access token required.")
	}

	return nil
}

// AvatarURL returns the URL to the user's avatar hosted by auth source.
// Returns an empty string if the auth source is invalid.
func (a *Auth) AvatarURL(size int) string {
	switch a.Source {
	case AuthSourceGitHub:
		return fmt.Sprintf("https://avatars1.githubusercontent.com/u/%s?s=%d", a.SourceID, size)
	default:
		return ""
	}
}

// AuthService represents a service for managin auths.
type AuthService interface {
	// Looks up an auth object by ID along with the associated user.
	// Returns ENOTFOUND if ID does not exist.
	FindAuthByID(ctx context.Context, id int64) (*Auth, error)

	// Retrieves authentication objects based on filter. Also returns the
	// total number of objects that match the filter. This may differ from the
	// returned object count if the Limit field is set.
	FindAuths(ctx context.Context, filter AuthFilter) ([]*Auth, int, error)

	// Creates a new auth object. If a user is attached to auth, then
	// the auth object is linked to an existing user. Otherwise, a new user
	// object is created.
	//
	// On success, the auth.ID is set to the new authentication ID.
	CreateAuth(ctx context.Context, auth *Auth) error

	// Permanently deletes an auth object from the system by ID.
	// The parent user object is not removed.
	DeleteAuth(ctx context.Context, id int64) error
}

// AuthFilter represents a filter accepted by FindAuths().
type AuthFilter struct {
	// Filtering fields
	ID       *int64  `json:"id"`
	UserID   *int64  `json:"userID"`
	Source   *string `json:"source"`
	SourceID *string `json:"sourceID"`

	// Reestricts results to a subset of the total range.
	// Can be used for pagination.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
