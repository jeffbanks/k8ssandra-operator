package telemetry

import (
	"errors"
	"strings"

	telemetry "github.com/k8ssandra/k8ssandra-operator/apis/telemetry/v1alpha1"
	"github.com/k8ssandra/k8ssandra-operator/pkg/cassandra"
	v1 "k8s.io/api/core/v1"
)

var (
	DefaultFilters = []string{"deny:org.apache.cassandra.metrics.Table",
		"deny:org.apache.cassandra.metrics.table",
		"allow:org.apache.cassandra.metrics.table.live_ss_table_count",
		"allow:org.apache.cassandra.metrics.Table.LiveSSTableCount",
		"allow:org.apache.cassandra.metrics.table.live_disk_space_used",
		"allow:org.apache.cassandra.metrics.table.LiveDiskSpaceUsed",
		"allow:org.apache.cassandra.metrics.Table.Pending",
		"allow:org.apache.cassandra.metrics.Table.Memtable",
		"allow:org.apache.cassandra.metrics.Table.Compaction",
		"allow:org.apache.cassandra.metrics.table.read",
		"allow:org.apache.cassandra.metrics.table.write",
		"allow:org.apache.cassandra.metrics.table.range",
		"allow:org.apache.cassandra.metrics.table.coordinator",
		"allow:org.apache.cassandra.metrics.table.dropped_mutations"}
)

// InjectCassandraTelemetryFilters adds MCAC filters to the cassandra container as an env variable.
// If filter list is set to nil, the default filters are used, otherwise the provided filters are used.
func InjectCassandraTelemetryFilters(telemetrySpec *telemetry.TelemetrySpec, dcConfig *cassandra.DatacenterConfig) error {
	filtersEnvVar := v1.EnvVar{}
	if telemetrySpec == nil || telemetrySpec.Mcac == nil || telemetrySpec.Mcac.MetricFilters == nil {
		// Default filters are applied
		filtersEnvVar = v1.EnvVar{Name: "METRIC_FILTERS", Value: strings.Join(DefaultFilters, " ")}
	} else {
		// Custom filters are applied
		filtersEnvVar = v1.EnvVar{Name: "METRIC_FILTERS", Value: strings.Join(*telemetrySpec.Mcac.MetricFilters, " ")}
	}
	if dcConfig.PodTemplateSpec == nil {
		return errors.New("no PodTemplateSpec in dcConfig to inject telemetry filtering env vars into")
	}
	containerIndex, containerFound := cassandra.FindContainer(dcConfig.PodTemplateSpec, "cassandra")
	if containerFound {
		dcConfig.PodTemplateSpec.Spec.Containers[containerIndex].Env = append(dcConfig.PodTemplateSpec.Spec.Containers[containerIndex].Env,
			filtersEnvVar)
	} else {
		dcConfig.PodTemplateSpec.Spec.Containers = append(dcConfig.PodTemplateSpec.Spec.Containers, v1.Container{
			Name: "cassandra",
			Env:  []v1.EnvVar{filtersEnvVar},
		})
	}
	return nil
}
