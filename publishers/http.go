package publishers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/peknur/ruuvibeacon"
)

func httpPublisher(ctx context.Context, e *ruuvibeacon.Envelope) {
	uri := os.Getenv("APP_PUBLISHER_HTTP_URI")
	if uri == "" {
		log.Print("httpPublisher: env variable 'APP_PUBLISHER_HTTP_URI' not set")
		return
	}
	js, err := json.Marshal(&e)
	if err != nil {
		log.Print(err)
		return
	}
	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewReader(js))
	if err != nil {
		log.Print(err)
		return
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
}

func init() {
	Add("http", httpPublisher)
}
