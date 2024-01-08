package file

import (
	"context"
)

type Saver interface {
	Save(context.Context, FileInformation) (*FileInformation, error)
}

type Deleter interface {
	Delete(context.Context, string) error
}

type Finder interface {
	FindById(context.Context, string) (*FileInformation, error)
}

type AllFinder interface {
	FindAll(context.Context) ([]*FileInformation, error)
}

// Storer defines the API for interacting with NO/SQL storage
type Storer interface {
	Saver
	Deleter
	Finder
	AllFinder
}

type FileInformationService struct {
	repo FileInformationRepository
}

func NewFileInformationService(repo FileInformationRepository) *FileInformationService {
	return &FileInformationService{repo: repo}
}

func (m *FileInformationService) Save(ctx context.Context, fileInfo FileInformation) (*FileInformation, error) {
	return m.repo.Create(ctx, fileInfo)
}

func (m *FileInformationService) Delete(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}

func (m *FileInformationService) FindById(ctx context.Context, id string) (*FileInformation, error) {
	return m.repo.FindById(ctx, id)
}

func (m *FileInformationService) FindAll(ctx context.Context) ([]*FileInformation, error) {
	return m.repo.FindAll(ctx)
}
