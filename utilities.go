package main

import (
	"errors"
	"log"

	"github.com/alexjlockwood/gcm"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func sendMessageToGCM(tokens []string, payload map[string]interface{}) (bool, error) {
	// At any exit, decrement pending
	defer func() {
		go decrementPending()
	}()

	if len(tokens) == 0 {
		errText := "No tokens were supplied, exiting"
		log.Println(errText)
		return false, errors.New(errText)
	}

	// All is well, make & send the message
	go appendAttempts(len(tokens))

	msg := gcm.NewMessage(payload, tokens...)
	sender := &gcm.Sender{ApiKey: settings.GCMAPIKey}
	response, err := sender.Send(msg, 2)
	if err != nil {
		log.Println("Failed to send message:")
		log.Println(err.Error())

		go appendFailures(1)
		return false, err
	}

	numCan := 0
	numErr := 0
	if response != nil {
		for i, result := range response.Results {
			// Canonicals
			if result.RegistrationID != "" {
				numCan++
				canonicalReplacements = append(canonicalReplacements, canonicalReplacement{tokens[i], result.RegistrationID})
			}
			if result.Error != "" {
				numErr++
				log.Printf("Error sending: %s", result.Error)

				if result.Error == "NotRegistered" {
					handleNotRegisteredError(tokens[i])
					go appendNotRegistered(1)
				}
			}
		}

		go appendCanonicals(numCan)
		go appendFailures(numErr)
	}

	log.Printf("Message sent. Attempts: %d, Errors: %d, Successful: %d (Canonicals: %d)", len(tokens), numErr, len(tokens)-numErr, numCan)

	return true, nil
}

func sendMessageToAPNS(tokens []string, payloadAsString []byte) (bool, error) {
	cert, pemErr := certificate.FromPemFile("conf/APNS2.pem", "")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	log.Println("token:", tokens[0])
	notification.DeviceToken = tokens[0] // "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.life360.enterpriseqa"
	notification.Payload = payloadAsString

	//client := apns.NewClient(cert).Development()
	client := apns.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return false, err
	}
	log.Println("Res:", res)

	return true, nil
}


func handleCanonicalsInResult(original string, results []gcm.Result) {
	for _, r := range results {
		canonicalReplacements = append(canonicalReplacements, canonicalReplacement{original, r.RegistrationID})
	}
}

func handleNotRegisteredError(original string) {
	notRegisteredMutex.Lock()
	notRegisteredKeys = append(notRegisteredKeys, original)
	notRegisteredMutex.Unlock()
}

func appendAttempts(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Attempts += numToAppend
}

func appendFailures(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Failures += numToAppend
}

func appendCanonicals(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Canonicals += numToAppend
}

func appendNotRegistered(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.NotRegistered += numToAppend
}

func incrementPending() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Pending++
}

func decrementPending() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Pending--
}
