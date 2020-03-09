// Copyright © 2020 Banzai Cloud
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

package process

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/banzaicloud/pipeline-internal-sdk/.gen/pipeline"
	"github.com/banzaicloud/pipeline-internal-sdk/transport/grpc"
)

type Client struct {
	grpcClient pb.ProcessClient
}

func NewClient(transport *grpc.Transport) (*Client, error) {

	grpcClient := pb.NewProcessClient(transport.ClientConn())

	return &Client{grpcClient: grpcClient}, nil
}

type ProcessEntry struct {
	ID           string
	ParentID     string
	OrgID        uint
	Name         string
	ResourceType ResourceType
	ResourceID   string
	Status       Status
	StartedAt    time.Time
	FinishedAt   *time.Time
}

type ProcessEvent struct {
	ProcessID string
	Name      string
	Log       string
	Timestamp time.Time
}

type ResourceType string
type Status string

const (
	ClusterResourceType ResourceType = "cluster"

	RunningStatus  Status = "running"
	FailedStatus   Status = "failed"
	FinishedStatus Status = "finished"
)

func (c *Client) LogProcess(ctx context.Context, e ProcessEntry) error {

	pe := pb.ProcessEntry{
		Id:           e.ID,
		ParentId:     e.ParentID,
		OrgId:        uint32(e.OrgID),
		Name:         e.Name,
		ResourceType: string(e.ResourceType),
		ResourceId:   e.ResourceID,
		Status:       string(e.Status),
	}

	pe.StartedAt, _ = ptypes.TimestampProto(e.StartedAt)

	if e.FinishedAt != nil {
		finishedAt, _ := ptypes.TimestampProto(*e.FinishedAt)
		pe.FinishedAt = finishedAt
	}

	_, err := c.grpcClient.LogProcess(ctx, &pe)

	return err
}

func (c *Client) LogEvent(ctx context.Context, e ProcessEvent) error {

	pe := pb.ProcessEvent{
		ProcessId: e.ProcessID,
		Name:      e.Name,
		Log:       e.Log,
	}

	pe.Timestamp, _ = ptypes.TimestampProto(e.Timestamp)

	_, err := c.grpcClient.LogEvent(ctx, &pe)

	return err
}
