package file

import (
	"context"
)

type Saver interface {
	Save(context.Context, AudioFile) (*AudioFile, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}

type Finder interface {
	FindById(context.Context, string) (*AudioFile, error)
}

type AllFinder interface {
	FindAll(context.Context) ([]*AudioFile, error)
}

// Storer defines the API for interacting with NO/SQL storage
type Storer interface {
	Saver
	Deleter
	Finder
	AllFinder
}

type AudioFileService struct {
	repo FileInformationRepository
}

func (m *AudioFileService) Save(ctx context.Context, audioFile AudioFile) (*AudioFile, error) {
	return m.repo.Create(ctx, audioFile)
}

func (m *AudioFileService) Delete(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}

func (m *AudioFileService) FindById(ctx context.Context, id string) (*AudioFile, error) {
	return m.repo.FindById(ctx, id)
}
func (m *AudioFileService) FindAll(ctx context.Context) ([]*AudioFile, error) {
	return m.repo.FindAll(ctx)
}
