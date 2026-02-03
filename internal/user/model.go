package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// model for starting-point
type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	PasswordHash string `json:"-" bson:"passwordhash"`
	Role string `json:"role" bson:"role"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type PublicUser struct{
	ID string `json:"id"`
	Email string `json:"email"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}


// Helper func()
func ToPublic(u User)PublicUser{
	return PublicUser{
		ID:u.ID.Hex(),
		Email: u.Email,
		Role: u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}