# packform


#### Dependent DB Services
1. PostgreSQL
2. MongoDB

#### Install GO Dependencies
From project root folder issue below command
```bash
go get ./...
```

## Data Loading
Data and schema will be dropped and re-created on every start of the server with the help of 5 csv files
```bash
orders.csv
order_items.csv
deliveries.csv
customers.csv
customer_companies.csv
```

#### Go Build and start
From project root folder issue below commnand to build
```bash
go build
```
This generates executable binary with the folder name in my case packform
run the binary to start the server at 8010
```bash
./packform
```

## End Points

### Companies
```bash
curl --location --request GET 'http://localhost:8010/companies'
```

>{
        "_id": "5f0b37f2573840eda146eca8",
        "companyid": "1",
        "companyname": "Roga & Kopyta"
 }

### Customers
```bash
curl --location --request GET 'http://localhost:8010/customers'
```

> {"_id":"5f0b37f2573840eda146eca6","company_id":"1","credit_cards":"[\"*****-1234\", \"*****-5678\"]","login":"ivan","name":"Ivan Ivanovich","password":"12345","userid":"ivan"}

### Total orders count
```bash
curl --location --request GET 'http://localhost:8010/total'
```

> 10 

with parameters
```bash
curl --location --request GET 'http://localhost:8010/total?search=4&sdate=2020-01-01&edate=2020-02-03'
```

### Total Amount
```bash
curl --location --request GET 'http://localhost:8010/totalamount'
```

> 27935.4241

with parameters
```bash
curl --location --request GET 'http://localhost:8010/totalamount?search=4&sdate=2020-01-01&edate=2020-02-03'
```

### Orders
```bash
curl --location --request GET 'http://localhost:8010/orders?start=0&count=5&sort=true'
```

> {"id":1,"created":"2020-01-02T15:34:12Z","order_name":"PO #001-I","totalamount":918.1220000000001,"deliveredamount":6.726999999999999,"customer_id":"ivan"}

With more parameters
```bash
curl --location --request GET 'http://localhost:8010/orders?search=4&sdate=2020-01-01&edate=2020-02-03&start=0&count=5&sort=true'
```
