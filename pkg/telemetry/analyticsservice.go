// Copyright 2023 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"context"

	"github.com/livekit/protocol/livekit"
	"github.com/livekit/protocol/logger"

	"github.com/livekit/livekit-server/pkg/config"
	"github.com/livekit/livekit-server/pkg/routing"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AnalyticsService
type AnalyticsService interface {
	SendStats(ctx context.Context, stats []*livekit.AnalyticsStat)
	SendEvent(ctx context.Context, events *livekit.AnalyticsEvent)
}

type analyticsService struct {
	analyticsKey string
	nodeID       string

	events livekit.AnalyticsRecorderService_IngestEventsClient
	stats  livekit.AnalyticsRecorderService_IngestStatsClient
}

func NewAnalyticsService(_ *config.Config, currentNode routing.LocalNode) AnalyticsService {
	return &analyticsService{
		analyticsKey: "", // TODO: conf.AnalyticsKey
		nodeID:       currentNode.Id,
	}
}

func (a *analyticsService) SendStats(_ context.Context, stats []*livekit.AnalyticsStat) {
	if a.stats == nil {
		return
	}

	for _, stat := range stats {
		stat.AnalyticsKey = a.analyticsKey
		stat.Node = a.nodeID
	}
	if err := a.stats.Send(&livekit.AnalyticsStats{Stats: stats}); err != nil {
		logger.Errorw("failed to send stats", err)
	}
}

func (a *analyticsService) SendEvent(_ context.Context, event *livekit.AnalyticsEvent) {
	if a.events == nil {
		return
	}

	event.AnalyticsKey = a.analyticsKey
	if err := a.events.Send(&livekit.AnalyticsEvents{
		Events: []*livekit.AnalyticsEvent{event},
	}); err != nil {
		logger.Errorw("failed to send event", err, "eventType", event.Type.String())
	}
}
