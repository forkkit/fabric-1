/*
	Copyright NetFoundry, Inc.

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

package handler_mgmt

import (
	"github.com/golang/protobuf/proto"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/fabric/controller/handler_common"
	"github.com/openziti/fabric/controller/network"
	"github.com/openziti/fabric/pb/mgmt_pb"
	"github.com/openziti/foundation/channel2"
)

type listSessionsHandler struct {
	network *network.Network
}

func newListSessionsHandler(network *network.Network) *listSessionsHandler {
	return &listSessionsHandler{network: network}
}

func (h *listSessionsHandler) ContentType() int32 {
	return int32(mgmt_pb.ContentType_ListSessionsRequestType)
}

func (h *listSessionsHandler) HandleReceive(msg *channel2.Message, ch channel2.Channel) {
	list := &mgmt_pb.ListSessionsRequest{}
	if err := proto.Unmarshal(msg.Body, list); err != nil {
		handler_common.SendFailure(msg, ch, err.Error())
		return
	}

	sessions := h.network.GetAllSessions()
	response := &mgmt_pb.ListSessionsResponse{
		Sessions: make([]*mgmt_pb.Session, 0),
	}
	for _, s := range sessions {
		responseSession := &mgmt_pb.Session{
			Id:        s.Id.Token,
			ClientId:  s.ClientId.Token,
			ServiceId: s.Service.Id,
			Circuit:   NewCircuit(s.Circuit),
		}
		response.Sessions = append(response.Sessions, responseSession)
	}
	body, err := proto.Marshal(response)
	if err != nil {
		handler_common.SendFailure(msg, ch, err.Error())
		return
	}

	responseMsg := channel2.NewMessage(int32(mgmt_pb.ContentType_ListSessionsResponseType), body)
	responseMsg.ReplyTo(msg)
	if err := ch.Send(responseMsg); err != nil {
		pfxlog.ContextLogger(ch.Label()).Errorf("unexpected error sending response (%s)", err)
	}
}

func NewCircuit(circuit *network.Circuit) *mgmt_pb.Circuit {
	mgmtCircuit := &mgmt_pb.Circuit{}
	for _, r := range circuit.Path {
		mgmtCircuit.Path = append(mgmtCircuit.Path, r.Id)
	}
	for _, l := range circuit.Links {
		mgmtCircuit.Links = append(mgmtCircuit.Links, l.Id.Token)
	}
	return mgmtCircuit
}
