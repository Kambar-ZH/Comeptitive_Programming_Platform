package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

var Logger, ErrorLogger *zap.Logger

func init() {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	var err error
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
