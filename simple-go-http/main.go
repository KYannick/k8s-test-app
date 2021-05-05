package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type podinfo struct {
	POD_NAME string ""
	POD_IP   string ""
}

var alive, ready bool = true, true
var info podinfo

func responseHandler(w http.ResponseWriter, status int, message interface{}) {
	logrus.Infof("Sending response. Status: \"%v\", Message: \"%+v\"", status, message)

	success := true
	if status/100 != 2 {
		success = false
	}

	data := struct {
		Success bool
		Status  int
		Message interface{}
	}{
		success,
		status,
		message,
	}

	response, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func app_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		responseHandler(w, http.StatusMethodNotAllowed, fmt.Sprintf("Received call using forbidden method (only GET is allowed): %v", r.Method))
		return
	}

	responsestring := "Hello from " + info.POD_NAME + ". I have ip " + info.POD_IP + ""

	responseHandler(w, http.StatusOK, responsestring)
}

func liveness_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		responseHandler(w, http.StatusMethodNotAllowed, fmt.Sprintf("Received call using forbidden method (only GET is allowed): %v", r.Method))
		return
	}

	if alive {
		responseHandler(w, http.StatusOK, "This app is alive")
		return
	} else {
		responseHandler(w, http.StatusInternalServerError, "This app is dead")
		return
	}
}

func readiness_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		responseHandler(w, http.StatusMethodNotAllowed, fmt.Sprintf("Received call using forbidden method (only GET is allowed): %v", r.Method))
		return
	}

	if ready {
		responseHandler(w, http.StatusOK, "This app is ready")
		return
	} else {
		responseHandler(w, http.StatusServiceUnavailable, "This app is not ready for traffic")
		return
	}
}

func set_probe_handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		responseHandler(w, http.StatusMethodNotAllowed, fmt.Sprintf("Received call using forbidden method (only POST is allowed): %v", r.Method))
		return
	}

	q := r.URL.Query()

	if q.Get("flag") == "alive" {
		alive = !alive
		responseHandler(w, http.StatusOK, fmt.Sprintf("alive bool toggled, new value is %v", alive))
		return
	} else if q.Get("flag") == "ready" {
		ready = !ready
		responseHandler(w, http.StatusOK, fmt.Sprintf("ready bool toggled, new value is %v", ready))
		return
	} else {
		responseHandler(w, http.StatusBadRequest, "You are asking the impossible!")
		return
	}

}

func main() {
	info.POD_IP = os.Getenv("POD_IP")
	info.POD_NAME = os.Getenv("POD_NAME")

	http.HandleFunc("/", app_handler)
	http.HandleFunc("/alive", liveness_handler)
	http.HandleFunc("/ready", readiness_handler)
	http.HandleFunc("/set-probe", set_probe_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
