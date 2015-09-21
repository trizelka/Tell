package models

import "gopkg.in/mgo.v2/bson"

type (
	// User represents the structure of our resource
	User struct {
		Id     		bson.ObjectId 	`json:"id" bson:"_id"`
		Regid		string		`json:"regid" bson:"regid"`
		Name   		string        	`json:"name" bson:"name"`
		Email 		string        	`json:"email" bson:"email"`
		Phone		string		`json:"phone" bson:"phone"`
	}
	
)
