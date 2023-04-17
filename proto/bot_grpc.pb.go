// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.6
// source: bot.proto

package botpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BotService_SignUp_FullMethodName       = "/bot.BotService/SignUp"
	BotService_CreateTask_FullMethodName   = "/bot.BotService/CreateTask"
	BotService_DeleteTask_FullMethodName   = "/bot.BotService/DeleteTask"
	BotService_AssignWorker_FullMethodName = "/bot.BotService/AssignWorker"
	BotService_CloseTask_FullMethodName    = "/bot.BotService/CloseTask"
	BotService_GetOpenTasks_FullMethodName = "/bot.BotService/GetOpenTasks"
)

// BotServiceClient is the client API for BotService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BotServiceClient interface {
	// SignUp user
	SignUp(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	CreateTask(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	DeleteTask(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	AssignWorker(ctx context.Context, in *AssignRequest, opts ...grpc.CallOption) (*AssignResponse, error)
	CloseTask(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error)
	GetOpenTasks(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type botServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBotServiceClient(cc grpc.ClientConnInterface) BotServiceClient {
	return &botServiceClient{cc}
}

func (c *botServiceClient) SignUp(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, BotService_SignUp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botServiceClient) CreateTask(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, BotService_CreateTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botServiceClient) DeleteTask(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, BotService_DeleteTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botServiceClient) AssignWorker(ctx context.Context, in *AssignRequest, opts ...grpc.CallOption) (*AssignResponse, error) {
	out := new(AssignResponse)
	err := c.cc.Invoke(ctx, BotService_AssignWorker_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botServiceClient) CloseTask(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error) {
	out := new(CloseResponse)
	err := c.cc.Invoke(ctx, BotService_CloseTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *botServiceClient) GetOpenTasks(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, BotService_GetOpenTasks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BotServiceServer is the server API for BotService service.
// All implementations must embed UnimplementedBotServiceServer
// for forward compatibility
type BotServiceServer interface {
	// SignUp user
	SignUp(context.Context, *RegisterRequest) (*RegisterResponse, error)
	CreateTask(context.Context, *CreateRequest) (*CreateResponse, error)
	DeleteTask(context.Context, *DeleteRequest) (*DeleteResponse, error)
	AssignWorker(context.Context, *AssignRequest) (*AssignResponse, error)
	CloseTask(context.Context, *CloseRequest) (*CloseResponse, error)
	GetOpenTasks(context.Context, *GetRequest) (*GetResponse, error)
	mustEmbedUnimplementedBotServiceServer()
}

// UnimplementedBotServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBotServiceServer struct {
}

func (UnimplementedBotServiceServer) SignUp(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedBotServiceServer) CreateTask(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedBotServiceServer) DeleteTask(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
func (UnimplementedBotServiceServer) AssignWorker(context.Context, *AssignRequest) (*AssignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignWorker not implemented")
}
func (UnimplementedBotServiceServer) CloseTask(context.Context, *CloseRequest) (*CloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseTask not implemented")
}
func (UnimplementedBotServiceServer) GetOpenTasks(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpenTasks not implemented")
}
func (UnimplementedBotServiceServer) mustEmbedUnimplementedBotServiceServer() {}

// UnsafeBotServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BotServiceServer will
// result in compilation errors.
type UnsafeBotServiceServer interface {
	mustEmbedUnimplementedBotServiceServer()
}

func RegisterBotServiceServer(s grpc.ServiceRegistrar, srv BotServiceServer) {
	s.RegisterService(&BotService_ServiceDesc, srv)
}

func _BotService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_SignUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).SignUp(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BotService_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_CreateTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).CreateTask(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BotService_DeleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).DeleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_DeleteTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).DeleteTask(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BotService_AssignWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).AssignWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_AssignWorker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).AssignWorker(ctx, req.(*AssignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BotService_CloseTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).CloseTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_CloseTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).CloseTask(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BotService_GetOpenTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BotServiceServer).GetOpenTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BotService_GetOpenTasks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BotServiceServer).GetOpenTasks(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BotService_ServiceDesc is the grpc.ServiceDesc for BotService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BotService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bot.BotService",
	HandlerType: (*BotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _BotService_SignUp_Handler,
		},
		{
			MethodName: "CreateTask",
			Handler:    _BotService_CreateTask_Handler,
		},
		{
			MethodName: "DeleteTask",
			Handler:    _BotService_DeleteTask_Handler,
		},
		{
			MethodName: "AssignWorker",
			Handler:    _BotService_AssignWorker_Handler,
		},
		{
			MethodName: "CloseTask",
			Handler:    _BotService_CloseTask_Handler,
		},
		{
			MethodName: "GetOpenTasks",
			Handler:    _BotService_GetOpenTasks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bot.proto",
}
