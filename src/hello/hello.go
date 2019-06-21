package main

import (
	"fmt"
	gocql "github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func insert_rows(numberofrows int, session *gocql.Session) {

	
	imsi_n := 132312321

	for i := 0; i < numberofrows; i++ {

		imsi := strconv.Itoa(imsi_n + i)
		opc := "adaddddddddadaasdass" + imsi
		msisdn := "7777" + imsi

		//Add an etry to the table
		if err := session.Query("INSERT INTO  sm.sim_inventory (imsi, msisdn, opc) VALUES (?, ?, ?)",
			imsi, msisdn, opc).Exec(); err != nil {
			log.Fatalf("Could not Insert to cassandra : %v", err)
		}

	}

}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func fetch_rows(numberofrows int, session *gocql.Session) {

	defer timeTrack(time.Now(), "fetch_rows")

	var imsi string
	var msisdn string
	var opc string

	// Fetch multiple rows and run process over them
	iter := session.Query("SELECT imsi, msisdn,opc FROM sm.sim_inventory LIMIT ? ", numberofrows).Iter()
	for iter.Scan(&imsi, &msisdn, &opc) {
		log.Printf("Iter imsi: %v", imsi)
		log.Printf("Iter msisdn: %v", msisdn)
		log.Printf("Iter opc: %v", opc)
	}
	if err := iter.Close(); err != nil {
		log.Fatalf("Could not Query table : %v", err)
	}

}
func main() {
	fmt.Println("Hello, World -Test  Cassandra")

	cluster := gocql.NewCluster("mycassandra")
	cluster.Keyspace = "sm"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		log.Fatalf("Could not connect to cassandra cluster: %v", err)
	} else {
		log.Info("Successfully connected")
	}

	insert_rows(10, session)
	fetch_rows(10, session)

	fmt.Println("Connected closed")
}
