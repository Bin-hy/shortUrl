package interfaces

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ShortenV2(c *gin.Context) {
	var req ShortenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := h.svc.GenerateShortUrlV2(c, req.LongUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	shortUrl = fmt.Sprintf("http://%s/v2/%s", c.Request.Host, shortUrl)
	c.JSON(http.StatusOK, ShortenResp{ShortUrl: shortUrl})
}

func (h *Handler) RedirectV2(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	longUrl, err := h.svc.GetLongUrlV2(c, shortUrl)
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
