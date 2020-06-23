package checks

import (
	"context"

	"github.com/influxdata/influxdb/v2"
	"github.com/influxdata/influxdb/v2/kit/feature"
)

var _ influxdb.CheckService = (*FeatureFlaggedService)(nil)

// FeatureFlaggedService is a feature flag driven proxy implementation
// of the influxdb.CheckService
type FeatureFlaggedService struct {
	influxdb.CheckService

	checksWithoutURMs *Service
}

// NewFlaggedService is a feature flag driven proxy for a check service.
// It enables the locally defined check service which does not use URMs to function when
// the feature flag is enabled on the request context.
func NewFlaggedService(original influxdb.CheckService, checksWithoutURMs *Service) *FeatureFlaggedService {
	return &FeatureFlaggedService{original, checksWithoutURMs}
}

func (f *FeatureFlaggedService) readEnabled(ctx context.Context) bool {
	return feature.ChecksReadDisableUrms().Enabled(ctx)
}

func (f *FeatureFlaggedService) writeEnabled(ctx context.Context) bool {
	return feature.ChecksWriteDisableUrms().Enabled(ctx)
}

// FindCheckByID returns a single check by ID.
func (f *FeatureFlaggedService) FindCheckByID(ctx context.Context, id influxdb.ID) (influxdb.Check, error) {
	if f.readEnabled(ctx) {
		return f.checksWithoutURMs.FindCheckByID(ctx, id)
	}

	return f.CheckService.FindCheckByID(ctx, id)
}

// FindCheck returns the first check that matches filter.
func (f *FeatureFlaggedService) FindCheck(ctx context.Context, filter influxdb.CheckFilter) (influxdb.Check, error) {
	if f.readEnabled(ctx) {
		return f.checksWithoutURMs.FindCheck(ctx, filter)
	}

	return f.CheckService.FindCheck(ctx, filter)
}

// FindChecks returns a list of checks that match filter and the total count of matching checks.
// Additional options provide pagination & sorting.
func (f *FeatureFlaggedService) FindChecks(ctx context.Context, filter influxdb.CheckFilter, opt ...influxdb.FindOptions) ([]influxdb.Check, int, error) {
	if f.readEnabled(ctx) {
		return f.checksWithoutURMs.FindChecks(ctx, filter, opt...)
	}

	return f.CheckService.FindChecks(ctx, filter, opt...)
}

// CreateCheck creates a new check and sets b.ID with the new identifier.
func (f *FeatureFlaggedService) CreateCheck(ctx context.Context, c influxdb.CheckCreate, userID influxdb.ID) error {
	if f.writeEnabled(ctx) {
		return f.checksWithoutURMs.CreateCheck(ctx, c, userID)
	}

	return f.CheckService.CreateCheck(ctx, c, userID)
}

// UpdateCheck updates the whole check.
// Returns the new check state after update.
func (f *FeatureFlaggedService) UpdateCheck(ctx context.Context, id influxdb.ID, c influxdb.CheckCreate) (influxdb.Check, error) {
	if f.writeEnabled(ctx) {
		return f.checksWithoutURMs.UpdateCheck(ctx, id, c)
	}

	return f.CheckService.UpdateCheck(ctx, id, c)
}

// PatchCheck updates a single bucket with changeset.
// Returns the new check state after update.
func (f *FeatureFlaggedService) PatchCheck(ctx context.Context, id influxdb.ID, upd influxdb.CheckUpdate) (influxdb.Check, error) {
	if f.writeEnabled(ctx) {
		return f.checksWithoutURMs.PatchCheck(ctx, id, upd)
	}

	return f.CheckService.PatchCheck(ctx, id, upd)
}

// DeleteCheck will delete the check by id.
func (f *FeatureFlaggedService) DeleteCheck(ctx context.Context, id influxdb.ID) error {
	if f.writeEnabled(ctx) {
		return f.checksWithoutURMs.DeleteCheck(ctx, id)
	}

	return f.CheckService.DeleteCheck(ctx, id)
}
