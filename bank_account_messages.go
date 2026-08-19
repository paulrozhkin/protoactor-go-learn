package main

import "errors"

var (
	AccountNotFound = errors.New("account not found")
	UnknownError    = errors.New("unknown error")
	InvariantError  = errors.New("invariant error")
)

type CreateAccount struct {
	RequestID string
	AccountId string
	FirstName string
	LastName  string
	Balance   float64
}

type GetAccountRequest struct {
	RequestID string
	AccountId string
}

type GetAccountResponse struct {
	RequestID string
	AccountId string
	FirstName string
	LastName  string
	Balance   float64
}

type DepositRequest struct {
	RequestID string
	AccountId string
	Amount    float64
}

type WithdrawRequest struct {
	RequestID string
	AccountId string
	Amount    float64
}

type OperationResponse struct {
	RequestID string
	AccountId string
	Success   bool
	Balance   float64
}

type InternalTransfer struct {
	RequestID string
	Recipient string
	Sender    string
	Amount    float64
}
