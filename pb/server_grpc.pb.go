// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.21.5
// source: server.proto

//protoc --go_out=./pb --go-grpc_out=./pb server.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	NotificationService_Subscribe_FullMethodName = "/notification.NotificationService/Subscribe"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 定义通知服务
type NotificationServiceClient interface {
	// 订阅更新，服务器端流式 RPC
	Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (NotificationService_SubscribeClient, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (NotificationService_SubscribeClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &NotificationService_ServiceDesc.Streams[0], NotificationService_Subscribe_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &notificationServiceSubscribeClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NotificationService_SubscribeClient interface {
	Recv() (*Notification, error)
	grpc.ClientStream
}

type notificationServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *notificationServiceSubscribeClient) Recv() (*Notification, error) {
	m := new(Notification)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility
//
// 定义通知服务
type NotificationServiceServer interface {
	// 订阅更新，服务器端流式 RPC
	Subscribe(*SubscriptionRequest, NotificationService_SubscribeServer) error
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServiceServer struct {
}

func (UnimplementedNotificationServiceServer) Subscribe(*SubscriptionRequest, NotificationService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NotificationServiceServer).Subscribe(m, &notificationServiceSubscribeServer{ServerStream: stream})
}

type NotificationService_SubscribeServer interface {
	Send(*Notification) error
	grpc.ServerStream
}

type notificationServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *notificationServiceSubscribeServer) Send(m *Notification) error {
	return x.ServerStream.SendMsg(m)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _NotificationService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "server.proto",
}

const (
	Service_Login_FullMethodName                 = "/notification.Service/Login"
	Service_Register_FullMethodName              = "/notification.Service/Register"
	Service_GetUserInfo_FullMethodName           = "/notification.Service/GetUserInfo"
	Service_ModUserInfo_FullMethodName           = "/notification.Service/ModUserInfo"
	Service_GetTaskListAll_FullMethodName        = "/notification.Service/GetTaskListAll"
	Service_GetTaskListOne_FullMethodName        = "/notification.Service/GetTaskListOne"
	Service_ImportXLSToTaskTable_FullMethodName  = "/notification.Service/ImportXLSToTaskTable"
	Service_DelTask_FullMethodName               = "/notification.Service/DelTask"
	Service_ModTask_FullMethodName               = "/notification.Service/ModTask"
	Service_AddTask_FullMethodName               = "/notification.Service/AddTask"
	Service_QueryTaskWithSQL_FullMethodName      = "/notification.Service/QueryTaskWithSQL"
	Service_QueryTaskWithField_FullMethodName    = "/notification.Service/QueryTaskWithField"
	Service_GetPatchsAll_FullMethodName          = "/notification.Service/GetPatchsAll"
	Service_GetOnePatchs_FullMethodName          = "/notification.Service/GetOnePatchs"
	Service_GetPatchsByState_FullMethodName      = "/notification.Service/GetPatchsByState"
	Service_DelPatch_FullMethodName              = "/notification.Service/DelPatch"
	Service_ImportXLSToPatchTable_FullMethodName = "/notification.Service/ImportXLSToPatchTable"
	Service_ModPatch_FullMethodName              = "/notification.Service/ModPatch"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 定义服务
type ServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error)
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoReply, error)
	ModUserInfo(ctx context.Context, in *ModUserInfoRequest, opts ...grpc.CallOption) (*ModUserInfoReply, error)
	// 修改单/任务
	GetTaskListAll(ctx context.Context, in *GetTaskListAllRequest, opts ...grpc.CallOption) (*GetTaskListAllReply, error)
	GetTaskListOne(ctx context.Context, in *GetTaskListOneRequest, opts ...grpc.CallOption) (*GetTaskListOneReply, error)
	ImportXLSToTaskTable(ctx context.Context, in *ImportToTaskListRequest, opts ...grpc.CallOption) (*ImportToTaskListReply, error)
	// CURD
	DelTask(ctx context.Context, in *DelTaskRequest, opts ...grpc.CallOption) (*DelTaskReply, error)
	ModTask(ctx context.Context, in *ModTaskRequest, opts ...grpc.CallOption) (*ModTaskReply, error)
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskReply, error)
	QueryTaskWithSQL(ctx context.Context, in *QueryTaskWithSQLRequest, opts ...grpc.CallOption) (*QueryTaskWithSQLReply, error)
	QueryTaskWithField(ctx context.Context, in *QueryTaskWithFieldRequest, opts ...grpc.CallOption) (*QueryTaskWithFieldReply, error)
	// 补丁
	GetPatchsAll(ctx context.Context, in *GetPatchsAllRequest, opts ...grpc.CallOption) (*GetPatchsAllReply, error)
	GetOnePatchs(ctx context.Context, in *GetOnePatchsRequest, opts ...grpc.CallOption) (*GetOnePatchsReply, error)
	GetPatchsByState(ctx context.Context, in *GetPatchsByStateRequest, opts ...grpc.CallOption) (*GetPatchsByStateReply, error)
	DelPatch(ctx context.Context, in *DelPatchRequest, opts ...grpc.CallOption) (*DelPatchReply, error)
	ImportXLSToPatchTable(ctx context.Context, in *ImportXLSToPatchRequest, opts ...grpc.CallOption) (*ImportXLSToPatchReply, error)
	ModPatch(ctx context.Context, in *ModPatchRequest, opts ...grpc.CallOption) (*ModPatchReply, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, Service_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterReply)
	err := c.cc.Invoke(ctx, Service_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInfoReply)
	err := c.cc.Invoke(ctx, Service_GetUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ModUserInfo(ctx context.Context, in *ModUserInfoRequest, opts ...grpc.CallOption) (*ModUserInfoReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModUserInfoReply)
	err := c.cc.Invoke(ctx, Service_ModUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetTaskListAll(ctx context.Context, in *GetTaskListAllRequest, opts ...grpc.CallOption) (*GetTaskListAllReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskListAllReply)
	err := c.cc.Invoke(ctx, Service_GetTaskListAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetTaskListOne(ctx context.Context, in *GetTaskListOneRequest, opts ...grpc.CallOption) (*GetTaskListOneReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskListOneReply)
	err := c.cc.Invoke(ctx, Service_GetTaskListOne_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ImportXLSToTaskTable(ctx context.Context, in *ImportToTaskListRequest, opts ...grpc.CallOption) (*ImportToTaskListReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ImportToTaskListReply)
	err := c.cc.Invoke(ctx, Service_ImportXLSToTaskTable_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DelTask(ctx context.Context, in *DelTaskRequest, opts ...grpc.CallOption) (*DelTaskReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DelTaskReply)
	err := c.cc.Invoke(ctx, Service_DelTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ModTask(ctx context.Context, in *ModTaskRequest, opts ...grpc.CallOption) (*ModTaskReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModTaskReply)
	err := c.cc.Invoke(ctx, Service_ModTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddTaskReply)
	err := c.cc.Invoke(ctx, Service_AddTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryTaskWithSQL(ctx context.Context, in *QueryTaskWithSQLRequest, opts ...grpc.CallOption) (*QueryTaskWithSQLReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryTaskWithSQLReply)
	err := c.cc.Invoke(ctx, Service_QueryTaskWithSQL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryTaskWithField(ctx context.Context, in *QueryTaskWithFieldRequest, opts ...grpc.CallOption) (*QueryTaskWithFieldReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryTaskWithFieldReply)
	err := c.cc.Invoke(ctx, Service_QueryTaskWithField_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetPatchsAll(ctx context.Context, in *GetPatchsAllRequest, opts ...grpc.CallOption) (*GetPatchsAllReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPatchsAllReply)
	err := c.cc.Invoke(ctx, Service_GetPatchsAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetOnePatchs(ctx context.Context, in *GetOnePatchsRequest, opts ...grpc.CallOption) (*GetOnePatchsReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOnePatchsReply)
	err := c.cc.Invoke(ctx, Service_GetOnePatchs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetPatchsByState(ctx context.Context, in *GetPatchsByStateRequest, opts ...grpc.CallOption) (*GetPatchsByStateReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPatchsByStateReply)
	err := c.cc.Invoke(ctx, Service_GetPatchsByState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DelPatch(ctx context.Context, in *DelPatchRequest, opts ...grpc.CallOption) (*DelPatchReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DelPatchReply)
	err := c.cc.Invoke(ctx, Service_DelPatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ImportXLSToPatchTable(ctx context.Context, in *ImportXLSToPatchRequest, opts ...grpc.CallOption) (*ImportXLSToPatchReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ImportXLSToPatchReply)
	err := c.cc.Invoke(ctx, Service_ImportXLSToPatchTable_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ModPatch(ctx context.Context, in *ModPatchRequest, opts ...grpc.CallOption) (*ModPatchReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModPatchReply)
	err := c.cc.Invoke(ctx, Service_ModPatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
//
// 定义服务
type ServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoReply, error)
	ModUserInfo(context.Context, *ModUserInfoRequest) (*ModUserInfoReply, error)
	// 修改单/任务
	GetTaskListAll(context.Context, *GetTaskListAllRequest) (*GetTaskListAllReply, error)
	GetTaskListOne(context.Context, *GetTaskListOneRequest) (*GetTaskListOneReply, error)
	ImportXLSToTaskTable(context.Context, *ImportToTaskListRequest) (*ImportToTaskListReply, error)
	// CURD
	DelTask(context.Context, *DelTaskRequest) (*DelTaskReply, error)
	ModTask(context.Context, *ModTaskRequest) (*ModTaskReply, error)
	AddTask(context.Context, *AddTaskRequest) (*AddTaskReply, error)
	QueryTaskWithSQL(context.Context, *QueryTaskWithSQLRequest) (*QueryTaskWithSQLReply, error)
	QueryTaskWithField(context.Context, *QueryTaskWithFieldRequest) (*QueryTaskWithFieldReply, error)
	// 补丁
	GetPatchsAll(context.Context, *GetPatchsAllRequest) (*GetPatchsAllReply, error)
	GetOnePatchs(context.Context, *GetOnePatchsRequest) (*GetOnePatchsReply, error)
	GetPatchsByState(context.Context, *GetPatchsByStateRequest) (*GetPatchsByStateReply, error)
	DelPatch(context.Context, *DelPatchRequest) (*DelPatchReply, error)
	ImportXLSToPatchTable(context.Context, *ImportXLSToPatchRequest) (*ImportXLSToPatchReply, error)
	ModPatch(context.Context, *ModPatchRequest) (*ModPatchReply, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) Login(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedServiceServer) Register(context.Context, *RegisterRequest) (*RegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedServiceServer) GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedServiceServer) ModUserInfo(context.Context, *ModUserInfoRequest) (*ModUserInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModUserInfo not implemented")
}
func (UnimplementedServiceServer) GetTaskListAll(context.Context, *GetTaskListAllRequest) (*GetTaskListAllReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskListAll not implemented")
}
func (UnimplementedServiceServer) GetTaskListOne(context.Context, *GetTaskListOneRequest) (*GetTaskListOneReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskListOne not implemented")
}
func (UnimplementedServiceServer) ImportXLSToTaskTable(context.Context, *ImportToTaskListRequest) (*ImportToTaskListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportXLSToTaskTable not implemented")
}
func (UnimplementedServiceServer) DelTask(context.Context, *DelTaskRequest) (*DelTaskReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelTask not implemented")
}
func (UnimplementedServiceServer) ModTask(context.Context, *ModTaskRequest) (*ModTaskReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModTask not implemented")
}
func (UnimplementedServiceServer) AddTask(context.Context, *AddTaskRequest) (*AddTaskReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTask not implemented")
}
func (UnimplementedServiceServer) QueryTaskWithSQL(context.Context, *QueryTaskWithSQLRequest) (*QueryTaskWithSQLReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTaskWithSQL not implemented")
}
func (UnimplementedServiceServer) QueryTaskWithField(context.Context, *QueryTaskWithFieldRequest) (*QueryTaskWithFieldReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTaskWithField not implemented")
}
func (UnimplementedServiceServer) GetPatchsAll(context.Context, *GetPatchsAllRequest) (*GetPatchsAllReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPatchsAll not implemented")
}
func (UnimplementedServiceServer) GetOnePatchs(context.Context, *GetOnePatchsRequest) (*GetOnePatchsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOnePatchs not implemented")
}
func (UnimplementedServiceServer) GetPatchsByState(context.Context, *GetPatchsByStateRequest) (*GetPatchsByStateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPatchsByState not implemented")
}
func (UnimplementedServiceServer) DelPatch(context.Context, *DelPatchRequest) (*DelPatchReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelPatch not implemented")
}
func (UnimplementedServiceServer) ImportXLSToPatchTable(context.Context, *ImportXLSToPatchRequest) (*ImportXLSToPatchReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportXLSToPatchTable not implemented")
}
func (UnimplementedServiceServer) ModPatch(context.Context, *ModPatchRequest) (*ModPatchReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModPatch not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetUserInfo(ctx, req.(*GetUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ModUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ModUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ModUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ModUserInfo(ctx, req.(*ModUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetTaskListAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskListAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetTaskListAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetTaskListAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetTaskListAll(ctx, req.(*GetTaskListAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetTaskListOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskListOneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetTaskListOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetTaskListOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetTaskListOne(ctx, req.(*GetTaskListOneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ImportXLSToTaskTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportToTaskListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ImportXLSToTaskTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ImportXLSToTaskTable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ImportXLSToTaskTable(ctx, req.(*ImportToTaskListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DelTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DelTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DelTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DelTask(ctx, req.(*DelTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ModTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ModTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ModTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ModTask(ctx, req.(*ModTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_AddTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryTaskWithSQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTaskWithSQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryTaskWithSQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_QueryTaskWithSQL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryTaskWithSQL(ctx, req.(*QueryTaskWithSQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryTaskWithField_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTaskWithFieldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryTaskWithField(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_QueryTaskWithField_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryTaskWithField(ctx, req.(*QueryTaskWithFieldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetPatchsAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPatchsAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetPatchsAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetPatchsAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetPatchsAll(ctx, req.(*GetPatchsAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetOnePatchs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOnePatchsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetOnePatchs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetOnePatchs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetOnePatchs(ctx, req.(*GetOnePatchsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetPatchsByState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPatchsByStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetPatchsByState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetPatchsByState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetPatchsByState(ctx, req.(*GetPatchsByStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DelPatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelPatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DelPatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_DelPatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DelPatch(ctx, req.(*DelPatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ImportXLSToPatchTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportXLSToPatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ImportXLSToPatchTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ImportXLSToPatchTable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ImportXLSToPatchTable(ctx, req.(*ImportXLSToPatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ModPatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModPatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ModPatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ModPatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ModPatch(ctx, req.(*ModPatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Service_Login_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Service_Register_Handler,
		},
		{
			MethodName: "GetUserInfo",
			Handler:    _Service_GetUserInfo_Handler,
		},
		{
			MethodName: "ModUserInfo",
			Handler:    _Service_ModUserInfo_Handler,
		},
		{
			MethodName: "GetTaskListAll",
			Handler:    _Service_GetTaskListAll_Handler,
		},
		{
			MethodName: "GetTaskListOne",
			Handler:    _Service_GetTaskListOne_Handler,
		},
		{
			MethodName: "ImportXLSToTaskTable",
			Handler:    _Service_ImportXLSToTaskTable_Handler,
		},
		{
			MethodName: "DelTask",
			Handler:    _Service_DelTask_Handler,
		},
		{
			MethodName: "ModTask",
			Handler:    _Service_ModTask_Handler,
		},
		{
			MethodName: "AddTask",
			Handler:    _Service_AddTask_Handler,
		},
		{
			MethodName: "QueryTaskWithSQL",
			Handler:    _Service_QueryTaskWithSQL_Handler,
		},
		{
			MethodName: "QueryTaskWithField",
			Handler:    _Service_QueryTaskWithField_Handler,
		},
		{
			MethodName: "GetPatchsAll",
			Handler:    _Service_GetPatchsAll_Handler,
		},
		{
			MethodName: "GetOnePatchs",
			Handler:    _Service_GetOnePatchs_Handler,
		},
		{
			MethodName: "GetPatchsByState",
			Handler:    _Service_GetPatchsByState_Handler,
		},
		{
			MethodName: "DelPatch",
			Handler:    _Service_DelPatch_Handler,
		},
		{
			MethodName: "ImportXLSToPatchTable",
			Handler:    _Service_ImportXLSToPatchTable_Handler,
		},
		{
			MethodName: "ModPatch",
			Handler:    _Service_ModPatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
