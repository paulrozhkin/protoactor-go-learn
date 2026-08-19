package main

import (
	"github.com/asynkron/protoactor-go/actor"
)

type BankAccount struct {
	AccountId string
	FirstName string
	LastName  string
	Balance   float64
}

func (b *BankAccount) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateAccount:
		b.FirstName = msg.FirstName
		b.LastName = msg.LastName
		b.Balance = msg.Balance
		b.AccountId = msg.AccountId
		ctx.Respond(&GetAccountResponse{
			RequestID: msg.RequestID,
			AccountId: b.AccountId,
			FirstName: b.FirstName,
			LastName:  b.LastName,
			Balance:   b.Balance,
		})
	case *GetAccountRequest:
		ctx.Respond(&GetAccountResponse{
			RequestID: msg.RequestID,
			AccountId: b.AccountId,
			FirstName: b.FirstName,
			LastName:  b.LastName,
			Balance:   b.Balance,
		})
	case *DepositRequest:
		b.Balance += msg.Amount
		ctx.Respond(&OperationResponse{
			RequestID: msg.RequestID,
			AccountId: b.AccountId,
			Success:   true,
			Balance:   b.Balance,
		})
	case *WithdrawRequest:
		response := &OperationResponse{RequestID: msg.RequestID,
			AccountId: b.AccountId}
		if b.Balance-msg.Amount > 0 {
			b.Balance -= msg.Amount
			response.Success = true
		}
		response.Balance = b.Balance
		ctx.Respond(response)
	}
}
