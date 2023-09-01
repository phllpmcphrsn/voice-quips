package metadata

import (
	"context"
)

type Saver interface {
	Save(context.Context, Metadata) (*Metadata, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}

type Finder interface {
	FindById(context.Context, string) (*Metadata, error)
}

type AllFinder interface {
	FindAll(context.Context) ([]*Metadata, error)
}

// MetadataStorer defines the API for interacting with NO/SQL storage
type MetadataStorer interface {
	Saver
	Deleter
	Finder
	AllFinder
}

type MetadataService struct {
	repo MetadataRepository
}

func (m *MetadataService) Save(ctx context.Context, metadata Metadata) (*Metadata, error) {
	return m.repo.Create(ctx, metadata)
}

func (m *MetadataService) Delete(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}

func (m *MetadataService) FindById(ctx context.Context, id string) (*Metadata, error) {
	return m.repo.FindById(ctx, id)
}
func (m *MetadataService) FindAll(ctx context.Context) ([]*Metadata, error) {
	return m.repo.FindAll(ctx)
}
