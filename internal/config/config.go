package config

import (
	"os"
)

func DSN() string {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://postgres:adminadmin@localhost:54320/codeforces"
	}
	return dsn
}

func KafkaConn() string {
	kafka_conn := os.Getenv("KAFKA_CONN")
	if kafka_conn == "" {
		kafka_conn = "localhost:29092"
	}
	return kafka_conn
}

func PeerName() string {
	peer := os.Getenv("PEER")
	if peer == "" {
		peer = "peer2"
	}
	return peer
}

func ServePort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return port
}

func TelebotToken() string {
	token := os.Getenv("TELEBOT_TOKEN")
	if token == "" {
		token = "1752339795:AAFomsaP4I3hr2Xh6QYi3s09Yyps6nEGGM4"
	}
	return token
}

func TelebotChannelName() string {
	channelName := os.Getenv("TELEBOT_CHANNEL_NAME")
	if channelName == "" {
		channelName = "@CompetitiveProgrammingPlatform"
	}
	return channelName
}
