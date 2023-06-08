package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/go-plugins-helpers/sdk"
	"io"
	"net/http"
	"os"
)

func main() {
	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)

	h.HandleFunc("/LogDriver.StartLogging", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		//var req StartLoggingRequest
		fmt.Fprintf(os.Stdout, "Start logging request was called for the container : %s", body)
		//json.Unmarshal(body, &req)
		//err := d.StartLogging(req.File, req.Info)
		sendResponse(w)
	})

	h.HandleFunc("/LogDriver.StopLogging", func(w http.ResponseWriter, r *http.Request) {
		//var req StopLoggingRequest
		//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}
		fmt.Fprintln(os.Stdout, "Stop logging request was called for the container")
		//err := d.StopLogging(req.File)
		sendResponse(w)
	})

	h.HandleFunc("/LogDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(os.Stdout, "Capabilities endpoint was called")
		json.NewEncoder(w).Encode(&CapabilitiesResponse{
			Cap: logger.Capability{ReadLogs: false},
		})
	})

	h.HandleFunc("/LogDriver.ReadLogs", func(w http.ResponseWriter, r *http.Request) {
		//var req ReadLogsRequest
		//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}

		//fmt.Fprintln(os.Stdout, "docker logs was called for the container : ", req.Info.ContainerID)
		fmt.Fprintln(os.Stdout, "docker logs was called for the container : ")
		//http.Error(w, "Not implemented", http.StatusNotImplemented)
	})

	if err := h.ServeUnix("customlogdriver", 0); err != nil {
		panic(err)
	}
}

func sendResponse(w http.ResponseWriter) {
	var res response
	json.NewEncoder(w).Encode(&res)
}
