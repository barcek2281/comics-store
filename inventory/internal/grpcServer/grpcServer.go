package grpcserver

import (
	"context"
	"fmt"

	"github.com/barcek2281/comics-store/inventory/internal/model"
	"github.com/barcek2281/comics-store/inventory/internal/storage/sqlite"
	inventoryv1 "github.com/barcek2281/proto/gen/go/inventory"
)

type GRPCserver struct {
	store *sqlite.Storage
	inventoryv1.UnimplementedInventoryServer
}

func New(store *sqlite.Storage) *GRPCserver {
	return &GRPCserver{
		store: store,
	}
}

func (g *GRPCserver) Create(ctx context.Context, in *inventoryv1.CreateRequest) (*inventoryv1.CreateResponce, error) {
	comic := model.Comics{
		Title:       in.GetTitle(),
		Author:      in.GetAuthor(),
		Description: in.GetDescription(),
		ReleaseDate: in.GetReleaseDate(),
		Price:       float32(in.GetPrice()),
		Quantity:    int32(in.GetQuantity()),
	}

	id, err := g.store.Create(comic)
	if err != nil {
		return nil, fmt.Errorf("failed to create comic: %w", err)
	}

	return &inventoryv1.CreateResponce{
		Id: id,
	}, nil
}

func (g *GRPCserver) Delete(ctx context.Context, in *inventoryv1.DeleteRequest) (*inventoryv1.DeleteResponce, error) {
	err := g.store.Delete(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete comic: %w", err)
	}

	return &inventoryv1.DeleteResponce{
		IsDeleted: true,
		Result: "",
	}, nil
}

func (g *GRPCserver) Get(ctx context.Context, in *inventoryv1.GetRequest) (*inventoryv1.Comics, error) {
	comic, err := g.store.Get(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get comic: %w", err)
	}

	return &inventoryv1.Comics{
		Id:          fmt.Sprint(comic.ID),
		Title:       comic.Title,
		Author:      comic.Author,
		Description: comic.Description,
		ReleaseDate: comic.ReleaseDate,
		Price:       float32(comic.Price),
		Quantity:    int32(comic.Quantity),
	}, nil
}

func (g *GRPCserver) List(ctx context.Context, in *inventoryv1.ListRequest) (*inventoryv1.ListResponse, error) {
	comics, err := g.store.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list comics: %w", err)
	}

	var list []*inventoryv1.Comics
	for _, c := range comics {
		list = append(list, &inventoryv1.Comics{
			Id:          fmt.Sprint(c.ID),
			Title:       c.Title,
			Author:      c.Author,
			Description: c.Description,
			ReleaseDate: c.ReleaseDate,
			Price:       float32(c.Price),
			Quantity:    int32(c.Quantity),
		})
	}

	return &inventoryv1.ListResponse{Comics: list}, nil
}

func (g *GRPCserver) Update(ctx context.Context, in *inventoryv1.UpdateRequest) (*inventoryv1.UpdateResponce, error) {
	comic := model.Comics{
		ID:          in.GetId(),
		Title:       in.GetTitle(),
		Author:      in.GetAuthor(),
		Description: in.GetDescription(),
		ReleaseDate: in.GetReleaseDate(),
		Price:       float32(in.GetPrice()),
		Quantity:    int32(in.GetQuantity()),
	}

	err := g.store.Update(comic)
	if err != nil {
		return nil, fmt.Errorf("failed to update comic: %w", err)
	}

	return &inventoryv1.UpdateResponce{
		Successfully: true,
		Result: "",
	}, nil
}
