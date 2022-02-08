package config

import "os"

func DSN() string {
	dsn := os.Getenv("DSN")
	if (dsn == "") {
		dsn = "postgres://postgres:adminadmin@localhost:54320/codeforces"
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
		peer = "peer2"
	}
	return peer
}

func PORT() string {
	port := os.Getenv("PORT")
	if (port == "") {
		port = ":8081"
	}
	return port
}