/*
Copyright 2018 Appvia Ltd <info@appvia.io>

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

package informer

import "github.com/prometheus/client_golang/prometheus"

var (
	addCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "informer_add_total",
			Help: "The number of adds which the informer has seen per resource",
		},
		[]string{"resource"},
	)
	deleteCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "informer_delete_total",
			Help: "The number of deletes which the informer has seen per resource",
		},
		[]string{"resource"},
	)
	errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "informer_error_total",
			Help: "The number of errors which the informer has seen per resource",
		},
		[]string{"resource"},
	)
	updateCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "informer_update_total",
			Help: "The number of updates which the informer has seen per resource",
		},
		[]string{"resource"},
	)
)

func init() {
	prometheus.MustRegister(addCounter)
	prometheus.MustRegister(deleteCounter)
	prometheus.MustRegister(errorCounter)
	prometheus.MustRegister(updateCounter)
}
