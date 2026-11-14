package models

import "testing"

func fsMetric(name, device, fstype, mount string, value float64) MetricData {
	return MetricData{MetricName: name, Value: value, Labels: map[string]string{
		"device": device, "fstype": fstype, "mountpoint": mount,
	}}
}

// A bootc node mounts its root partition at /sysroot, /var and /etc; it must
// be counted once, and tmpfs not at all.
func TestDiskCountsEachDeviceOnce(t *testing.T) {
	um := NewUsageMetrics()
	const gb = BytesInGB
	for _, m := range []string{"/sysroot", "/var", "/etc"} {
		um.UpdateUsageMetrics(fsMetric("node_filesystem_size_bytes", "/dev/vda4", "ext4", m, 20*gb))
		um.UpdateUsageMetrics(fsMetric("node_filesystem_free_bytes", "/dev/vda4", "ext4", m, 15*gb))
	}
	um.UpdateUsageMetrics(fsMetric("node_filesystem_size_bytes", "/dev/vda3", "ext4", "/boot", 1*gb))
	um.UpdateUsageMetrics(fsMetric("node_filesystem_free_bytes", "/dev/vda3", "ext4", "/boot", 1*gb))
	um.UpdateUsageMetrics(fsMetric("node_filesystem_size_bytes", "tmpfs", "tmpfs", "/tmp", 1*gb))

	if got := convertToGB(um.nodeDiskTotal); got != 21 {
		t.Fatalf("total %v GB, want 21", got)
	}
	if got := calculateDiskUsage(um.nodeDiskTotal, um.nodeDiskUsage); got != 5 {
		t.Fatalf("used %v GB, want 5", got)
	}

	um.ResetAllReadedMetrics()
	um.UpdateUsageMetrics(fsMetric("node_filesystem_size_bytes", "/dev/vda4", "ext4", "/var", 20*gb))
	if got := convertToGB(um.nodeDiskTotal); got != 20 {
		t.Fatalf("after reset: total %v GB, want 20", got)
	}
}
