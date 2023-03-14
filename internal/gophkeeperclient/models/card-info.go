package models

// BankCard structure represents a bank card and its associated information.
// It includes the card number, the expiry month and year, the CVV, and the cardholder's name.
// The CardNumber field represents the card number, the ExpiryMonth and ExpiryYear
// fields represent the month and year the card expires, respectively, the CVV field represents the card's security code,
// and the CardHolderName field represents the name of the cardholder.
type BankCard struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
	CardHolderName string `json:"card_holder_name"`
}
