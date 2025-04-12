package domain

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateAPIKey() string {
	formattedKey := make([]byte, 16)

	rand.Read(formattedKey)

	return hex.EncodeToString(formattedKey)
}

func NewAccount(name, email string) *Account {
	return &Account {
		ID: uuid.New().String(),
		Name: name,
		Email: email,
		APIKey: generateAPIKey(),
		Balance: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (props *Account) AddBalance(amount float64) {
	props.mu.Lock()
	defer props.mu.Unlock()
	
	props.Balance += amount
	props.UpdatedAt = time.Now()
}