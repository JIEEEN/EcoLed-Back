package controllers

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/Eco-Led/EcoLed-Back_test/utils"

	"github.com/gin-gonic/gin"
)

type PostControllers struct{}

var postService = new(services.PostService)

// TODO: 제목, 내용은 작성완료해서 DB에 넣었으나, 이미지가 실패했을 때 처리하기.
func (ctr PostControllers) CreatePost(c *gin.Context) {
	//Get body by PostForm
	title := c.PostForm("title")
	body := c.PostForm("body")
	var postForm = forms.PostForm{
		Title: title,
		Body:  body,
	}

	//Get userID from token & Change type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	//Create (service)
	err = postService.CreatePost(userID, postForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Upload image
	var imageService services.ImageService
	//By form-data type, file is uploaded (in Controller)
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Open file
	filename := filepath.Base(file.Filename)
	filecontent, _ := file.Open()
	defer filecontent.Close()

	//Get imageURL (in Service)
	imageURL, err := imageService.UploadPostImage(context.Background(), filecontent, userID, filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{
		"Post created successfully with image!": imageURL,
	})

}

func (ctr PostControllers) GetUserPosts(c *gin.Context) {
	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	//Get User's all posts (service)
	posts, err := postService.GetUserPosts(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return posts
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})

}

func (ctr PostControllers) GetOnePost(c *gin.Context) {
	// Get postID from param
	postIDstring := c.Param("postID")
	postIDint64, err := strconv.ParseUint(postIDstring, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get postID",
		})
		return
	}
	postID := uint(postIDint64)

	//Get One Post (service)
	post, err := postService.GetOnePost(postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})

}

// TODO: 제목, 내용은 작성완료해서 DB에 넣었으나, 이미지가 실패했을 때 처리하기.
func (ctr PostControllers) UpdatePost(c *gin.Context) {
	//Get body by PostForm
	title := c.PostForm("title")
	body := c.PostForm("body")
	var postForm = forms.PostForm{
		Title: title,
		Body:  body,
	}

	//Get userID from token & Change type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	// Get postID from param
	postIDstring := c.Param("postID")
	postIDint64, err := strconv.ParseUint(postIDstring, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get postID",
		})
		return
	}
	postID := uint(postIDint64)

	//Update post (in DB)
	err = postService.UpdatePost(userID, postID, postForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Upload image
	var imageService services.ImageService
	//By form-data type, file is uploaded
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Open file
	filename := filepath.Base(file.Filename)
	filecontent, _ := file.Open()
	defer filecontent.Close()

	//Get imageURL
	imageURL, err := imageService.UploadPostImage(context.Background(), filecontent, userID, filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return imageURL
	c.JSON(http.StatusOK, gin.H{
		"Post Updated successfully with image!": imageURL,
	})

}

func (ctr PostControllers) DeletePost(c *gin.Context) {
	// Get postID from param
	postIDstring := c.Param("postID")
	postIDint64, err := strconv.ParseUint(postIDstring, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to get postID",
		})
		return
	}
	postID := uint(postIDint64)

	// Get userID from token & Chage type to uint
	userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	//Delete image(in google cloud storage & DB)
	var imageService services.ImageService
	err = imageService.DeleteImage(context.Background(), userID, postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Delete post(in DB)
	err = postService.DeletePost(userID, postID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Return message
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Success",
	})

}
