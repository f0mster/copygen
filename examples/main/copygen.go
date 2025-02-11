// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	c "strconv"

	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

/* Define the function and field this converter is applied to using regex. */
// Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return c.Itoa(i)
}

// ModelsToDomain copies a Account, User to a Account.
func ModelsToDomain(tA *domain.Account, fA *models.Account, fU *models.User) {
	// Account fields
	tA.Name = fA.Name
	tA.UserID = Itoa(fU.UserID)
	tA.ID = fA.ID

}
