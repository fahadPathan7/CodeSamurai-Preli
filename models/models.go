package models

type User struct {
	User_id int `json:"user_id" bson:"user_id"`
	User_name string `json:"user_name" bson:"user_name"`
	Balance int `json:"balance" bson:"balance"`
}

type Station struct {
	Station_id int `json:"station_id" bson:"station_id"`
	Station_name string `json:"station_name" bson:"station_name"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude float64 `json:"latitude" bson:"latitude"`
}

type Stop struct {
    Station_id int `json:"station_id" bson:"station_id"`
    Arrival_time string `json:"arrival_time" bson:"arrival_time"`
    Departure_time string `json:"departure_time" bson:"departure_time"`
    Fare int `json:"fare" bson:"fare"`
}

type Train struct {
    Train_id int `json:"train_id" bson:"train_id"`
    Train_name string `json:"train_name" bson:"train_name"`
    Capacity int `json:"capacity" bson:"capacity"`
    Stops []Stop `json:"stops" bson:"stops"`
}

type TrainResponse struct {
    Train_id int `json:"train_id"`
    Train_name string `json:"train_name"`
    Capacity int `json:"capacity"`
    Service_start string `json:"service_start"`
    Service_ends string `json:"service_ends"`
    Num_stations int `json:"num_stations"`
}

type PrintResponseTrainStation struct {
	Train_id int `json:"train_id"`
	Arrival_time *string `json:"arrival_time"`
	Departure_time *string `json:"departure_time"`
}

type Wallet struct {
	Wallet_id int `json:"wallet_id" bson:"wallet_id"`
    Balance int `json:"wallet_balance" bson:"wallet_balance"`
    Wallet_user UserWalletInfo `json:"wallet_user" bson:"wallet_user"`
}

type UserWalletInfo struct {
    User_id int `json:"user_id" bson:"user_id"`
	User_name string `json:"user_name" bson:"user_name"`
}

type WalletMoneyAdd struct {
    Recharge int `json:"recharge" bson:"recharge"`
}

type PurchaseTicketRequest struct {
    WalletID     int    `json:"wallet_id"`
    TimeAfter    string `json:"time_after"`
    StationFrom  int    `json:"station_from"`
    StationTo    int    `json:"station_to"`
}

type TicketResponse struct {
    TicketID  int       `json:"ticket_id"`
    WalletID  int       `json:"wallet_id"`
    Balance   int       `json:"balance"`
    Stations  []StationsTicket `json:"stations"`
}

type StationsTicket struct {
    Station_id int `json:"station_id" bson:"station_id"`
    Train_id int `json:"train_id" bson:"train_id"`
    Arrival_time string `json:"arrival_time" bson:"arrival_time"`
    Departure_time string `json:"departure_time" bson:"departure_time"`
}