package match

import (
	"github.com/timshannon/bolthold"
	"github.com/twistedogic/doom/pkg/fetch"
	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/doom/pkg/store"
	"github.com/twistedogic/doom/pkg/target"
)

const (
	maxOffset   = -180
	concurrency = 5
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
	items, err := a.client.GetMatch(offset)
	if err != nil {
		return err
	}
	if err := a.store.BulkUpsert(items); err != nil {
		return err
	}
	loop := helper.NewLoop(-1)
	for i := range items {
		loop.Add(1)
		id := items[i].ID
		go func() {
			loop.Done(a.UpsertDetail(id))
		}()
	}
	return loop.Wait()
}

func (a *App) UpsertDetail(matchID int) error {
	items, err := a.client.GetDetail(matchID)
	if err != nil {
		return err
	}
	return a.store.BulkUpsert(items)
}

func (a *App) GetLastest(i interface{}) error {
	query := bolthold.Where("ID").MatchFunc(func(*bolthold.RecordAccess) (bool, error) {
		return true, nil
	}).SortBy("ID").Reverse()
	return a.store.FindOne(i, query)
}

func (a *App) Backfill() error {
	loop := helper.NewLoop(-1)
	for i := 0; i >= a.maxOffset; i-- {
		loop.Add(1)
		go func() {
			loop.Done(a.UpsertMatch(i))
		}()
	}
	return loop.Wait()
}

func (a *App) Update(t target.Target) error {
	panic("not implemented")
	return nil
}
