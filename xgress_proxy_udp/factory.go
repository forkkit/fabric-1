/*
	Copyright 2020 NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package xgress_proxy_udp

import (
	"fmt"
	"github.com/netfoundry/ziti-fabric/xgress"
)

func NewFactory(ctrl xgress.CtrlChannel) xgress.Factory {
	return &factory{ctrl: ctrl}
}

func (f *factory) CreateListener(optionsData xgress.OptionsData) (xgress.Listener, error) {
	options := xgress.LoadOptions(optionsData)
	service := ""
	if value, found := optionsData["service"]; found {
		service = value.(string)
	} else {
		return nil, fmt.Errorf("missing 'service' configuration option")
	}
	return newListener(service, f.ctrl, options), nil
}

func (f *factory) CreateDialer(_ xgress.OptionsData) (xgress.Dialer, error) {
	return nil, fmt.Errorf("not implemented")
}

type factory struct {
	ctrl    xgress.CtrlChannel
	options *xgress.Options
}