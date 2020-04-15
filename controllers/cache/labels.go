package cache

const (
	AppLabel = "harbor-cluster-operator"
)

func (d *defaultCache) mergeLabels() map[string]string {
	dynLabels := map[string]string{
		AppLabel:                       d.Harbor.Name,
		"app.kubernetes.io/name":       "cache",
		"app.kubernetes.io/instance":   d.Harbor.Name,
		"app.kubernetes.io/managed-by": "harbor-cluster-operator",
		"app.kubernetes.io/part-of":    "harbor-cluster",
	}

	return MergeLabels(d.Labels, dynLabels, d.Harbor.Labels)
}

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

func generateLabels(component, role string) map[string]string {
	return map[string]string{
		"app":       AppLabel,
		"component": component,
		component:   role,
	}
}
