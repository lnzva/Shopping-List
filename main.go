package main

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/gorilla/mux"
)

type shoppingItem struct {
	Id int	`json:"Id, omitempty"`
	Name string `json:"Name"`
	Price float64 `json:"Price, omitempty"`
}

type Response struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	ShoppingItem   []shoppingItem `json:"shoppingItem, omitempty"`
}

//var shoppingItemList = make(map[shoppingItem]bool)
var shoppingItemList []shoppingItem
var itemCount int

func getShoppingList (w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Shopping List", ShoppingItem: shoppingItemList})
}

func addShoppingItem (w http.ResponseWriter, r *http.Request) {
	var shoppingItemInput shoppingItem
	
	err := json.NewDecoder(r.Body).Decode(&shoppingItemInput)

	if err == nil {
		if shoppingItemInput.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
			return
		}
		itemCount++
		shoppingItemInput.Id = itemCount
		shoppingItemList = append(shoppingItemList, shoppingItemInput)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "Item added successfully"})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
	}
}

func main() {
//	shoppingItemList[shoppingItem{1, "Kolagach", 34.5}]
//	shoppingItemList[0] = shoppingItem{1, "kolagach", 34.5};
	itemCount = 0
	shoppingItemList = append(shoppingItemList, shoppingItem{1, "Kolagach", 34.5})

	m := mux.NewRouter()

	m.HandleFunc("/shoppinglist/list",  getShoppingList).Methods("GET")

	m.HandleFunc("/shoppingList/list", addShoppingItem).Methods("POST")
//	m.HandleFunc("/shoppingList/list", updateShoppingItem).Methods("PUT")
//	m.HandleFunc("/shoppingList/list/{id}", deleteShoppingItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", m))
}
