package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DebateMap struct {
    ID               primitive.ObjectID     `bson:"_id,omitempty"`
    RegistrationDate time.Time              `bson:"registration_date"`
	RootNodeTopic    string                 `bson:"root_node_topic" json:"root_node_topic"`
    NodesJSON        map[string]interface{} `bson:"nodes_json" json:"nodes_json"`
    UserID           primitive.ObjectID     `bson:"user_id"`
}