package oauth2x_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/rxnew/go-oauth2x"
	"github.com/rxnew/go-oauth2x/mock"
)

func TestPreloadTokenSource_Token(t *testing.T) {
	t.Run("static token", func(t *testing.T) {
		token := &oauth2.Token{
			AccessToken: "test",
			TokenType:   "bearer",
		}
		actual, err := oauth2x.PreloadTokenSource(context.Background(), nil, oauth2.StaticTokenSource(token)).Token()
		require.NoError(t, err)
		assert.Equal(t, token.AccessToken, actual.AccessToken)
		assert.Equal(t, token.TokenType, actual.TokenType)
	})

	t.Run("reuse token", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		s := mock.NewMockTokenSource(mockCtrl)
		s.EXPECT().Token().Return(&oauth2.Token{
			AccessToken: "test",
			TokenType:   "bearer",
			Expiry:      time.Now().Add(12 * time.Hour),
		}, nil).Times(1)
		sut := oauth2x.PreloadTokenSource(context.Background(), nil, s)
		_, err := sut.Token()
		require.NoError(t, err)
		_, err = sut.Token()
		require.NoError(t, err)
	})

	t.Run("auto fetch token", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		s := mock.NewMockTokenSource(mockCtrl)
		s.EXPECT().Token().DoAndReturn(func() (*oauth2.Token, error) {
			return &oauth2.Token{
				AccessToken: "test",
				TokenType:   "bearer",
				Expiry:      time.Now().Add(1 * time.Nanosecond),
			}, nil
		}).Times(3)
		sut := oauth2x.PreloadTokenSource(context.Background(), nil, s)
		time.Sleep(10 * time.Millisecond)
		_, err := sut.Token()
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
		_, err = sut.Token()
		require.NoError(t, err)
	})
}
