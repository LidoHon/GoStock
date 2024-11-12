package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/LidoHon/GoStock/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type response struct{
	ID int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB{
	err :=godotenv.Load()

	if err != nil{
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil{
		panic(err)
	} 
    err = db.Ping()
	if err !=nil{
		panic(err)
	}
	fmt.Println("Successfully created connection to database")
	return db

}


// create stock

func CreateStock(w http.ResponseWriter, r *http.Request){
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err !=nil {
		log.Fatal("unable to decode the request ",err)
	}


	insterID := insertStock(stock)
	res := response{
		ID: insterID,
		Message: "stock created successfully",
	}
	json.NewEncoder(w).Encode(res)
}


// get stock
func GetStock(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, err :=strconv.Atoi(params["id"])
	if err !=nil{
		log.Fatal("unable to get id", err)
	}

	stock, err := getStock(int64(id))
	if err !=nil{
		log.Fatal("unable to get stock", err)
	}
	json.NewEncoder(w).Encode(stock)

}


// get all stock
func GetAllStock(w http.ResponseWriter, r *http.Request){
	stocks, err := getAllStock()
	if err !=nil{
		log.Fatal("unable to get stock", err)
	}
	json.NewEncoder(w).Encode(stocks)
}

// delete stock

func DeleteStock(w http.ResponseWriter, r *http.Request){
	params :=mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err !=nil{
		log.Fatal("unable to get id", err)
	}
	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("stock deleted successfully. total rows/record affected %v", deletedRows)

	res := response{
		ID: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
	
}


// update stock	
func UpdateStock(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err !=nil{
		log.Fatal("unable to get id", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err !=nil{
		log.Fatal("unable to decode the request", err)
	}
	updatedRows := updateStock(int64(id), stock)
	msg:=fmt.Sprintf("stock updated successfully. total rows/record affected %v", updatedRows)

	res := response{
    ID: int64(id),
    Message: msg,
}

	json.NewEncoder(w).Encode(res)

}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query: %v", err)
	}

	fmt.Println("New stock created successfully, ID:", id)
	return id
}



func getStock(id int64)(models.Stock, error){

	db := createConnection()
	defer db.Close()

	sqlStatement :=`SELECT * FROM stocks WHERE stockid = $1`
	var stock models.Stock
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err{
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil
	case  nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row %v", err)
	}
	return stock, err
}

func getAllStock()([]models.Stock, error){
	db := createConnection()

	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err != nil{
		log.Fatalf("unable to execute the query %v", err)
	}
	defer rows.Close()

	for rows.Next(){
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err !=nil{
			log.Fatalf("unable to scan the row %v", err)
		}
		stocks =append(stocks, stock)
	}
	return stocks, err

}

func updateStock(id int64, stock models.Stock) int64{

	db := createConnection()
	defer db.Close()

	sqlStatment := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlStatment, id, stock.Name, stock.Price, stock.Company)
	if err != nil{
		log.Fatalf("unable to excute a query: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil{
		log.Fatalf("error occured while checking the affected rows %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected

}

func deleteStock(id int64)int64{
	db := createConnection()
	defer db.Close()

	sqlStatement :=`DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)
	if err !=nil{
		log.Fatalf("unable to delete the stock %v", err)
	}

	rowsAffected, err:= res.RowsAffected()
	if err !=nil{
		log.Fatalf("error occured while checking the affected rows %v", err)

	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected

}
