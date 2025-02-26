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
package admin

import (
	"airman.com/airfk/pkg/server"
	"airman.com/airfk/pkg/service"
	"airman.com/airfk/pkg/types"

	"airman.com/fkExample/node/conf"
)

// Backend interface provides the common API services.
type Backend interface {
	// General API
	IsRunning() bool
	StartWS(c *conf.Config, apis []types.API) error
	StopWS()
	DataDir() string
	WSHandle() *server.Server
	WSEndpoint() string
	Version() string
	Name() string
	Config() interface{}
	RpcAPIs() []types.API
	Services() []service.Service
}
