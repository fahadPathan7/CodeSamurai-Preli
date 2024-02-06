package models

type Book struct {
	ID     int     `json:"id" bson:"id"`
	Title  string  `json:"title" bson:"title"`
	Author string  `json:"author" bson:"author"`
	Genre  string  `json:"genre" bson:"genre"`
	Price  float64 `json:"price" bson:"price"`
}