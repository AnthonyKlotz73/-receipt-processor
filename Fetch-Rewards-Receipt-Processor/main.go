package main

import (
	"FetchRewardsChallenge/receiptstructs"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var dummyRecept []receiptstructs.Receipt
var items1 []receiptstructs.Item
var items2 []receiptstructs.Item

func main() {
	items1 = append(items1, receiptstructs.Item{ShortDescription: "Apple Barrel Pewter Gray", Price: "1.25"}, receiptstructs.Item{ShortDescription: "Max Hot Soda", Price: "12.75"})
	items2 = append(items2, receiptstructs.Item{ShortDescription: "DeWalt 20v-60A Battery Drill", Price: "100.99"})
	r := mux.NewRouter()

	dummyRecept = append(dummyRecept, receiptstructs.Receipt{Retailer: "Wal-Mart", PurchaseDate: "2025-02-23", PurchaseTime: "14:26", Items: items1, Total: "14.00"})
	dummyRecept = append(dummyRecept, receiptstructs.Receipt{Retailer: "Home-Depo", PurchaseDate: "2025-01-22", PurchaseTime: "12:24", Items: items2, Total: "100.99"})
	/*jFileArray := []string{"examples/Ex1.json", "examples/Ex2.json", "examples/Ex3.json", "examples/Ex4.json"}
	for _, i := range jFileArray {
		var i1 int32 = receiptstructs.PrintReceiptFJson(i)
		if i1 > 0 {
			dummyId := uuid.New().String()
			InMemoryReceptMap[dummyId] = i1
			fmt.Println(dummyId)
		}
	}*/
	//fmt.Println(dummyId)
	r.HandleFunc("/receipts/process", GetProcessReceipts).Methods("Get")
	r.HandleFunc("/receipts/", GetIds).Methods("Get")
	r.HandleFunc("/receipts/process", PostProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPointsReceipts).Methods("GET")
	//log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8000", r))

	//receiptstructs.PrintReceiptFJson("examples/Ex1.json")
	//receiptstructs.PrintReceiptFJson("examples/Ex2.json")
	//receiptstructs.PrintReceiptFJson("examples/Ex3.json")
	//receiptstructs.PrintReceiptFJson("examples/Ex4.json")

}

// string is the id and int32 is the points
var InMemoryReceptMap = make(map[string]int32)

/*
A Get that returns the ids of receipts that have been processed already
*/
func GetIds(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(InMemoryReceptMap)
}

/*
A get process that get the processed receipts
*/
func GetProcessReceipts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(dummyRecept)
}

/*
A post process that sends a json file to add a new id and points
*/
func PostProcessReceipts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	//var dummyId string
	var c receiptstructs.Receipt
	_ = json.NewDecoder(req.Body).Decode(&c)
	dummyString, points := receiptstructs.ReceiptRewards(c)
	//dummyString, points := receiptstructs.PrintBreakDown(ReceptRewards(itemsJ1))
	//var i1 int32 = receiptstructs.PrintReceiptFJson(_)
	if points > 0 {
		dummyId := uuid.New().String()
		InMemoryReceptMap[dummyId] = points
		fmt.Println(dummyId)
		receiptstructs.PrintBreakDown(dummyString, points)
	}
}

/*
a Get process that gets the amount of points that an id from a receipt has earned
*/
func GetPointsReceipts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var id string = params["id"]
	receiptPointsById, ok := InMemoryReceptMap[id]
	if !ok {
		json.NewEncoder(res).Encode("Not Found")
		return
	}
	json.NewEncoder(res).Encode(receiptPointsById)
}
