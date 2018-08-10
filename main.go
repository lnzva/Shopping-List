package main

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/gorilla/mux"
)

type shoppingItem struct {
	Id int	`json:"Id"`
	Name string `json:"Name"`
	Price float64 `json:"Price"`
}

type Response struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	ShoppingItem   []shoppingItem `json:"shoppingItem"`
}

//var shoppingItemList = make(map[shoppingItem]bool)
var shoppingItemList []shoppingItem

func getShoppingList (w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Shopping List", ShoppingItem: shoppingItemList})
}

func addShoppingItem (w http.ResponseWriter, r *http.Request) {
	

func main() {

//	shoppingItemList[shoppingItem{1, "Kolagach", 34.5}]
//	shoppingItemList[0] = shoppingItem{1, "kolagach", 34.5};
	shoppingItemList = append(shoppingItemList, shoppingItem{1, "Kolagach", 34.5})

	m := mux.NewRouter()

	m.HandleFunc("/shoppinglist/list",  getShoppingList).Methods("GET")

//	m.HandleFunc("/shoppingList/list", addShoppingItem).Methods("POST")
//	m.HandleFunc("/shoppingList/list", updateShoppingItem).Methods("PUT")
//	m.HandleFunc("/shoppingList/list/{id}", deleteShoppingItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", m))

}
