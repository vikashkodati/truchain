package vote

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Vote errors reserve 1100 ~ 1199.
const (
	DefaultCodespace sdk.CodespaceType = 12

	CodeNotFound       sdk.CodeType = 1201
	CodeDuplicate      sdk.CodeType = 1202
	CodeGameNotStarted sdk.CodeType = 1203
	CodeUnknownVote    sdk.CodeType = 1204
	CodeInvalidVote    sdk.CodeType = 1205
)

// ErrNotFound creates an error when the searched entity is not found
func ErrNotFound(id int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeNotFound,
		"Vote with id "+fmt.Sprintf("%d", id)+" not found")
}

// ErrDuplicateVoteForGame throws when a vote already has been cast
func ErrDuplicateVoteForGame(
	gameID int64, user sdk.AccAddress) sdk.Error {

	return sdk.NewError(
		DefaultCodespace,
		CodeDuplicate,
		"Vote with for game "+fmt.Sprintf("%d", gameID)+" has already been cast by user "+user.String())
}

// ErrGameNotStarted is thrown when a vote is attempted on a story
// that hasn't begun the validation game yet
func ErrGameNotStarted(storyID int64) sdk.Error {
	return sdk.NewError(
		DefaultCodespace,
		CodeGameNotStarted,
		"Validation game not started for story: "+
			fmt.Sprintf("%d", storyID))
}

// ErrVoteHandler returns an unknown Vote type error
func ErrVoteHandler(vote interface{}) sdk.Error {
	if mType := reflect.TypeOf(vote); mType != nil {
		errMsg := "Unrecognized Vote type: " + mType.Name()
		return sdk.NewError(DefaultCodespace, CodeUnknownVote, errMsg)
	}

	return sdk.NewError(DefaultCodespace, CodeUnknownVote, "Unknown Vote type")
}