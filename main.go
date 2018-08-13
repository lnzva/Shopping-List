package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"fmt"

	"github.com/gorilla/mux"
)

type shoppingItem struct {
	Id int	`json:"Id, omitempty"`
	Name string `json:"Name, omitempty"`
	Price float64 `json:"Price, omitempty"`
	Count int `json:"Count, omitempty"`
}

type Response struct {
	Success int    `json:"success"`
	Message string `json:"message"`
	ShoppingItem   []shoppingItem `json:"shoppingItem, omitempty"`
}

var shoppingItemList []shoppingItem
var idCount int

func getShoppingList (w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Shopping List", ShoppingItem: shoppingItemList})
}

func addShoppingItem (w http.ResponseWriter, r *http.Request) {
	var shoppingItemInput shoppingItem

	err := json.NewDecoder(r.Body).Decode(&shoppingItemInput)

	if err == nil {
		if shoppingItemInput.Name == "" || shoppingItemInput.Price < 0 || shoppingItemInput.Count < 0  {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
			return
		}
		idCount++
		shoppingItemInput.Id = idCount
		shoppingItemList = append(shoppingItemList, shoppingItemInput)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Success: 1, Message: "Item added successfully"})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
	}
}

func updateShoppingItem (w http.ResponseWriter, r *http.Request) {
	var shoppingItemInput shoppingItem
	tmpVar := mux.Vars(r)
	reqIdx, _ := strconv.Atoi(tmpVar["id"])
	
	json.NewDecoder(r.Body).Decode(&shoppingItemInput)

	if shoppingItemInput.Name == "" || shoppingItemInput.Price < 0 || shoppingItemInput.Count < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid information"})
		return
	}

	for i, item := range shoppingItemList {
			if item.Id == reqIdx {
				fmt.Println(shoppingItemInput.Id)
				shoppingItemList[i] = shoppingItemInput
				json.NewEncoder(w).Encode(Response{Success: 1, Message: "Updated information", ShoppingItem: shoppingItemList})
				return
			}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success :0, Message: "Invalid item ID or item not found"})
}

func deleteShoppingItem (w http.ResponseWriter, r *http.Request) {
	tmpVar := mux.Vars(r)
	reqIdx, _ := strconv.Atoi(tmpVar["id"])

	for i, item := range shoppingItemList {
		if item.Id == reqIdx {
			shoppingItemList = append(shoppingItemList[:i], shoppingItemList[i + 1:]...)
			json.NewEncoder(w).Encode(Response{Success: 1, Message: "Deleted item successfully", ShoppingItem: shoppingItemList})
			return
		}
	}
}


func main() {
	idCount = 0
	fmt.Println("DOING")
	shoppingItemList = append(shoppingItemList, shoppingItem{1, "Kolagach", 34.5, 1})

	m := mux.NewRouter()

	m.HandleFunc("/shoppinglist/list",  getShoppingList).Methods("GET")

	m.HandleFunc("/shoppinglist/list", addShoppingItem).Methods("POST")
	m.HandleFunc("/shoppinglist/list/{id}", updateShoppingItem).Methods("PUT")
	m.HandleFunc("/shoppinglist/list/{id}", deleteShoppingItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", m))
}
