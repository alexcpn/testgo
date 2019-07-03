package main

//alexcpn@gmail.com

import (
	"flag"
	"fmt"
	gocql "github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func createTable(session *gocql.Session) {

	/*
		createNamespaceStatement := `create keyspace sm with
		replication = {'class':'SimpleStrategy','replication_factor':1};`

		if err := session.Query(createNamespaceStatement).Exec(); err != nil {
			log.Warning("Could not Create Table : %v", err)

		}*/

	createTableStatement := `CREATE TABLE IF NOT EXISTS sm.sim_inventory (    
		imsi text PRIMARY KEY,
		msisdn text,
		opc text);`

	if err := session.Query(createTableStatement).Exec(); err != nil {
		log.Fatalf("Could not Create Table : %v", err)
		return
	}
	log.Debug("Created Table")

}
func insertRows(numberofrows int, session *gocql.Session) {

	defer timeTrack(time.Now(), "insertRows Time")
	imsiN := 132312321

	for i := 0; i < numberofrows; i++ {

		imsi := strconv.Itoa(imsiN + i)
		opc := "adaddddddddadaasdass" + imsi
		msisdn := "7777" + imsi

		//Add an etry to the table
		if err := session.Query("INSERT INTO  sm.sim_inventory (imsi, msisdn, opc) VALUES (?, ?, ?)",
			imsi, msisdn, opc).Exec(); err != nil {
			log.Fatalf("Could not Insert to cassandra : %v", err)
		}

	}

}

func fetchRows(numberofrows int, session *gocql.Session) {

	defer timeTrack(time.Now(), "fetchRows Time")

	// Fetch multiple rows and run process over them
	iter := session.Query("SELECT imsi, msisdn,opc FROM sm.sim_inventory LIMIT ? ", numberofrows).Iter()

	var imsi string
	var msisdn string
	var opc string

	counter := 0

	for iter.Scan(&imsi, &msisdn, &opc) {
		//log.Printf("Iter imsi: %v", imsi)
		//log.Printf("Iter msisdn: %v", msisdn)
		//log.Printf("Iter opc: %v", opc)
		counter++
	}
	log.Printf("Number of rows read  %v", counter)
	if err := iter.Close(); err != nil {
		log.Fatalf("Could not Query table : %v", err)
	}

}
func dropTable(session *gocql.Session) {

	removeTableStmt := fmt.Sprintf("DROP TABLE IF EXISTS %v", "sm.sim_inventory")
	if err := session.Query(removeTableStmt).Exec(); err != nil {
		log.Fatalf("Could not Drop Keyspace : %v", err)
	}
	log.Debug("Dropped Keyspace")

}

func alterTable(session *gocql.Session) {

	alterTableStatement := fmt.Sprintf(`ALTER TABLE sm.sim_inventory
	ADD  page_no  int ;`)
	if err := session.Query(alterTableStatement).Exec(); err != nil {
		log.Warningf("Could not alterTableStatement Table : %v", err)
	}
	log.Info("Altered Table")

}

func createIndex(session *gocql.Session) {

	createIndex := fmt.Sprintf(`CREATE INDEX pageindx ON sm.sim_inventory(page_no)`)
	if err := session.Query(createIndex).Exec(); err != nil {
		log.Warningf("Could not alterTableStatement Table : %v", err)
	}
	log.Info("Created Index")

}

func updateRows(numberofrows int, session *gocql.Session) {

	defer timeTrack(time.Now(), "updateRows Time")

	// Fetch multiple rows and run process over them
	iter := session.Query("SELECT imsi FROM sm.sim_inventory LIMIT ? ", numberofrows).Iter()

	var imsi string

	counter := 0
	pageNum := 1

	for iter.Scan(&imsi) {
		counter++
		if counter % 1000 == 0 {
			pageNum++
		}
		if err := session.Query("UPDATE sm.sim_inventory SET page_no = ? where imsi = ?",
			pageNum, imsi).Exec(); err != nil {
			log.Fatalf("Could not Update page_no : %v", err)
		}

	}
	log.Printf("Number of rows read  %v", counter)
	if err := iter.Close(); err != nil {
		log.Fatalf("Could not Query table : %v", err)
	}

}

func fetchRowsChunks(pageno int, session *gocql.Session) {

	defer timeTrack(time.Now(), "fetchRows By Page Time")

	// Fetch multiple rows and run process over them
	iter := session.Query("SELECT imsi, msisdn,opc,page_no FROM sm.sim_inventory where page_no = ? ", pageno).Iter()

	var imsi string
	var msisdn string
	var opc string
	var page int

	counter := 0

	for iter.Scan(&imsi, &msisdn, &opc, &page) {
		log.Printf("Iter imsi: %v %v %v   Page=%v", imsi,msisdn,opc,page)
	
	}
	log.Printf("Number of rows read  %v", counter)
	if err := iter.Close(); err != nil {
		log.Fatalf("Could not Query table : %v", err)
	}

}

func main() {
	//coding/testgo# go install hello && ./bin/hello -ni 10 -nq 1000000 -del=false -pno=10
	fmt.Println("Hello, World -Test  Cassandra Pagination")
	numberOfRowsToInsert := flag.Int("ni", 1, "Number of rows to insert")
	numberOfRowsToQuery := flag.Int("nq", 1, "Number of rows to Query")
	deleteTable := flag.Bool("del", false, "Delete table after test")
	pageNumberToGet := flag.Int("pno", 1, "Page to get")

	flag.Parse()
	log.Printf("Number of Rows to Insert %v", *numberOfRowsToInsert)
	log.Printf("Number of Rows to Query %v", *numberOfRowsToQuery)
	log.Printf("Delete Table after tests %v", *deleteTable)
	log.Printf("Page Number to read %v", *pageNumberToGet)


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

	createTable(session)
	insertRows(*numberOfRowsToInsert, session)
	fetchRows(*numberOfRowsToQuery, session)
	//Do the below only once - if column is not created
	//alterTable(session)
	//createIndex(session)
	//updateRows(*numberOfRowsToQuery, session)
	fetchRowsChunks(*pageNumberToGet, session)
	
	if *deleteTable {
		log.Info("Going to Drop the table")
		dropTable(session)
	}

	fmt.Println("Connected closed")
}
