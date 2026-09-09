package provider

import (
	"context"
	"strings"
	"testing"
	"time"

	"terraform-provider-noozle/internal/client"
)

type outletStatusReaderStub struct {
	status client.OutletStatus
	err    error
}

func (s outletStatusReaderStub) GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error) {
	return s.status, s.err
}

func TestOutletLifecycleObserverReportsDisabledAndDeletedConvergence(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target outletLifecycleObservationTarget
		status client.OutletStatus
	}{
		"disabled": {
			target: outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}},
			status: client.OutletStatus{RequestedGeneration: 7, DesiredStatus: "disabled", RuntimeStatus: "disabled", Realization: client.OutletRuntimeRealization{Phase: "stopped"}},
		},
		"deleted": {
			target: outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Deleting: true},
			status: client.OutletStatus{RequestedGeneration: 7, DesiredStatus: "deleted", RuntimeStatus: "deleted"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			observer := newOutletLifecycleObserver(outletStatusReaderStub{status: test.status}, time.Minute, time.Millisecond)
			if err := observer.Observe(context.Background(), test.target); err != nil {
				t.Fatalf("observe outlet lifecycle: %v", err)
			}
		})
	}
}

func TestOutletLifecycleObserverRejectsSupersededGeneration(t *testing.T) {
	t.Parallel()

	observer := newOutletLifecycleObserver(outletStatusReaderStub{status: client.OutletStatus{RequestedGeneration: 7, Superseded: true}}, time.Minute, time.Millisecond)
	err := observer.Observe(context.Background(), outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Enabled: true})
	if err == nil || !strings.Contains(err.Error(), "superseded") {
		t.Fatalf("observe superseded outlet lifecycle error = %v", err)
	}
}

func TestOutletLifecycleObserverRejectsTerminalRuntimeFailure(t *testing.T) {
	t.Parallel()

	observer := newOutletLifecycleObserver(outletStatusReaderStub{status: client.OutletStatus{
		RequestedGeneration: 7,
		Realization:         client.OutletRuntimeRealization{Phase: "failed_start"},
	}}, time.Minute, time.Millisecond)
	err := observer.Observe(context.Background(), outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Enabled: true})
	if err == nil || !strings.Contains(err.Error(), "failed_start") {
		t.Fatalf("observe failed outlet lifecycle error = %v", err)
	}
}

func TestOutletLifecycleObserverRejectsDegradedRuntimeWithError(t *testing.T) {
	t.Parallel()

	observer := newOutletLifecycleObserver(outletStatusReaderStub{status: client.OutletStatus{
		RequestedGeneration: 7,
		Realization: client.OutletRuntimeRealization{
			Degraded:  true,
			LastError: &client.OutletRuntimeError{Detail: "broker unavailable"},
		},
	}}, time.Minute, time.Millisecond)
	err := observer.Observe(context.Background(), outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Enabled: true})
	if err == nil || !strings.Contains(err.Error(), "degraded") {
		t.Fatalf("observe degraded outlet lifecycle error = %v", err)
	}
}

func TestOutletLifecycleObserverReportsCanceledObservation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	observer := newOutletLifecycleObserver(outletStatusReaderStub{}, time.Minute, time.Millisecond)
	err := observer.Observe(ctx, outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Enabled: true})
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("observe canceled outlet lifecycle error = %v", err)
	}
}

func TestOutletLifecycleObserverAcceptsNotFoundAfterDeletion(t *testing.T) {
	t.Parallel()

	observer := newOutletLifecycleObserver(outletStatusReaderStub{err: &client.APIError{StatusCode: 404}}, time.Minute, time.Millisecond)
	err := observer.Observe(context.Background(), outletLifecycleObservationTarget{OutletType: "ntfy", OutletID: 99, Operation: client.OutletOperation{DesiredGeneration: 7}, Deleting: true})
	if err != nil {
		t.Fatalf("observe deleted outlet lifecycle: %v", err)
	}
}

func TestOutletLifecycleObserverReportsEnabledConvergence(t *testing.T) {
	t.Parallel()

	observer := newOutletLifecycleObserver(outletStatusReaderStub{status: client.OutletStatus{
		RequestedGeneration: 7,
		DesiredStatus:       "enabled",
		RuntimeStatus:       "active",
		Realization: client.OutletRuntimeRealization{
			Phase: "running",
			Ready: true,
		},
	}}, time.Minute, time.Millisecond)

	err := observer.Observe(context.Background(), outletLifecycleObservationTarget{
		OutletType: "ntfy",
		OutletID:   99,
		Operation:  client.OutletOperation{DesiredGeneration: 7},
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("observe enabled outlet lifecycle: %v", err)
	}
}

func TestObserveTypedOutletLifecycleUsesFinalGeneration(t *testing.T) {
	t.Parallel()

	err := observeTypedOutletLifecycle(context.Background(), outletStatusReaderStub{status: client.OutletStatus{
		RequestedGeneration: 11,
		DesiredStatus:       "enabled",
		RuntimeStatus:       "active",
		Realization: client.OutletRuntimeRealization{
			Phase: "running",
			Ready: true,
		},
	}}, "ntfy", 99, 11, true, false)
	if err != nil {
		t.Fatalf("observe final outlet generation: %v", err)
	}
}
