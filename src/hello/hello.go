package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

func main() {
	fmt.Println("Hello, World -Test cassandra")

	cluster := gocql.NewCluster("mycassandra")
	cluster.Keyspace = "sm"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()
	fmt.Println("Connected closed")
}
