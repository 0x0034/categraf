package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

type reloadTestCollector struct {
	id int
}

func (c *reloadTestCollector) Update(ch chan<- prometheus.Metric) error {
	return nil
}

func TestNewNodeCollectorReloadsCollectorParameters(t *testing.T) {
	const collectorName = "reloadtest"

	previousFactory, hadFactory := factories[collectorName]
	previousState, hadState := collectorState[collectorName]
	t.Cleanup(func() {
		if hadFactory {
			factories[collectorName] = previousFactory
		} else {
			delete(factories, collectorName)
		}
		if hadState {
			collectorState[collectorName] = previousState
		} else {
			delete(collectorState, collectorName)
		}
	})

	created := 0
	registerCollector(collectorName, false, func() (Collector, error) {
		created++
		return &reloadTestCollector{id: created}, nil
	})

	first, err := NewNodeCollector(false, "--collector.reloadtest")
	if err != nil {
		t.Fatalf("create first node collector: %v", err)
	}
	firstCollector, ok := first.Collectors[collectorName].(*reloadTestCollector)
	if !ok {
		t.Fatalf("first collector type = %T", first.Collectors[collectorName])
	}

	second, err := NewNodeCollector(false,
		"--collector.reloadtest",
		"--collector.reloadtest.fields=changed",
	)
	if err != nil {
		t.Fatalf("create second node collector: %v", err)
	}
	secondCollector, ok := second.Collectors[collectorName].(*reloadTestCollector)
	if !ok {
		t.Fatalf("second collector type = %T", second.Collectors[collectorName])
	}
	if firstCollector == secondCollector {
		t.Fatal("new node collector reused collector instance after parameters changed")
	}
	if secondCollector.id != 2 {
		t.Fatalf("second collector id = %d, want 2", secondCollector.id)
	}
}
