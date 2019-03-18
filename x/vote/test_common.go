package vote

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/argument"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/trubank"

	"github.com/TruStory/truchain/x/challenge"

	app "github.com/TruStory/truchain/types"
	c "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, c.Keeper) {

	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	argumentKey := sdk.NewKVStoreKey(argument.StoreKey)
	storyListKey := sdk.NewKVStoreKey(story.PendingListStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiringQueueStoreKey)
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	pendingGameListKey := sdk.NewKVStoreKey("pendingGameList")
	votingStoryQueueKey := sdk.NewKVStoreKey("gameQueue")
	voteKey := sdk.NewKVStoreKey("vote")
	backingKey := sdk.NewKVStoreKey("backing")
	paramsKey := sdk.NewKVStoreKey(sdkparams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkparams.TStoreKey)
	truBankKey := sdk.NewKVStoreKey(trubank.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(argumentKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(truBankKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	pk := sdkparams.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	ck := c.NewKeeper(catKey, codec)
	sk := story.NewKeeper(
		storyKey,
		storyListKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		ck,
		pk.Subspace(story.StoreKey),
		codec)

	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	truBankKeeper := trubank.NewKeeper(
		truBankKey,
		bankKeeper,
		ck,
		codec)

	stakeKeeper := stake.NewKeeper(
		sk,
		truBankKeeper,
		pk.Subspace(stake.StoreKey),
	)
	stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())

	argumentKeeper := argument.NewKeeper(
		argumentKey,
		sk,
		codec)

	backingKeeper := backing.NewKeeper(
		backingKey,
		argumentKeeper,
		stakeKeeper,
		sk,
		bankKeeper,
		truBankKeeper,
		ck,
		codec,
	)

	challengeKeeper := challenge.NewKeeper(
		challengeKey,
		argumentKeeper,
		stakeKeeper,
		backingKeeper,
		truBankKeeper,
		bankKeeper,
		sk,
		pk.Subspace(challenge.StoreKey),
		codec,
	)
	challenge.InitGenesis(ctx, challengeKeeper, challenge.DefaultGenesisState())

	k := NewKeeper(
		voteKey,
		votingStoryQueueKey,
		argumentKeeper,
		stakeKeeper,
		am,
		backingKeeper,
		challengeKeeper,
		sk,
		bankKeeper,
		truBankKeeper,
		pk.Subspace(StoreKey),
		codec)
	InitGenesis(ctx, k, DefaultGenesisState())

	return ctx, k, ck
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper, ck c.WriteKeeper, st story.Status) int64 {
	body := "TruStory validators can be bootstrapped with a single genesis file."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	storyID, _ := sk.Create(ctx, body, cat.ID, creator, source, storyType)
	s, _ := sk.Story(ctx, storyID)
	s.Status = st
	sk.UpdateStory(ctx, s)

	return storyID
}

func createFakeCategory(ctx sdk.Context, ck c.WriteKeeper) c.Category {
	existing, err := ck.GetCategory(ctx, 1)
	if err == nil {
		return existing
	}
	id := ck.Create(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}
