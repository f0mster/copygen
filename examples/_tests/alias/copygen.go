// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	service "github.com/switchupcb/copygen/examples/main/domain"
	data "github.com/switchupcb/copygen/examples/main/models"
)

// ModelsToDomain copies a Account, User to a Account.
func ModelsToDomain(tA *service.Account, fA *data.Account, fU *data.User) {
	// Account fields
	tA.Name = fA.Name
	tA.ID = fA.ID

}
