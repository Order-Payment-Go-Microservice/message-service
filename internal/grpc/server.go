package grpc

import (
	"context"
	"log"

	"github.com/Order-Payment-Go-Microservice/message-service/internal/model"
	"github.com/Order-Payment-Go-Microservice/message-service/internal/service"
	pb "github.com/Order-Payment-Go-Microservice/proto-generation/gen/message/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	svc service.MessageService
}

func NewMessageServer(svc service.MessageService) *MessageServer {
	return &MessageServer{svc: svc}
}

func toProtoMessage(m model.Message) *pb.MessageResponse {
	return &pb.MessageResponse{
		Id: m.ID.String(), ChatId: m.ChatID.String(),
		SenderId: m.SenderID.String(), ReceiverId: m.ReceiverID.String(),
		Content: m.Content, MessageType: m.MessageType,
		IsRead: m.IsRead, IsDelivered: m.IsDelivered,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
}

func (s *MessageServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.MessageResponse, error) {
	chatID, _ := uuid.Parse(req.ChatId)
	senderID, _ := uuid.Parse(req.SenderId)
	receiverID, _ := uuid.Parse(req.ReceiverId)
	msg, err := s.svc.SendMessage(chatID, senderID, receiverID, req.Content, req.MessageType)
	if err != nil {
		return nil, err
	}
	return toProtoMessage(*msg), nil
}

func (s *MessageServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	chatID, _ := uuid.Parse(req.ChatId)
	limit := int(req.Limit)
	offset := int(req.Offset)
	messages, err := s.svc.GetChatHistory(chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.MessageResponse, len(messages))
	for i, m := range messages {
		list[i] = toProtoMessage(m)
	}
	return &pb.GetChatHistoryResponse{Messages: list}, nil
}

func (s *MessageServer) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MessageResponse, error) {
	id, _ := uuid.Parse(req.MessageId)
	if err := s.svc.MarkRead(id); err != nil {
		return nil, err
	}
	msg, _ := s.svc.GetMessage(id)
	if msg != nil {
		return toProtoMessage(*msg), nil
	}
	return &pb.MessageResponse{Id: id.String(), IsRead: true}, nil
}

func (s *MessageServer) MarkAsDelivered(ctx context.Context, req *pb.MarkAsDeliveredRequest) (*pb.MessageResponse, error) {
	id, _ := uuid.Parse(req.MessageId)
	if err := s.svc.MarkDelivered(id); err != nil {
		return nil, err
	}
	msg, _ := s.svc.GetMessage(id)
	if msg != nil {
		return toProtoMessage(*msg), nil
	}
	return &pb.MessageResponse{Id: id.String(), IsDelivered: true}, nil
}

func (s *MessageServer) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*emptypb.Empty, error) {
	id, _ := uuid.Parse(req.MessageId)
	if err := s.svc.DeleteMessage(id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageServer) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.MessageResponse, error) {
	id, _ := uuid.Parse(req.MessageId)
	msg, err := s.svc.EditMessage(id, req.NewContent)
	if err != nil {
		return nil, err
	}
	return toProtoMessage(*msg), nil
}

func (s *MessageServer) SearchMessages(ctx context.Context, req *pb.SearchMessagesRequest) (*pb.GetChatHistoryResponse, error) {
	chatID, _ := uuid.Parse(req.ChatId)
	messages, err := s.svc.SearchMessages(chatID, req.Query)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.MessageResponse, len(messages))
	for i, m := range messages {
		list[i] = toProtoMessage(m)
	}
	return &pb.GetChatHistoryResponse{Messages: list}, nil
}

func (s *MessageServer) GetMessageStatus(ctx context.Context, req *pb.GetMessageStatusRequest) (*pb.MessageStatusResponse, error) {
	id, _ := uuid.Parse(req.MessageId)
	msg, err := s.svc.GetMessage(id)
	if err != nil {
		return nil, err
	}
	return &pb.MessageStatusResponse{
		IsRead: msg.IsRead, IsDelivered: msg.IsDelivered,
		UpdatedAt: timestamppb.New(msg.UpdatedAt),
	}, nil
}

func (s *MessageServer) PinMessage(ctx context.Context, req *pb.PinMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Pin message %s", req.MessageId)
	return &emptypb.Empty{}, nil
}

func (s *MessageServer) UnpinMessage(ctx context.Context, req *pb.UnpinMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Unpin message %s", req.MessageId)
	return &emptypb.Empty{}, nil
}

func (s *MessageServer) GetPinnedMessages(ctx context.Context, req *pb.GetPinnedMessagesRequest) (*pb.GetChatHistoryResponse, error) {
	return &pb.GetChatHistoryResponse{}, nil
}

func (s *MessageServer) ExportChatToEmail(ctx context.Context, req *pb.ExportChatRequest) (*emptypb.Empty, error) {
	log.Printf("Exporting chat %s to email %s", req.ChatId, req.Email)
	return &emptypb.Empty{}, nil
}
