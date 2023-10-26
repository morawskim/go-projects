package main

import (
	"encoding/json"
	"flag"
	ics "github.com/arran4/golang-ical"
	"net/http"
	"time"
)

type event struct {
	Summary     string
	Description string
	StartAt     *time.Time
	EndAt       *time.Time
}

func main() {
	var url = flag.String("url", "", "The URL to iCal")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/calendar", func(w http.ResponseWriter, r *http.Request) {
		get, err := http.Get(*url)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))

			return
		}

		defer get.Body.Close()
		cal, err := ics.ParseCalendar(get.Body)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))

			return
		}

		events := make([]event, len(cal.Events()), len(cal.Events()))

		for i, val := range cal.Events() {
			var startAt, endAt *time.Time

			tmp, err := val.GetStartAt()
			startAt = timeOrNull(tmp, err)
			tmp, err = val.GetEndAt()
			endAt = timeOrNull(tmp, err)

			events[i] = event{
				Summary:     val.GetProperty(ics.ComponentPropertySummary).Value,
				Description: val.GetProperty(ics.ComponentPropertyDescription).Value,
				StartAt:     startAt,
				EndAt:       endAt,
			}
		}
		json.NewEncoder(w).Encode(events)
	})

	if err := http.ListenAndServe(":8000", mux); err != nil {
		panic(err)
	}
}

func timeOrNull(time time.Time, err error) *time.Time {
	if err == nil {
		return &time
	} else {
		return nil
	}
}
