package main

import (
   // "unicode"
	"fmt"
	"log"
	"net/http"
	//"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	//"github.com/satori/go.uuid"
	//"strconv"
	"time"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// MongoDB Config
var mongodb_server = "admin:cmpe281@ds115664.mlab.com:15664"
var mongodb_database = "starbucks"
var mongodb_collection = "starbucks"
var MongoDBHosts = "ds115664.mlab.com:15664"
var AuthDatabase = "starbucks"
var AuthUserName = "admin"
var AuthPassword = "cmpe281"

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}


//var _init_ctx sync.Once
var _instance *DB

/*func New() *mgo.Database {
    _init_ctx.Do(func() {
        _instance = new(DB)

        
        // Create a session which maintains a pool of socket connections
        // to our MongoDB.
        

        if err != nil {
            fmt.Printf("Error en mongo: %+v\n", err)
            os.Exit(1)
        }
      _instance.Database = session.DB(AuthDatabase)  
    })
    return _instance.Database
}
*/

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
    mx.HandleFunc("/payment/{id}", starbucksPaymentHandler(formatter)).Methods("POST")
	mx.HandleFunc("/payment/{userid}/{orderid}/{amt}/{orderstatus}", starbucksNewOrderHandler(formatter)).Methods("POST")
//	mx.HandleFunc("/paymentInfo/{id}/{cardnumber}/{Date}/{CVV}", starbucksPaymentValidateHandler(formatter)).Methods("GET")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Starbucks version 1.0 alive!"})
	}
}

func starbucksPaymentHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println( "Order ID: ", uuid )

		// Open MongoDB Session
		// Open MongoDB Session
		mongoDBDialInfo := &mgo.DialInfo{
            Addrs:    []string{MongoDBHosts},
            Timeout:  600 * time.Second,
            Database: AuthDatabase,
            Username: AuthUserName,
            Password: AuthPassword,
        }

		session, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                panic(err)
        }
        c := session.DB(AuthDatabase).C(mongodb_collection)

      //Get Starbucks Order
        var result starbucks
        //err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
        err = c.Find(bson.M{"OrderId" : uuid}).One(&result)
        if err != nil {
                log.Fatal(err)
                fmt.Println( "Order not found: ", uuid)
        } else{
			fmt.Println( "Order found: ", uuid)
			// Update Starbucks Order Status
			query := bson.M{"OrderId" : uuid}
		    change := bson.M{"$set": bson.M{"OrderStatus" : "Order Paid"}}
		    err = c.Update(query, change)
		    if err != nil {
		                log.Fatal(err)
		    } else {
		    	formatter.JSON(w, http.StatusOK, result)
		    }

		}
		
	}
}

/*

// Create Order without Mognodb
func starbucksNewOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uuid, _ := uuid.NewV4()
		var ord = order{
			Id:          uuid.String(),
			OrderStatus: "Order Placed",
		}
		if orders == nil {
			orders = make(map[string]order)
		}
		orders[uuid.String()] = ord
		fmt.Println("Orders: ", orders)
		formatter.JSON(w, http.StatusOK, ord)
		formatter.JSON(w, http.StatusOK, struct{ Payment string }{"Order is ready for Payment"})
	}
}
*/

// Create Order with Mongodb
func starbucksNewOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)

		var uuid string = params["userid"]
		fmt.Println( "User ID: ", uuid )

		var orderid string = params["orderid"]
		fmt.Println( "Order ID: ", orderid )

		var orderstatus string = params["orderstatus"]
		fmt.Println( "Order Status: ", orderstatus )

		var amt string = params["amt"]
		fmt.Println("Order amount:", amt)


		// Open MongoDB Session
		mongoDBDialInfo := &mgo.DialInfo{
            Addrs:    []string{MongoDBHosts},
            Timeout:  600 * time.Second,
            Database: AuthDatabase,
            Username: AuthUserName,
            Password: AuthPassword,
        }

		session, err := mgo.DialWithInfo(mongoDBDialInfo)
        if err != nil {
                panic(err)
        }
        c := session.DB(mongodb_database).C(mongodb_collection)
        err = c.Insert(&starbucks{UserId : uuid, OrderId : orderid, OrderAmount: amt, OrderStatus : orderstatus})
          if err != nil {
                log.Fatal(err)
        } else {

        	fmt.Println("Successfull Order: ", orderid)
        	formatter.JSON(w, http.StatusOK, orderid)
        	formatter.JSON(w, http.StatusOK, struct{ Payment string }{"Order is ready for Payment"})
        }
    }
}
		

/*

// API  Order Payment
func starbucksPaymentHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println("Order ID: ", uuid)
		for key, value := range orders {
			//fmt.Println("Order ID:", uuid, "Value:", "Order Paid")
			fmt.Println("Key:", key, "Value:", value)
			var ord = orders[uuid]
			ord.OrderStatus = "Order Paid"
			orders[uuid] = ord
			fmt.Println("Order: ", ord)
			formatter.JSON(w, http.StatusOK, ord)
		}
		//formatter.JSON(w, http.StatusOK, "Orders Paid!")						
	}
}
*/

/*
// API Process Orders
func starbucksPaymentValidateHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		///paymentInfo/{id}/{cardnumber}/{Date}/{CVV}
		params := mux.Vars(req)
		var good = 0
		var uuid string = params["id"]
		cardNumber := params["cardnumber"]
		x :=  params["Date"]
		Date := time.New()
		var Date time.time = date
		/*if unicode.IsDigit(params["CVV"]){
			fmt.Println("Invalid CVV")
		}
		cvv_str := params["CVV"]
		CVV, _ := strconv.Atoi(cvv_str)
		b := make([]int, CVV)
		fmt.Println("Order ID: ", uuid)
		fmt.Println("Card Number: ", cardNumber)
		fmt.Println("Date: ", Date)
		fmt.Println("CVV: ", b)
		fmt.Println("length of b ", len(b))
		if len(b) == 3 {
         	good = 1
		}
		if good ==1 {
			for key, value := range orders {
				fmt.Println("Key:", key, "Value:", value)
				var ord = orders[uuid]
				ord.OrderStatus = "Payment Information Validated"
				orders[uuid] = ord
				fmt.Println("Order: ", ord)
				formatter.JSON(w, http.StatusOK, ord)
			}
		} else {
			formatter.JSON(w, http.StatusOK, struct{ Payment string }{"Invalid Payment Information"})
		}			
	}
}
*/