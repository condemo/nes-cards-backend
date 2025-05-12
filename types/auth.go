package types

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int64  `bun:",pk,autoincrement" json:"id"`
	Username string `bun:",notnull" json:"username" validate:"required,min=3,max=12,alphanum"`
	Password string `bun:",notnull" json:"-" validate:"required,min=6"`
}

func (u *User) Validate() error {
	err := validate.Struct(u)
	return err
}
