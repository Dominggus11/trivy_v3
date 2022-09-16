package controllers

import (
	"net/http"
	"os"
	"trivy_v3/models"

	"github.com/gin-gonic/gin"
)

func HelloUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Developer": "Roy Dominggus Andornov Malau",
		"Project":   "Trivy Misconfiguration",
	})
}

func FindProjects(c *gin.Context) {
	var projects []models.Projects
	models.DB.Find(&projects)

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

func PostProject(c *gin.Context) {
	db := models.DB
	var input models.Projects
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := db.Where("project_name = ?", input.ProjectName).First(&input).Error; err != nil {
		project := models.Projects{
			ProjectName: input.ProjectName,
		}
		db.Create(&project)
		c.JSON(http.StatusOK, gin.H{
			"data": project})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Name Sudah Ada"})
		return
	}

}

func FindProject(c *gin.Context) {
	db := models.DB
	// Get model if exist
	var project models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Tidak Tersedia Bos Q"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": project})
}

func UpdateProject(c *gin.Context) {
	db := models.DB
	// Get model if exist
	var input models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Tidak Tersedia Bos Q"})
		return
	}

	if err := db.Where("project_name = ?", input.ProjectName).First(&input).Error; err != nil {
		project := models.Projects{
			ProjectName: input.ProjectName,
		}
		//db.Updates(&dockerfile)
		db.Model(&input).Updates(project)
		c.JSON(http.StatusOK, gin.H{
			"data": project,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Name Sudah Ada"})
		return
	}

}

func DeleteProject(c *gin.Context) {
	db := models.DB
	// Get model if exist
	var input models.Projects
	if err := db.Where("id = ?", c.Param("id")).First(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Tidak Tersedia Bos Q"})
		return
	}

	var dockerfile []models.Dockerfiles
	err := db.Where("project_id = ?", c.Param("id")).Find(&dockerfile).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Tidak Tersedia Bos Q"})
		return
	}

	for _, d := range dockerfile {
		pathFile := d.Pathfile
		pathJson := d.PathJson
		os.RemoveAll(pathFile)
		os.RemoveAll(pathJson)
		db.Delete(&d)
	}

	db.Delete(&input)
	c.JSON(http.StatusOK, gin.H{
		"data": "Deleted",
	})
}
