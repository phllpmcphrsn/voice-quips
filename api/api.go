package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	log "log/slog"

	"github.com/gin-gonic/gin"
	"github.com/phllpmcphrsn/voice-quips/config"
	"github.com/phllpmcphrsn/voice-quips/file"
	"github.com/phllpmcphrsn/voice-quips/s3"
)

const MegaByte int64 = 10 << 10

type APIServer struct {
	// API properties
	basePath   string
	listenAddr string
	env        string

	// (Dependency) Inject the services
	s3Service   s3.DownloadUploader
	fileService file.Storer
}

func NewAPIServer(apiConfig config.APIConfig, s3Service s3.DownloadUploader, fileService file.Storer) *APIServer {
	return &APIServer{
		basePath:    apiConfig.Path,
		listenAddr:  apiConfig.Address,
		env:         apiConfig.Env,
		s3Service:   s3Service,
		fileService: fileService,
	}
}

// Ping godoc
func (a *APIServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}

// GET /api/v1/audio
// This endpoint will simply return the metadata of all files; no audio will be returned here
func (a *APIServer) getAudio(c *gin.Context) {
	// call to fileService to get a list of filenames (or perhaps s3links)
	metadatum, err := a.fileService.FindAll(c)
	if err != nil {
		log.Error("Could not retrieve entries", "err", err)
		c.AbortWithError(http.StatusInternalServerError, InternalServerError(""))
		return
	}

	c.IndentedJSON(http.StatusOK, metadatum)
}

// GET /audio/{id, name}
func (a *APIServer) getAudioByIdAndName(c *gin.Context) {
	var err error

	id := c.Query("id")
	name := c.Query("name")

	if id == "" && name == "" {
		err = errors.New("neither id or name params given")
		log.Error("request failed", "err", err, "request", c.Request.RequestURI)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	file_info, err := a.fileService.FindById(c, id)
	if err != nil {
		log.Error("request for the following ID was not found", "err", err, "id", id, "request", c.Request.RequestURI)
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	file, err := a.s3Service.DownloadObject(c, name, file_info.S3Link)
	print(file)

}

// POST /audio
func (a *APIServer) createAudio(c *gin.Context) {
	
	err := c.Request.ParseMultipartForm(MegaByte) // 1 MB
	if err != nil {
		log.Error("could not parse form in request", "err", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("Could not retrieve upload file from request", "err", err)
		c.AbortWithError(http.StatusInternalServerError, InternalServerError("Could not retrieve upload file from request"))
		return
	}

	defer file.Close()

	// content goes to S3
	content, err := io.ReadAll(file)
	if err != nil {
		log.Error("could not read file", "err", err)
	}

	if len(content) > int(MegaByte) {
		err = errors.New("file size too large")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
	}

	// make call to fileInfo DB, returns required fileInfo object
	a.fileService.Save(c, file)
	// make call to s3Service, get bucket and filename from returned fileInfo object, supply content
	// a.s3Service.UploadObject(c, )
	log.Info("file related stuff", "file", string(content), "header", header)
	c.IndentedJSON(http.StatusCreated, nil)
}

// DELETE /api/v1/audio/{id}

// StartRouter starts up a Gin router for the API
func (a *APIServer) StartRouter() {
	r := gin.Default()
	if os.Getenv(gin.EnvGinMode) == "" {
		mode := ginEnvMode(a.env)
		gin.SetMode(mode)
	}

	// setup v1 routes
	v1 := r.Group(a.basePath)
	{
		v1.GET("/ping", a.ping)
		v1.GET("/audio/", a.getAudio)
		// v1.GET("/audio/:id", a.getAudioById)
		v1.POST("/audio", a.createAudio)
		// v1.DELETE("/audio/:id", a.deleteAudio)
	}

	r.Run(a.listenAddr)
}

func ginEnvMode(env string) string {
	switch strings.ToLower(env) {
	case "prod":
		return gin.ReleaseMode
	case "dev":
		return gin.DebugMode
	default:
		return gin.TestMode
	}
}
