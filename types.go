package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//OrderItem - Holding Order Items
type OrderItem struct {
	ID			int			`json:"id"`
	orderID 	int			`json:"order_id"`
    priceUnit	float64		`json:"price_per_unit"`
	quantity	int			`json:"quantity"`
	product		string		`json:"product"`
}

//Order - Holding the Order info
type Order struct {
	ID			int			`json:"id"`
	created 	string		`json:"created"`
    orderName  	string		`json:"order_name"`
    customerID 	string		`json:"customer_id"`
}

//Deliveries - Holding the delivery info
type Deliveries struct {
	ID					int		`json:"id"`
	orderItemID 		int		`json:"order_item_id"`
    deliveredQuantity  	int		`json:"delivered_quantity"`
}

//Customer - Holding the customer info
type Customer struct {
	userID			string				`json:"userid,omitempty" bson:"userid,omitempty"`
	login			string				`json:"login,omitempty" bson:"login,omitempty"`
	password		string				`json:"password,omitempty" bson:"password,omitempty"`
	name			string				`json:"name,omitempty" bson:"name,omitempty"`
	companyID  		string				`json:"company_id,omitempty" bson:"company_id,omitempty"`
	creditCards		string				`json:"credit_cards,omitempty" bson:"credit_cards,omitempty"`
}

//CustomerCompany - Holding customer company info
type CustomerCompany struct {
	ID					primitive.ObjectID 	`json:"_id" bson:"_id,omitempty"`
	companyID			string					`json:"companyid" bson:"companyid,omitempty"`
	companyName			string				`json:"companyname" bson:"companyname,omitempty"`
}

//Sorder - Hybrid Order Struct to store query output
type Sorder struct {
	ID				int			`json:"id"`
	Created 		string		`json:"created"`
    OrderName  		string		`json:"order_name"`
	TotalAmount		float64		`json:"totalamount"`
	DeliveredAmount	float64		`json:"deliveredamount"`
	CustomerID		string		`json:"customer_id"`
}