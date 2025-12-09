package controllers

import (
	"net/http"
	"strconv"

	"devplus-backend/internal/services"
	"devplus-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RepoController struct {
	Service *services.RepoService
}

func NewRepoController() *RepoController {
	return &RepoController{
		Service: &services.RepoService{},
	}
}

func (ctrl *RepoController) ListRepos(c *gin.Context) {
	repos, err := ctrl.Service.ListRepos()
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to list repositories", err)
		return
	}
	c.JSON(http.StatusOK, repos)
}

func (ctrl *RepoController) SyncRepo(c *gin.Context) {
	idStr := c.Param("id")
	repoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid repository ID", err)
		return
	}

	token := c.GetHeader("Authorization")
	// Expected format: Bearer <token>
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	} else {
		token = ""
	}

	if err := ctrl.Service.SyncRepo(uint(repoID), token); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to sync repository", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

func (ctrl *RepoController) ListPRs(c *gin.Context) {
	idStr := c.Param("id")
	repoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid repository ID", err)
		return
	}

	prs, err := ctrl.Service.ListPRs(uint(repoID))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to list PRs", err)
		return
	}
	c.JSON(http.StatusOK, prs)
}

func (ctrl *RepoController) GetPR(c *gin.Context) {
	idStr := c.Param("id")
	repoID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid repository ID", err)
		return
	}

	prNumStr := c.Param("pr_number")
	prNumber, err := strconv.Atoi(prNumStr)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid PR number", err)
		return
	}

	pr, err := ctrl.Service.GetPR(uint(repoID), prNumber)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get PR", err)
		return
	}
	c.JSON(http.StatusOK, pr)
}
