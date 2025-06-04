package user

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        interface{} `bson:"_id,omitempty" json:"id"`
	Name      string      `bson:"name,omitempty" json:"name"`
	Email     string      `bson:"email,omitempty" json:"email"`
	Password  string      `bson:"password,omitempty" json:"-"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
}

func StructToBsonM(u User) bson.M {
	bsonBytes, err := bson.Marshal(u)
	if err != nil {
		log.Fatal("Marshal error:", err)
	}

	var result bson.M
	err = bson.Unmarshal(bsonBytes, &result)
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}
	return result
}
