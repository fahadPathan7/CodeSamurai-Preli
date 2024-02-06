package controller

import (
	"context"
	"encoding/json"
	"log"
	"samurai/models"
	"strconv"

	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://fahadpathan56:fahadpathan@cluster0.oxoqi6z.mongodb.net/"

const dbName = "samurai"
const colName = "books"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Connected to MongoDB!")
}


// insering book
func insertBook(book models.Book) {
	// Inserting data into database
	_, err := collection.InsertOne(context.Background(), book)

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}
}


func InsertBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	insertBook(book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}


// updating a book
func updateABook(bookID int, book models.Book) {
	filter := bson.M{"id": bookID}

	update := bson.M{ "$set": bson.M{ "title": book.Title, "author": book.Author, "genre": book.Genre, "price": book.Price } }

	updateResult, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

// controller function to update a book
func UpdateABook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	// get book from body
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	// check if book exists
	if !bookExists(id) {
		w.WriteHeader(http.StatusNotFound)
		message := "Book with id: " + strconv.Itoa(id) + " was not found."
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	updateABook(id, book)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// check if a book exists
func bookExists(id int) bool {
	filter := bson.M{"id": id}

	var book models.Book
	err := collection.FindOne(context.Background(), filter).Decode(&book)

	if err != nil {
		return false
	}

	return book.ID != 0
}

// get operation
// get a book
func getABook(id int) models.Book {
	filter := bson.M{"id": id}

	var book models.Book
	err := collection.FindOne(context.Background(), filter).Decode(&book)

	if err != nil {
		log.Fatal(err)
	}

	return book
}

// controller function to get a book
func GetABook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get id from url
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])
	// id is string type

	book := getABook(id)

	// if book is empty
	if book.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		message := "Book with id: " + strconv.Itoa(id) + " was not found."
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(book)
}


// get all the Books
func getAllBooks() []models.Book {
	var books []models.Book

	cur, err := collection.Find(context.Background(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var book models.Book
		err := cur.Decode(&book)

		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}

	// if there is an error getting the books, handle it
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return books
}

// controller function to get all books
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books := getAllBooks()

	response := map[string][]models.Book{"books": books}

	json.NewEncoder(w).Encode(response)
}

// search books