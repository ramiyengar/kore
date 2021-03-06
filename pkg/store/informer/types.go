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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

// Config is the configuration for the resource informer
type Config struct {
	// Factories is the shared informer factory
	Factories []informers.SharedInformerFactory
	// Resource the kind we are watching
	Resource string
	// AddFunc is called on newly created object
	AddFunc func(schema.GroupVersionResource, metav1.Object)
	// DeleteFunc is called when an object is being removed
	DeleteFunc func(schema.GroupVersionResource, metav1.Object)
	// ErrorFunc is called on an error from the informer
	ErrorFunc func(schema.GroupVersionResource, error)
	// UpdateFunc is called when an object has been updated - old / new
	UpdateFunc func(schema.GroupVersionResource, metav1.Object, metav1.Object)
}

// Informer is the contract for a resource informer
type Informer interface {
	// Stop stops and releases any resources from the informer
	Stop() error
}
