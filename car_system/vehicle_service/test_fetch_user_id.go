package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	client := &http.Client{}

	// Replace <actual-session-value> with the actual session ID obtained after login
	sessionCookie := "user-session=<actual-session-value>"

	req, err := http.NewRequest("GET", "http://localhost:8080/get-session-user-id", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set the session cookie
	req.Header.Set("Cookie", sessionCookie)

	// Log request headers
	fmt.Printf("Request Headers: %v\n", req.Header)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Error: Unauthorized. Check if the session is valid.")
		return
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return
	}

	// Decode the response body
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	fmt.Printf("Response: %v\n", response)
}
