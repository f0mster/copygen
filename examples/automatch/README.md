# Example: Automatch

The automatch example uses the automatcher to match three models with varying level of depth:
```go
// 0
type Account // domain
    // 1
    ID      int
    Name    string
    Email   string
            // 2
    User    domain.DomainUser
                UserID   int
                Username string    
// 0
type User    // models
    // 1
    UserID    int
    Username  string
              // 2
    UserData  models.UserData
                  Options map[string]interface{}
                  // 3
                  Data    models.Data
                        ID      int
// 0            
type Account // models
    // 1
    ID       int
    Name     string
    Password string
    Email    string
```

## YML

```yml
# Define where the code will be generated.
generated:
  setup: ./setup.go
  output: ../copygen.go

# Templates and custom options aren't used for this example.
```

## Go

Specify a depth-level of two for the subfields of `domain.Account`. Specify a depth-level of 1 for the `models.User` field. Keep in mind that not specifying a depth-level for `models.Account` would result in the same outcome, as Copygen uses a maximum depth by default.

```go
// Copygen defines the functions that will be generated.
type Copygen interface {
	// depth domain.Account 2
	// depth models.User 1
	ModelsToDomain(*models.Account, *models.User) *domain.Account
}
```

_Use pointers to avoid allocations._

## Output

`copygen -yml path/to/yml`

```go
// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.

// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/automatch/domain"
	"github.com/switchupcb/copygen/examples/automatch/models"
)

// ModelsToDomain copies a Account, User to a Account.
func ModelsToDomain(tA *domain.Account, fA *models.Account, fU *models.User) {
	// Account fields
	tA.User.Username = fU.Username
	tA.User.UserID = fU.UserID
	tA.Email = fA.Email
	tA.Name = fA.Name
	tA.ID = fA.ID

}
```