package main

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/require"
)

func TestBankCreateAccountSendDepositMessageAndGetCorrectResult(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &Bank{} })
	pid := system.Root.Spawn(props)
	future := system.Root.RequestFuture(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
		RequestID: "request-id-create-account",
	}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	accountId := response.(*GetAccountResponse).AccountId
	require.Equal(t, &GetAccountResponse{
		RequestID: "request-id-create-account",
		AccountId: accountId,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
	}, response)

	future = system.Root.RequestFuture(pid, &DepositRequest{
		Amount:    5,
		RequestID: "request-id-deposit",
		AccountId: accountId}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &OperationResponse{}, response)
	require.Equal(t, &OperationResponse{
		Balance:   5,
		RequestID: "request-id-deposit",
		AccountId: accountId,
		Success:   true,
	}, response)

	future = system.Root.RequestFuture(pid, &GetAccountRequest{
		RequestID: "request-id-get-account",
		AccountId: accountId,
	}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   5,
		RequestID: "request-id-get-account",
		AccountId: accountId,
	}, response)
}

func TestBankCreateAccountSendWithdrawMessageAndGetCorrectResult(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &Bank{} })
	pid := system.Root.Spawn(props)
	future := system.Root.RequestFuture(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
		RequestID: "request-id-create-account",
	}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	accountId := response.(*GetAccountResponse).AccountId
	require.Equal(t, &GetAccountResponse{
		RequestID: "request-id-create-account",
		AccountId: accountId,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
	}, response)

	future = system.Root.RequestFuture(pid, &DepositRequest{
		Amount:    5,
		RequestID: "request-id-deposit",
		AccountId: accountId}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &OperationResponse{}, response)
	require.Equal(t, &OperationResponse{
		Balance:   5,
		RequestID: "request-id-deposit",
		AccountId: accountId,
		Success:   true,
	}, response)

	future = system.Root.RequestFuture(pid, &WithdrawRequest{
		Amount:    2,
		RequestID: "request-id-withdraw",
		AccountId: accountId}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &OperationResponse{}, response)
	require.Equal(t, &OperationResponse{
		Balance:   3,
		RequestID: "request-id-withdraw",
		AccountId: accountId,
		Success:   true,
	}, response)

	future = system.Root.RequestFuture(pid, &GetAccountRequest{
		RequestID: "request-id-get-account",
		AccountId: accountId,
	}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   3,
		RequestID: "request-id-get-account",
		AccountId: accountId,
	}, response)
}

func TestBankCreateMultipleAccountSendDepositMessageAndMakeTransfer(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &Bank{} })
	pid := system.Root.Spawn(props)
	future := system.Root.RequestFuture(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   15,
		RequestID: "request-id-create-account",
	}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	accountIdSender := response.(*GetAccountResponse).AccountId
	require.Equal(t, &GetAccountResponse{
		RequestID: "request-id-create-account",
		AccountId: accountIdSender,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   15,
	}, response)

	future = system.Root.RequestFuture(pid, &CreateAccount{
		FirstName: "Petr",
		LastName:  "Petrovich",
		Balance:   9,
		RequestID: "request-id-create-account",
	}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	accountIdRecipient := response.(*GetAccountResponse).AccountId
	require.Equal(t, &GetAccountResponse{
		RequestID: "request-id-create-account",
		AccountId: accountIdRecipient,
		FirstName: "Petr",
		LastName:  "Petrovich",
		Balance:   9,
	}, response)

	future = system.Root.RequestFuture(pid, &InternalTransfer{
		Amount:    2,
		RequestID: "request-id-transfer",
		Recipient: accountIdRecipient,
		Sender:    accountIdSender}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &OperationResponse{}, response)
	require.Equal(t, &OperationResponse{
		Balance:   11,
		RequestID: "request-id-transfer",
		AccountId: accountIdRecipient,
		Success:   true,
	}, response)

	future = system.Root.RequestFuture(pid, &GetAccountRequest{
		RequestID: "request-id-get-account",
		AccountId: accountIdSender,
	}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   13,
		RequestID: "request-id-get-account",
		AccountId: accountIdSender,
	}, response)

	future = system.Root.RequestFuture(pid, &GetAccountRequest{
		RequestID: "request-id-get-account",
		AccountId: accountIdRecipient,
	}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Petr",
		LastName:  "Petrovich",
		Balance:   11,
		RequestID: "request-id-get-account",
		AccountId: accountIdRecipient,
	}, response)
}

func TestBankCreateInvariantAccountSendDepositMessageAndHandleSupervision(t *testing.T) {
	system := actor.NewActorSystem()
	invariantReceived := make(chan any, 1)
	strategy := actor.NewOneForOneStrategy(
		3,
		10*time.Second,
		func(reason any) actor.Directive {
			switch reason {
			case InvariantError:
				invariantReceived <- reason
				t.Logf("invariant error received: %v", reason)
				return actor.StopDirective
			}
			t.Logf("supervisor received: %v", reason)
			return actor.RestartDirective
		},
	)

	props := actor.PropsFromProducer(func() actor.Actor { return &Bank{} }, actor.WithSupervisor(strategy))
	pid := system.Root.Spawn(props)
	future := system.Root.RequestFuture(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   -10,
		RequestID: "request-id-create-account",
	}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	accountId := response.(*GetAccountResponse).AccountId
	require.Equal(t, &GetAccountResponse{
		RequestID: "request-id-create-account",
		AccountId: accountId,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   -10,
	}, response)

	future = system.Root.RequestFuture(pid, &WithdrawRequest{
		Amount:    5,
		RequestID: "request-id-withdraw",
		AccountId: accountId}, time.Second)
	select {
	case reason := <-invariantReceived:
		require.Equal(t, InvariantError, reason)

	case <-time.After(time.Second):
		t.Fatal("supervisor did not receive invariant error")
	}
}
