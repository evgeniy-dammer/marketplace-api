package handlers

import (
	"context"

	itemservice "github.com/evgeniy-dammer/emenu-api/pkg/proto"
)

type ItemServiceServer struct {
	itemservice.UnimplementedItemServiceServer
}

func (*ItemServiceServer) FindAll(ctx context.Context, in *itemservice.FindAllRequest) (*itemservice.FindAllResponse, error) {
	return &itemservice.FindAllResponse{
		Items: []*itemservice.Item{
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
