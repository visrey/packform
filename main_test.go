package main

import (
    "os"
    "testing"

    "net/http"
    "net/http/httptest"
    "context"
    "time"
    "log"
)

var mainApp App
var initDB InitDB

func TestMain(m *testing.M) {
	initDB.dataInitDB("postgres","123456","packform")
    // mainApp.Initialize(
    //     os.Getenv("APP_DB_USERNAME"),
    //     os.Getenv("APP_DB_PASSWORD"),
	//     os.Getenv("APP_DB_NAME"))
	mainApp.Initialize("postgres","123456","packform")
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func clearTable() {
	mainApp.DB.Exec("DELETE FROM orders")
	mainApp.DB.Exec("DELETE FROM order_items")
    mainApp.DB.Exec("DELETE FROM deliveries")
    customerCollection := mainApp.Mongo.Database("packform").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := customerCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}
	companyCollection := mainApp.Mongo.Database("packform").Collection("customercompanies")
	ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := companyCollection.Drop(ctx2); err != nil {
		log.Fatal(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    mainApp.Router.ServeHTTP(rr, req)

    return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestCustomers(t *testing.T) {
    req, _ := http.NewRequest("GET", "/customers", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestOrders(t *testing.T) {
    req, _ := http.NewRequest("GET", "/orders?start=0&count=5&sort=false", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCompanies(t *testing.T) {
    req, _ := http.NewRequest("GET", "/companies", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestTotalOrders(t *testing.T) {
    req, _ := http.NewRequest("GET", "/total", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func TestTotalAmount(t *testing.T) {
    req, _ := http.NewRequest("GET", "/totalamount", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}