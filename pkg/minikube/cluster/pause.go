/*
Copyright 2019 The Kubernetes Authors All rights reserved.

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

package cluster

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/minikube/pkg/minikube/command"
	"k8s.io/minikube/pkg/minikube/cruntime"
	"k8s.io/minikube/pkg/minikube/kubelet"
)

// Pause pauses a Kubernetes cluster
func Pause(cr cruntime.Manager, r command.Runner, namespaces []string) ([]string, error) {
	ids := []string{}
	// Disable the kubelet so it does not attempt to restart paused pods
	if err := kubelet.Disable(r); err != nil {
		return ids, errors.Wrap(err, "kubelet disable")
	}
	if err := kubelet.Stop(r); err != nil {
		return ids, errors.Wrap(err, "kubelet stop")
	}
	ids, err := cr.ListContainers(cruntime.ListOptions{State: cruntime.Running, Namespaces: namespaces})
	if err != nil {
		return ids, errors.Wrap(err, "list running")
	}
	if len(ids) == 0 {
		glog.Warningf("no running containers to pause")
		return ids, nil
	}
	return ids, cr.PauseContainers(ids)

}

// Unpause unpauses a Kubernetes cluster
func Unpause(cr cruntime.Manager, r command.Runner, namespaces []string) ([]string, error) {
	ids, err := cr.ListContainers(cruntime.ListOptions{State: cruntime.Paused, Namespaces: namespaces})
	if err != nil {
		return ids, errors.Wrap(err, "list paused")
	}

	if len(ids) == 0 {
		glog.Warningf("no paused containers found")
	} else if err := cr.UnpauseContainers(ids); err != nil {
		return ids, errors.Wrap(err, "unpause")
	}

	if err := kubelet.Enable(r); err != nil {
		return ids, errors.Wrap(err, "kubelet enable")
	}
	if err := kubelet.Start(r); err != nil {
		return ids, errors.Wrap(err, "kubelet start")
	}
	return ids, nil
}
