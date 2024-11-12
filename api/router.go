package api

import (
	_ "food/api/docs"
	"food/api/handler"
	"food/config"
	"food/pkg/logger"
	"food/service"
	"food/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewApi(r *gin.Engine, cfg *config.Config, storage storage.IStorage, logger logger.LoggerI, service service.IServiceManager) {
	h := handler.NewStrg(logger, storage, cfg, service)
	r.Use(customCORSMiddleware())
	// r.Use(authMiddleware)
	v1 := r.Group("/food/api/v1")

	v1.POST("/uploadfiles", h.UploadFiles)
	v1.DELETE("/deletefiles", h.DeleteFile)

	v1.POST("/sendcode", h.UserRegister)
	v1.POST("/user/register", h.Register)
	// v1.POST("/user/login", h.UserLogin)
	v1.POST("/user/byphoneconfirm", h.UserLoginByPhoneConfirm)

	// v1.POST("/admin/sendcode", h.AdminRegister)
	// v1.POST("/admin/verifycode", h.AdminRegisterConfirm)
	v1.POST("/admin/login", h.AdminLogin)

	v1.POST("/combo", h.CreateCombo)
	v1.GET("/getallcombos", h.GetAllCombos)
	v1.GET("/getcombo/:id", h.GetCombo)
	v1.PUT("/updatecombo/:id", h.UpdateCombo)

	v1.POST("/category", h.CreateCategory)
	v1.GET("/getbycategory/:id", h.GetCategoryByID)
	v1.GET("/getallcategory", h.GetAllCategories)
	v1.PUT("/category/:id", h.UpdateCategory)
	v1.DELETE("/deletecategory", h.DeleteCustomer)

	v1.POST("/order", h.CreateOrder)
	v1.GET("/getorder/:id", h.GetOrder)
	v1.GET("/getallorders", h.GetAllOrders)
	v1.PUT("/updateorder", h.UpdateOrder)
	v1.DELETE("/deleteorder/:id", h.DeleteOrder)
	v1.PATCH("/orderStatus/:id", h.ChangeOrderStatus)

	v1.POST("/createadmin", h.CreateAdmin)
	v1.GET("/getbyidadmin/:id", h.GetAdminByID)
	v1.GET("/getalladmins", h.GetAllAdmins)
	v1.PUT("/updateadmin/:id", h.UpdateAdmin)
	v1.DELETE("/deleteadmin/:id", h.DeleteAdmin)

	v1.POST("/createuser", h.CreateUser)
	v1.GET("/getbyiduser/:id", h.GetUserByID)
	v1.GET("/getallusers", h.GetAllUsers)
	v1.PUT("/updateuser/:id", h.UpdateUser)
	v1.DELETE("/deleteuser/:id", h.DeleteUser)

	v1.POST("/createproduct", h.CreateProduct)
	v1.GET("/getproduct/:id", h.GetProductByID)
	v1.GET("/getallproducts", h.GetAllProducts)
	v1.PUT("/updateproduct/:id", h.UpdateProduct)
	v1.DELETE("/deleteproduct/:id", h.DeleteProduct)

	v1.POST("/createbranch", h.CreateBranch)
	v1.GET("/getbranch/:id", h.GetBranchByID)
	v1.GET("/getallbranches", h.GetAllBranches)
	v1.PUT("/updatebranch/:id", h.UpdateBranch)
	v1.DELETE("/deletebranch/:id", h.DeleteBranch)

	v1.POST("/createbanner", h.CreateBanner)
	v1.GET("/getallbanners", h.GetAllBanners)
	v1.DELETE("/deletebanner", h.DeleteBanner)

	url := ginSwagger.URL("swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// func authMiddleware(c *gin.Context) {
// 	auth := c.GetHeader("Authorization")
// 	if auth == "" {
// 		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
// 	}

// 	c.Next()
// }
