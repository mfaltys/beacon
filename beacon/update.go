package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v3"
)

func update(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	r.ParseForm()
	clientId := strings.TrimSpace(r.FormValue("id"))
	clientSec := strings.TrimSpace(r.FormValue("sec"))
	clientAddress := strings.TrimSpace(r.FormValue("address"))

	// sha3:512 hash the id and sec
	clientIdHash := sha3.Sum512([]byte(clientId))
	clientSecHash := sha3.Sum512([]byte(clientSec))

	// check if id exists
	storedSecHash, err := client.Get(fmt.Sprintf("sec:%x", clientIdHash)).Result()
	if err != redis.Nil {
		// id exists, make sure clientSecHash is the same as the stored version
		if fmt.Sprintf("%x", clientSecHash) == storedSecHash {
			// client is authed, update clientAddress
			w.WriteHeader(http.StatusOK)
			client.Set(fmt.Sprintf("ip:%x", clientIdHash), clientAddress, 0).Err()
		} else {
			// client auth failed
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		// id does not exist
	}
}
