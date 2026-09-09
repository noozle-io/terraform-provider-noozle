package provider

import (
	"context"
	"fmt"
	"time"

	"terraform-provider-noozle/internal/client"
)

const (
	defaultOutletLifecycleObservationTimeout  = 10 * time.Minute
	defaultOutletLifecycleObservationInterval = 2 * time.Second
)

type outletLifecycleStatusReader interface {
	GetTypedOutletStatus(context.Context, string, int64, int64) (client.OutletStatus, error)
}

type outletLifecycleObservationTarget struct {
	OutletType string
	OutletID   int64
	Operation  client.OutletOperation
	Enabled    bool
	Deleting   bool
}

type outletLifecycleObserver struct {
	statusReader outletLifecycleStatusReader
	timeout      time.Duration
	interval     time.Duration
}

func newOutletLifecycleObserver(statusReader outletLifecycleStatusReader, timeout, interval time.Duration) outletLifecycleObserver {
	if timeout <= 0 {
		timeout = defaultOutletLifecycleObservationTimeout
	}
	if interval <= 0 {
		interval = defaultOutletLifecycleObservationInterval
	}

	return outletLifecycleObserver{
		statusReader: statusReader,
		timeout:      timeout,
		interval:     interval,
	}
}

func observeTypedOutletLifecycle(ctx context.Context, statusReader outletLifecycleStatusReader, outletType string, outletID, generation int64, enabled, deleting bool) error {
	return newOutletLifecycleObserver(statusReader, 0, 0).Observe(ctx, outletLifecycleObservationTarget{
		OutletType: outletType,
		OutletID:   outletID,
		Operation:  client.OutletOperation{DesiredGeneration: generation},
		Enabled:    enabled,
		Deleting:   deleting,
	})
}

func observeFinalTypedOutletOperation(ctx context.Context, statusReader outletLifecycleStatusReader, recorder *client.OutletOperationRecorder, outletType string, outletID int64, enabled bool) error {
	operation, ok := recorder.FinalOperation()
	if !ok {
		return fmt.Errorf("no durable outlet operation was returned by the successful mutation sequence")
	}

	return observeTypedOutletLifecycle(ctx, statusReader, outletType, outletID, operation.DesiredGeneration, enabled, false)
}

func (o outletLifecycleObserver) Observe(ctx context.Context, target outletLifecycleObservationTarget) error {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("outlet lifecycle observation timed out: %w", err)
		}

		status, err := o.statusReader.GetTypedOutletStatus(ctx, target.OutletType, target.OutletID, target.Operation.DesiredGeneration)
		if err != nil {
			if target.Deleting && client.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("get typed outlet lifecycle status: %w", err)
		}

		if err := outletLifecycleStatusMatches(status, target); err != nil {
			return err
		}
		if outletLifecycleStatusConverged(status, target) {
			return nil
		}

		timer := time.NewTimer(o.interval)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return fmt.Errorf("outlet lifecycle observation timed out: %w", ctx.Err())
		case <-timer.C:
		}
	}
}

func outletLifecycleStatusMatches(status client.OutletStatus, target outletLifecycleObservationTarget) error {
	if status.Superseded {
		return fmt.Errorf("outlet lifecycle generation %d was superseded", target.Operation.DesiredGeneration)
	}
	if status.RequestedGeneration != target.Operation.DesiredGeneration {
		return fmt.Errorf("outlet lifecycle returned requested generation %d, want %d", status.RequestedGeneration, target.Operation.DesiredGeneration)
	}
	if status.Realization.Phase == "failed_start" || status.Realization.Phase == "stuck" {
		return fmt.Errorf("outlet lifecycle reached terminal realization phase %q", status.Realization.Phase)
	}
	if (status.Realization.Degraded || status.RuntimeStatus == "degraded") && status.Realization.LastError != nil {
		return fmt.Errorf("outlet lifecycle is degraded: %s", status.Realization.LastError.Detail)
	}

	return nil
}

func outletLifecycleStatusConverged(status client.OutletStatus, target outletLifecycleObservationTarget) bool {
	if target.Deleting {
		return status.DesiredStatus == "deleted" && status.RuntimeStatus == "deleted"
	}
	if target.Enabled {
		return status.DesiredStatus == "enabled" && status.RuntimeStatus == "active" && status.Realization.Phase == "running" && status.Realization.Ready
	}

	return status.DesiredStatus == "disabled" && status.RuntimeStatus == "disabled" && status.Realization.Phase == "stopped"
}
