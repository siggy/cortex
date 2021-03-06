package tsdb

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cortexproject/cortex/pkg/storage/bucket"
)

func TestUsersScanner_ScanUsers_ShouldReturnedOwnedUsersOnly(t *testing.T) {
	bucketClient := &bucket.ClientMock{}
	bucketClient.MockIter("", []string{"user-1", "user-2", "user-3"}, nil)

	isOwned := func(userID string) (bool, error) {
		return userID == "user-1" || userID == "user-3", nil
	}

	s := NewUsersScanner(bucketClient, isOwned, log.NewNopLogger())
	actual, err := s.ScanUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"user-1", "user-3"}, actual)

}

func TestUsersScanner_ScanUsers_ShouldReturnUsersForWhichOwnerCheckFailed(t *testing.T) {
	expected := []string{"user-1", "user-2"}

	bucketClient := &bucket.ClientMock{}
	bucketClient.MockIter("", expected, nil)

	isOwned := func(userID string) (bool, error) {
		return false, errors.New("failed to check if user is owned")
	}

	s := NewUsersScanner(bucketClient, isOwned, log.NewNopLogger())
	actual, err := s.ScanUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
