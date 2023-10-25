package main

import (
	"flag"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"net/http"
)

func main() {
	var url = flag.String("url", "", "The URL to iCal")
	flag.Parse()

	get, err := http.Get(*url)
	if err != nil {
		panic(err)
	}

	defer get.Body.Close()

	cal, err := ics.ParseCalendar(get.Body)

	if err != nil {
		panic(err)
	}

	for _, val := range cal.Events() {
		fmt.Println(val.GetProperty(ics.ComponentPropertySummary).Value)
		fmt.Println(val.GetProperty(ics.ComponentPropertyDescription).Value)

		time, _ := val.GetStartAt()
		fmt.Println(time.Local())
	}
}
