module cassandratest

go 1.12

require (
	github.com/gocql/gocql v0.0.0-20190629212933-1335d3dd7fe2
	github.com/sirupsen/logrus v1.4.2
	internal/util v0.0.0
)

replace internal/util => ../util
