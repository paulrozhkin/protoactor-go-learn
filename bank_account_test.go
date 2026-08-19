package main

import (
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/require"
)

var timeout = 5 * time.Second

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
