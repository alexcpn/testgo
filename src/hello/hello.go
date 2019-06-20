package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gocql/gocql"
)

func main() {
	fmt.Println("Hello, World -Test cassandra")

	cluster := gocql.NewCluster("mycassandra")
	cluster.Keyspace = "sm"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
        log.Fatalf("Could not connect to cassandra cluster: %v", err)
    }else{
		log.Info("Successfully connected")
	}
	
	defer session.Close()
	fmt.Println("Connected closed")
}
