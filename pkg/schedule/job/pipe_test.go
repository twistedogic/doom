package job_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/twistedogic/doom/pkg/model"
	"github.com/twistedogic/doom/pkg/model/odd"
	"github.com/twistedogic/doom/pkg/schedule/job"
	"github.com/twistedogic/doom/pkg/tap"
	"github.com/twistedogic/doom/pkg/tap/jockey"
	"github.com/twistedogic/doom/testutil"
)

func TestPipe(t *testing.T) {
	transformers := []model.TransformFunc{
		odd.Transform,
	}
	taps := []tap.Tap{
		jockey.New(jockey.Base, "fha", 5),
	}
	s := testutil.NewMockStore(t, make(map[string][]byte), false)
	dst := model.New(s, transformers...)
	wg := new(sync.WaitGroup)
	for i := range taps {
		wg.Add(1)
		go func(src tap.Tap) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
			defer cancel()
			j := job.New("test", src, dst, time.Millisecond)
			if err := j.Execute(ctx); err != nil {
				t.Fatal(err)
			}
		}(taps[i])
	}
	wg.Wait()
	if len(s.Content()) == 0 {
		t.Fatal("no data")
	}
	for k, v := range s.Content() {
		t.Logf("%s: %s\n", k, v)
	}
}
