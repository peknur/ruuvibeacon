package publishers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/peknur/ruuvibeacon"
)

func logPublisher(ctx context.Context, e *ruuvibeacon.Envelope) {
	js, err := json.Marshal(&e)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println(string(js))
}

func init() {
	Add("log", logPublisher)
}
