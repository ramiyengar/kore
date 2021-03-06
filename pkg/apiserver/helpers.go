/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/appvia/kore/pkg/kore"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
func returnNotImplemented(req *restful.Request, wr *restful.Response) {
	wr.WriteHeader(http.StatusNotImplemented)
}
*/

// newList provides an api list type
func newList() *List {
	return &List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
	}
}

func makeListWithSize(size int) *List {
	l := newList()
	l.Items = make([]string, size)

	return l
}

// handleErrors is a generic wrapper for handling the error from downstream kore brigde
func handleErrors(req *restful.Request, resp *restful.Response, handler func() error) {
	if err := handler(); err != nil {
		code := http.StatusInternalServerError
		switch err {
		case kore.ErrNotFound:
			code = http.StatusNotFound
		case kore.ErrNotAllowed{}:
			code = http.StatusNotAcceptable
		case kore.ErrUnauthorized:
			code = http.StatusForbidden
		case kore.ErrRequestInvalid:
			code = http.StatusBadRequest
		case io.EOF:
			code = http.StatusBadRequest
		}
		if strings.Contains(err.Error(), "record not found") {
			code = http.StatusNotFound
			err = errors.New("resource not found")
		}
		if kerrors.IsNotFound(err) {
			code = http.StatusNotFound
		}

		e := newError(err.Error()).
			WithCode(code).
			WithVerb(req.Request.Method).
			WithURI(req.Request.RequestURI)

		if err := resp.WriteHeaderAndEntity(code, e); err != nil {
			log.WithError(err).Error("failed to respond to request")
		}
	}
}
