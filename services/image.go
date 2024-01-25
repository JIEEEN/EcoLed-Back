package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"gorm.io/gorm"
)

type ImageService struct{}

func (srv ImageService) UploadProfileImage(ctx context.Context, file io.Reader, userID uint, fileName string) (imageURL string, err error) {
	//Get filename
	uniqueFilename := time.Now().Format("20060102150405") + "_" + fileName

	// Upload image on Google Cloud Storage
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS1")))
	if err != nil {
		err = errors.New("failed to create client")
		return "", err
	}
	defer client.Close()

	// Upload an object with storage.Writer.
	bucketName := "ecoled_test_profile_images"
	wc := client.Bucket(bucketName).Object(uniqueFilename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		err = errors.New("failed to upload image1")
		return "", err
	}
	if err := wc.Close(); err != nil {
		err = errors.New("failed to upload image2")
		return "", err
	}

	//Get imageURL
	imageURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, uniqueFilename)

	//Save imageURL to DB
	result := initializers.DB.Model(&models.Profiles{}).
		Where("user_id = ?", userID).
		Update("profile_image", imageURL)
	if result.Error != nil {
		err = errors.New("failed to upload profileImage in db")
		return imageURL, err
	}

	//return imageURL
	return imageURL, nil

}

func (srv ImageService) UploadPostImage(tx *gorm.DB, ctx context.Context, file io.Reader, userID uint, fileName string) (imageURL string, err error) {
	//Get filename (in Service)
	uniqueFilename := time.Now().Format("20060102150405") + "_" + fileName

	// Upload image on Google Cloud Storage
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS1")))
	if err != nil {
		err = errors.New("failed to create client")
		return "", err
	}
	defer client.Close()

	// Upload an object with storage.Writer.
	bucketName := "ecoled_test_post_images"
	wc := client.Bucket(bucketName).Object(uniqueFilename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		err = errors.New("failed to upload image1")
		return "", err
	}
	if err := wc.Close(); err != nil {
		err = errors.New("failed to upload image2")
		return "", err
	}

	//Get imageURL
	imageURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, uniqueFilename)

	//Save imageURL to DB
	var post models.Posts
	result := tx.Where("user_id = ?", userID).Order("updated_at DESC").First(&post)
	if result.Error != nil {
		err = errors.New("failed to get last uploaded post")
		return imageURL, err
	}
	post.Image = imageURL
	result = tx.Save(&post)
	if result.Error != nil {
		err = errors.New("failed to upload image in db")
		return imageURL, err
	}

	//return imageURL
	return imageURL, nil

}

func (svc ImageService) DeleteProfileImage(ctx context.Context, userID uint) (err error) {
	//Delete image on Google Cloud Storage
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS1")))
	if err != nil {
		err = errors.New("failed to create client")
		return err
	}
	defer client.Close()

	//Find filename
	var profile models.Profiles
	result := initializers.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		err = errors.New("failed to get profile")
		return err
	}
	fileNameStr := profile.Profile_image
	fileName := path.Base(fileNameStr)

	//Delete an object with storage.Writer.
	bucketName := "ecoled_test_profile_images"
	if err := client.Bucket(bucketName).Object(fileName).Delete(ctx); err != nil {
		err = errors.New("failed to delete image")
		return err
	}

	//Delete imageURL in DB
	profile.Profile_image = ""
	result = initializers.DB.Save(&profile)
	if result.Error != nil {
		err = errors.New("failed to delete image in db")
		return err
	}

	return nil

}

func (svc ImageService) DeletePostImage(ctx context.Context, userID uint, postID uint) (err error) {
	//Delete image on Google Cloud Storage
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS1")))
	if err != nil {
		err = errors.New("failed to create client")
		return err
	}
	defer client.Close()

	//Find filename
	var post models.Posts
	result := initializers.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		err = errors.New("failed to get post")
		return err
	}
	fileNameStr := post.Image
	fileName := path.Base(fileNameStr)

	//Delete an object with storage.Writer.
	bucketName := "ecoled_test_post_images"
	if err := client.Bucket(bucketName).Object(fileName).Delete(ctx); err != nil {
		err = errors.New("failed to delete image")
		return err
	}

	//Delete imageURL in DB
	post.Image = ""
	result = initializers.DB.Save(&post)
	if result.Error != nil {
		err = errors.New("failed to delete image in db")
		return err
	}

	return nil

}
