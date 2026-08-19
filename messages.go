package main

type CreateAccount struct {
	FirstName string
	LastName  string
	Balance   float64
}

type GetAccountRequest struct {
}

type GetAccountResponse struct {
	FirstName string
	LastName  string
	Balance   float64
}

type DepositRequest struct {
	Amount float64
}

type WithdrawRequest struct {
	Amount float64
}

type OperationResponse struct {
	Success bool
	Balance float64
}
