package api

import (
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

// GET /api/v1/audio/{id, name}
// POST /api/v1/audio
func (a *APIServer) createAudio(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 10) // 1 MB
	if err != nil {
		log.Error("could not parse form in request", "err", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		log.Error("Could not retrieve upload file from request", "err", err)
		c.AbortWithError(http.StatusInternalServerError, InternalServerError(""))
		return
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Error("could not read file", "err", err)
	}

	log.Info("file related stuff", "file", string(content), "header", header)
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
