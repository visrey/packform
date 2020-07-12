package main

import (
    "database/sql"
    "fmt"
	"log"
	"context"

    "net/http"
	"strconv"
	"time"
    "encoding/json"

    "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/handlers"
)

//App Struct 
type App struct {
    Router *mux.Router
	DB     *sql.DB
	Mongo  *mongo.Client
}

//Initialize Application DB
func (app *App) Initialize(user, password, dbname string) {
    connectionString :=
        fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

    var err error
    app.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	app.Mongo, _ = mongo.Connect(ctx, clientOptions)

    app.Router = mux.NewRouter()

    app.initializeRoutes()
}

//Run Start Server
func (app *App) Run(addr string) {
    log.Fatal(http.ListenAndServe(":8010", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(app.Router)))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (app *App) getOrders(w http.ResponseWriter, r *http.Request) {
	qcount, ok := r.URL.Query()["count"]
	if !ok || len(qcount[0]) < 1 {
		fmt.Println("missing counts parameter")
		return
	}
	qstart, ok := r.URL.Query()["start"]
	if !ok || len(qstart[0]) < 1 {
		fmt.Println("missing start parameter")
		return
	}
	qsort, ok := r.URL.Query()["sort"]
	if !ok || len(qsort[0]) < 1 {
		fmt.Println("missing sort parameter")
		return
    }
    search, ok := r.URL.Query()["search"]
	if !ok || len(search[0]) < 1 {
        fmt.Println("missing search parameter")
        search = append(search, "")
		// return
    }
    sdate, ok := r.URL.Query()["sdate"]
	if !ok || len(sdate[0]) < 1 {
        fmt.Println("missing sdate parameter")
        sdate = append(sdate, "")
		// return
    }
    edate, ok := r.URL.Query()["edate"]
	if !ok || len(edate[0]) < 1 {
        fmt.Println("missing edate parameter")
        edate = append(edate, "")
		// return
	}
    count, _ := strconv.Atoi(qcount[0])
	start, _ := strconv.Atoi(qstart[0])
    sort, _ := strconv.ParseBool(qsort[0])
    searchQuery := search[0]
    startDate := sdate[0]
    endDate := edate[0]

    if count > 5 || count < 1 {
        count = 5
    }
    if start < 0 {
        start = 0
	}

    orders, err := getOrders(app.DB, start, count, sort, searchQuery, startDate, endDate)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, orders)
}
func (app *App) getTotalAmount(w http.ResponseWriter, r *http.Request) {
    search, ok := r.URL.Query()["search"]
	if !ok || len(search[0]) < 1 {
        fmt.Println("missing search parameter")
        search = append(search, "")
		// return
    }
    sdate, ok := r.URL.Query()["sdate"]
	if !ok || len(sdate[0]) < 1 {
        fmt.Println("missing sdate parameter")
        sdate = append(sdate, "")
		// return
    }
    edate, ok := r.URL.Query()["edate"]
	if !ok || len(edate[0]) < 1 {
        fmt.Println("missing edate parameter")
        edate = append(edate, "")
		// return
	}
    searchQuery := search[0]
    startDate := sdate[0]
    endDate := edate[0]

    totalAmount, err := getTotalAmount(app.DB, searchQuery, startDate, endDate)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, totalAmount)
}

func (app *App) getTotal(w http.ResponseWriter, r *http.Request) {
    search, ok := r.URL.Query()["search"]
	if !ok || len(search[0]) < 1 {
        fmt.Println("missing search parameter")
        search = append(search, "")
		// return
    }
    sdate, ok := r.URL.Query()["sdate"]
	if !ok || len(sdate[0]) < 1 {
        fmt.Println("missing sdate parameter")
        sdate = append(sdate, "")
		// return
    }
    edate, ok := r.URL.Query()["edate"]
	if !ok || len(edate[0]) < 1 {
        fmt.Println("missing edate parameter")
        edate = append(edate, "")
		// return
	}
    searchQuery := search[0]
    startDate := sdate[0]
    endDate := edate[0]
    
    totalCount, err := getTotal(app.DB, searchQuery, startDate, endDate)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, totalCount)
}

func (app *App) getCustomers(w http.ResponseWriter, r *http.Request) {

    customers, err := getCustomers(app.Mongo)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, customers)
}

func (app *App) getCompanies(w http.ResponseWriter, r *http.Request) {

    companies, err := getCompanies(app.Mongo)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, companies)
}


func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/total", app.getTotal).Methods("GET")
    app.Router.HandleFunc("/orders", app.getOrders).Methods("GET")
    app.Router.HandleFunc("/totalamount", app.getTotalAmount).Methods("GET")
	app.Router.HandleFunc("/customers", app.getCustomers).Methods("GET")
	app.Router.HandleFunc("/companies", app.getCompanies).Methods("GET")
}