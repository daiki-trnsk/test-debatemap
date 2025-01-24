package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	Nickname        *string            `json:"nickname"   validate:"required,min=2,max=10"`
	Password        *string            `json:"password"   validate:"required,min=6"`
	Email           *string            `json:"email"      validate:"email,required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `josn:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updtaed_at"`
	User_ID         string             `json:"user_id"`
}

type DebateMap struct {
    ID               primitive.ObjectID     `bson:"_id,omitempty"`
    RegistrationDate time.Time              `bson:"registration_date"`
	RootNodeTopic    string                 `bson:"root_node_topic" json:"root_node_topic"`
    NodesJSON        map[string]interface{} `bson:"nodes_json" json:"nodes_json"`
    UserID           primitive.ObjectID     `bson:"user_id"`
}