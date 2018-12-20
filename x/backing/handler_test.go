package backing

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestBackStoryMsg_FailBasicValidation(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 5 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration, validEvidence())
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidBackingPeriod := strings.Contains(res.Log, "901")
	assert.True(t, hasInvalidBackingPeriod, "should return err code")
}

func TestBackStoryMsg_FailInsufficientFunds(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration, validEvidence())
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInsufficientFunds := strings.Contains(res.Log, "65541")
	assert.True(t, hasInsufficientFunds, "should return err code")
}

func TestBackStoryMsg(t *testing.T) {
	ctx, bk, sk, ck, _, am := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := createFakeFundedAccount(ctx, am, sdk.Coins{amount})
	duration := 99 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration, validEvidence())
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID, "incorrect result backing id")
}
