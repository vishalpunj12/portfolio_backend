package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type ContactRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Message   string `json:"message"`
}

func main() {
	http.HandleFunc("/send", handleSend)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or specific origin
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Handle preflight (OPTIONS) request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("Request received!!!")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ContactRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = sendEmail(req)
	if err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Request completed !!!")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func sendEmail(data ContactRequest) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASS")
	to := os.Getenv("EMAIL_TO")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: New Contact Message received from portfolio site\r\n" +
			"\r\n" +
			"Name: " + data.FirstName + " " + data.LastName + "\n" +
			"Email: " + data.Email + "\n" +
			"Phone: " + data.Phone + "\n" +
			"Message: " + data.Message + "\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}
