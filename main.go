package main

import (
	"log"
	"net/http"
	"os"

	"./src/app"
	"github.com/julienschmidt/httprouter"

	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

func main() {
	// get env parameters
	dbPath := os.Getenv("LOCATION_DB_PATH")
	if dbPath == "" {
		log.Fatal("Not found LOCATION_DB_PATH params")
	}
	
	uaParserPath := os.Getenv("UA_PARSER_REGEXP_PATH")
	if uaParserPath == "" {
		log.Fatal("Not found UA_PARSER_REGEXP_PATH params")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	// init DBs and server
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	parser, err := uaparser.New(uaParserPath)
	if err != nil {
		log.Fatal(err)
	}

	br := &app.BidRequest{
		LocationDB: db,
		UAParser:   parser,
	}

	router := httprouter.New()
	router.POST("/api/bidrequest/v1/", br.HTTPHandle)

	log.Fatal(http.ListenAndServe(port, router))
}
