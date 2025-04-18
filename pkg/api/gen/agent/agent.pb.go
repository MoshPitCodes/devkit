// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: agent.proto

package agent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	mi := &file_agent_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{0}
}

type PingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	mi := &file_agent_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StartContainerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TemplateName  string                 `protobuf:"bytes,1,opt,name=template_name,json=templateName,proto3" json:"template_name,omitempty"` // Name of the template being used (for naming)
	Image         string                 `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	Ports         []string               `protobuf:"bytes,3,rep,name=ports,proto3" json:"ports,omitempty"`                                                                                       // Format: ["host:container", ...]
	Volumes       []string               `protobuf:"bytes,4,rep,name=volumes,proto3" json:"volumes,omitempty"`                                                                                   // Format: ["host:container", ...]
	Environment   map[string]string      `protobuf:"bytes,5,rep,name=environment,proto3" json:"environment,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // Key-value pairs
	Command       []string               `protobuf:"bytes,6,rep,name=command,proto3" json:"command,omitempty"`                                                                                   // Command to run
	NetworkMode   string                 `protobuf:"bytes,7,opt,name=network_mode,json=networkMode,proto3" json:"network_mode,omitempty"`                                                        // e.g., "bridge", "host", "none"
	Resources     *Resources             `protobuf:"bytes,8,opt,name=resources,proto3" json:"resources,omitempty"`                                                                               // Resource limits
	RestartPolicy string                 `protobuf:"bytes,9,opt,name=restart_policy,json=restartPolicy,proto3" json:"restart_policy,omitempty"`                                                  // e.g., "no", "always", "on-failure"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartContainerRequest) Reset() {
	*x = StartContainerRequest{}
	mi := &file_agent_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartContainerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartContainerRequest) ProtoMessage() {}

func (x *StartContainerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartContainerRequest.ProtoReflect.Descriptor instead.
func (*StartContainerRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{2}
}

func (x *StartContainerRequest) GetTemplateName() string {
	if x != nil {
		return x.TemplateName
	}
	return ""
}

func (x *StartContainerRequest) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *StartContainerRequest) GetPorts() []string {
	if x != nil {
		return x.Ports
	}
	return nil
}

func (x *StartContainerRequest) GetVolumes() []string {
	if x != nil {
		return x.Volumes
	}
	return nil
}

func (x *StartContainerRequest) GetEnvironment() map[string]string {
	if x != nil {
		return x.Environment
	}
	return nil
}

func (x *StartContainerRequest) GetCommand() []string {
	if x != nil {
		return x.Command
	}
	return nil
}

func (x *StartContainerRequest) GetNetworkMode() string {
	if x != nil {
		return x.NetworkMode
	}
	return ""
}

func (x *StartContainerRequest) GetResources() *Resources {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *StartContainerRequest) GetRestartPolicy() string {
	if x != nil {
		return x.RestartPolicy
	}
	return ""
}

// Nested message for resource limits
type Resources struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Memory        string                 `protobuf:"bytes,1,opt,name=memory,proto3" json:"memory,omitempty"` // e.g., "512m", "1g"
	Cpu           string                 `protobuf:"bytes,2,opt,name=cpu,proto3" json:"cpu,omitempty"`       // e.g., "0.5", "1"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Resources) Reset() {
	*x = Resources{}
	mi := &file_agent_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Resources) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resources) ProtoMessage() {}

func (x *Resources) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resources.ProtoReflect.Descriptor instead.
func (*Resources) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{3}
}

func (x *Resources) GetMemory() string {
	if x != nil {
		return x.Memory
	}
	return ""
}

func (x *Resources) GetCpu() string {
	if x != nil {
		return x.Cpu
	}
	return ""
}

type StartContainerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ContainerId   string                 `protobuf:"bytes,1,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // e.g., "Container devkit-redis started successfully"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartContainerResponse) Reset() {
	*x = StartContainerResponse{}
	mi := &file_agent_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartContainerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartContainerResponse) ProtoMessage() {}

func (x *StartContainerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartContainerResponse.ProtoReflect.Descriptor instead.
func (*StartContainerResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{4}
}

func (x *StartContainerResponse) GetContainerId() string {
	if x != nil {
		return x.ContainerId
	}
	return ""
}

func (x *StartContainerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ListContainersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListContainersRequest) Reset() {
	*x = ListContainersRequest{}
	mi := &file_agent_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListContainersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListContainersRequest) ProtoMessage() {}

func (x *ListContainersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListContainersRequest.ProtoReflect.Descriptor instead.
func (*ListContainersRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{5}
}

type ContainerInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`         // Short container ID
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`     // Container name (e.g., devkit-redis)
	Image         string                 `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`   // Image used
	Status        string                 `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"` // Status (e.g., running, exited)
	Ports         []string               `protobuf:"bytes,5,rep,name=ports,proto3" json:"ports,omitempty"`   // Published ports (e.g., "host:container/proto")
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ContainerInfo) Reset() {
	*x = ContainerInfo{}
	mi := &file_agent_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ContainerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerInfo) ProtoMessage() {}

func (x *ContainerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerInfo.ProtoReflect.Descriptor instead.
func (*ContainerInfo) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{6}
}

func (x *ContainerInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ContainerInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ContainerInfo) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *ContainerInfo) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ContainerInfo) GetPorts() []string {
	if x != nil {
		return x.Ports
	}
	return nil
}

type ListContainersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Containers    []*ContainerInfo       `protobuf:"bytes,1,rep,name=containers,proto3" json:"containers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListContainersResponse) Reset() {
	*x = ListContainersResponse{}
	mi := &file_agent_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListContainersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListContainersResponse) ProtoMessage() {}

func (x *ListContainersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListContainersResponse.ProtoReflect.Descriptor instead.
func (*ListContainersResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{7}
}

func (x *ListContainersResponse) GetContainers() []*ContainerInfo {
	if x != nil {
		return x.Containers
	}
	return nil
}

type StopContainerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // ID or Name of the container to stop
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopContainerRequest) Reset() {
	*x = StopContainerRequest{}
	mi := &file_agent_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopContainerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopContainerRequest) ProtoMessage() {}

func (x *StopContainerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopContainerRequest.ProtoReflect.Descriptor instead.
func (*StopContainerRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{8}
}

func (x *StopContainerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type StopContainerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // e.g., "Container devkit-redis stopped successfully"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopContainerResponse) Reset() {
	*x = StopContainerResponse{}
	mi := &file_agent_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopContainerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopContainerResponse) ProtoMessage() {}

func (x *StopContainerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopContainerResponse.ProtoReflect.Descriptor instead.
func (*StopContainerResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{9}
}

func (x *StopContainerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type RemoveContainerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`        // ID or Name of the container to remove
	Force         bool                   `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"` // Whether to force remove a running container
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveContainerRequest) Reset() {
	*x = RemoveContainerRequest{}
	mi := &file_agent_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveContainerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveContainerRequest) ProtoMessage() {}

func (x *RemoveContainerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveContainerRequest.ProtoReflect.Descriptor instead.
func (*RemoveContainerRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{10}
}

func (x *RemoveContainerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RemoveContainerRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

type RemoveContainerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"` // e.g., "Container devkit-redis removed successfully"
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveContainerResponse) Reset() {
	*x = RemoveContainerResponse{}
	mi := &file_agent_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveContainerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveContainerResponse) ProtoMessage() {}

func (x *RemoveContainerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveContainerResponse.ProtoReflect.Descriptor instead.
func (*RemoveContainerResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{11}
}

func (x *RemoveContainerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type LogsContainerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                  // ID or Name of the container
	Follow        bool                   `protobuf:"varint,2,opt,name=follow,proto3" json:"follow,omitempty"`         // Stream logs continuously
	Tail          string                 `protobuf:"bytes,3,opt,name=tail,proto3" json:"tail,omitempty"`              // Number of lines to show from end (e.g., "100", "all")
	Timestamps    bool                   `protobuf:"varint,4,opt,name=timestamps,proto3" json:"timestamps,omitempty"` // Show timestamps
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogsContainerRequest) Reset() {
	*x = LogsContainerRequest{}
	mi := &file_agent_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogsContainerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsContainerRequest) ProtoMessage() {}

func (x *LogsContainerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsContainerRequest.ProtoReflect.Descriptor instead.
func (*LogsContainerRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{12}
}

func (x *LogsContainerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LogsContainerRequest) GetFollow() bool {
	if x != nil {
		return x.Follow
	}
	return false
}

func (x *LogsContainerRequest) GetTail() string {
	if x != nil {
		return x.Tail
	}
	return ""
}

func (x *LogsContainerRequest) GetTimestamps() bool {
	if x != nil {
		return x.Timestamps
	}
	return false
}

type LogsContainerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Line          []byte                 `protobuf:"bytes,1,opt,name=line,proto3" json:"line,omitempty"` // A single line of log output (stdout or stderr)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogsContainerResponse) Reset() {
	*x = LogsContainerResponse{}
	mi := &file_agent_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogsContainerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsContainerResponse) ProtoMessage() {}

func (x *LogsContainerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsContainerResponse.ProtoReflect.Descriptor instead.
func (*LogsContainerResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{13}
}

func (x *LogsContainerResponse) GetLine() []byte {
	if x != nil {
		return x.Line
	}
	return nil
}

type RefreshTemplatesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTemplatesRequest) Reset() {
	*x = RefreshTemplatesRequest{}
	mi := &file_agent_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTemplatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTemplatesRequest) ProtoMessage() {}

func (x *RefreshTemplatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTemplatesRequest.ProtoReflect.Descriptor instead.
func (*RefreshTemplatesRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{14}
}

type RefreshTemplatesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTemplatesResponse) Reset() {
	*x = RefreshTemplatesResponse{}
	mi := &file_agent_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTemplatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTemplatesResponse) ProtoMessage() {}

func (x *RefreshTemplatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTemplatesResponse.ProtoReflect.Descriptor instead.
func (*RefreshTemplatesResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{15}
}

func (x *RefreshTemplatesResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ShutdownRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShutdownRequest) Reset() {
	*x = ShutdownRequest{}
	mi := &file_agent_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShutdownRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownRequest) ProtoMessage() {}

func (x *ShutdownRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownRequest.ProtoReflect.Descriptor instead.
func (*ShutdownRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{16}
}

type ShutdownResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShutdownResponse) Reset() {
	*x = ShutdownResponse{}
	mi := &file_agent_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShutdownResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownResponse) ProtoMessage() {}

func (x *ShutdownResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownResponse.ProtoReflect.Descriptor instead.
func (*ShutdownResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{17}
}

func (x *ShutdownResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Request message for stopping the agent.
type StopAgentRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// If true, perform an immediate shutdown rather than graceful.
	Force         bool `protobuf:"varint,1,opt,name=force,proto3" json:"force,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopAgentRequest) Reset() {
	*x = StopAgentRequest{}
	mi := &file_agent_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopAgentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopAgentRequest) ProtoMessage() {}

func (x *StopAgentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopAgentRequest.ProtoReflect.Descriptor instead.
func (*StopAgentRequest) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{18}
}

func (x *StopAgentRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

// Response message for stopping the agent.
type StopAgentResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Confirmation message.
	Message       string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopAgentResponse) Reset() {
	*x = StopAgentResponse{}
	mi := &file_agent_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopAgentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopAgentResponse) ProtoMessage() {}

func (x *StopAgentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_agent_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopAgentResponse.ProtoReflect.Descriptor instead.
func (*StopAgentResponse) Descriptor() ([]byte, []int) {
	return file_agent_proto_rawDescGZIP(), []int{19}
}

func (x *StopAgentResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_agent_proto protoreflect.FileDescriptor

const file_agent_proto_rawDesc = "" +
	"\n" +
	"\vagent.proto\x12\x05agent\"\r\n" +
	"\vPingRequest\"(\n" +
	"\fPingResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"\xa7\x03\n" +
	"\x15StartContainerRequest\x12#\n" +
	"\rtemplate_name\x18\x01 \x01(\tR\ftemplateName\x12\x14\n" +
	"\x05image\x18\x02 \x01(\tR\x05image\x12\x14\n" +
	"\x05ports\x18\x03 \x03(\tR\x05ports\x12\x18\n" +
	"\avolumes\x18\x04 \x03(\tR\avolumes\x12O\n" +
	"\venvironment\x18\x05 \x03(\v2-.agent.StartContainerRequest.EnvironmentEntryR\venvironment\x12\x18\n" +
	"\acommand\x18\x06 \x03(\tR\acommand\x12!\n" +
	"\fnetwork_mode\x18\a \x01(\tR\vnetworkMode\x12.\n" +
	"\tresources\x18\b \x01(\v2\x10.agent.ResourcesR\tresources\x12%\n" +
	"\x0erestart_policy\x18\t \x01(\tR\rrestartPolicy\x1a>\n" +
	"\x10EnvironmentEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"5\n" +
	"\tResources\x12\x16\n" +
	"\x06memory\x18\x01 \x01(\tR\x06memory\x12\x10\n" +
	"\x03cpu\x18\x02 \x01(\tR\x03cpu\"U\n" +
	"\x16StartContainerResponse\x12!\n" +
	"\fcontainer_id\x18\x01 \x01(\tR\vcontainerId\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"\x17\n" +
	"\x15ListContainersRequest\"w\n" +
	"\rContainerInfo\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x14\n" +
	"\x05image\x18\x03 \x01(\tR\x05image\x12\x16\n" +
	"\x06status\x18\x04 \x01(\tR\x06status\x12\x14\n" +
	"\x05ports\x18\x05 \x03(\tR\x05ports\"N\n" +
	"\x16ListContainersResponse\x124\n" +
	"\n" +
	"containers\x18\x01 \x03(\v2\x14.agent.ContainerInfoR\n" +
	"containers\"&\n" +
	"\x14StopContainerRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"1\n" +
	"\x15StopContainerResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\">\n" +
	"\x16RemoveContainerRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05force\x18\x02 \x01(\bR\x05force\"3\n" +
	"\x17RemoveContainerResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"r\n" +
	"\x14LogsContainerRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x16\n" +
	"\x06follow\x18\x02 \x01(\bR\x06follow\x12\x12\n" +
	"\x04tail\x18\x03 \x01(\tR\x04tail\x12\x1e\n" +
	"\n" +
	"timestamps\x18\x04 \x01(\bR\n" +
	"timestamps\"+\n" +
	"\x15LogsContainerResponse\x12\x12\n" +
	"\x04line\x18\x01 \x01(\fR\x04line\"\x19\n" +
	"\x17RefreshTemplatesRequest\"4\n" +
	"\x18RefreshTemplatesResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"\x11\n" +
	"\x0fShutdownRequest\",\n" +
	"\x10ShutdownResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"(\n" +
	"\x10StopAgentRequest\x12\x14\n" +
	"\x05force\x18\x01 \x01(\bR\x05force\"-\n" +
	"\x11StopAgentResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage2\x9b\x05\n" +
	"\fAgentService\x12/\n" +
	"\x04Ping\x12\x12.agent.PingRequest\x1a\x13.agent.PingResponse\x12M\n" +
	"\x0eStartContainer\x12\x1c.agent.StartContainerRequest\x1a\x1d.agent.StartContainerResponse\x12M\n" +
	"\x0eListContainers\x12\x1c.agent.ListContainersRequest\x1a\x1d.agent.ListContainersResponse\x12J\n" +
	"\rStopContainer\x12\x1b.agent.StopContainerRequest\x1a\x1c.agent.StopContainerResponse\x12P\n" +
	"\x0fRemoveContainer\x12\x1d.agent.RemoveContainerRequest\x1a\x1e.agent.RemoveContainerResponse\x12L\n" +
	"\rLogsContainer\x12\x1b.agent.LogsContainerRequest\x1a\x1c.agent.LogsContainerResponse0\x01\x12S\n" +
	"\x10RefreshTemplates\x12\x1e.agent.RefreshTemplatesRequest\x1a\x1f.agent.RefreshTemplatesResponse\x12;\n" +
	"\bShutdown\x12\x16.agent.ShutdownRequest\x1a\x17.agent.ShutdownResponse\x12>\n" +
	"\tStopAgent\x12\x17.agent.StopAgentRequest\x1a\x18.agent.StopAgentResponseB\tZ\a.;agentb\x06proto3"

var (
	file_agent_proto_rawDescOnce sync.Once
	file_agent_proto_rawDescData []byte
)

func file_agent_proto_rawDescGZIP() []byte {
	file_agent_proto_rawDescOnce.Do(func() {
		file_agent_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_agent_proto_rawDesc), len(file_agent_proto_rawDesc)))
	})
	return file_agent_proto_rawDescData
}

var file_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 21)
var file_agent_proto_goTypes = []any{
	(*PingRequest)(nil),              // 0: agent.PingRequest
	(*PingResponse)(nil),             // 1: agent.PingResponse
	(*StartContainerRequest)(nil),    // 2: agent.StartContainerRequest
	(*Resources)(nil),                // 3: agent.Resources
	(*StartContainerResponse)(nil),   // 4: agent.StartContainerResponse
	(*ListContainersRequest)(nil),    // 5: agent.ListContainersRequest
	(*ContainerInfo)(nil),            // 6: agent.ContainerInfo
	(*ListContainersResponse)(nil),   // 7: agent.ListContainersResponse
	(*StopContainerRequest)(nil),     // 8: agent.StopContainerRequest
	(*StopContainerResponse)(nil),    // 9: agent.StopContainerResponse
	(*RemoveContainerRequest)(nil),   // 10: agent.RemoveContainerRequest
	(*RemoveContainerResponse)(nil),  // 11: agent.RemoveContainerResponse
	(*LogsContainerRequest)(nil),     // 12: agent.LogsContainerRequest
	(*LogsContainerResponse)(nil),    // 13: agent.LogsContainerResponse
	(*RefreshTemplatesRequest)(nil),  // 14: agent.RefreshTemplatesRequest
	(*RefreshTemplatesResponse)(nil), // 15: agent.RefreshTemplatesResponse
	(*ShutdownRequest)(nil),          // 16: agent.ShutdownRequest
	(*ShutdownResponse)(nil),         // 17: agent.ShutdownResponse
	(*StopAgentRequest)(nil),         // 18: agent.StopAgentRequest
	(*StopAgentResponse)(nil),        // 19: agent.StopAgentResponse
	nil,                              // 20: agent.StartContainerRequest.EnvironmentEntry
}
var file_agent_proto_depIdxs = []int32{
	20, // 0: agent.StartContainerRequest.environment:type_name -> agent.StartContainerRequest.EnvironmentEntry
	3,  // 1: agent.StartContainerRequest.resources:type_name -> agent.Resources
	6,  // 2: agent.ListContainersResponse.containers:type_name -> agent.ContainerInfo
	0,  // 3: agent.AgentService.Ping:input_type -> agent.PingRequest
	2,  // 4: agent.AgentService.StartContainer:input_type -> agent.StartContainerRequest
	5,  // 5: agent.AgentService.ListContainers:input_type -> agent.ListContainersRequest
	8,  // 6: agent.AgentService.StopContainer:input_type -> agent.StopContainerRequest
	10, // 7: agent.AgentService.RemoveContainer:input_type -> agent.RemoveContainerRequest
	12, // 8: agent.AgentService.LogsContainer:input_type -> agent.LogsContainerRequest
	14, // 9: agent.AgentService.RefreshTemplates:input_type -> agent.RefreshTemplatesRequest
	16, // 10: agent.AgentService.Shutdown:input_type -> agent.ShutdownRequest
	18, // 11: agent.AgentService.StopAgent:input_type -> agent.StopAgentRequest
	1,  // 12: agent.AgentService.Ping:output_type -> agent.PingResponse
	4,  // 13: agent.AgentService.StartContainer:output_type -> agent.StartContainerResponse
	7,  // 14: agent.AgentService.ListContainers:output_type -> agent.ListContainersResponse
	9,  // 15: agent.AgentService.StopContainer:output_type -> agent.StopContainerResponse
	11, // 16: agent.AgentService.RemoveContainer:output_type -> agent.RemoveContainerResponse
	13, // 17: agent.AgentService.LogsContainer:output_type -> agent.LogsContainerResponse
	15, // 18: agent.AgentService.RefreshTemplates:output_type -> agent.RefreshTemplatesResponse
	17, // 19: agent.AgentService.Shutdown:output_type -> agent.ShutdownResponse
	19, // 20: agent.AgentService.StopAgent:output_type -> agent.StopAgentResponse
	12, // [12:21] is the sub-list for method output_type
	3,  // [3:12] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_agent_proto_init() }
func file_agent_proto_init() {
	if File_agent_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_agent_proto_rawDesc), len(file_agent_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   21,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_proto_goTypes,
		DependencyIndexes: file_agent_proto_depIdxs,
		MessageInfos:      file_agent_proto_msgTypes,
	}.Build()
	File_agent_proto = out.File
	file_agent_proto_goTypes = nil
	file_agent_proto_depIdxs = nil
}
