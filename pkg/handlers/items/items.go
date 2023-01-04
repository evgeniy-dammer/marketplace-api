package handlers

import (
	"context"

	protos "github.com/evgeniy-dammer/emenu-api/pkg/protos"
)

type ItemServiceServer struct {
	protos.UnimplementedItemServiceServer
}

func (*ItemServiceServer) FindAll(ctx context.Context, in *protos.FindAllRequest) (*protos.FindAllResponse, error) {
	return &protos.FindAllResponse{
		Items: []*protos.Item{
			{
				Id:       "p01",
				Name:     "name 1",
				Price:    4.5,
				Quantity: 4,
				Status:   true,
			},
			{
				Id:       "p02",
				Name:     "name 2",
				Price:    22,
				Quantity: 3,
				Status:   false,
			},
			{
				Id:       "p03",
				Name:     "name 3",
				Price:    27,
				Quantity: 3,
				Status:   true,
			},
		},
	}, nil
}
