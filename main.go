package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Movie representa a estrutura de dados de um filme.
type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

var db *sql.DB

func main() {
	// Conexão com o banco de dados PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getRequiredEnv("DB_HOST"),
		getRequiredEnv("DB_PORT"),
		getRequiredEnv("DB_USER"),
		getRequiredEnv("DB_PASSWORD"),
		getRequiredEnv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria a tabela de filmes se ela não existir
	createTable()

	// Inicializa o roteador
	router := mux.NewRouter()

	// Define as rotas da API
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Inicia o servidor
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		director TEXT NOT NULL,
		year INT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// --- Handlers da API ---

func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Year); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var m Movie
	err = db.QueryRow("SELECT id, title, director, year FROM movies WHERE id = $1", id).Scan(&m.ID, &m.Title, &m.Director, &m.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Filme não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO movies (title, director, year) VALUES ($1, $2, $3) RETURNING id", m.Title, m.Director, m.Year).Scan(&m.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE movies SET title = $1, director = $2, year = $3 WHERE id = $4", m.Title, m.Director, m.Year, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Filme não encontrado", http.StatusNotFound)
		return
	}

	m.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Filme não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão.
func getRequiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Erro fatal: Variável de ambiente não encontrada: %s", key)
	}
	return value
}
