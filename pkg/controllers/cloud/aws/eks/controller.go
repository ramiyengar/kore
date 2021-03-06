/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package eks

import (
	"context"
	"time"

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type eksCtrl struct {
	kore.Interface
	// mgr is the controller manager
	mgr manager.Manager
	// stopCh is the stop channel
	stopCh chan struct{}
}

func init() {
	if err := controllers.Register(&eksCtrl{}); err != nil {
		log.WithError(err).Fatal("failed to register controller")
	}
}

// Name returns the name of the controller
func (t *eksCtrl) Name() string {
	return "gke.compute.kore.appvia.io"
}

// Run starts the controller
func (t *eksCtrl) Run(ctx context.Context, cfg *rest.Config, hubi kore.Interface) error {
	logger := log.WithFields(log.Fields{
		"controller": t.Name(),
	})

	mgr, err := manager.New(cfg, controllers.DefaultManagerOptions(t))
	if err != nil {
		logger.WithError(err).Error("trying to create the manager")

		return err
	}
	t.mgr = mgr
	t.Interface = hubi

	// @step: create the controller
	_, err = controllers.NewController(
		"eks-credentials.kore.appvia.io", t.mgr,
		&source.Kind{Type: &aws.AWSCredentials{}},
		&controllers.ReconcileHandler{
			HandlerFunc: t.ReconcileCredentials,
		},
	)
	if err != nil {
		logger.WithError(err).Error("trying to create the aws credentials controller")

		return err
	}

	go func() {
		logger.Info("starting the controller loop")

		for {
			t.stopCh = make(chan struct{})

			if err := mgr.Start(t.stopCh); err != nil {
				logger.WithError(err).Error("failed to start the controller")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// @step: use a routine to catch the stop channel
	go func() {
		<-ctx.Done()
		logger.Info("stopping the teams controller")

		close(t.stopCh)
	}()

	return nil
}

// Stop is responsible for calling a halt on the controller
func (t *eksCtrl) Stop(_ context.Context) error {
	log.WithFields(log.Fields{
		"controller": t.Name(),
	}).Info("attempting to stop the controller")

	return nil
}
