package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"log"
)

// Send a message to GCM or APNS
func send(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokens := r.PostForm["tokens"]
	payloadAsString := r.PostFormValue("payload")

	if payloadAsString == "" {
		errText := "Payload was empty, exiting"
		log.Println(errText)
		//return false, errors.New(errText)
	}

	// Unpack the JSON payload
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(payloadAsString), &payload)
	if err != nil {
		log.Println("Can't unmarshal the json: " + err.Error())
		log.Println("Original: " + payloadAsString)
		//return false, err
	}

	log.Println("Original: " + payloadAsString)
	if payload["aps"] != nil {
		log.Println("Sending apns")
		go func() {
			incrementPending()
			byteArray := []byte(payloadAsString)
			sendMessageToAPNS(tokens, byteArray)
		}()
	} else {
		log.Println("Sending gcm")
		go func() {
			incrementPending()
			sendMessageToGCM(tokens, payload)
		}()
	}

	// Return immediately
	output := "ok\n"
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(len(output)))
	io.WriteString(w, output)
}

// Return a run report for this process
func getReport(w http.ResponseWriter, r *http.Request) {
	runReportMutex.Lock()
	a, _ := json.Marshal(runReport)
	runReportMutex.Unlock()
	b := string(a)
	io.WriteString(w, b)
}

// Return all currently collected canonical reports from GCM
func getCanonicalReport(w http.ResponseWriter, r *http.Request) {
	canonicalReplacementsMutex.Lock()
	ids := map[string][]canonicalReplacement{"canonical_replacements": canonicalReplacements}
	a, _ := json.Marshal(ids)
	canonicalReplacementsMutex.Unlock()

	b := string(a)
	io.WriteString(w, b)

	// Clear out canonicals
	go func() {
		canonicalReplacementsMutex.Lock()
		defer canonicalReplacementsMutex.Unlock()
		canonicalReplacements = nil
	}()
}

// Return all tokens that need to be unregistered
func getNotRegisteredReport(w http.ResponseWriter, r *http.Request) {
	notRegisteredMutex.Lock()
	ids := map[string][]string{"tokens": notRegisteredKeys}
	a, _ := json.Marshal(ids)
	notRegisteredMutex.Unlock()

	b := string(a)
	io.WriteString(w, b)

	// Clear ids
	go func() {
		notRegisteredMutex.Lock()
		defer notRegisteredMutex.Unlock()
		notRegisteredKeys = nil
	}()
}
