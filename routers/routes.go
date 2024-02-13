package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	
	"github.com/Eco-Led/EcoLed-Back_test/controllers"
	"github.com/Eco-Led/EcoLed-Back_test/middlewares"

	"github.com/gin-gonic/gin"
)

var userController = new(controllers.UserControllers)
var AuthController = new(controllers.AuthControllers)
var oauthController = new(controllers.OauthControllers)
var profileController = new(controllers.ProfileControllers)
var profileimageController = new(controllers.ProfileImageControllers)
var accountController = new(controllers.AccountControllers)
var paylogController = new(controllers.PaylogControllers)
var rankingController = new(controllers.RankingControllers)
var postController = new(controllers.PostControllers)


func UserRoutes(router *gin.Engine, apiVersion string) {
	router.POST(apiVersion+"/login", userController.Login)
	router.POST(apiVersion+"/register", userController.Register)
	router.POST(apiVersion+"/logout", userController.Logout)
	router.POST(apiVersion+"/email", userController.FindEmail)
	router.POST(apiVersion+"/password", userController.FindPassword)
	router.POST(apiVersion+"/code", userController.VerifyCode)
	router.PUT(apiVersion+"/password", userController.UpdatePassword)
}

func AuthRoutes(router *gin.Engine, apiVersion string) {
	router.POST(apiVersion+"/refresh", AuthController.RefreshToken)
}

func OauthRoutes(router *gin.Engine, apiVersion string) {
	// 세션 스토어 설정
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("oauth_session", store))

	router.GET(apiVersion+"/oauth/google", oauthController.GoogleLogin)
	router.GET(apiVersion+"/oauth/google/callback", oauthController.GoogleCallback)
	router.POST(apiVersion+"/oauth/google/register", oauthController.GoogleRegister)
}

func ProfileRoutes(router *gin.Engine, apiVersion string) {
	router.Use(middlewares.AuthToken())
	router.PUT(apiVersion+"/profile", profileController.UpdateProfile)
	router.GET(apiVersion+"/profile/:userID", profileController.GetUserProfile)
	router.GET(apiVersion+"/profile", profileController.GetMyProfile)
	router.GET(apiVersion+"/profile/search", profileController.GetProfileByNickname)
}

func ProfileImageRoutes(router *gin.Engine, apiVersion string) {
	router.Use(middlewares.AuthToken())
	router.POST(apiVersion+"/profileimage", profileimageController.UploadProfileImage)
	router.DELETE(apiVersion+"/profileimage", profileimageController.DeleteProfileImage)
}

func AccountRoutes(router *gin.Engine, apiVersion string) {
	router.Use(middlewares.AuthToken())
	router.GET(apiVersion+"/account", accountController.GetAccount)
}

func PaylogRoutes(router *gin.Engine, apiVersion string) {
	router.Use(middlewares.AuthToken())
	router.POST(apiVersion+"/paylog", paylogController.CreatePaylog)
	router.GET(apiVersion+"/paylog/:paylogID", paylogController.GetPaylog)
	router.GET(apiVersion+"/paylog", paylogController.GetPaylogs)
	router.PUT(apiVersion+"/paylog/:paylogID", paylogController.UpdatePaylog)
	router.DELETE(apiVersion+"/paylog/:paylogID", paylogController.DeletePaylog)
}

func RankingRoutes(router *gin.Engine, apiVersion string) {
	router.GET(apiVersion+"/ranking", rankingController.GetRanking)
}

func PostRoutes(router *gin.Engine, apiVersion string) {
	router.Use(middlewares.AuthToken())
	router.POST(apiVersion+"/post", postController.CreatePost)
	router.GET(apiVersion+"/post/:postID", postController.GetPost)
	router.GET(apiVersion+"/post", postController.GetMyPost)
	router.GET(apiVersion+"/posts/:userID", postController.GetUserPost)
	router.PUT(apiVersion+"/post/:postID", postController.UpdatePost)
	router.DELETE(apiVersion+"/post/:postID", postController.DeletePost)
}

func RouterSetupV1() *gin.Engine {
	r := gin.Default()

	

	apiVersion := "/api/v1"
	r.Group(apiVersion)
	{
		UserRoutes(r, apiVersion)
		AuthRoutes(r, apiVersion)
		OauthRoutes(r, apiVersion)
		ProfileRoutes(r, apiVersion)
		ProfileImageRoutes(r, apiVersion)
		AccountRoutes(r, apiVersion)
		PaylogRoutes(r, apiVersion)
		RankingRoutes(r, apiVersion)
		PostRoutes(r, apiVersion)
	}

	return r
}
