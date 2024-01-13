package file

import (
	"context"
	"errors"
	log "log/slog"
	"mime/multipart"

	"github.com/dhowden/tag"
)

type Saver interface {
	Save(context.Context, multipart.File) (*FileInformation, error)
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

func (m *FileInformationService) Save(ctx context.Context, file multipart.File) (*FileInformation, error) {
	var fileInfo FileInformation

	metadata, err := GetMetadata(file)
	if err != nil {
		return nil, err
	}
	fileInfo.Metadata = metadata

	fileInfo.Filename = ctx.
	return m.repo.Create(ctx, fileInfo)
}

func GetMetadata(file multipart.File) (Metadata, error) {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		log.Error("could not parse metadata from file", "err", err)
		return Metadata{}, err
	}
	
	return Metadata{
		Title: metadata.Title(), 
		Artist: metadata.Artist(), 
		Album: metadata.Album(), 
		Year: metadata.Year(),
		}, 
		nil
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
