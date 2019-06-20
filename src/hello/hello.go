package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gocql/gocql"
)

func main() {
	fmt.Println("Hello, World -Test  Cassandra")

	cluster := gocql.NewCluster("mycassandra")
	cluster.Keyspace = "sm"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
        log.Fatalf("Could not connect to cassandra cluster: %v", err)
    }else{
		log.Info("Successfully connected")
	}
	var imsi string
	var msisdn string
	var opc string
	
	//Add an etry to the table

	if err := session.Query("INSERT INTO  sm.sim_inventory (imsi, msisdn, opc) VALUES (?, ?, ?)",
	"132312321", "7777777", "adsdasdas1").Exec(); err != nil {
		log.Fatalf("Could not Insert to cassandra : %v", err)
	}

	// Fetch multiple rows and run process over them
    iter := session.Query("SELECT imsi, msisdn,opc FROM sm.sim_inventory").Iter()
    for iter.Scan(&imsi, &msisdn, &opc) {
        log.Printf("Iter imsi: %v", imsi)
		log.Printf("Iter msisdn: %v", msisdn)
		log.Printf("Iter opc: %v", opc)
	}
	if err := iter.Close(); err != nil {
		log.Fatalf("Could not Query table : %v", err)
	}
	
 
	
	fmt.Println("Connected closed")
}
