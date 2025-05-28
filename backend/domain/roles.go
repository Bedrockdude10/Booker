package domain

import "github.com/Bedrockdude10/Booker/backend/utils"

const (
	RolePromoter = "promoter"
	RoleArtist   = "artist"
	RoleAdmin    = "admin"
)

var ValidRoles = utils.NewSet(
	RolePromoter, RoleArtist, RoleAdmin,
)

// Returns if the provided role exists in set of valid roles
func HasRole(role string) bool {
	return ValidRoles.Has(role)
}
