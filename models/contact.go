package models

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

const contactLogPath = "./storage/contact_messages.jsonl"

var contactMu sync.Mutex

type ContactMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	SentAt  string `json:"sentAt"`
}

func SaveContactMessage(name, email, subject, message string) error {
	contactMu.Lock()
	defer contactMu.Unlock()

	if err := os.MkdirAll("./storage", 0755); err != nil {
		return err
	}

	m := ContactMessage{
		Name:    name,
		Email:   email,
		Subject: subject,
		Message: message,
		SentAt:  time.Now().Format(time.RFC3339),
	}
	line, err := json.Marshal(m)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(contactLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(line, '\n'))
	return err
}
