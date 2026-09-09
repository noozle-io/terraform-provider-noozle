package client

import (
	"context"
	"encoding/json"
)

type outletOperationRecorderKey struct{}

// OutletOperationRecorder retains the most recent durable typed outlet mutation.
type OutletOperationRecorder struct {
	operation OutletOperation
	ok        bool
}

// WithOutletOperationRecorder adds a recorder to ctx for one sequential outlet
// mutation sequence. The final recorded operation is the generation to observe.
func WithOutletOperationRecorder(ctx context.Context) (context.Context, *OutletOperationRecorder) {
	recorder := &OutletOperationRecorder{}
	return context.WithValue(ctx, outletOperationRecorderKey{}, recorder), recorder
}

// FinalOperation returns the final durable mutation operation, if one was recorded.
func (r *OutletOperationRecorder) FinalOperation() (OutletOperation, bool) {
	if r == nil {
		return OutletOperation{}, false
	}
	return r.operation, r.ok
}

func recordOutletMutationOperation(ctx context.Context, method string, payload []byte) {
	if method != httpMethodPost && method != httpMethodPut && method != httpMethodPatch && method != httpMethodDelete {
		return
	}

	recorder, ok := ctx.Value(outletOperationRecorderKey{}).(*OutletOperationRecorder)
	if !ok || recorder == nil || len(payload) == 0 {
		return
	}

	var response struct {
		Operation *OutletOperation `json:"operation"`
	}
	if err := json.Unmarshal(payload, &response); err != nil || response.Operation == nil || response.Operation.DesiredGeneration == 0 {
		return
	}

	RecordOutletOperation(ctx, *response.Operation)
}

// RecordOutletOperation records operation as the latest mutation in ctx.
// It lets typed-client adapters preserve the same contract when their
// transport implementation is supplied by a test or alternate boundary.
func RecordOutletOperation(ctx context.Context, operation OutletOperation) {
	recorder, ok := ctx.Value(outletOperationRecorderKey{}).(*OutletOperationRecorder)
	if !ok || recorder == nil || operation.DesiredGeneration == 0 {
		return
	}

	recorder.operation = operation
	recorder.ok = true
}
