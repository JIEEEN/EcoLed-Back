package services

import (
	"errors"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"
	"gorm.io/gorm"
)

type PostService struct{}

func (srv PostService) CreatePost(tx *gorm.DB, userID uint, postForm forms.PostForm) error {
	// Create post
	result := tx.Create(&models.Posts{
		Title:   postForm.Title,
		Body:    postForm.Body,
		User_id: userID,
	})
	if result.Error != nil {
		err := errors.New("failed to create post")
		return err
	}

	return nil

}

func (srv PostService) GetUserPost(userID uint) ([]models.Posts, error) {
	//Get all posts
	var posts []models.Posts
	result := initializers.DB.Where("user_id =?", userID).
		Where("deleted_at is NULL").
		Find(&posts)
	if result.Error != nil {
		err := errors.New("failed to get all posts")
		return nil, err
	}
	if result.RowsAffected == 0 {
		err := errors.New("there are no posts")
		return nil, err
	}

	return posts, nil

}

func (srv PostService) GetPost(postID uint) (models.Posts, error) {
	//Get one post
	var post models.Posts
	result := initializers.DB.
		Where("deleted_at is NULL").
		First(&post, postID)
	if result.Error != nil {
		err := errors.New("failed to get post")
		return post, err
	}
	if result.RowsAffected == 0 {
		err := errors.New("there are no post")
		return post, err
	}

	return post, nil

}

func (srv PostService) UpdatePost(tx *gorm.DB, userID uint, postID uint, postForm forms.PostForm) error {
	//Check whether post is
	var post models.Posts
	result := tx.First(&post, postID)
	if result.Error != nil {
		err := errors.New("there are no post that match postID")
		return err
	}

	//Check whether post is user's post
	if post.User_id != userID {
		err := errors.New("you are not this post's creater")
		return err
	}

	//Update post
	post.Title = postForm.Title
	post.Body = postForm.Body
	result = tx.Save(&post)
	if result.Error != nil {
		err := errors.New("failed to update post")
		return err
	}

	return nil

}

func (srv PostService) DeletePost(userID uint, postID uint) error {
	//Check whether post is
	var post models.Posts
	result := initializers.DB.First(&post, postID)
	if result.Error != nil {
		err := errors.New("there are no post that match postID")
		return err
	}

	//Check whether post is user's post
	if post.User_id != userID {
		err := errors.New("you are not this post's creater")
		return err
	}

	//Delete post
	result = initializers.DB.Delete(&post, postID)
	if result.Error != nil {
		err := errors.New("failed to delete post")
		return err
	}

	return nil

}
