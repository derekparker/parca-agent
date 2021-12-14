package agent

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/parca-dev/parca-agent/pkg/debuginfo"
	"github.com/parca-dev/parca-agent/pkg/ksym"
)

func TestCgroupProfiler(t *testing.T) {
	var (
		unit           = "syslog.service"
		logger         = log.NewNopLogger()
		ksymCache      = ksym.NewKsymCache(logger)
		ctx            = context.Background()
		errCh          = make(chan error)
		externalLabels = map[string]string{"systemdunit": unit}
	)

	f, err := os.CreateTemp(os.TempDir(), "test.tmp")
	if err != nil {
		t.Fatal(err)
	}
	p := NewCgroupProfiler(
		logger,
		externalLabels,
		ksymCache,
		NewNoopProfileStoreClient(),
		debuginfo.NewNoopClient(),
		&SystemdUnitTarget{
			Name:     unit,
			NodeName: "testnode",
		},
		10*time.Second,
		sink,
		f.Name(),
	)
	if p == nil {
		t.Fatal("expected a non-nil profiler")
	}

	// Start the profiler. Run in separate goroutine so we can
	// assert since this operation blocks.
	go func(errc chan error) { errc <- p.Run(ctx) }(errCh)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for profiler to run")
	}
}

func sink(r Record) {

}
