package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAppAccount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	appAccount, err := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	assert.NoError(t, err)
	t.Log(appAccount)

	assert.Equal(t, appAccount.PrimaryAddress(), address)
	acc, err := keeper.PrimaryAccount(ctx, address)
	assert.NoError(t, err)
	assert.Equal(t, acc.GetPubKey(), publicKey)

	assert.Equal(t, false, appAccount.IsJailed)
}

func TestJailUntil_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	isJailed, err := keeper.IsJailed(ctx, createdAppAccount.PrimaryAddress())
	assert.Nil(t, err)
	assert.Equal(t, false, isJailed)

	err = keeper.JailUntil(ctx, createdAppAccount.PrimaryAddress(), time.Now().AddDate(0, 0, 10))
	assert.NoError(t, err)
	isJailed, err = keeper.IsJailed(ctx, createdAppAccount.PrimaryAddress())
	assert.Nil(t, err)
	assert.Equal(t, true, isJailed)

	accounts, err := keeper.JailedAccountsAfter(ctx, time.Now().AddDate(0, 0, 10))
	assert.NoError(t, err)
	assert.Len(t, accounts, 0)

	err = keeper.JailUntil(ctx, createdAppAccount.PrimaryAddress(), time.Now().AddDate(0, 0, 10))
	accounts, _ = keeper.JailedAccountsAfter(ctx, time.Now().AddDate(0, 0, 110))
	assert.Len(t, accounts, 0)

	accounts, err = keeper.JailedAccountsAfter(ctx, time.Now())
	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
}

func TestIncrementSlashCount_Success(t *testing.T) {
	ctx, keeper := mockDB()

	_, publicKey, address, coins := getFakeAppAccountParams()

	createdAppAccount, _ := keeper.CreateAppAccount(ctx, address, coins, publicKey)
	assert.Equal(t, createdAppAccount.SlashCount, 0)

	// incrementing once
	keeper.IncrementSlashCount(ctx, createdAppAccount.PrimaryAddress())
	returnedAppAccount, ok := keeper.getAppAccount(ctx, address)
	assert.True(t, ok)
	assert.Equal(t, returnedAppAccount.SlashCount, 1)

	// incrementing again
	keeper.IncrementSlashCount(ctx, createdAppAccount.PrimaryAddress())
	returnedAppAccount, ok = keeper.getAppAccount(ctx, address)
	assert.True(t, ok)
	assert.Equal(t, returnedAppAccount.SlashCount, 2)
}
