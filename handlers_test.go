package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

// --- Função auxiliar para criar o mock do DB ---
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	var mock sqlmock.Sqlmock
	var dbMock *sql.DB
	var err error

	dbMock, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock do banco: %v", err)
	}

	// substitui o db global pelo mock
	db = dbMock

	// função para fechar o banco ao final
	return mock, func() {
		dbMock.Close()
	}
}

// --- Teste GET /movies ---
func TestGetMovies(t *testing.T) {
	mock, close := setupMockDB(t)
	defer close()

	rows := sqlmock.NewRows([]string{"id", "title", "director", "year"}).
		AddRow(1, "Matrix", "Wachowski", 1999).
		AddRow(2, "Interestelar", "Christopher Nolan", 2014)

	mock.ExpectQuery("SELECT \\* FROM movies").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/movies", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status esperado: %v, obtido: %v", http.StatusOK, status)
	}

	var movies []Movie
	if err := json.Unmarshal(rr.Body.Bytes(), &movies); err != nil {
		t.Errorf("erro ao decodificar resposta: %v", err)
	}

	if len(movies) != 2 {
		t.Errorf("esperado 2 movies, obtido: %d", len(movies))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

// --- Teste GET /movies/{id} ---
func TestGetMovieByID(t *testing.T) {
	mock, close := setupMockDB(t)
	defer close()

	row := sqlmock.NewRows([]string{"id", "title", "director", "year"}).
		AddRow(1, "Matrix", "Wachowski", 1999)

	mock.ExpectQuery("SELECT id, title, director, year FROM movies WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(row)

	req, _ := http.NewRequest("GET", "/movies/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status esperado %d, recebido %d", http.StatusOK, rr.Code)
	}
}

// --- Teste POST /movies ---
func TestCreateMovie(t *testing.T) {
	mock, close := setupMockDB(t)
	defer close()

	mock.ExpectQuery("INSERT INTO movies").
		WithArgs("Matrix", "Wachowski", 1999).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	body := `{"title":"Matrix", "director":"Wachowski", "year":1999}`

	req, _ := http.NewRequest("POST", "/movies", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("status esperado %d, recebido %d", http.StatusCreated, rr.Code)
	}
}

// --- Teste DELETE /movies/{id} ---
func TestDeleteMovie(t *testing.T) {
	mock, close := setupMockDB(t)
	defer close()

	mock.ExpectExec("DELETE FROM movies WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req, _ := http.NewRequest("DELETE", "/movies/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("esperado %d, recebido %d", http.StatusNoContent, rr.Code)
	}
}
