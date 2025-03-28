package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Golang-Personal-Projects/Go-Projects/06-Go-Postgres-API/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Create DB connection
func createConnection() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!!!")
	}
	// Open connection to SQL
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic("Error connecting to the postgres database😱")
	}

	err = db.Ping()
	if err != nil {
		panic("Error could not ping the database🤯")
	}

	return db
}

// CreateStock create a new stock
func CreateStock(w http.ResponseWriter, req *http.Request) {
	var stock models.Stock

	if err := json.NewDecoder(req.Body).Decode(&stock); err != nil {
		log.Fatalf("Unable to decode the request body %v", err)
	} // decode the json-values into the address of models.Stock

	insertID := insertStock(stock)

	response := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	// return the response by encoding
	json.NewEncoder(w).Encode(response)
}

// GetStock Get a stock item by id
func GetStock(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to int %v", err)
	}

	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to find the stock %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

// DeleteStock Delete a stock item
func DeleteStock(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to int %v", err)
	}

	deleteRows, err := deleteStock(int64(id))
	if err != nil {
		log.Fatalf("unable to delete item %v", err)
	}

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records deleted %v", deleteRows)

	res := response{ID: int64(id), Message: msg}

	json.NewEncoder(w).Encode(res)

}

// UpdateStock Update a stock item
func UpdateStock(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert string to interger %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(req.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body, %v", err)
	}

	updateRows, err := updateStock(int64(id), stock)
	if err != nil {
		log.Fatalf("Unable to update the stock %v", err)
	}

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updateRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// GetAllStocks Get all stock items
func GetAllStocks(w http.ResponseWriter, req *http.Request) {

	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("Unable to get all the stocks, %v", err)
	}
	json.NewEncoder(w).Encode(stocks)
}

// insert stock function

func insertStock(s models.Stock) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	if err := db.QueryRow(sqlStatement, s.Name, s.Price, s.Company).Scan(&id); err != nil {
		log.Fatalf("Unable to insert item into row, %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	err := db.QueryRow(sqlStatement, id).Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch {

	case errors.Is(err, sql.ErrNoRows):
		fmt.Println("No rows were returned")
		return stock, nil
	case err == nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the rows %v", err)
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error) {

	db := createConnection()

	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) (int64, error) {
	db := createConnection()

	defer db.Close()
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	result, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking rows affected %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected, err
}

func deleteStock(id int64) (int64, error) {
	db := createConnection()

	defer db.Close()
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking rows affected %v", err)
	}
	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected, err
}
