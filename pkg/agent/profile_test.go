package agent

import (
	"os"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/parca-dev/parca-agent/pkg/debuginfo"
	"github.com/parca-dev/parca-agent/pkg/ksym"
)

func TestCgroupProfiler(t *testing.T) {
	var (
		unit           = "test.service"
		logger         = log.NewNopLogger()
		ksymCache      = ksym.NewKsymCache(logger)
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
}

func sink(r Record) {

}
