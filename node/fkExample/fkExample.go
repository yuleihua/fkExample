// Copyright 2018 The huayulei_2003@hotmail.com Authors
// This file is part of the airfk library.
//
// The airfk library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The airfk library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the airfk library. If not, see <http://www.gnu.org/licenses/>.
package fkExample

import (
	"context"
	"path/filepath"
	"time"

	"airman.com/airfk/pkg/common"
	"airman.com/airfk/pkg/server"

	cmn "airman.com/fkExample/node/common"
)

// Backend interface provides the common API services.
type Backend interface {
	// General API
	DataDir() string
}

// PrivateFkExampleAPI is the collection of service API methods exposed only
// over a secure RPC channel.
type PrivateFkExampleAPI struct {
	manager *Manager
}

// NewPrivateFkExampleAPI creates a new API definition for the private admin methods
// of the node itself.
func NewPrivateFkExampleAPI(manager *Manager) *PrivateFkExampleAPI {
	return &PrivateFkExampleAPI{manager: manager}
}

// RootDir return node root path.
func (api *PrivateFkExampleAPI) RootDir() (string, error) {
	return api.manager.backend.DataDir(), nil
}

// PublicFkExampleAPI
type PublicFkExampleAPI struct {
	manager *Manager
}

// NewPublicFkExampleAPI create PublicTaskAPI.
func NewPublicFkExampleAPI(manager *Manager) *PublicFkExampleAPI {
	return &PublicFkExampleAPI{
		manager: manager,
	}
}

// WriteFile create file and write file.
func (api *PublicFkExampleAPI) WriteFile(file string, text string) error {
	if api.manager == nil || file == "" {
		return nil
	}

	// create file name
	fileName := filepath.Join(api.manager.backend.DataDir(), file)
	if err := common.WriteFile(fileName, []byte(text)); err != nil {
		return err
	}
	r := cmn.NewResultWithEnd(time.Now().String(), time.Now().Unix(), time.Now().Unix(), "write ok", []byte(file))
	api.manager.resultsFeed.Send([]cmn.Result{*r})
	return nil
}

// ReadFile read file content.
func (api *PublicFkExampleAPI) ReadFile(file string) (string, error) {
	if api.manager == nil || file == "" {
		return "", cmn.ErrInvalidParameter
	}

	// create file name
	fileName := filepath.Join(api.manager.backend.DataDir(), file)
	text, err := common.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	r := cmn.NewResultWithEnd(time.Now().String(), time.Now().Unix(), time.Now().Unix(), "read ok", []byte(file))
	api.manager.resultsFeed.Send([]cmn.Result{*r})
	return string(text), nil
}

// ResultsSubscription creates a subscription that is result of task.
func (api *PublicFkExampleAPI) ResultsSubscription(ctx context.Context) (*server.Subscription, error) {
	notifier, supported := server.NotifierFromContext(ctx)
	if !supported {
		return &server.Subscription{}, server.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		results := make(chan []cmn.Result, 128)
		resultsSub := api.manager.es.SubscribeResultTask(results)

		for {
			select {
			case rs := <-results:
				for _, h := range rs {
					notifier.Notify(rpcSub.ID, h)
				}
			case <-rpcSub.Err():
				resultsSub.Unsubscribe()
				return
			case <-notifier.Closed():
				resultsSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}
