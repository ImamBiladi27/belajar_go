package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

const (
    dbDriver = "mysql"
    dbUser   = "root"
    dbPass   = ""
    dbName   = "gocrud"
)

func main() {
    
    r := mux.NewRouter()
    r.HandleFunc("/all-user", getAllUserHandler).Methods("GET")
    r.HandleFunc("/user", createUserHandler).Methods("POST")
    r.HandleFunc("/user/{id}", getUserHandler).Methods("GET")
    r.HandleFunc("/user/{id}", updateUserHandler).Methods("PUT")
    r.HandleFunc("/user/{id}", deleteUserHandler).Methods("DELETE")
    r.HandleFunc("/login",loginUserHandler).Methods("POST")
    handler := cors.Default().Handler(r)
    fmt.Println("Server started at :8090")
    log.Fatal(http.ListenAndServe(":8090", handler))
}

// func loginUserHandler(w http.ResponseWriter, r *http.Request){
//     db := connectDB() //memanggil db
//     defer db.Close() //menutup db

//     var input struct {
//         Email string `json:"email"`
//         Password string `json:"password"`
//     }
//     if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
//         http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
//     }
//     email := input.Email
//     password := input.Password
//     // password digunakan untuk autentikasi, misal cek ke database
//     fmt.Println("email:", email,password)

//     var user User
//     err := db.QueryRow("SELECT id, name, email FROM users WHERE email = ? AND password = ?", email, password).Scan(&user.ID, &user.Name, &user.Email)
//     if err != nil {
//         log.Println("Query error:", err)
//         http.Error(w, "Credential salah, coba lagi", http.StatusUnauthorized)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(user)
// }
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
    db := connectDB()
    defer db.Close()

    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var user User
    var hashedPassword string
    err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", input.Email).
        Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)
    if err != nil {
        log.Println("Query error:", err)
        http.Error(w, "Credential salah, coba lagi", http.StatusUnauthorized)
        return
    }
    fmt.Println("hashedPassword:", []byte(hashedPassword), "input.Password:", []byte(input.Password))
    // Compare hashed password
    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
        http.Error(w, "Credential salah, coba lagi", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// CREATE
// func createDataHandler(w http.ResponseWriter, r *http.Request) {
//     db := connectDB()
//     defer db.Close()

//     var user User
//     json.NewDecoder(r.Body).Decode(&user)

//     stmt, err := db.Prepare("INSERT INTO users(name, email ,password) VALUES(?, ?)")
//     if err != nil {
//         panic(err.Error())
//     }
//     _, err = stmt.Exec(user.Name, user.Email)
//     if err != nil {
//         panic(err.Error())
//     }

//     w.WriteHeader(http.StatusCreated)
//     fmt.Fprintln(w, "User created successfully")
// }

// CREATE
// Handler untuk membuat user baru
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    db := connectDB()
    defer db.Close()

    var input struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Hash password sebelum insert
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    stmt, err := db.Prepare("INSERT INTO users(name, email, password) VALUES(?, ?, ?)")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    res, err := stmt.Exec(input.Name, input.Email, string(hashedPassword))
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
        return
    }

    user := User{
        ID:    int(lastID),
        Name:  input.Name,
        Email: input.Email,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}


// READ
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    db := connectDB()
    defer db.Close()

    params := mux.Vars(r)
    id := params["id"]

    var user User
    
    err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)

    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// Read All
func getAllUserHandler(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	fmt.Println("rows",rows)
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UPDATE
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
    db := connectDB()
    defer db.Close()

    params := mux.Vars(r)
    id := params["id"]

    var user User
    json.NewDecoder(r.Body).Decode(&user)

    stmt, err := db.Prepare("UPDATE users SET name = ?, email = ? WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    _, err = stmt.Exec(user.Name, user.Email, id)
    if err != nil {
        panic(err.Error())
    }

    fmt.Fprintln(w, "User updated successfully")
}

// DELETE
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
    db := connectDB()
    defer db.Close()

    params := mux.Vars(r)
    id := params["id"]

    stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    _, err = stmt.Exec(id)
    if err != nil {
        panic(err.Error())
    }

    fmt.Fprintln(w, "User deleted successfully")
}

// Koneksi ke DB
func connectDB() *sql.DB {
    // db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    // if err != nil {
    //     panic(err.Error())
    // }
    // return db
	dsn := "root:@tcp(127.0.0.1:3306)/gocrud"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err.Error())
    }
    return db
}
