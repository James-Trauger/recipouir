package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	First    *string            `json:"first" validate:"required, min=2, max=50"`
	Last     *string            `json:"last" validate:"required, min=2, max=50"`
	Pass     *string            `json:"password" validate:"required, min=8"`
	Email    *string            `json:"email" validate:"email, required"`
	Token    *string            `json:"token"`
	UserType *string            `json:"UserType" validate:"required, eq=ADMIN|eq=USER"`
	Refresh  *string            `json:"refresh"`
}
