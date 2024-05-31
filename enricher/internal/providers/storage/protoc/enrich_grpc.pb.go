// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: enrich.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EnrichStorageClient is the client API for EnrichStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnrichStorageClient interface {
	IsFIOPresents(ctx context.Context, in *FIO, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
	Update(ctx context.Context, in *Enriched, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type enrichStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewEnrichStorageClient(cc grpc.ClientConnInterface) EnrichStorageClient {
	return &enrichStorageClient{cc}
}

func (c *enrichStorageClient) IsFIOPresents(ctx context.Context, in *FIO, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/EnrichStorage/IsFIOPresents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *enrichStorageClient) Update(ctx context.Context, in *Enriched, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/EnrichStorage/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnrichStorageServer is the server API for EnrichStorage service.
// All implementations must embed UnimplementedEnrichStorageServer
// for forward compatibility
type EnrichStorageServer interface {
	IsFIOPresents(context.Context, *FIO) (*wrapperspb.BoolValue, error)
	Update(context.Context, *Enriched) (*emptypb.Empty, error)
	mustEmbedUnimplementedEnrichStorageServer()
}

// UnimplementedEnrichStorageServer must be embedded to have forward compatible implementations.
type UnimplementedEnrichStorageServer struct {
}

func (UnimplementedEnrichStorageServer) IsFIOPresents(context.Context, *FIO) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFIOPresents not implemented")
}
func (UnimplementedEnrichStorageServer) Update(context.Context, *Enriched) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedEnrichStorageServer) mustEmbedUnimplementedEnrichStorageServer() {}

// UnsafeEnrichStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnrichStorageServer will
// result in compilation errors.
type UnsafeEnrichStorageServer interface {
	mustEmbedUnimplementedEnrichStorageServer()
}

func RegisterEnrichStorageServer(s grpc.ServiceRegistrar, srv EnrichStorageServer) {
	s.RegisterService(&EnrichStorage_ServiceDesc, srv)
}

func _EnrichStorage_IsFIOPresents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FIO)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrichStorageServer).IsFIOPresents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EnrichStorage/IsFIOPresents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrichStorageServer).IsFIOPresents(ctx, req.(*FIO))
	}
	return interceptor(ctx, in, info, handler)
}

func _EnrichStorage_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Enriched)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnrichStorageServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EnrichStorage/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnrichStorageServer).Update(ctx, req.(*Enriched))
	}
	return interceptor(ctx, in, info, handler)
}

// EnrichStorage_ServiceDesc is the grpc.ServiceDesc for EnrichStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnrichStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "EnrichStorage",
	HandlerType: (*EnrichStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsFIOPresents",
			Handler:    _EnrichStorage_IsFIOPresents_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _EnrichStorage_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "enrich.proto",
}
