package mdb

//MongoDB business logic

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type StockData struct {
	Symbol     string        `json:"symbol"`
	Price      float32       `json:"Price"`
	Quantity   int           `json:"Quantity"`
	InsertDate string        `json:"InsertDate"`
	Operation  string        `json:"Operation"`
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

type StockDataTables struct {
	Data [][]string `json:"data"`
}

//MongoDBConn connection details for the DB
type MongoDBConn struct {
	host string
	db   string
	coll string
}

func CreateMongoDBConn(host string, db string, coll string) *MongoDBConn {
	return &MongoDBConn{host, db, coll}
}

//Stock returns a slice of stocks for the given symbol from MongoDB
func (mconn *MongoDBConn) Stock(symbol string, data *[]StockData) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	q := c.Find(bson.M{"symbol": symbol})
	if q == nil {
		log.Fatal(err)
	}

	err = q.All(data)
	if err != nil {
		log.Fatal(err)
	}
}

//StockDataTables returns a slice of stocks for the given symbol from MongoDB
func (mconn *MongoDBConn) StockDataTables(symbol string, data *StockDataTables) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	var sf = bson.M{}
	if symbol != "" {
		sf = bson.M{"symbol": symbol}
	}
	q := c.Find(sf)
	if q == nil {
		log.Fatal(err)
	}
	d := []StockData{}
	err = q.All(&d)
	if err != nil {
		log.Fatal(err)
	}
	dlen := len(d)
	data.Data = make([][]string, 0, dlen-1)
	for i := 0; i < dlen; i++ {
		cd := []string{d[i].Symbol, fmt.Sprintf("%f", d[i].Price), fmt.Sprintf("%d", d[i].Quantity), d[i].InsertDate, d[i].Operation, fmt.Sprintf("<button id='%s' class='delstock btn-danger btn-xs'>Delete</button>", d[i].Id.Hex())}
		//cd := coldata{d[i].Symbol, d[i].Price, d[i].Quantity, d[i].InsertDate, d[i].Operation, fmt.Sprintf("<button id='%s' class='delstock btn-danger btn-xs'>Delete</button>", d[i].Id.Hex())}
		data.Data = append(data.Data, cd)
		//data.Data = append(data.Data, fmt.Sprintf("%s,%f,%d,%s,%s,<button id='%s' class='delstock btn-danger btn-xs'>Delete</button>", d[i].Symbol, d[i].Price, d[i].Quantity, d[i].InsertDate, d[i].Operation, d[i].Id.Hex()))
	}

}

//Stock returns a slice of stocks for the given symbol from MongoDB
func (mconn *MongoDBConn) StockUniqueSymbols(data *[]string) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	q := c.Find(nil)
	if q == nil {
		log.Fatal(err)
	}
	var n int
	n, err = q.Count()
	if err != nil {
		log.Fatal("Failed to count result", err)
	}

	log.Println("Count = ", n)
	err = q.Distinct("symbol", data)
	if err != nil {
		log.Fatal(err)
	}
}

//StockByID returns a slice of stocks for the given Object ID from MongoDB
func (mconn *MongoDBConn) StockByID(id string, data *[]StockData) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Failed to connect to DB:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	q := c.Find(bson.M{"_id": bson.ObjectIdHex(id)})
	if q == nil {
		log.Println("Failed to find by ID:", err)
		return err
	}

	err = q.All(data)
	if err != nil {
		log.Println("Failed to parse result:", err)
		return err
	}
	return nil
}

//RemoveById removes a document by ID from MongoDB
func (mconn *MongoDBConn) RemoveByID(id string) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Failed to connec to DB:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Failed to delete:", err)
		return err
	}
	return nil
}

//Place adds a stock to DB
func (mconn *MongoDBConn) Place(symbol string, price float32, quantity int, oper string, inDate time.Time) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	// sqr := StockQueryResult{}
	// err = YahooStockData("GSK.L", &sqr)
	// if err != nil {
	// 	log.Fatal("test failed")
	// }
	data := StockData{
		Symbol:     symbol,
		Price:      price,
		Quantity:   quantity,
		Operation:  oper,
		InsertDate: "",
	}
	err = c.Insert(&data)
	if err != nil {
		log.Fatal(err)
	}
}

//Place2 adds a stock to DB
func (mconn *MongoDBConn) Place2(sd StockData) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Error Place2, DB connection:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	//sd.InsertDate = ""
	err = c.Insert(&sd)
	if err != nil {
		log.Println("Error Place2, insert operation:", err)
		return err
	}
	return nil
}
