package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"context"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

//InitDB Struct for Initiaing 
type InitDB struct {
	postgresClient		*sql.DB
	mongoClient			*mongo.Client
}

func (initDB *InitDB) dataInitDB(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	initDB.postgresClient, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	initDB.mongoClient, _ = mongo.Connect(ctx, clientOptions)
	
	initDB.tableCheck()
	initDB.clearTable()

	//Read CSV and Initialize DB
	orders := readCSV("orders.csv")
	initDB.writeOrderData(orders)
	orderItems := readCSV("order_items.csv")
	initDB.writeOrderItemsData(orderItems)
	deliveries := readCSV("deliveries.csv")
	initDB.writeDeliveriesData(deliveries)
	customers := readCSV("customers.csv")
	initDB.writeCustomersData(customers)
	customerCompanies := readCSV("customer_companies.csv")
	initDB.writeCustomerCompaniesData(customerCompanies)
}

func (initDB *InitDB) tableCheck() {
    if _, err := initDB.postgresClient.Exec(tableCreationQuery); err != nil {
        log.Fatal(err.Error())
    }
}

func (initDB *InitDB) clearTable() {
	initDB.postgresClient.Exec("DELETE FROM orders")
	initDB.postgresClient.Exec("DELETE FROM order_items")
	initDB.postgresClient.Exec("DELETE FROM deliveries")
	customerCollection := initDB.mongoClient.Database("packform").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := customerCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	companyCollection := initDB.mongoClient.Database("packform").Collection("customercompanies")
	ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := companyCollection.Drop(ctx2); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS orders
(
	id int PRIMARY KEY,
	created timestamp NOT NULL,
    order_name TEXT NOT NULL,
    customer_id TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS order_items
(
	id Integer PRIMARY KEY,
	order_id Integer NOT NULL,
	price_per_unit float8 NOT NULL,
	quantity Integer NOT NULL,
	product TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS deliveries
(
	id Integer PRIMARY KEY,
	order_item_id Integer NOT NULL,
	delivered_quantity Integer NOT NULL
);
`

func readCSV(name string) [][]string {
	csvFile, err := os.Open(name)
	if err != nil {
		log.Fatal("Unable to open file:", name, err.Error())
	}

	defer csvFile.Close()
	fileReader := csv.NewReader(csvFile)
	fileReader.Comma = ','

	rows, err := fileReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to Read File Properly ::\n", name, err.Error())
	}

	return rows
}

func (initDB *InitDB) writeOrderData(dataRows [][]string) {
	queryStmt, err := initDB.postgresClient.Prepare("INSERT INTO orders(id, created, order_name, customer_id) VALUES($1, $2, $3, $4)")

	if err != nil {
		log.Fatal(err.Error())
	}
	for _, row := range dataRows[1:] {
		_, err = queryStmt.Exec(row[0], row[1], row[2], row[3])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (initDB *InitDB) writeOrderItemsData(dataRows [][]string) {
	queryStmt, err := initDB.postgresClient.Prepare("INSERT INTO order_items(id, order_id, price_per_unit, quantity, product) VALUES($1, $2, $3, $4, $5)")

	if err != nil {
		log.Fatal(err.Error())
	}
	for _, row := range dataRows[1:] {
		if row[2] == "" {
			_, err = queryStmt.Exec(row[0], row[1], float64(0.0), row[3], row[4])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err = queryStmt.Exec(row[0], row[1], row[2], row[3], row[4])
			if err != nil {
				log.Fatal(err)
			}
		}
		
	}
}

func (initDB *InitDB) writeDeliveriesData(dataRows [][]string) {
	queryStmt, err := initDB.postgresClient.Prepare("INSERT INTO deliveries(id, order_item_id, delivered_quantity) VALUES($1, $2, $3)")

	if err != nil {
		log.Fatal(err.Error())
	}
	for _, row := range dataRows[1:] {
		_, err = queryStmt.Exec(row[0], row[1], row[2])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (initDB *InitDB) writeCustomersData(dataRows [][]string) {
	collection := initDB.mongoClient.Database("packform").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	for _, row := range dataRows[1:] {
		_, err := collection.InsertOne(ctx, bson.D{{"userid",row[0]}, {"login", row[1]}, {"password",row[2]}, {"name", row[3]}, {"company_id", row[4]},{"credit_cards", row[5]}})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (initDB *InitDB) writeCustomerCompaniesData(dataRows [][]string) {
	collection := initDB.mongoClient.Database("packform").Collection("customercompanies")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	for _, row := range dataRows[1:] {
		_, err := collection.InsertOne(ctx, bson.D{{"companyid",row[0]},{"companyname",row[1]}})
		if err != nil {
			log.Fatal(err)
		}
	}
}