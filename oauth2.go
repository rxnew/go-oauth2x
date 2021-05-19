package oauth2x

//go:generate mockgen -destination mock/oauth2.go -package mock golang.org/x/oauth2 TokenSource

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	expiryDelta        = 15 * time.Second
	minPreloadInterval = 20 * time.Second
)

func NewClient(ctx context.Context, src oauth2.TokenSource) *http.Client {
	if src == nil {
		return oauth2.NewClient(ctx, nil)
	}
	return &http.Client{
		Transport: &oauth2.Transport{
			Base:   oauth2.NewClient(ctx, nil).Transport,
			Source: PreloadTokenSource(ctx, nil, src),
		},
	}
}

func PreloadTokenSource(ctx context.Context, t *oauth2.Token, src oauth2.TokenSource) oauth2.TokenSource {
	if pt, ok := src.(*preloadTokenSource); ok {
		if t == nil {
			return pt
		}
		src = pt.new
	}

	s := &preloadTokenSource{new: src, t: t}
	s.start(ctx)
	return s
}

type preloadTokenSource struct {
	new oauth2.TokenSource
	mu  sync.Mutex
	r   bool
	t   *oauth2.Token
	err error
}

func (s *preloadTokenSource) Token() (*oauth2.Token, error) {
	t, err := s.t, s.err
	if t.Valid() {
		return t, nil
	}
	if err != nil {
		return nil, err
	}
	return s.fetch()
}

func (s *preloadTokenSource) fetch() (*oauth2.Token, error) {
	s.r = true

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.r {
		s.t, s.err = s.new.Token()
		s.r = false
	}
	return s.t, s.err
}

func (s *preloadTokenSource) start(ctx context.Context) {
	if !s.t.Valid() {
		s.t, s.err = s.new.Token()
	}
	if s.t.Valid() && s.t.Expiry.IsZero() {
		return
	}

	t := s.t
	go func() {
		for {
			d := maxDuration(time.Until(t.Expiry.Round(0).Add(-expiryDelta)), minPreloadInterval)
			select {
			case <-time.After(d):
			case <-ctx.Done():
				return
			}
			t, _ = s.fetch()
		}
	}()
}

func maxDuration(x, y time.Duration) time.Duration {
	if x > y {
		return x
	}
	return y
}
