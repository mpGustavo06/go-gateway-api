package dto

import "github.com/mpGustavo06/go-gateway-api/go-gateway/internal/domain"

type CreateAccountInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountOutput struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	APIKey    string  `json:"api_key,omitempty"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func ToAccount(input CreateAccountInput) *domain.Account {
	return domain.NewAccount(input.Name, input.Email)
}

func FromAccount(account *domain.Account) AccountOutput {
	return AccountOutput{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		APIKey:    account.APIKey,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z"),	
		UpdatedAt: account.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}