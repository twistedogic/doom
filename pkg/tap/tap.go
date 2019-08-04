package tap

import (
	"github.com/twistedogic/doom/pkg/fetch"
	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/doom/pkg/store"
)

const (
	maxOffset          = -180
	concurrency        = 5
	rateLimitPerSecond = 10
)

type App struct {
	maxOffset int
	store     *store.Store
	client    *fetch.Fetcher
}

func New(filename, base string, maxOffset, rate int) (*App, error) {
	s, err := store.New(filename)
	if err != nil {
		return nil, err
	}
	f := fetch.New(base, rate)
	return &App{maxOffset: maxOffset, store: s, client: f}, nil
}

func (a *App) UpsertMatch(offset int) error {
	out, err := a.client.GetMatch(offset)
	if err != nil {
		return err
	}
	loop := helper.NewLoop(-1)
	for i := range out {
		loop.Add(2)
		id := out[i].ID
		go func() {
			loop.Done(a.store.UpsertItem(out[i]))
			loop.Done(a.UpsertDetail(id))
		}()
	}
	return loop.Wait()
}

func (a *App) UpsertDetail(matchID int) error {
	out, err := a.client.GetDetail(matchID)
	if err != nil {
		return err
	}
	loop := helper.NewLoop(-1)
	for i := range out {
		loop.Add(1)
		go func() {
			loop.Done(a.store.UpsertItem(out[i]))
		}()
	}
	return loop.Wait()
}

func (a *App) update() error {
	loop := helper.NewLoop(-1)
	for i := 0; i >= a.maxOffset; i-- {
		loop.Add(1)
		go func() {
			loop.Done(a.UpsertMatch(i))
		}()
	}
	return loop.Wait()
}
