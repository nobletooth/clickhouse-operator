// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package chi

import (
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/altinity/clickhouse-operator/pkg/announcer"
	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/controller"
	chiLabeler "github.com/altinity/clickhouse-operator/pkg/model/chi/tags/labeler"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

// getPodsIPs gets all pod IPs
func (c *Controller) getPodsIPs(obj interface{}) (ips []string) {
	l := log.V(3).M(obj).F()

	l.S().Info("looking for pods IPs")
	defer l.E().Info("looking for pods IPs")

	for _, pod := range c.kube.Pod().GetAll(obj) {
		if ip := pod.Status.PodIP; ip == "" {
			l.Warning("Pod NO IP address found. Pod: %s", util.NamespacedName(pod))
		} else {
			ips = append(ips, ip)
			l.Info("Pod IP address found. Pod: %s IP: %s", util.NamespacedName(pod), ip)
		}
	}
	return ips
}

// GetCHI gets CHI by any object from a CHI
func (c *Controller) GetCHI(obj meta.Object) (*api.ClickHouseInstallation, error) {
	switch obj.(type) {
	case *api.ClickHouseInstallation:
		cr, err := c.kube.CR().Get(controller.NewContext(), obj.GetNamespace(), obj.GetName())
		if cr == nil {
			return nil, err
		}
		return cr.(*api.ClickHouseInstallation), err
	default:
		return c.GetCHIByObject(obj)
	}
}

// GetCHIByObject gets CHI by namespaced name
func (c *Controller) GetCHIByObject(obj meta.Object) (*api.ClickHouseInstallation, error) {
	crName, err := chiLabeler.New(nil).GetCRNameFromObjectMeta(obj)
	if err != nil {
		return nil, fmt.Errorf("unable to find CR by name: '%s'. More info: %v", obj.GetName(), err)
	}

	cr, err := c.kube.CR().Get(controller.NewContext(), obj.GetNamespace(), crName)
	if cr == nil {
		return nil, err
	}
	return cr.(*api.ClickHouseInstallation), err
}
