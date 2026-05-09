package notification

import (
	"context"
	"google.golang.org/grpc"
)

type NotificationRequest struct {
	UserId  string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Title   string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *NotificationRequest) Reset()         {}
func (m *NotificationRequest) String() string { return "" }
func (*NotificationRequest) ProtoMessage()    {}
func (x *NotificationRequest) GetUserId() string {
	if x != nil { return x.UserId }
	return ""
}
func (x *NotificationRequest) GetTitle() string {
	if x != nil { return x.Title }
	return ""
}
func (x *NotificationRequest) GetMessage() string {
	if x != nil { return x.Message }
	return ""
}

type NotificationResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (m *NotificationResponse) Reset()         {}
func (m *NotificationResponse) String() string { return "" }
func (*NotificationResponse) ProtoMessage()    {}
func (x *NotificationResponse) GetStatus() string {
	if x != nil { return x.Status }
	return ""
}

type NotificationServiceClient interface {
	SendNotification(ctx context.Context, in *NotificationRequest, opts ...grpc.CallOption) (*NotificationResponse, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) SendNotification(ctx context.Context, in *NotificationRequest, opts ...grpc.CallOption) (*NotificationResponse, error) {
	out := new(NotificationResponse)
	err := c.cc.Invoke(ctx, "/notification.NotificationService/SendNotification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type NotificationServiceServer interface {
	SendNotification(context.Context, *NotificationRequest) (*NotificationResponse, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) SendNotification(context.Context, *NotificationRequest) (*NotificationResponse, error) {
	return nil, nil
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendNotification",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(NotificationRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(NotificationServiceServer).SendNotification(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/notification.NotificationService/SendNotification",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(NotificationServiceServer).SendNotification(ctx, req.(*NotificationRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/notification.proto",
}
