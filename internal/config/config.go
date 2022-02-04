package config

import "os"

func DSN() string {
	dsn := os.Getenv("DSN")
	if (dsn == "") {
		dsn = "postgres://postgres:adminadmin@localhost:5432/codeforces"
	}
	return dsn
}

func KAFKA_CONN() string {
	kafka_conn := os.Getenv("KAFKA_CONN")
	if (kafka_conn == "") {
		kafka_conn = "localhost:29092"
	}
	return kafka_conn
}

func PEER() string {
	peer := os.Getenv("PEER")
	if (peer == "") {
		peer = "peer1"
	}
	return peer
}