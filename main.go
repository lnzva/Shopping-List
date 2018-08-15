package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

/*****************Required Structs***************/

//struct for shopping item
type shoppingItem struct {
	Id    int     `json:"Id,omitempty"`
	Name  string  `json:"Name,omitempty"`
	Price float64 `json:"Price,omitempty"`
	Count int     `json:"Count,omitempty"`
}

//struct for responding through json object
type Response struct {
	Success      int            `json:"success"`
	Message      string         `json:"message"`
	ShoppingItem []shoppingItem `json:"shoppingItem,omitempty"`
}

//struct that represents the user
type User struct {
	UserName string `json:"Username,omitempty"`
	Password string `json:"Password,omitmepty"`
}

/*****************Required global variables***************/

var (
	shoppingItemList []shoppingItem
	userList         = make(map[string]User)
	idCount          int

	access sync.Mutex
)

/*****************Basic Authentication Function***************/

//function to check if user is already logged in
func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("UserName")

	if err == nil {
		access.Lock()
		defer access.Unlock()

		_, flag := userList[cookie.Value]
		return flag
	} else {
		return false
	}
}

//function to handle user registration
func registerUser(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "User already logged in; Please logout to register new user!"})
		return
	}

	var tmpUser User
	var flag bool

	tmpUser.UserName, tmpUser.Password, flag = r.BasicAuth()

	access.Lock()
	defer access.Unlock()

	if flag == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
	} else {
		if _, found := userList[tmpUser.UserName]; found == true {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Username already exists"})
		} else if tmpUser.UserName == "" || tmpUser.Password == "" {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid username or password"})
		} else {
			userList[tmpUser.UserName] = tmpUser
			json.NewEncoder(w).Encode(Response{Success: 1, Message: "User registered successfully"})
		}
	}
}

//function to handle user login
func loginUser(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "User already logged in; Please logout to login another user!"})
		return
	}

	var tmpUser User
	var flag bool

	tmpUser.UserName, tmpUser.Password, flag = r.BasicAuth()

	access.Lock()
	defer access.Unlock()

	if flag == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
	} else {
		val, found := userList[tmpUser.UserName]
		if found == true && val.Password == tmpUser.Password {
			cookie := http.Cookie{Name: "UserName", Value: tmpUser.UserName, Path: "/shoppinglist"}
			http.SetCookie(w, &cookie)
			json.NewEncoder(w).Encode(Response{Success: 1, Message: "User logged-in successfully"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid username or password"})
		}
	}
}

//function to logout user
func logoutUser(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "No user already logged in to logout"})
		return
	}

	cookie := http.Cookie{Name: "UserName", Value: "", Path: "/shoppinglist", Expires: time.Now()}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(Response{Success: 1, Message: "Logged out"})
}

/*****************Handling Shopping List and Items***************/

//get all items in list
func getShoppingList(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please login first"})
		return
	}

	access.Lock()
	defer access.Unlock()

	json.NewEncoder(w).Encode(Response{Success: 1, Message: "The Shopping List", ShoppingItem: shoppingItemList})
}

//add an item to list
func addShoppingItem(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please login first"})
		return
	}

	var shoppingItemInput shoppingItem

	err := json.NewDecoder(r.Body).Decode(&shoppingItemInput)

	if err == nil {
		if shoppingItemInput.Name == "" || shoppingItemInput.Price <= 0 || shoppingItemInput.Count <= 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
			return
		}

		access.Lock()
		defer access.Unlock()

		for _, item := range shoppingItemList {
			if item.Name == shoppingItemInput.Name {
				w.WriteHeader(http.StatusNotAcceptable)
				json.NewEncoder(w).Encode(Response{Success: 0, Message: "Item exists, please update the item"})
				return
			}
		}

		idCount++
		shoppingItemInput.Id = idCount
		shoppingItemList = append(shoppingItemList, shoppingItemInput)

		json.NewEncoder(w).Encode(Response{Success: 1, Message: "Item added successfully"})
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
	}
}

//update an entry in list
func updateShoppingItem(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please login first"})
		return
	}

	var shoppingItemInput shoppingItem
	tmpVar := mux.Vars(r)
	reqIdx, err := strconv.Atoi(tmpVar["id"])

	if err == nil {
		err = json.NewDecoder(r.Body).Decode(&shoppingItemInput)
		shoppingItemInput.Id = reqIdx

		if err == nil {
			if shoppingItemInput.Name == "" || shoppingItemInput.Price <= 0 || shoppingItemInput.Count <= 0 {
				w.WriteHeader(http.StatusNotAcceptable)
				json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid information"})
				return
			}

			access.Lock()
			defer access.Unlock()

			tmpIdx := -1

			for i, item := range shoppingItemList {
				if shoppingItemInput.Name == item.Name && item.Id != reqIdx {
					w.WriteHeader(http.StatusNotAcceptable)
					json.NewEncoder(w).Encode(Response{Success: 0, Message: "Item of same name exists in another index, please choose another name or update that item"})
					return
				}
				if item.Id == reqIdx {
					tmpIdx = i
				}
			}

			if tmpIdx != -1 {
				shoppingItemList[tmpIdx] = shoppingItemInput
				json.NewEncoder(w).Encode(Response{Success: 1, Message: "Updated information", ShoppingItem: shoppingItemList})
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid item ID or item not found"})
}

//delete an item from list
func deleteShoppingItem(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Please login first"})
		return
	}

	tmpVar := mux.Vars(r)
	reqIdx, err := strconv.Atoi(tmpVar["id"])

	if err == nil {
		access.Lock()
		defer access.Unlock()

		for i, item := range shoppingItemList {
			if item.Id == reqIdx {
				shoppingItemList = append(shoppingItemList[:i], shoppingItemList[i+1:]...)
				json.NewEncoder(w).Encode(Response{Success: 1, Message: "Deleted item successfully", ShoppingItem: shoppingItemList})
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid or insufficient information"})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Success: 0, Message: "Invalid item ID or item not found"})
}

/***************Le main function***************/

func main() {
	idCount = 0

	m := mux.NewRouter()

	m.HandleFunc("/shoppinglist/register/", registerUser).Methods("POST")
	m.HandleFunc("/shoppinglist/login/", loginUser).Methods("POST")
	m.HandleFunc("/shoppinglist/logout/", logoutUser).Methods("GET")

	m.HandleFunc("/shoppinglist/list/", getShoppingList).Methods("GET")
	m.HandleFunc("/shoppinglist/list/", addShoppingItem).Methods("POST")
	m.HandleFunc("/shoppinglist/list/{id}", updateShoppingItem).Methods("PUT")
	m.HandleFunc("/shoppinglist/list/{id}", deleteShoppingItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":12345", m))
}
