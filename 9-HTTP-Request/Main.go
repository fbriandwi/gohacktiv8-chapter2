// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Role     string `json:"role"`
// }

// type Product struct {
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Price int    `json:"price"`
// }

// var users = []User{
// 	{ID: 1, Username: "user1", Password: "password1", Role: "user"},
// 	{ID: 2, Username: "user2", Password: "password2", Role: "admin"},
// }

// var products = []Product{
// 	{ID: 1, Name: "Product 1", Price: 1000},
// 	{ID: 2, Name: "Product 2", Price: 2000},
// 	{ID: 3, Name: "Product 3", Price: 3000},
// }

// var mySigningKey = []byte("mysupersecretkey")

// func main() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/register", register).Methods("POST")
// 	r.HandleFunc("/login", login).Methods("POST")

// 	// Middleware Authentication
// 	r.Use(authentication)

// 	// Middleware Authorization Multi Level User
// 	r.Use(authorization)

// 	// API Product
// 	r.HandleFunc("/products", getProducts).Methods("GET")
// 	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
// 	r.HandleFunc("/products", createProduct).Methods("POST")
// 	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
// 	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

// 	// Middleware Authorization Access Product by ID
// 	r.Use(authorizationByID)

// 	log.Fatal(http.ListenAndServe(":8000", r))
// }

// func register(w http.ResponseWriter, r *http.Request) {
// 	var newUser User
// 	err := json.NewDecoder(r.Body).Decode(&newUser)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Check if user already exists
// 	for _, user := range users {
// 		if newUser.Username == user.Username {
// 			http.Error(w, "Username already taken", http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	// Add new user
// 	newUser.ID = len(users) + 1
// 	users = append(users, newUser)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newUser)
// }



// func login(w http.ResponseWriter, r *http.Request) {
// 	var credentials User
// 	err := json.NewDecoder(r.Body).Decode(&credentials)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Find user by username and password
// 	var user User
// 	for _, u := range users {
// 		if u.Username == credentials.Username && u.Password == credentials.Password {
// 			user = u
// 			break
// 		}
// 	}

// 	// Add new login
// 	credentials User.ID = len(user) + 1
// 	users = append(users, credentials User)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(credentials User)
// }



package main

import (
    "fmt"
    "net/http"
    // "strconv"

    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
)

type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price int    `json:"price"`
}

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

var products []Product
var users []User
var jwtKey = []byte("my_secret_key")

func createToken(user User) (string, error) {
    claims := jwt.MapClaims{}
    claims["authorized"] = true
    claims["user_id"] = user.ID
    claims["user_email"] = user.Email
    claims["user_role"] = user.Role
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprint(w, "Missing Authorization Token")
            return
        }
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })
        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                w.WriteHeader(http.StatusUnauthorized)
                fmt.Fprint(w, "Invalid Authorization Token Signature")
                return
            }
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprint(w, "Invalid Authorization Token")
            return
        }
        if !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprint(w, "Invalid Authorization Token")
            return
        }
        next.ServeHTTP(w, r)
    }
}


