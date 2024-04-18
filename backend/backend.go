package backend

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

type WoofEvent struct {
	EventType   int     `json:"event_type"`
	Lat         float32 `json:"lat"`
	Lng         float32 `json:"lng"`
	TimestampMs uint64  `json:"timestamp_ms"`
}

type WoofListResp struct {
	Events []WoofEvent `json:"events"`
}

func Run() {
	certFile := flag.String("cert", "", "certificate file")
	keyFile := flag.String("key", "", "key file")
	flag.Parse()

	events := make([]WoofEvent, 0)

	http.HandleFunc("/woof", func(w http.ResponseWriter, r *http.Request) {
		var event WoofEvent
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			log.Println("Error decoding JSON", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Processing event request %+v\n", event)
		events = append(events, event)
	})

	http.HandleFunc("/woof_list", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(WoofListResp{
			Events: events,
		})
	})

	http.ListenAndServeTLS(":443", *certFile, *keyFile, http.DefaultServeMux)
}
