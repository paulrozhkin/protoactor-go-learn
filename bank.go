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

		command := *msg
		command.AccountId = accountId

		ctx.RequestWithCustomSender(
			pid,
			&command,
			ctx.Sender(),
		)
	case *GetAccountRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
			return
		}
		ctx.Forward(accountPid)
	case *DepositRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
			return
		}
		ctx.Forward(accountPid)
	case *WithdrawRequest:
		accountPid, ok := b.accounts[msg.AccountId]
		if !ok {
			ctx.Respond(AccountNotFound)
			return
		}
		ctx.Forward(accountPid)
	case *InternalTransfer:
		b.makeInternalTransferInBehavior(ctx, msg)
	}
}

func (b *Bank) makeInternalTransferInBehavior(ctx actor.Context, msg *InternalTransfer) {
	senderPid, ok := b.accounts[msg.Sender]
	if !ok {
		ctx.Respond(AccountNotFound)
		return
	}
	recipientPid, ok := b.accounts[msg.Recipient]
	if !ok {
		ctx.Respond(AccountNotFound)
		return
	}
	transferPID := ctx.Spawn(
		actor.PropsFromProducer(func() actor.Actor {
			return NewTransferActor(senderPid, recipientPid)
		}),
	)

	// Сохраняет отправителя исходного запроса.
	ctx.Forward(transferPID)
}
