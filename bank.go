package main

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

var timeout = 5 * time.Second

type Bank struct {
	accounts map[string]*actor.PID
}

func (b *Bank) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		b.accounts = make(map[string]*actor.PID)
	case *CreateAccount:
		pid := ctx.Spawn(actor.PropsFromProducer(func() actor.Actor { return &BankAccount{} }))
		accountId := uuid.NewString()
		b.accounts[accountId] = pid
		msg.AccountId = accountId
		ctx.Forward(pid)
	case *GetAccountRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
		}
		ctx.Forward(accountPid)
	case *DepositRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
		}
		ctx.Forward(accountPid)
	case *WithdrawRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
		}
		ctx.Forward(accountPid)
	case *InternalTransfer:
		b.makeInternalTransferInBlockingMode(ctx, msg)
	}
}

func (b *Bank) makeInternalTransferInBlockingMode(ctx actor.Context, msg *InternalTransfer) {
	senderPid, ok := b.accounts[msg.Sender]
	if !ok {
		ctx.Respond(AccountNotFound)
	}
	recipientPid, ok := b.accounts[msg.Recipient]
	if !ok {
		ctx.Respond(AccountNotFound)
	}
	// Simple logic without compensation and with blocking mode

	// Withdraw from sender
	response := ctx.RequestFuture(senderPid, &WithdrawRequest{
		RequestID: msg.RequestID,
		AccountId: msg.Sender,
		Amount:    msg.Amount,
	}, timeout)
	withdrawFeature, err := response.Result()
	if err != nil {
		ctx.Respond(UnknownError)
	}
	withdrawResponse := withdrawFeature.(*OperationResponse)
	if !withdrawResponse.Success {
		ctx.Respond(withdrawResponse)
	}

	// Deposit to recipient
	response = ctx.RequestFuture(recipientPid, &DepositRequest{
		RequestID: msg.RequestID,
		AccountId: msg.Recipient,
		Amount:    msg.Amount,
	}, timeout)
	depositFeature, err := response.Result()
	if err != nil {
		// simple compensation
		ctx.Request(senderPid, &DepositRequest{
			RequestID: msg.RequestID,
			AccountId: msg.Recipient,
			Amount:    msg.Amount,
		})
		ctx.Respond(&OperationResponse{
			Success:   false,
			RequestID: msg.RequestID,
			AccountId: msg.Recipient,
		})
	}
	depositResponse := depositFeature.(*OperationResponse)
	if !depositResponse.Success {
		// simple compensation
		ctx.Request(senderPid, &DepositRequest{
			RequestID: msg.RequestID,
			AccountId: msg.Recipient,
			Amount:    msg.Amount,
		})
	}
	ctx.Respond(depositResponse)
}
