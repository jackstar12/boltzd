// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: boltzrpc.proto

package boltzrpc

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
	Boltz_GetInfo_FullMethodName                = "/boltzrpc.Boltz/GetInfo"
	Boltz_GetServiceInfo_FullMethodName         = "/boltzrpc.Boltz/GetServiceInfo"
	Boltz_ListSwaps_FullMethodName              = "/boltzrpc.Boltz/ListSwaps"
	Boltz_GetSwapInfo_FullMethodName            = "/boltzrpc.Boltz/GetSwapInfo"
	Boltz_Deposit_FullMethodName                = "/boltzrpc.Boltz/Deposit"
	Boltz_CreateSwap_FullMethodName             = "/boltzrpc.Boltz/CreateSwap"
	Boltz_CreateChannel_FullMethodName          = "/boltzrpc.Boltz/CreateChannel"
	Boltz_CreateReverseSwap_FullMethodName      = "/boltzrpc.Boltz/CreateReverseSwap"
	Boltz_GetSwapRecommendations_FullMethodName = "/boltzrpc.Boltz/GetSwapRecommendations"
)

// BoltzClient is the client API for Boltz service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BoltzClient interface {
	// Gets general information about the daemon like the chain of the LND node it is connected to
	// and the IDs of pending swaps.
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
	// Fetches the latest limits and fees from the Boltz backend API it is connected to.
	GetServiceInfo(ctx context.Context, in *GetServiceInfoRequest, opts ...grpc.CallOption) (*GetServiceInfoResponse, error)
	// Returns a list of all swaps, reverse swaps and channel creations in the database.
	ListSwaps(ctx context.Context, in *ListSwapsRequest, opts ...grpc.CallOption) (*ListSwapsResponse, error)
	// Gets all available information about a swap from the database.
	GetSwapInfo(ctx context.Context, in *GetSwapInfoRequest, opts ...grpc.CallOption) (*GetSwapInfoResponse, error)
	// This is a wrapper for channel creation swaps. The daemon only returns the ID, timeout block height and lockup address.
	// The Boltz backend takes care of the rest. When an amount of onchain coins that is in the limits is sent to the address
	// before the timeout block height, the daemon creates a new lightning invoice, sends it to the Boltz backend which
	// will try to pay it and if that is not possible, create a new channel to make the swap succeed.
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error)
	// Creates a new swap from onchain to lightning.
	CreateSwap(ctx context.Context, in *CreateSwapRequest, opts ...grpc.CallOption) (*CreateSwapResponse, error)
	// Create a new swap from onchain to a new lightning channel. The daemon will only accept the invoice payment if the HTLCs
	// is coming trough a new channel channel opened by Boltz.
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateSwapResponse, error)
	// Creates a new reverse swap from lightning to onchain. If `accept_zero_conf` is set to true in the request, the daemon
	// will not wait until the lockup transaction from Boltz is confirmed in a block, but will claim it instantly.
	CreateReverseSwap(ctx context.Context, in *CreateReverseSwapRequest, opts ...grpc.CallOption) (*CreateReverseSwapResponse, error)
	GetSwapRecommendations(ctx context.Context, in *GetSwapRecommendationsRequest, opts ...grpc.CallOption) (*GetSwapRecommendationsResponse, error)
}

type boltzClient struct {
	cc grpc.ClientConnInterface
}

func NewBoltzClient(cc grpc.ClientConnInterface) BoltzClient {
	return &boltzClient{cc}
}

func (c *boltzClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := c.cc.Invoke(ctx, Boltz_GetInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) GetServiceInfo(ctx context.Context, in *GetServiceInfoRequest, opts ...grpc.CallOption) (*GetServiceInfoResponse, error) {
	out := new(GetServiceInfoResponse)
	err := c.cc.Invoke(ctx, Boltz_GetServiceInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) ListSwaps(ctx context.Context, in *ListSwapsRequest, opts ...grpc.CallOption) (*ListSwapsResponse, error) {
	out := new(ListSwapsResponse)
	err := c.cc.Invoke(ctx, Boltz_ListSwaps_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) GetSwapInfo(ctx context.Context, in *GetSwapInfoRequest, opts ...grpc.CallOption) (*GetSwapInfoResponse, error) {
	out := new(GetSwapInfoResponse)
	err := c.cc.Invoke(ctx, Boltz_GetSwapInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, Boltz_Deposit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) CreateSwap(ctx context.Context, in *CreateSwapRequest, opts ...grpc.CallOption) (*CreateSwapResponse, error) {
	out := new(CreateSwapResponse)
	err := c.cc.Invoke(ctx, Boltz_CreateSwap_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateSwapResponse, error) {
	out := new(CreateSwapResponse)
	err := c.cc.Invoke(ctx, Boltz_CreateChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) CreateReverseSwap(ctx context.Context, in *CreateReverseSwapRequest, opts ...grpc.CallOption) (*CreateReverseSwapResponse, error) {
	out := new(CreateReverseSwapResponse)
	err := c.cc.Invoke(ctx, Boltz_CreateReverseSwap_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boltzClient) GetSwapRecommendations(ctx context.Context, in *GetSwapRecommendationsRequest, opts ...grpc.CallOption) (*GetSwapRecommendationsResponse, error) {
	out := new(GetSwapRecommendationsResponse)
	err := c.cc.Invoke(ctx, Boltz_GetSwapRecommendations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BoltzServer is the server API for Boltz service.
// All implementations must embed UnimplementedBoltzServer
// for forward compatibility
type BoltzServer interface {
	// Gets general information about the daemon like the chain of the LND node it is connected to
	// and the IDs of pending swaps.
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	// Fetches the latest limits and fees from the Boltz backend API it is connected to.
	GetServiceInfo(context.Context, *GetServiceInfoRequest) (*GetServiceInfoResponse, error)
	// Returns a list of all swaps, reverse swaps and channel creations in the database.
	ListSwaps(context.Context, *ListSwapsRequest) (*ListSwapsResponse, error)
	// Gets all available information about a swap from the database.
	GetSwapInfo(context.Context, *GetSwapInfoRequest) (*GetSwapInfoResponse, error)
	// This is a wrapper for channel creation swaps. The daemon only returns the ID, timeout block height and lockup address.
	// The Boltz backend takes care of the rest. When an amount of onchain coins that is in the limits is sent to the address
	// before the timeout block height, the daemon creates a new lightning invoice, sends it to the Boltz backend which
	// will try to pay it and if that is not possible, create a new channel to make the swap succeed.
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
	// Creates a new swap from onchain to lightning.
	CreateSwap(context.Context, *CreateSwapRequest) (*CreateSwapResponse, error)
	// Create a new swap from onchain to a new lightning channel. The daemon will only accept the invoice payment if the HTLCs
	// is coming trough a new channel channel opened by Boltz.
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateSwapResponse, error)
	// Creates a new reverse swap from lightning to onchain. If `accept_zero_conf` is set to true in the request, the daemon
	// will not wait until the lockup transaction from Boltz is confirmed in a block, but will claim it instantly.
	CreateReverseSwap(context.Context, *CreateReverseSwapRequest) (*CreateReverseSwapResponse, error)
	GetSwapRecommendations(context.Context, *GetSwapRecommendationsRequest) (*GetSwapRecommendationsResponse, error)
	mustEmbedUnimplementedBoltzServer()
}

// UnimplementedBoltzServer must be embedded to have forward compatible implementations.
type UnimplementedBoltzServer struct {
}

func (UnimplementedBoltzServer) GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedBoltzServer) GetServiceInfo(context.Context, *GetServiceInfoRequest) (*GetServiceInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServiceInfo not implemented")
}
func (UnimplementedBoltzServer) ListSwaps(context.Context, *ListSwapsRequest) (*ListSwapsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSwaps not implemented")
}
func (UnimplementedBoltzServer) GetSwapInfo(context.Context, *GetSwapInfoRequest) (*GetSwapInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSwapInfo not implemented")
}
func (UnimplementedBoltzServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (UnimplementedBoltzServer) CreateSwap(context.Context, *CreateSwapRequest) (*CreateSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSwap not implemented")
}
func (UnimplementedBoltzServer) CreateChannel(context.Context, *CreateChannelRequest) (*CreateSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedBoltzServer) CreateReverseSwap(context.Context, *CreateReverseSwapRequest) (*CreateReverseSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReverseSwap not implemented")
}
func (UnimplementedBoltzServer) GetSwapRecommendations(context.Context, *GetSwapRecommendationsRequest) (*GetSwapRecommendationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSwapRecommendations not implemented")
}
func (UnimplementedBoltzServer) mustEmbedUnimplementedBoltzServer() {}

// UnsafeBoltzServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BoltzServer will
// result in compilation errors.
type UnsafeBoltzServer interface {
	mustEmbedUnimplementedBoltzServer()
}

func RegisterBoltzServer(s grpc.ServiceRegistrar, srv BoltzServer) {
	s.RegisterService(&Boltz_ServiceDesc, srv)
}

func _Boltz_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_GetInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).GetInfo(ctx, req.(*GetInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_GetServiceInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).GetServiceInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_GetServiceInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).GetServiceInfo(ctx, req.(*GetServiceInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_ListSwaps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSwapsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).ListSwaps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_ListSwaps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).ListSwaps(ctx, req.(*ListSwapsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_GetSwapInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSwapInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).GetSwapInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_GetSwapInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).GetSwapInfo(ctx, req.(*GetSwapInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_Deposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_CreateSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).CreateSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_CreateSwap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).CreateSwap(ctx, req.(*CreateSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_CreateChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_CreateReverseSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReverseSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).CreateReverseSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_CreateReverseSwap_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).CreateReverseSwap(ctx, req.(*CreateReverseSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Boltz_GetSwapRecommendations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSwapRecommendationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoltzServer).GetSwapRecommendations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Boltz_GetSwapRecommendations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoltzServer).GetSwapRecommendations(ctx, req.(*GetSwapRecommendationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Boltz_ServiceDesc is the grpc.ServiceDesc for Boltz service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Boltz_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "boltzrpc.Boltz",
	HandlerType: (*BoltzServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Boltz_GetInfo_Handler,
		},
		{
			MethodName: "GetServiceInfo",
			Handler:    _Boltz_GetServiceInfo_Handler,
		},
		{
			MethodName: "ListSwaps",
			Handler:    _Boltz_ListSwaps_Handler,
		},
		{
			MethodName: "GetSwapInfo",
			Handler:    _Boltz_GetSwapInfo_Handler,
		},
		{
			MethodName: "Deposit",
			Handler:    _Boltz_Deposit_Handler,
		},
		{
			MethodName: "CreateSwap",
			Handler:    _Boltz_CreateSwap_Handler,
		},
		{
			MethodName: "CreateChannel",
			Handler:    _Boltz_CreateChannel_Handler,
		},
		{
			MethodName: "CreateReverseSwap",
			Handler:    _Boltz_CreateReverseSwap_Handler,
		},
		{
			MethodName: "GetSwapRecommendations",
			Handler:    _Boltz_GetSwapRecommendations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "boltzrpc.proto",
}
