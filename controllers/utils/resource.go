package utils

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func GetLimits(limCpu string, limitMem string) corev1.ResourceList {
	return generateResourceList(limCpu, limitMem)
}

func GetRequests(reqCpu string, reqMem string) corev1.ResourceList {
	return generateResourceList(reqCpu, reqMem)
}

func generateResourceList(cpu string, memory string) corev1.ResourceList {
	resources := corev1.ResourceList{}
	if cpu != "" {
		resources[corev1.ResourceCPU], _ = resource.ParseQuantity(cpu)
	}
	if memory != "" {
		resources[corev1.ResourceMemory], _ = resource.ParseQuantity(memory)
	}
	return resources
}
