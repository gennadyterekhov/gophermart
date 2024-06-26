package tests

import (
	"context"
	"testing"

	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestCanRegisterUsingMock(t *testing.T) {
	repo := repositories.NewRepositoryMock()

	user, err := repo.AddUser(context.Background(), "a", "a")
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Login)
	assert.Equal(t, "a", user.Password)

	user, err = repo.GetUserByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Login)
	assert.Equal(t, "a", user.Password)

	user, err = repo.GetUserByLogin(context.Background(), user.Login)
	assert.NoError(t, err)
	assert.Equal(t, "a", user.Login)
	assert.Equal(t, "a", user.Password)
}
