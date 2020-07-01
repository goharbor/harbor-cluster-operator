// +build !ignore_autogenerated

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Azure) DeepCopyInto(out *Azure) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Azure.
func (in *Azure) DeepCopy() *Azure {
	if in == nil {
		return nil
	}
	out := new(Azure)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChartMuseum) DeepCopyInto(out *ChartMuseum) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChartMuseum.
func (in *ChartMuseum) DeepCopy() *ChartMuseum {
	if in == nil {
		return nil
	}
	out := new(ChartMuseum)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Clair) DeepCopyInto(out *Clair) {
	*out = *in
	if in.VulnerabilitySources != nil {
		in, out := &in.VulnerabilitySources, &out.VulnerabilitySources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Clair.
func (in *Clair) DeepCopy() *Clair {
	if in == nil {
		return nil
	}
	out := new(Clair)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Database) DeepCopyInto(out *Database) {
	*out = *in
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(PostgresSQL)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Database.
func (in *Database) DeepCopy() *Database {
	if in == nil {
		return nil
	}
	out := new(Database)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Gcs) DeepCopyInto(out *Gcs) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Gcs.
func (in *Gcs) DeepCopy() *Gcs {
	if in == nil {
		return nil
	}
	out := new(Gcs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborCluster) DeepCopyInto(out *HarborCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborCluster.
func (in *HarborCluster) DeepCopy() *HarborCluster {
	if in == nil {
		return nil
	}
	out := new(HarborCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborClusterCondition) DeepCopyInto(out *HarborClusterCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborClusterCondition.
func (in *HarborClusterCondition) DeepCopy() *HarborClusterCondition {
	if in == nil {
		return nil
	}
	out := new(HarborClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborClusterList) DeepCopyInto(out *HarborClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HarborCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborClusterList.
func (in *HarborClusterList) DeepCopy() *HarborClusterList {
	if in == nil {
		return nil
	}
	out := new(HarborClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarborClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborClusterSpec) DeepCopyInto(out *HarborClusterSpec) {
	*out = *in
	out.CertificateIssuerRef = in.CertificateIssuerRef
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int32)
		**out = **in
	}
	if in.ImageSource != nil {
		in, out := &in.ImageSource, &out.ImageSource
		*out = new(ImageSource)
		**out = **in
	}
	if in.JobService != nil {
		in, out := &in.JobService, &out.JobService
		*out = new(JobService)
		**out = **in
	}
	if in.Clair != nil {
		in, out := &in.Clair, &out.Clair
		*out = new(Clair)
		(*in).DeepCopyInto(*out)
	}
	if in.Trivy != nil {
		in, out := &in.Trivy, &out.Trivy
		*out = new(Trivy)
		**out = **in
	}
	if in.ChartMuseum != nil {
		in, out := &in.ChartMuseum, &out.ChartMuseum
		*out = new(ChartMuseum)
		**out = **in
	}
	if in.Notary != nil {
		in, out := &in.Notary, &out.Notary
		*out = new(Notary)
		**out = **in
	}
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(Redis)
		(*in).DeepCopyInto(*out)
	}
	if in.Database != nil {
		in, out := &in.Database, &out.Database
		*out = new(Database)
		(*in).DeepCopyInto(*out)
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(Storage)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborClusterSpec.
func (in *HarborClusterSpec) DeepCopy() *HarborClusterSpec {
	if in == nil {
		return nil
	}
	out := new(HarborClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarborClusterStatus) DeepCopyInto(out *HarborClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]HarborClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarborClusterStatus.
func (in *HarborClusterStatus) DeepCopy() *HarborClusterStatus {
	if in == nil {
		return nil
	}
	out := new(HarborClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Hosts) DeepCopyInto(out *Hosts) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Hosts.
func (in *Hosts) DeepCopy() *Hosts {
	if in == nil {
		return nil
	}
	out := new(Hosts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageSource) DeepCopyInto(out *ImageSource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageSource.
func (in *ImageSource) DeepCopy() *ImageSource {
	if in == nil {
		return nil
	}
	out := new(ImageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InCluster) DeepCopyInto(out *InCluster) {
	*out = *in
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(MinIOSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InCluster.
func (in *InCluster) DeepCopy() *InCluster {
	if in == nil {
		return nil
	}
	out := new(InCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobService) DeepCopyInto(out *JobService) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobService.
func (in *JobService) DeepCopy() *JobService {
	if in == nil {
		return nil
	}
	out := new(JobService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MinIOSpec) DeepCopyInto(out *MinIOSpec) {
	*out = *in
	in.VolumeClaimTemplate.DeepCopyInto(&out.VolumeClaimTemplate)
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MinIOSpec.
func (in *MinIOSpec) DeepCopy() *MinIOSpec {
	if in == nil {
		return nil
	}
	out := new(MinIOSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Notary) DeepCopyInto(out *Notary) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Notary.
func (in *Notary) DeepCopy() *Notary {
	if in == nil {
		return nil
	}
	out := new(Notary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Oss) DeepCopyInto(out *Oss) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Oss.
func (in *Oss) DeepCopy() *Oss {
	if in == nil {
		return nil
	}
	out := new(Oss)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresSQL) DeepCopyInto(out *PostgresSQL) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresSQL.
func (in *PostgresSQL) DeepCopy() *PostgresSQL {
	if in == nil {
		return nil
	}
	out := new(PostgresSQL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redis) DeepCopyInto(out *Redis) {
	*out = *in
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(RedisSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redis.
func (in *Redis) DeepCopy() *Redis {
	if in == nil {
		return nil
	}
	out := new(Redis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedisServer) DeepCopyInto(out *RedisServer) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisServer.
func (in *RedisServer) DeepCopy() *RedisServer {
	if in == nil {
		return nil
	}
	out := new(RedisServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedisSpec) DeepCopyInto(out *RedisSpec) {
	*out = *in
	if in.Server != nil {
		in, out := &in.Server, &out.Server
		*out = new(RedisServer)
		(*in).DeepCopyInto(*out)
	}
	if in.Sentinel != nil {
		in, out := &in.Sentinel, &out.Sentinel
		*out = new(Sentinel)
		**out = **in
	}
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]Hosts, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisSpec.
func (in *RedisSpec) DeepCopy() *RedisSpec {
	if in == nil {
		return nil
	}
	out := new(RedisSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *S3) DeepCopyInto(out *S3) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new S3.
func (in *S3) DeepCopy() *S3 {
	if in == nil {
		return nil
	}
	out := new(S3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Sentinel) DeepCopyInto(out *Sentinel) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sentinel.
func (in *Sentinel) DeepCopy() *Sentinel {
	if in == nil {
		return nil
	}
	out := new(Sentinel)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Storage) DeepCopyInto(out *Storage) {
	*out = *in
	if in.InCluster != nil {
		in, out := &in.InCluster, &out.InCluster
		*out = new(InCluster)
		(*in).DeepCopyInto(*out)
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(Azure)
		**out = **in
	}
	if in.Gcs != nil {
		in, out := &in.Gcs, &out.Gcs
		*out = new(Gcs)
		**out = **in
	}
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(S3)
		**out = **in
	}
	if in.Swift != nil {
		in, out := &in.Swift, &out.Swift
		*out = new(Swift)
		**out = **in
	}
	if in.Oss != nil {
		in, out := &in.Oss, &out.Oss
		*out = new(Oss)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Storage.
func (in *Storage) DeepCopy() *Storage {
	if in == nil {
		return nil
	}
	out := new(Storage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Swift) DeepCopyInto(out *Swift) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Swift.
func (in *Swift) DeepCopy() *Swift {
	if in == nil {
		return nil
	}
	out := new(Swift)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Trivy) DeepCopyInto(out *Trivy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Trivy.
func (in *Trivy) DeepCopy() *Trivy {
	if in == nil {
		return nil
	}
	out := new(Trivy)
	in.DeepCopyInto(out)
	return out
}
