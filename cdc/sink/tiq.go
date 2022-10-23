// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package sink

import (
	"context"
	"sync/atomic"

	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"go.uber.org/zap"
)

// newTiQSink creates a tiq sink
func newTiQSink(ctx context.Context) *tiQSink {
	return &tiQSink{
		statistics: NewStatistics(ctx, "tiq"),
	}
}

type tiQSink struct {
	statistics      *Statistics
	accumulated     uint64
	lastAccumulated uint64
}

// Init table sink resources
func (b *tiQSink) Init(tableID model.TableID) error {
	return nil
}

func (b *tiQSink) EmitRowChangedEvents(ctx context.Context, rows ...*model.RowChangedEvent) error {
	for _, row := range rows {
		log.Debug("TiQsSink: EmitRowChangedEvents", zap.Any("row", row))
		log.Warn("TiQsSink: EmitRowChangedEvents", zap.Any("row", row))
	}
	rowsCount := len(rows)
	atomic.AddUint64(&b.accumulated, uint64(rowsCount))
	b.statistics.AddRowsCount(rowsCount)
	return nil
}

func (b *tiQSink) FlushRowChangedEvents(ctx context.Context, _ model.TableID, resolvedTs uint64) (uint64, error) {
	log.Debug("TiQSink: FlushRowChangedEvents", zap.Uint64("resolvedTs", resolvedTs))
	err := b.statistics.RecordBatchExecution(func() (int, error) {
		// TODO: add some random replication latency
		accumulated := atomic.LoadUint64(&b.accumulated)
		batchSize := accumulated - b.lastAccumulated
		b.lastAccumulated = accumulated
		return int(batchSize), nil
	})
	b.statistics.PrintStatus(ctx)
	return resolvedTs, err
}

func (b *tiQSink) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	log.Debug("TiQSink: Checkpoint Event", zap.Uint64("ts", ts))
	return nil
}

func (b *tiQSink) EmitDDLEvent(ctx context.Context, ddl *model.DDLEvent) error {
	log.Debug("TiQSink: DDL Event", zap.Any("ddl", ddl))
	return nil
}

func (b *tiQSink) Close(ctx context.Context) error {
	return nil
}

func (b *tiQSink) Barrier(ctx context.Context, tableID model.TableID) error {
	return nil
}
