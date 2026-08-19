package main

import (
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountSendDepositMessageAndGetCorrectResult(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &BankAccount{} })
	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
	})

	system.Root.Send(pid, &DepositRequest{Amount: 5})

	future := system.Root.RequestFuture(pid, &GetAccountRequest{}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   5,
	}, response)
}

func TestCreateAccountFeaturedDepositMessageAndGetCorrectResult(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &BankAccount{} })
	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   0,
	})

	future := system.Root.RequestFuture(pid, &DepositRequest{Amount: 5}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &OperationResponse{}, response)
	require.Equal(t, &OperationResponse{
		Success: true,
		Balance: 5,
	}, response)
}

func TestCreateMultipleAccount(t *testing.T) {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &BankAccount{} })
	firstAccountPid := system.Root.Spawn(props)
	system.Root.Send(firstAccountPid, &CreateAccount{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   3,
	})

	secondAccountPid := system.Root.Spawn(props)
	system.Root.Send(secondAccountPid, &CreateAccount{
		FirstName: "Petr",
		LastName:  "Petrovic",
		Balance:   5,
	})

	future := system.Root.RequestFuture(firstAccountPid, &GetAccountRequest{}, timeout)
	response, err := future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Balance:   3,
	}, response)

	future = system.Root.RequestFuture(secondAccountPid, &GetAccountRequest{}, timeout)
	response, err = future.Result()
	require.NoError(t, err)
	require.IsType(t, &GetAccountResponse{}, response)
	require.Equal(t, &GetAccountResponse{
		FirstName: "Petr",
		LastName:  "Petrovic",
		Balance:   5,
	}, response)
}
