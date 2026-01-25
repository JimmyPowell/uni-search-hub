package handler

import (
	"net/http"
	"strconv"
	"uni-search-hub/internal/model"
	"uni-search-hub/internal/server"

	"github.com/gin-gonic/gin"
)

type ChannelUpsertRequest struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ApiKey   string `json:"api_key,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
	Weight   int    `json:"weight,omitempty"`
}

func GetChannels(c *gin.Context) {
	pageInfo := model.GetPageQuery(c)
	channels, total, err := model.GetChannels(pageInfo)
	if err != nil {
		server.ApiError(c, err)
		return
	}
	for _, ch := range channels {
		ch.Clean()
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(channels)
	server.ApiSuccess(c, pageInfo)
}

func CreateChannel(c *gin.Context) {
	var req ChannelUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	ch := &model.Channel{
		Name:     req.Name,
		Provider: req.Provider,
		ApiKey:   req.ApiKey,
		Enabled:  enabled,
		Weight:   req.Weight,
	}
	if ch.Weight == 0 {
		ch.Weight = 1
	}
	if err := ch.Insert(); err != nil {
		server.ApiError(c, err)
		return
	}
	ch.Clean()
	server.ApiSuccess(c, ch)
}

func UpdateChannel(c *gin.Context) {
	var req ChannelUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的参数"})
		return
	}

	existing, err := model.GetChannelById(req.Id)
	if err != nil {
		server.ApiError(c, err)
		return
	}

	// 仅更新传入字段（ApiKey 允许留空表示不修改）
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Provider != "" {
		existing.Provider = req.Provider
	}
	if req.Enabled != nil {
		existing.Enabled = *req.Enabled
	}
	if req.Weight != 0 {
		existing.Weight = req.Weight
	}
	if req.ApiKey != "" {
		existing.ApiKey = req.ApiKey
	}

	if err := existing.Update(); err != nil {
		server.ApiError(c, err)
		return
	}
	existing.Clean()
	server.ApiSuccess(c, existing)
}

func DeleteChannel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的参数"})
		return
	}
	ch := &model.Channel{Id: id}
	if err := ch.Delete(); err != nil {
		server.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

