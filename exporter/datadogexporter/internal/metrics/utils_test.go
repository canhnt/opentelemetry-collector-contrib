// Copyright The OpenTelemetry Authors
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

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
)

func TestNewMetric(t *testing.T) {
	name := "test.metric"
	ts := uint64(1e9)
	value := 2.0
	tags := []string{"tag:value"}

	metric := newMetric(name, ts, value, tags)

	assert.Equal(t, "test.metric", *metric.Metric)
	// Assert timestamp conversion from uint64 ns to float64 s
	assert.Equal(t, 1.0, *metric.Points[0][0])
	// Assert value
	assert.Equal(t, 2.0, *metric.Points[0][1])
	// Assert tags
	assert.Equal(t, []string{"tag:value"}, metric.Tags)
}

func TestNewType(t *testing.T) {
	name := "test.metric"
	ts := uint64(1e9)
	value := 2.0
	tags := []string{"tag:value"}

	gauge := NewGauge(name, ts, value, tags)
	assert.Equal(t, gauge.GetType(), string(Gauge))

	count := NewCount(name, ts, value, tags)
	assert.Equal(t, count.GetType(), string(Count))

}

func TestDefaultMetrics(t *testing.T) {
	buildInfo := component.BuildInfo{
		Version: "1.0",
		Command: "otelcontribcol",
	}

	ms := DefaultMetrics("metrics", "test-host", uint64(2e9), buildInfo)

	assert.Equal(t, "otel.datadog_exporter.metrics.running", *ms[0].Metric)
	// Assert metrics list length (should be 1)
	assert.Equal(t, 1, len(ms))
	// Assert timestamp
	assert.Equal(t, 2.0, *ms[0].Points[0][0])
	// Assert value (should always be 1.0)
	assert.Equal(t, 1.0, *ms[0].Points[0][1])
	// Assert hostname tag is set
	assert.Equal(t, "test-host", *ms[0].Host)
	// Assert no other tags are set
	assert.ElementsMatch(t, []string{"version:1.0", "command:otelcontribcol"}, ms[0].Tags)
}
