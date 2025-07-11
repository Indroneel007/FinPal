package db

import (
	"context"
	"examples/SimpleBankProject/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	s := util.RandomString(6)
	password, err := util.HashPassword(s)

	require.NoError(t, err)
	require.NotEmpty(t, password)
	require.True(t, util.CheckPasswordHash(s, password))

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: password,
		FullName:       util.RandomOwner(),
		Email:          util.RandomString(8),
	}
	user, err := testQueries.CreateUser(context.Background(), args)

	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	return user
}

func TestUser(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)
	fetchedUser, err := store.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.HashedPassword, fetchedUser.HashedPassword)
	require.Equal(t, user.FullName, fetchedUser.FullName)
	require.Equal(t, user.Email, fetchedUser.Email)

}
