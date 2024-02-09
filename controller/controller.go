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
const colNameUser = "user"
const colNameStation = "station"
const colNameTrain = "train"

var collectionUser *mongo.Collection
var collectionStation *mongo.Collection
var collectionTrain *mongo.Collection

func init() {
	fmt.Println("Taking some time for deleting the existing database and creating a new one...")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// delete the database before new connection to avoid duplicate data
	err = client.Database(dbName).Drop(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	collectionUser = client.Database(dbName).Collection(colNameUser)
	collectionStation = client.Database(dbName).Collection(colNameStation)
	collectionTrain = client.Database(dbName).Collection(colNameTrain)

	fmt.Println("Connected to MongoDB!")
}

// inserting user
func insertUser(user models.User) {
	// Inserting data into database
	_, err := collectionUser.InsertOne(context.Background(), user)

	// if there is an error inserting, handle it
	if err != nil {
		fmt.Println(err)
	}
}

// Controller function to insert a user
func InsertUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	insertUser(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// insert a station
func insertStation(station models.Station) {
	_, err := collectionStation.InsertOne(context.Background(), station)

	// if there is an error inserting, handle it
	if err != nil {
		fmt.Println(err)
	}
}

// Controller function to insert a station
func InsertStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var station models.Station
	_ = json.NewDecoder(r.Body).Decode(&station)

	insertStation(station)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(station)
}

// insert a train
func insertTrain(train models.Train) {
	_, err := collectionTrain.InsertOne(context.Background(), train)

	// if there is an error inserting, handle it
	if err != nil {
		fmt.Println(err)
	}
}

// Controller function to insert a train
func InsertTrain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var train models.Train
    _ = json.NewDecoder(r.Body).Decode(&train)

    insertTrain(train)

    // Create the response
    response := models.TrainResponse{
        Train_id: train.Train_id,
        Train_name: train.Train_name,
        Capacity: train.Capacity,
        Service_start: train.Stops[0].Departure_time,
        Service_ends: train.Stops[len(train.Stops)-1].Arrival_time,
        Num_stations: len(train.Stops),
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

// list all stations and sort by id ascending order
func listAllStations() []models.Station {
	var stations []models.Station

	opts := options.Find()

	opts.SetSort(map[string]int{"station_id": 1})

	cur, err := collectionStation.Find(context.Background(), bson.D{{}}, opts)

	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var station models.Station
		err := cur.Decode(&station)

		if err != nil {
			fmt.Println(err)
		}

		stations = append(stations, station)
	}

	// if there is an error getting the stations, handle it
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}

	return stations
}

// Controller function to list all stations
func ListAllStations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stations := listAllStations()

	// If stations is nil, return an empty array instead to ensure JSON consistency
	if stations == nil {
		stations = []models.Station{}
	}

	response := map[string][]models.Station{"stations": stations}

	json.NewEncoder(w).Encode(response)
}

// list all trains at a station
func listTrainsAtStation(stationID int) []models.Train {
	var trains []models.Train

	filter := bson.M{"stops.station_id": stationID}

	cur, err := collectionTrain.Find(context.Background(), filter)

	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var train models.Train
		err := cur.Decode(&train)

		if err != nil {
			fmt.Println(err)
		}

		trains = append(trains, train)
	}

	// if there is an error getting the trains, handle it
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}

	return trains

}

// Controller function to list all trains at a station
func ListTrainsAtStation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	stationID := params["station_id"]

	// convert string to int
	id, _ := strconv.Atoi(stationID)

	// search if station exists
	if !stationExists(id) {
		w.WriteHeader(http.StatusNotFound)
		message := "station with id: " + strconv.Itoa(id) + " was not found"
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	trains := listTrainsAtStation(id)

	// If trains is nil, return an empty array instead to ensure JSON consistency
	if trains == nil {
		trains = []models.Train{}
	}

	for i, train := range trains {
		var stops []models.Stop
		for _, stop := range train.Stops {
			if stop.Station_id == id {
				stops = append(stops, stop)
			}
		}
		trains[i].Stops = stops
	}

	var printResponseTrainStation []models.PrintResponseTrainStation

	for _, train := range trains {
		for _, stop := range train.Stops {
			var arrivalTime, departureTime *string
        	if stop.Arrival_time != "" {
            	arrivalTime = &stop.Arrival_time
       		}
        	if stop.Departure_time != "" {
            	departureTime = &stop.Departure_time
        	}
			response := models.PrintResponseTrainStation{
				Train_id: train.Train_id,
				Arrival_time: arrivalTime,
				Departure_time: departureTime,
			}
			printResponseTrainStation = append(printResponseTrainStation, response)
		}
	}

	// now create the final response
	finalResponse := map[string]interface{}{"station_id": id, "trains": printResponseTrainStation}

	json.NewEncoder(w).Encode(finalResponse)
}

// check if a station exists
func stationExists(id int) bool {
	filter := bson.M{"station_id": id}

	var station models.Station
	err := collectionStation.FindOne(context.Background(), filter).Decode(&station)

	if err != nil {
		return false
	}

	return station.Station_id != 0
}

// print the wallet of a user
func printWallet(userID int) models.User {
	filter := bson.M{"user_id": userID}

	var user models.User
	err := collectionUser.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		fmt.Println(err)
	}

	return user
}

// Controller function to print the wallet of a user
func PrintWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userID := params["wallet_id"]

	// convert string to int
	id, _ := strconv.Atoi(userID)

	// search if user exists
	if !userExists(id) {
		w.WriteHeader(http.StatusNotFound)
		message := "wallet with id: " + strconv.Itoa(id) + " was not found"
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	user := printWallet(id)

	var wallet models.Wallet

	wallet.Wallet_id = user.User_id
	wallet.Balance = user.Balance

	var walletUserInfo models.UserWalletInfo
	walletUserInfo.User_id = user.User_id
	walletUserInfo.User_name = user.User_name

	wallet.Wallet_user = walletUserInfo

	json.NewEncoder(w).Encode(wallet)
}

// check if a user exists
func userExists(id int) bool {
	filter := bson.M{"user_id": id}

	var user models.User
	err := collectionUser.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return false
	}

	return user.User_id != 0
}

// insert money into wallet
func insertMoneyIntoWallet(userID int, recharge int) {
	filter := bson.M{"user_id": userID}

	update := bson.M{"$inc": bson.M{"balance": recharge}}

	_, err := collectionUser.UpdateOne(context.Background(), filter, update)

	if err != nil {
		fmt.Println(err)
	}
}

// Controller function to insert money into wallet
// valid range is 100 to 10000
func InsertMoneyIntoWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userID := params["wallet_id"]

	// convert string to int
	id, _ := strconv.Atoi(userID)

	var walletMoneyAdd models.WalletMoneyAdd
	_ = json.NewDecoder(r.Body).Decode(&walletMoneyAdd)

	// search if user exists
	if !userExists(id) {
		w.WriteHeader(http.StatusNotFound)
		message := "wallet with id: " + strconv.Itoa(id) + " was not found"
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if recharge is in valid range
	if walletMoneyAdd.Recharge < 100 || walletMoneyAdd.Recharge > 10000 {
		w.WriteHeader(http.StatusBadRequest)
		// invalid amount: {recharge_amount}
		message := "invalid amount: " + strconv.Itoa(walletMoneyAdd.Recharge)
		response := map[string]string{"message": message}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	insertMoneyIntoWallet(id, walletMoneyAdd.Recharge)

	user := printWallet(id)

	// it should print the wallet. struct will be Wallet

	var wallet models.Wallet

	wallet.Wallet_id = user.User_id
	wallet.Balance = user.Balance

	var walletUserInfo models.UserWalletInfo
	walletUserInfo.User_id = user.User_id
	walletUserInfo.User_name = user.User_name

	wallet.Wallet_user = walletUserInfo

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// // get users wallet balance
// func checkWalletBalance(userID int) (int, error) {
// 	filter := bson.M{"user_id": userID}

// 	var user models.User
// 	err := collectionUser.FindOne(context.Background(), filter).Decode(&user)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return user.Balance, nil
// }

// // find optimal route and calculate total fare
// func findOptimalRoute(stationFrom int, stationTo int, timeAfter string) []models.StationsTicket {

// }

// // deduct fare from wallet
// func deductFromWallet(userID int, fare int) (int, error) {
// 	filter := bson.M{"user_id": userID}

// 	update := bson.M{"$inc": bson.M{"balance": -fare}}

// 	_, err := collectionUser.UpdateOne(context.Background(), filter, update)

// 	if err != nil {
// 		return 0, err
// 	}

// 	user := printWallet(userID)

// 	return user.Balance, nil
// }

// controller function to purchase a ticket
func PurchaseTicket(w http.ResponseWriter, r *http.Request) {
	// json format
	w.Header().Set("Content-Type", "application/json")
    // Parse and validate request
    var req models.PurchaseTicketRequest
    _ = json.NewDecoder(r.Body).Decode(&req)

	// status code forbidden
	w.WriteHeader(http.StatusForbidden)

	// no ticket available for station: {station_from} to station: {station_to}
	message := "no ticket available for station: " + strconv.Itoa(req.StationFrom) + " to station: " + strconv.Itoa(req.StationTo)

	response := map[string]string{"message": message}

	json.NewEncoder(w).Encode(response)
}

// controller function to purchase a ticket
func BestTicket(w http.ResponseWriter, r *http.Request) {
	// json format
	w.Header().Set("Content-Type", "application/json")
    // Parse and validate request

	args := r.URL.Query()

    var fromm, too string

    if from := args.Get("from"); from != "" {
        fromm = from
    }
	if to := args.Get("to"); to != "" {
        too = to
    }


	// 	args := r.URL.Query()

//     var searchKey, searchVal, sortKey, sortOrder string

//     if title := args.Get("title"); title != "" {
//         searchKey, searchVal = "title", title
//     } else if author := args.Get("author"); author != "" {
//         searchKey, searchVal = "author", author
//     } else if genre := args.Get("genre"); genre != "" {
//         searchKey, searchVal = "genre", genre
//     }

	// status code forbidden
	w.WriteHeader(http.StatusForbidden)

	// no ticket available for station: {station_from} to station: {station_to}
	message := "no ticket available for station: " + fromm + " to station: " + too

	response := map[string]string{"message": message}

	json.NewEncoder(w).Encode(response)
}


// func InsertBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var book models.Book
// 	_ = json.NewDecoder(r.Body).Decode(&book)

// 	insertBook(book)
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(book)
// }

// // updating a book
// func updateABook(bookID int, book models.Book) {
// 	filter := bson.M{"id": bookID}

// 	update := bson.M{ "$set": bson.M{ "title": book.Title, "author": book.Author, "genre": book.Genre, "price": book.Price } }

// 	updateResult, err := collection.UpdateOne(context.Background(), filter, update)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
// }

// // controller function to update a book
// func UpdateABook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	params := mux.Vars(r)

// 	id, _ := strconv.Atoi(params["id"])

// 	// get book from body
// 	var book models.Book
// 	_ = json.NewDecoder(r.Body).Decode(&book)

// 	// check if book exists
// 	if !bookExists(id) {
// 		w.WriteHeader(http.StatusNotFound)
// 		message := "book with id: " + strconv.Itoa(id) + " was not found"
// 		response := map[string]string{"message": message}
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	book.ID = id

// 	updateABook(id, book)
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(book)
// }

// // check if a book exists
// func bookExists(id int) bool {
// 	filter := bson.M{"id": id}

// 	var book models.Book
// 	err := collection.FindOne(context.Background(), filter).Decode(&book)

// 	if err != nil {
// 		return false
// 	}

// 	return book.ID != 0
// }

// // get operation
// // get a book
// func getABook(id int) models.Book {
// 	filter := bson.M{"id": id}

// 	var book models.Book
// 	err := collection.FindOne(context.Background(), filter).Decode(&book)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return book
// }

// // controller function to get a book
// func GetABook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	// get id from url
// 	params := mux.Vars(r)

// 	id, _ := strconv.Atoi(params["id"])
// 	// id is string type

// 	if !bookExists(id) {
// 		w.WriteHeader(http.StatusNotFound)
// 		message := "book with id: " + strconv.Itoa(id) + " was not found"
// 		response := map[string]string{"message": message}
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	book := getABook(id)

// 	// if book is empty

// 	json.NewEncoder(w).Encode(book)
// }

// // get all the Books
// func getAllBooks() []models.Book {
// 	var books []models.Book

// 	cur, err := collection.Find(context.Background(), bson.D{})

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer cur.Close(context.Background())

// 	for cur.Next(context.Background()) {
// 		var book models.Book
// 		err := cur.Decode(&book)

// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		books = append(books, book)
// 	}

// 	// if there is an error getting the books, handle it
// 	if err := cur.Err(); err != nil {
// 		fmt.Println(err)
// 	}

// 	return books
// }

// // controller function to get all books
// func GetAllBooks(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	books := getAllBooks()

// 	// If books is nil, return an empty array instead to ensure JSON consistency
//     if books == nil {
//         books = []models.Book{}
//     }

// 	response := map[string][]models.Book{"books": books}

// 	json.NewEncoder(w).Encode(response)
// }

// // search books
// // api: /api/books?{search_field}={value}&sort={sorting_field}&order={sorting_order}
// func searchBooksByField(searchField string, searchValue string, sortField string, sortOrder string) []models.Book {
// 	var books []models.Book

// 	filter := bson.M{searchField: searchValue}

// 	if searchField == "" {
// 		filter = bson.M{}
// 	}

// 	opts := options.Find()

// 	if sortField != "" {
// 		sortValue := 1
// 		if sortOrder == "DESC" {
// 			sortValue = -1
// 		}
// 		opts.SetSort(bson.D{{ Key: sortField, Value: sortValue }})
// 	}

// 	cur, err := collection.Find(context.Background(), filter, opts)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer cur.Close(context.Background())

// 	for cur.Next(context.Background()) {
// 		var book models.Book
// 		err := cur.Decode(&book)

// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		books = append(books, book)
// 	}

// 	// if there is an error getting the books, handle it
// 	if err := cur.Err(); err != nil {
// 		fmt.Println(err)
// 	}

// 	return books
// }

// func SearchBooks(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	args := r.URL.Query()

//     var searchKey, searchVal, sortKey, sortOrder string

//     if title := args.Get("title"); title != "" {
//         searchKey, searchVal = "title", title
//     } else if author := args.Get("author"); author != "" {
//         searchKey, searchVal = "author", author
//     } else if genre := args.Get("genre"); genre != "" {
//         searchKey, searchVal = "genre", genre
//     }

//     if sort := args.Get("sort"); sort != "" {
//         sortKey = sort
//     } else {
//         sortKey = "id"
//     }

//     if order := args.Get("order"); order != "" {
//         sortOrder = order
//     } else {
//         sortOrder = "ASC"
//     }

//     // get books
//     books := searchBooksByField(searchKey, searchVal, sortKey, sortOrder)

// 	// If books is nil, return an empty array instead to ensure JSON consistency
//     if books == nil {
//         books = []models.Book{}
//     }

// 	response := map[string][]models.Book{"books": books}

// 	json.NewEncoder(w).Encode(response)
// }

// // delete a book
// func deleteABook(id int) {
// 	filter := bson.M{"id": id}

// 	deleteResult, err := collection.DeleteOne(context.Background(), filter)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Printf("Deleted %v documents in the books collection\n", deleteResult.DeletedCount)
// }

// // controller function to delete a book
// func DeleteABook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	params := mux.Vars(r)

// 	id, _ := strconv.Atoi(params["id"])

// 	// check if book exists
// 	if !bookExists(id) {
// 		w.WriteHeader(http.StatusNotFound)
// 		message := "book with id: " + strconv.Itoa(id) + " was not found"
// 		response := map[string]string{"message": message}
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	deleteABook(id)

// 	w.WriteHeader(http.StatusOK)
// 	message := "book with id: " + strconv.Itoa(id) + " was deleted"
// 	response := map[string]string{"message": message}
// 	json.NewEncoder(w).Encode(response)
// }