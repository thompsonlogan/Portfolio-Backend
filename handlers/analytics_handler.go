package handlers

import (
	"api/internal/data/model"
	"api/internal/util"
	"api/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddVisitRequest defines the JSON body for recording a visit
type AddVisitRequest struct {
	Source  string `json:"source" binding:"required"` // Source of the visit (required)
}

// AddVisitResponse defines the JSON response for recording a visit
type AddVisitResponse struct {
	Message string `json:"message"` // Success message
}

type AnalyticsHandler struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
  return &AnalyticsHandler{service: service}
}

type VisitJob struct {
	Visit   *model.PortfolioVisit
	Handler func(*model.PortfolioVisit) error
}

var writeQueue = make(chan VisitJob, 100)

// AddVisit godoc
// @Summary Record a portfolio visit
// @Description Adds a visit record for the given source. Automatically records IP, user-agent, and referrer.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param visit body AddVisitRequest true "Visit data"
// @Success 200 {object} AddVisitResponse
// @Failure 400 {object} map[string]string "error: bad request"
// @Failure 500 {object} map[string]string "error: internal server error"
// @Router /visit [post]
func (h *AnalyticsHandler) AddVisit(c *gin.Context) {
	h.handleVisit(c, h.service.AddVisit)
}


// AddGithubVisit godoc
// @Summary Record a github visit
// @Description Adds a visit record for the given src and IP. If the record exists, increments github visit count.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param visit body AddVisitRequest true "Visit data"
// @Success 200 {object} AddVisitResponse
// @Failure 400 {object} map[string]string "error: bad request"
// @Failure 500 {object} map[string]string "error: internal server error"
// @Router /visit/github [post]
func (h *AnalyticsHandler) AddGithubVisit(c *gin.Context) {
	h.handleVisit(c, h.service.AddGithubVisit)
}

// AddLinkedinVisit godoc
// @Summary Record a Linkedin visit
// @Description Adds a visit record for the given src and IP. If the record exists, increments Linkedin visit count.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param visit body AddVisitRequest true "Visit data"
// @Success 200 {object} AddVisitResponse
// @Failure 400 {object} map[string]string "error: bad request"
// @Failure 500 {object} map[string]string "error: internal server error"
// @Router /visit/linkedin [post]
func (h *AnalyticsHandler) AddLinkedinVisit(c *gin.Context) {
	h.handleVisit(c, h.service.AddLinkedinVisit)
}

// AddLinkedinVisit godoc
// @Summary Record a resume download
// @Description Adds a visit record for the given src and IP. If the record exists, increments resume download count.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param visit body AddVisitRequest true "Visit data"
// @Success 200 {object} AddVisitResponse
// @Failure 400 {object} map[string]string "error: bad request"
// @Failure 500 {object} map[string]string "error: internal server error"
// @Router /visit/resume [post]
func (h *AnalyticsHandler) AddResumeDownload(c *gin.Context) {
	h.handleVisit(c, h.service.AddResumeDownload)
}

func (h *AnalyticsHandler) handleVisit(c *gin.Context, handlerFunc func(*model.PortfolioVisit) error) {
	var req AddVisitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.Request.UserAgent()
	if util.IsBot(userAgent) {
		c.JSON(http.StatusOK, AddVisitResponse{Message: "ignored bot visit"})
		return
	}

	visit := &model.PortfolioVisit{
		Source:    req.Source,
		Referrer:  c.Request.Referer(),
		IP:        util.HashIP(c.ClientIP()),
		UserAgent: userAgent,
	}

	// Enqueue asynchronously
	select {
	case writeQueue <- VisitJob{Visit: visit, Handler: handlerFunc}:
		c.JSON(http.StatusOK, AddVisitResponse{Message: "visit queued successfully"})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server busy, try again later"})
	}
}

// StartAnalyticsQueueWorker runs a background worker that processes visit events
func StartAnalyticsQueueWorker() {
	for job := range writeQueue {
		if err := job.Handler(job.Visit); err != nil {
			log.Printf("failed to process visit for source=%s: %v", job.Visit.Source, err)
		}
	}
}