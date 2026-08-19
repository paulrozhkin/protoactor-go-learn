package main

import "github.com/asynkron/protoactor-go/actor"

type TransferActor struct {
	behavior actor.Behavior

	senderAccount    *actor.PID
	recipientAccount *actor.PID

	replyTo *actor.PID
	command InternalTransfer
}

func NewTransferActor(
	senderAccount *actor.PID,
	recipientAccount *actor.PID,
) actor.Actor {
	a := &TransferActor{
		senderAccount:    senderAccount,
		recipientAccount: recipientAccount,
	}

	a.behavior.Become(a.idle)

	return a
}

func (a *TransferActor) Receive(ctx actor.Context) {
	a.behavior.Receive(ctx)
}

func (a *TransferActor) idle(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		// Ждём InternalTransfer.

	case *InternalTransfer:
		a.command = *msg
		a.replyTo = ctx.Sender()

		ctx.SetReceiveTimeout(timeout)
		a.behavior.Become(a.waitingWithdraw)

		ctx.Request(a.senderAccount, &WithdrawRequest{
			RequestID: msg.RequestID,
			AccountId: msg.Sender,
			Amount:    msg.Amount,
		})
	}
}

func (a *TransferActor) waitingWithdraw(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *OperationResponse:
		if !msg.Success {
			a.finish(ctx, msg)
			return
		}

		a.behavior.Become(a.waitingDeposit)

		ctx.Request(a.recipientAccount, &DepositRequest{
			RequestID: a.command.RequestID,
			AccountId: a.command.Recipient,
			Amount:    a.command.Amount,
		})

	case *actor.ReceiveTimeout:
		a.finish(ctx, UnknownError)
	}
}

func (a *TransferActor) waitingDeposit(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *OperationResponse:
		if msg.Success {
			a.finish(ctx, msg)
			return
		}

		a.startCompensation(ctx)

	case *actor.ReceiveTimeout:
		// Результат зачисления неизвестен.
		a.finish(ctx, UnknownError)
	}
}

func (a *TransferActor) startCompensation(ctx actor.Context) {
	a.behavior.Become(a.waitingCompensation)

	ctx.Request(a.senderAccount, &DepositRequest{
		RequestID: a.command.RequestID + ":compensation",
		AccountId: a.command.Sender,
		Amount:    a.command.Amount,
	})
}

func (a *TransferActor) waitingCompensation(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *OperationResponse:
		if !msg.Success {
			a.finish(ctx, UnknownError)
			return
		}

		a.finish(ctx, &OperationResponse{
			RequestID: a.command.RequestID,
			AccountId: a.command.Sender,
			Success:   false,
		})

	case *actor.ReceiveTimeout:
		a.finish(ctx, UnknownError)
	}
}

func (a *TransferActor) finish(
	ctx actor.Context,
	response any,
) {
	ctx.CancelReceiveTimeout()

	if a.replyTo != nil {
		ctx.Send(a.replyTo, response)
	}

	ctx.Stop(ctx.Self())
}
