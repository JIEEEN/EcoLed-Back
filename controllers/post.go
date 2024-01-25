package controllers

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"

	"github.com/gin-gonic/gin"
)

type PostControllers struct{}

var postService = new(services.PostService)

func (ctr PostControllers) CreatePost(c *gin.Context) {
	//Get body by PostForm (form)
	title := c.PostForm("title")
	body := c.PostForm("body")
	var postForm = forms.PostForm{
		Title: title,
		Body:  body,
	}

	//Get userID from token & Change type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Strat transaction
	tx := initializers.DB.Begin()
	defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

	//Create (service)
	err = postService.CreatePost(tx, userID, postForm)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Upload image
	var imageService services.ImageService
	//By form-data type, file is uploaded 
	file, err := c.FormFile("file")
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Open file
	filename := filepath.Base(file.Filename)
	filecontent, _ := file.Open()
	defer filecontent.Close()

	//Get imageURL (in Service)
	imageURL, err := imageService.UploadPostImage(tx, context.Background(), filecontent, userID, filename)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{"Post created successfully with image!": imageURL})

}

func (ctr PostControllers) GetUserPost(c *gin.Context) {
	// Get userID from param & Change type to uint (util)
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get User's all posts (service)
	posts, err := postService.GetUserPost(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return posts
	c.JSON(http.StatusOK, gin.H{"posts": posts})

}


func (ctr PostControllers) GetMyPost(c *gin.Context) {
	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get My all posts (service)
	posts, err := postService.GetUserPost(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return posts
	c.JSON(http.StatusOK, gin.H{"posts": posts})

}


func (ctr PostControllers) GetPost(c *gin.Context) {
	// Get postID from param & Change type to uint (util)
	postID, err := utils.GetPostID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Get One Post (service)
	post, err := postService.GetPost(postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return post
	c.JSON(http.StatusOK, gin.H{"post": post})

}


func (ctr PostControllers) UpdatePost(c *gin.Context) {
	//Get body by PostForm (form)
	title := c.PostForm("title")
	body := c.PostForm("body")
	var postForm = forms.PostForm{
		Title: title,
		Body:  body,
	}

	//Get userID from token & Change type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get postID from param & Change type to uint (util)
	postID, err := utils.GetPostID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Start transaction
	tx := initializers.DB.Begin()

	//Update post (service)
	err = postService.UpdatePost(tx, userID, postID, postForm)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Upload image
	var imageService services.ImageService
	//By form-data type, file is uploaded
	file, err := c.FormFile("file")
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Open file
	filename := filepath.Base(file.Filename)
	filecontent, _ := file.Open()
	defer filecontent.Close()

	//Get imageURL
	imageURL, err := imageService.UploadPostImage(tx, context.Background(), filecontent, userID, filename)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{"Post Updated successfully with image!": imageURL})

}


func (ctr PostControllers) DeletePost(c *gin.Context) {
	// Get postID from param & Change type to uint (util)
	postID, err := utils.GetPostID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from token & Chage type to uint (util)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Delete image(in google cloud storage & DB) (service)
	var imageService services.ImageService
	err = imageService.DeletePostImage(context.Background(), userID, postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Delete post (service)
	err = postService.DeletePost(userID, postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Return message
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Success",
	})

}
