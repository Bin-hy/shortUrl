package interfaces

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenV3(c *gin.Context) {
	var req ShortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.svc.GenerateShortUrlV3(c, req.LongUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	shortUrl = fmt.Sprintf("http://%s/v3/%s", c.Request.Host, shortUrl)
	c.JSON(http.StatusOK, ShortenResp{ShortUrl: shortUrl})
}

func (h *Handler) RedirectV3(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longUrl, err := h.svc.GetLongUrlV3(c, shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if longUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusFound, longUrl)
}
