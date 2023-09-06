package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phllpmcphrsn/voice-quips/metadata"
	"github.com/phllpmcphrsn/voice-quips/s3"
	"golang.org/x/text/cases"
)

const basePath = "/api/v1"
type APIServer struct {
	// API properties
	listenAddr string
	env        string

	// (Dependency) Inject the services
	s3Service       s3.S3DownloadUploader
	metadataService metadata.MetadataStorer
}

func NewAPIServer(address, env string, s3Service s3.S3DownloadUploader, metadataService metadata.MetadataStorer) *APIServer{
	return &APIServer{
		listenAddr: address,
		env: env,
		s3Service: s3Service,
		metadataService: metadataService,
	}
}

// Ping godoc
func (a *APIServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "PONG")
}

// GET /api/v1/audio
// GET /api/v1/audio/{id, name}
// POST /api/v1/audio
// DELETE /api/v1/audio/{id}

// StartRouter starts up a Gin router for the API
func (a *APIServer) StartRouter() {
	r := gin.Default()
	if os.Getenv(gin.EnvGinMode) == "" {
		mode := ginEnvMode(a.env)
		gin.SetMode(mode)
	}

	v1 := r.Group(basePath) 
	{
		v1.GET("/ping", a.ping)
		v1.GET("/audio/", a.getAudio)
		v1.GET("/audio/:id", a.getAudioById)
		v1.POST("/audio", a.createAudio)
		v1.DELETE("/audio/:id", a.deleteAudio)
	}
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