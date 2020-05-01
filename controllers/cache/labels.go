package cache

const (
	AppLabel = "harbor-cluster-operator"
)

// NewLabels returns new labels
func (redis *RedisReconciler) NewLabels() map[string]string {
	dynLabels := map[string]string{
		AppLabel:                       redis.Name,
		"app.kubernetes.io/name":       "cache",
		"app.kubernetes.io/instance":   redis.Namespace,
		"app.kubernetes.io/managed-by": "harbor-cluster-operator",
		"app.kubernetes.io/part-of":    "harbor-cluster",
	}

	return MergeLabels(redis.Labels, dynLabels, redis.HarborCluster.Labels)
}

// MergeLabels merge new label to existing labels
func MergeLabels(allLabels ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, labels := range allLabels {
		if labels != nil {
			for k, v := range labels {
				res[k] = v
			}
		}
	}
	return res
}

// generateLabels returns labels
func generateLabels(component, role string) map[string]string {
	return map[string]string{
		"app":       AppLabel,
		"component": component,
		component:   role,
	}
}

