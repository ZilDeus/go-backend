package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	_ "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"github.com/rs/cors"
)

type Item struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex"`
	Cratio float32
	Pratio float32
	Unit   string
}

type User struct {
	gorm.Model
	Email    string
	Password string
	Meals    datatypes.JSONSlice[Meal]
}

type Meal struct {
	Name        string
	Description string
	Dishes      []Dish
}

type Dish struct {
	Item   int32
	Amount float32
}

type ReturnUser struct {
	Email string
	Meals []ReturnMeal
}

type ReturnMeal struct {
	Name        string
	Description string
	Dishes      []ReturnDish
}

type ReturnDish struct {
	Name   string
	Cratio float32
	Pratio float32
	Unit   string
	Amount float32
}

var db *gorm.DB

func startupServer() {
	r := mux.NewRouter()

	r.HandleFunc("/start", handleStart).Methods("POST")
	r.HandleFunc("/sign-up", handleSignup).Methods("POST")
	r.HandleFunc("/sign-in", handleSignin).Methods("POST")
	r.HandleFunc("/get-user", handleGetUser).Methods("POST")
	r.HandleFunc("/get-item", handleGetItem).Methods("POST")
	r.HandleFunc("/update-meal", handleUpdateMeal).Methods("POST")
	r.HandleFunc("/add-meal", handleAddMeal).Methods("POST")
	r.HandleFunc("/rem-meal", handleRemMeal).Methods("POST")
	r.HandleFunc("/get-items", handleGetAllItems).Methods("POST")
	r.HandleFunc("/add_item_8", handleAddItem).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:      []string{"*","https://svelte-k59b4wquf-zildeus.vercel.app"},
		AllowedMethods:      []string{"GET", "POST"},
		AllowedHeaders:      []string{"*"},
		AllowPrivateNetwork: true,
		AllowCredentials:    true})

	srv := c.Handler(r)

	http.ListenAndServe(":8080", srv)
}
func GetDB() *gorm.DB {
	var err error
  //dsn := "user=postgres password=H0e39EytMYVB12lV host=db.daagzkqbsqqvbecdjtda.supabase.co port=5432 dbname=postgres"
	//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Print("in opening database:", err.Error())
		panic("")
	}
	return db
}

func ValidateUser(w *http.ResponseWriter, req *http.Request, endPoint string) bool {
	id := req.Header.Get("Id")
	if id == "0" || id == "" {
		fmt.Println("in", endPoint)
		fmt.Println("user id not valid")
		http.Error(*w, "user id not valid", http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateKey(w *http.ResponseWriter, req *http.Request, endPoint string) bool {
	Key := req.Header.Get("Key")
	if Key != "1202" {
		fmt.Println("in", endPoint)
		fmt.Println("someone tried to use db with a bad api-key", Key)
		http.Error(*w, "uncorrect api-key", http.StatusBadRequest)
		return false
	}
	return true
}

func startupDatabase() {
	db = GetDB()
	db.AutoMigrate(&User{}, &Item{})
	db.Create(&Item{Name: "صدر دجاج", Cratio: 1.65, Pratio: 0.31, Unit: "g"})
	db.Create(&Item{Name: "بتيتة", Cratio: 0.92, Pratio: 0.02, Unit: "g"})
	db.Create(&Item{Name: "تمن", Cratio: 1.33, Pratio: 0.03, Unit: "g"})
	db.Create(&Item{Name: "خيار", Cratio: 0.13, Pratio: 0.01, Unit: "g"})
	db.Create(&Item{Name: "خس", Cratio: 0.17, Pratio: 0.01, Unit: "g"})
	db.Create(&Item{Name: "لهانة", Cratio: 0.32, Pratio: 0.01, Unit: "g"})
}
func main() {
  fmt.Println("testing")
	startupDatabase()
  fmt.Println("testing")
	startupServer()
  fmt.Println("testing")
}

func GetItemByIdI(id int32) (Item, error) {
	var item Item
	result := db.First(&item, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return item, errors.New("item not found")
	}
	return item, nil
}
func GetItemById(id string) (Item, error) {
	var item Item
	result := db.First(&item, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return item, errors.New("item not found")
	}
	return item, nil
}

func GetItemByName(name string) (Item, error) {
	var item Item
	result := db.First(&item, "Name = ?", name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return item, errors.New("item not found")
	}
	return item, nil
}

func GetUserById(id string) (User, error) {
	var user User
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}
	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	result := db.First(&user, "Email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}
	return user, nil
}

func GetUserByEmailAndPassword(email string, password string) (User, error) {
	var user User
	result := db.First(&user, "Email = ? AND Password = ?", email, password)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}
	return user, nil
}

func enableCors(w *http.ResponseWriter) {
}
func handleAddMeal(w http.ResponseWriter, req *http.Request) {

	if !ValidateKey(&w, req, "AddMeal") {
		return
	}
	if !ValidateUser(&w, req, "AddMeal") {
		return
	}

	var newMeal struct {
		Name        string
		Description string
	}
	err := json.NewDecoder(req.Body).Decode(&newMeal)
	if err != nil {
		fmt.Println("in add-meal", err.Error())
		fmt.Fprint(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _ := GetUserById(req.Header.Get("Id"))
	meals := user.Meals

	for _, meal := range meals {
		if newMeal.Name == meal.Name {
			fmt.Println("meal", newMeal.Name, "already exsist")
			http.Error(w, "meal with the same name already exist", http.StatusBadRequest)
			return
		}
	}
	meals = append(meals, Meal{Name: newMeal.Name, Description: newMeal.Description, Dishes: []Dish{}})
	db.Model(&user).Updates(User{Meals: datatypes.NewJSONSlice(meals)})
	w.WriteHeader(200)
	fmt.Println("add user:", user.ID, "meal", newMeal.Name)
	fmt.Fprintf(w, "succesfully add users meal %s", newMeal.Name)
}
func handleRemMeal(w http.ResponseWriter, req *http.Request) {
  fmt.Println("remove Meal request");
	if !ValidateKey(&w, req, "AddMeal") {
		return
	}
	if !ValidateUser(&w, req, "AddMeal") {
		return
	}

	var removeMeal struct {
    Name string
  };

  json.NewDecoder(req.Body).Decode(&removeMeal)
	fmt.Println("removnig", removeMeal.Name)
	user, _ := GetUserById(req.Header.Get("Id"))
	meals := user.Meals

	for i, meal := range meals {
		if removeMeal.Name == meal.Name {

			copy(meals[i:], meals[i+1:])
			meals[len(meals)-1] = Meal{}
			meals = meals[:len(meals)-1]
			db.Model(&user).Updates(User{Meals: datatypes.NewJSONSlice(meals)})
			w.WriteHeader(200)
			fmt.Println("removed user:", user.ID, "meal", removeMeal.Name)
			fmt.Fprintf(w, "succesfully removed users meal %s", removeMeal.Name)
			return
		}
	}
}
func handleUpdateMeal(w http.ResponseWriter, req *http.Request) {
	if !ValidateKey(&w, req, "AddMeal") {
		return
	}
	if !ValidateUser(&w, req, "AddMeal") {
		return
	}

	user, _ := GetUserById(req.Header.Get("Id"))

	var meals []Meal = user.Meals

	fmt.Printf("old:%+v\n", meals)

	var reqRMeal ReturnMeal

	err := json.NewDecoder(req.Body).Decode(&reqRMeal)
	if err != nil {
		fmt.Println("in update-meal", err.Error())
		fmt.Fprint(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("change:%+v\n", reqRMeal)

	reqMeal := GetMeal(reqRMeal)
	fmt.Printf("change:%+v\n", reqMeal)

	for i, meal := range meals {
		if reqMeal.Name == meal.Name {
			meals[i] = reqMeal
			db.Model(&user).Updates(User{Meals: datatypes.NewJSONSlice(meals)})
			fmt.Printf("new:%+v", meals)
			w.WriteHeader(200)
			fmt.Println("updated user:", user.ID, "meal")
			fmt.Fprintf(w, "succesfully updated users meals")
			return
		}
	}

	meals = append(meals, reqMeal)
	db.Model(&user).Updates(User{Meals: datatypes.NewJSONSlice(meals)})
	w.WriteHeader(200)
	fmt.Println("created user:", user.ID, "meal")
	fmt.Fprintf(w, "succesfully created user meal")
}
func handleGetAllItems(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	if !ValidateKey(&w, req, "GetAllItems") {
		return
	}
	var items []Item
	res := db.Find(&items)

	fmt.Println("found", res.RowsAffected, "items")
	json.NewEncoder(w).Encode(items)
}
func handleGetItem(w http.ResponseWriter, req *http.Request) {

	itemId := req.Header.Get("Item")

	if !ValidateKey(&w, req, "GetItem") {
		return
	}

	item, _ := GetItemById(itemId)

	fmt.Println(item.ID, item.Name, "c:", item.Cratio, "p:", item.Pratio, item.Unit)

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(item)
}
func handleStart(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "server is up")
}
func handleAddItem(w http.ResponseWriter, req *http.Request) {
	if !ValidateKey(&w, req, "AddItem") {
		return
	}
	var item Item

	err := json.NewDecoder(req.Body).Decode(&item)

	if err != nil {
		fmt.Println("error decoing jSON")
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Fprintf(w, "geting Item ERROR")
		return
	}

	_, err = GetItemByName(item.Name)

	if err != nil {
		db.Create(&item)
		fmt.Println(item.ID, item.Name, "c:", item.Cratio, "p:", item.Pratio, item.Unit)
		w.WriteHeader(200)
		fmt.Fprintf(w, "succesfully created item: %s , to ratio of (c:%f,p:%f) , mesured in %s", item.Name, item.Cratio, item.Pratio, item.Unit)
	} else {
		GetItemByName(item.Name)
		fmt.Println(item.ID, item.Name, "c:", item.Cratio, "p:", item.Pratio, item.Unit)
		db.Model(&Item{}).Where("name = ?", item.Name).Updates(&Item{Cratio: item.Cratio, Pratio: item.Pratio, Unit: item.Unit})
		w.WriteHeader(200)
		fmt.Fprintf(w, "succesfully updated item: %s , to ratio of (c:%f,p:%f) , mesured in %s", item.Name, item.Cratio, item.Pratio, item.Unit)
	}
}
func handleGetUser(w http.ResponseWriter, req *http.Request) {
	userId := req.Header.Get("Id")

	if !ValidateKey(&w, req, "GetUser") {
		return
	}

	fmt.Println("the id is ", userId)

	user, err := GetUserById(userId)

	if err != nil {
		fmt.Println("Error getting user, user with ID", userId, "is not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Fprintf(w, "error getting user")
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ReturnUser{Email: user.Email, Meals: GetRMeals(user.Meals)})
}

func GetMeal(rMeal ReturnMeal) Meal {
	var meal Meal = Meal{Name: rMeal.Name, Description: rMeal.Description}
	for _, dish := range rMeal.Dishes {
		item, _ := GetItemByName(dish.Name)
		meal.Dishes = append(meal.Dishes, Dish{Item: int32(item.ID), Amount: dish.Amount})
	}
	return meal
}
func GetRMeals(meals []Meal) []ReturnMeal {
	var rMeals []ReturnMeal
	for _, meal := range meals {
		var rMeal ReturnMeal = ReturnMeal{Name: meal.Name, Description: meal.Description}
		for _, dish := range meal.Dishes {
			item, _ := GetItemByIdI(dish.Item)
			rMeal.Dishes = append(rMeal.Dishes, ReturnDish{Amount: dish.Amount, Unit: item.Unit, Name: item.Name, Cratio: item.Cratio, Pratio: item.Pratio})
		}
		rMeals = append(rMeals, rMeal)
	}
	return rMeals
}

func handleSignin(w http.ResponseWriter, req *http.Request) {
	fmt.Println("sign in request")

	var user User

	if !ValidateKey(&w, req, "SignIn") {
		return
	}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		fmt.Fprintf(w, "ERROR in signin")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(user.ID, "::", user.Email, "::", user.Password)

	user, err = GetUserByEmailAndPassword(user.Email, user.Password)

	fmt.Println(user.ID, "::", user.Email, "::", user.Password)

	if err != nil {
		fmt.Println("ERROR account or password are not correct")
		w.WriteHeader(301)
		fmt.Fprintf(w, "ERROR account or password are not correct")
	} else {
		w.WriteHeader(200)
		fmt.Println("user with email : ", user.Email, " and password : ", user.Password)
		fmt.Fprintf(w, "%d", user.ID)
	}
}
func handleSignup(w http.ResponseWriter, req *http.Request) {
	fmt.Println("sign up request")

	var user User

	if !ValidateKey(&w, req, "Signup") {
		return
	}

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		fmt.Fprintf(w, "ERROR in signup")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = GetUserByEmail(user.Email)

	if err != nil {
		fmt.Println("creating user with email : ", user.Email, " and password : ", user.Password)
		var Meals = []Meal{{Name: "صدر دجاج و بتيتة", Dishes: []Dish{{Item: 1, Amount: 200}, {Item: 2, Amount: 100}}, Description: "بروتين+سعرات+طيب"}}
		user.Meals = datatypes.NewJSONSlice(Meals)
		db.Create(&user)
		fmt.Println("id:", user.ID)
		w.WriteHeader(200)
		fmt.Fprintf(w, "%d", user.ID)
		return
	} else {
		fmt.Println("email:", user.Email, " is  already in use")
		w.WriteHeader(301)
		fmt.Fprintf(w,"email already in use")
		return
	}
}
