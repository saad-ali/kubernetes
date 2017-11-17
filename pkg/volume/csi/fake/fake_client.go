/*
Copyright 2017 The Kubernetes Authors.

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

package fake

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"

	csipb "github.com/container-storage-interface/spec/lib/go/csi"
	grpctx "golang.org/x/net/context"
)

type FakeIdentityClient struct {
	nextErr error
}

func NewIdentityClient() *FakeIdentityClient {
	return &FakeIdentityClient{}
}

func (f *FakeIdentityClient) SetNextError(err error) {
	f.nextErr = err
}

func (f *FakeIdentityClient) GetSupportedVersions(ctx grpctx.Context, req *csipb.GetSupportedVersionsRequest, opts ...grpc.CallOption) (*csipb.GetSupportedVersionsResponse, error) {
	// short circuit with an error
	if f.nextErr != nil {
		return nil, f.nextErr
	}

	rsp := &csipb.GetSupportedVersionsResponse{
		SupportedVersions: []*csipb.Version{
			{Major: 0, Minor: 0, Patch: 1},
			{Major: 0, Minor: 1, Patch: 0},
			{Major: 1, Minor: 0, Patch: 0},
			{Major: 1, Minor: 0, Patch: 1},
			{Major: 1, Minor: 1, Patch: 1},
		},
	}
	return rsp, nil
}

func (f *FakeIdentityClient) GetPluginInfo(ctx context.Context, in *csipb.GetPluginInfoRequest, opts ...grpc.CallOption) (*csipb.GetPluginInfoResponse, error) {
	return nil, nil
}

type FakeNodeClient struct {
	nodePublishedVolumes map[string]string
	nextErr              error
}

func NewNodeClient() *FakeNodeClient {
	return &FakeNodeClient{nodePublishedVolumes: make(map[string]string)}
}

func (f *FakeNodeClient) SetNextError(err error) {
	f.nextErr = err
}

func (f *FakeNodeClient) GetNodePublishedVolumes() map[string]string {
	return f.nodePublishedVolumes
}

func (f *FakeNodeClient) NodePublishVolume(ctx grpctx.Context, req *csipb.NodePublishVolumeRequest, opts ...grpc.CallOption) (*csipb.NodePublishVolumeResponse, error) {

	if f.nextErr != nil {
		return nil, f.nextErr
	}

	if req.GetVolumeId() == "" {
		return nil, errors.New("missing volume id")
	}
	if req.GetTargetPath() == "" {
		return nil, errors.New("missing target path")
	}
	fsTypes := "ext4|xfs|zfs"
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	if !strings.Contains(fsTypes, fsType) {
		return nil, errors.New("invlid fstype")
	}
	f.nodePublishedVolumes[req.GetVolumeId()] = req.GetTargetPath()
	return &csipb.NodePublishVolumeResponse{}, nil
}

func (f *FakeNodeClient) NodeUnpublishVolume(ctx context.Context, req *csipb.NodeUnpublishVolumeRequest, opts ...grpc.CallOption) (*csipb.NodeUnpublishVolumeResponse, error) {
	if f.nextErr != nil {
		return nil, f.nextErr
	}

	if req.GetVolumeId() == "" {
		return nil, errors.New("missing volume id")
	}
	if req.GetTargetPath() == "" {
		return nil, errors.New("missing target path")
	}
	delete(f.nodePublishedVolumes, req.GetVolumeId())
	return &csipb.NodeUnpublishVolumeResponse{}, nil
}

func (f *FakeNodeClient) GetNodeID(ctx context.Context, in *csipb.GetNodeIDRequest, opts ...grpc.CallOption) (*csipb.GetNodeIDResponse, error) {
	return nil, nil
}
func (f *FakeNodeClient) NodeProbe(ctx context.Context, in *csipb.NodeProbeRequest, opts ...grpc.CallOption) (*csipb.NodeProbeResponse, error) {
	return nil, nil
}
func (f *FakeNodeClient) NodeGetCapabilities(ctx context.Context, in *csipb.NodeGetCapabilitiesRequest, opts ...grpc.CallOption) (*csipb.NodeGetCapabilitiesResponse, error) {
	return nil, nil
}

type FakeControllerClient struct {
	nextCapabilities []*csipb.ControllerServiceCapability
	nextErr          error
}

func NewControllerClient() *FakeControllerClient {
	return &FakeControllerClient{}
}

func (f *FakeControllerClient) SetNextError(err error) {
	f.nextErr = err
}

func (f *FakeControllerClient) SetNextCapabilities(caps []*csipb.ControllerServiceCapability) {
	f.nextCapabilities = caps
}

func (f *FakeControllerClient) ControllerGetCapabilities(ctx context.Context, in *csipb.ControllerGetCapabilitiesRequest, opts ...grpc.CallOption) (*csipb.ControllerGetCapabilitiesResponse, error) {
	if f.nextErr != nil {
		return nil, f.nextErr
	}

	if f.nextCapabilities == nil {
		f.nextCapabilities = []*csipb.ControllerServiceCapability{
			{&csipb.ControllerServiceCapability_Rpc{&csipb.ControllerServiceCapability_RPC{csipb.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME}}},
		}
	}
	return &csipb.ControllerGetCapabilitiesResponse{
		Capabilities: f.nextCapabilities,
	}, nil
}

func (f *FakeControllerClient) CreateVolume(ctx context.Context, in *csipb.CreateVolumeRequest, opts ...grpc.CallOption) (*csipb.CreateVolumeResponse, error) {
	return nil, nil
}

func (f *FakeControllerClient) DeleteVolume(ctx context.Context, in *csipb.DeleteVolumeRequest, opts ...grpc.CallOption) (*csipb.DeleteVolumeResponse, error) {
	return nil, nil
}

func (f *FakeControllerClient) ControllerPublishVolume(ctx context.Context, in *csipb.ControllerPublishVolumeRequest, opts ...grpc.CallOption) (*csipb.ControllerPublishVolumeResponse, error) {
	return nil, nil
}
func (f *FakeControllerClient) ControllerUnpublishVolume(ctx context.Context, in *csipb.ControllerUnpublishVolumeRequest, opts ...grpc.CallOption) (*csipb.ControllerUnpublishVolumeResponse, error) {
	return nil, nil
}

func (f *FakeControllerClient) ValidateVolumeCapabilities(ctx context.Context, in *csipb.ValidateVolumeCapabilitiesRequest, opts ...grpc.CallOption) (*csipb.ValidateVolumeCapabilitiesResponse, error) {
	return nil, nil
}
func (f *FakeControllerClient) ListVolumes(ctx context.Context, in *csipb.ListVolumesRequest, opts ...grpc.CallOption) (*csipb.ListVolumesResponse, error) {
	return nil, nil
}

func (f *FakeControllerClient) GetCapacity(ctx context.Context, in *csipb.GetCapacityRequest, opts ...grpc.CallOption) (*csipb.GetCapacityResponse, error) {
	return nil, nil
}
func (f *FakeControllerClient) ControllerProbe(ctx context.Context, in *csipb.ControllerProbeRequest, opts ...grpc.CallOption) (*csipb.ControllerProbeResponse, error) {
	return nil, nil
}
