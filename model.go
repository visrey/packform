package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
)

const (
	ISO = "2006-01-02"
)

func getTotal(db *sql.DB, searchString string, startDate string, endDate string) (int, error) {
	var total int
	locale, _ := time.LoadLocation("Australia/Melbourne")
	sTimeStamp, _ := time.Parse(ISO, startDate)
	startDateF := sTimeStamp.In(locale).Format(ISO)
	eTimeStamp, _ := time.Parse(ISO, endDate)
	endDateF := eTimeStamp.In(locale).Format(ISO)

	query := `select count(*) as total from orders ord `

	nameQuery := ` ((ord.id in (
					select ord2.id from orders ord2 where ord2.order_name like '%` + searchString + `%' 
					)) or (ord.id in (
					select oi2.order_id from order_items oi2 where oi2.product like '%` + searchString + `%'
					)))
				`
	dateQuery := ` (ord.created > date('` + startDateF + `') and ord.created < date('` + endDateF + `'))`

	if (searchString != "") || (startDate != "" && endDate != "") {
		query += " where "
	}

	if searchString != "" {
		query += nameQuery
	}

	if startDate != "" && endDate != "" {
		if searchString != "" {
			query += " and "
		}
		query += dateQuery
	}

	row := db.QueryRow(query)
	err := row.Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return total, err
}

func getTotalAmount(db *sql.DB, searchString string, startDate string, endDate string) (float64, error) {
	var totalamount float64
	locale, _ := time.LoadLocation("Australia/Melbourne")
	sTimeStamp, _ := time.Parse(ISO, startDate)
	startDateF := sTimeStamp.In(locale).Format(ISO)
	eTimeStamp, _ := time.Parse(ISO, endDate)
	endDateF := eTimeStamp.In(locale).Format(ISO)

	query := `select coalesce(sum(ois.totalamount),0) as totalamount from orders ord inner join (
			select oi.order_id, sum(oi.price_per_unit * quantity) as totalamount from order_items oi group by oi.order_id
			) as ois on ord.id = ois.order_id
			left outer join (
			select oi.order_id, sum(del.delivered_quantity * oi.price_per_unit) as deliveredamount from 
			order_items oi, deliveries del where oi.id = del.order_item_id group by oi.order_id
			) as dis on ord.id = dis.order_id `

	nameQuery := ` ((ord.id in (
						select ord2.id from orders ord2 where ord2.order_name like '%` + searchString + `%' 
						)) or (ord.id in (
						select oi2.order_id from order_items oi2 where oi2.product like '%` + searchString + `%'
						)))
					`
	dateQuery := ` (ord.created > date('` + startDateF + `') and ord.created < date('` + endDateF + `'))`

	if (searchString != "") || (startDate != "" && endDate != "") {
		query += " where "
	}

	if searchString != "" {
		query += nameQuery
	}

	if startDate != "" && endDate != "" {
		if searchString != "" {
			query += " and "
		}
		query += dateQuery
	}

	row := db.QueryRow(query)

	err := row.Scan(&totalamount)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return totalamount, err
}

func getOrders(db *sql.DB, start, count int, sorted bool, searchString string, startDate string, endDate string) ([]Sorder, error) {

	sort := ""
	if sorted {
		sort = "asc"
	} else {
		sort = "desc"

	}
	locale, _ := time.LoadLocation("Australia/Melbourne")
	sTimeStamp, _ := time.Parse(ISO, startDate)
	startDateF := sTimeStamp.In(locale).Format(ISO)
	eTimeStamp, _ := time.Parse(ISO, endDate)
	endDateF := eTimeStamp.In(locale).Format(ISO)

	query := `select ord.id as id, ord.created as created, ord.order_name as order_name, ord.customer_id as customer_id, ois.totalamount as totalamount,coalesce(dis.deliveredamount, 0) as deliveredamount from orders ord inner join (
		select oi.order_id, sum(oi.price_per_unit * quantity) as totalamount from order_items oi group by oi.order_id
		) as ois on ord.id = ois.order_id
		left outer join (
		select oi.order_id, sum(del.delivered_quantity * oi.price_per_unit) as deliveredamount from 
		order_items oi, deliveries del where oi.id = del.order_item_id group by oi.order_id
		) as dis on ord.id = dis.order_id `

	nameQuery := ` ((ord.id in (
					select ord2.id from orders ord2 where ord2.order_name like '%` + searchString + `%' 
					)) or (ord.id in (
					select oi2.order_id from order_items oi2 where oi2.product like '%` + searchString + `%'
					)))
				`
	dateQuery := ` (ord.created > date('` + startDateF + `') and ord.created < date('` + endDateF + `'))`

	if (searchString != "") || (startDate != "" && endDate != "") {
		query += " where "
	}

	if searchString != "" {
		query += nameQuery
	}

	if startDate != "" && endDate != "" {
		if searchString != "" {
			query += " and "
		}
		query += dateQuery
	}

	query += " order by created " + sort + " LIMIT $1 OFFSET $2"

	rows, err := db.Query(query, count, start)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	sorders := []Sorder{}

	for rows.Next() {
		var sorder Sorder
		if err := rows.Scan(&sorder.ID, &sorder.Created, &sorder.OrderName, &sorder.CustomerID, &sorder.TotalAmount, &sorder.DeliveredAmount); err != nil {
			return nil, err
		}
		sorders = append(sorders, sorder)
	}

	return sorders, nil
}

func getCustomers(client *mongo.Client) ([]bson.M, error) {
	var customers []bson.M

	customersCollection := client.Database("packform").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := customersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var customer bson.M
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return customers, err
}

func getCompanies(client *mongo.Client) ([]bson.M, error) {
	var companies []bson.M
	companiesCollection := client.Database("packform").Collection("customercompanies")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := companiesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var company bson.M
		err = cursor.Decode(&company)
		companies = append(companies, company)
	}
	return companies, err
}
