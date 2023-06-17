package httpserver

import (
	"fmt"
	"net/http"
	"robinhood/config"
	"robinhood/internal/dto"
	"robinhood/internal/handlers/bloghdl"
	"robinhood/internal/handlers/userhdl"
	"robinhood/pkg/auth"
	"robinhood/pkg/meta"
	"strings"

	_ "robinhood/docs"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewHTTPServer(
	bh *bloghdl.Handler,
	uh *userhdl.Handler,
) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler

	// auth middleware
	authMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Get().JWT.Secret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTCustomClaims)
		},
	})

	// swagger
	if config.Get().App.EnableSwagger {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		fmt.Printf("Enable swager local: http://localhost:%s/swagger/index.html \n", config.Get().Endpoint.Port)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is running...")
	})

	// add prefix for all routes below
	v1 := e.Group("/api/v1")
	user := v1.Group("/user")
	user.POST("/register", uh.Register)
	user.POST("/login", uh.Login)
	user.PUT("", uh.UpdateUser, authMiddleware)

	blog := v1.Group("/blog", authMiddleware)
	blog.GET("", bh.ListBlog)
	blog.GET("/:blogId", bh.GetBlogByID)
	blog.POST("", bh.CreateBlog)
	blog.PUT("/:blogId", bh.UpdateBlogStatus)
	blog.DELETE("/:blogId", bh.ArchiveBlog)

	comment := v1.Group("/comment", authMiddleware)
	comment.GET("/:blogId", bh.ListComment)
	comment.POST("/:blogId", bh.CreateComment)

	return e
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var m *meta.MetaError

	if metaErr, ok := meta.IsError(err); ok {
		m = metaErr
	} else if he, ok := err.(*echo.HTTPError); ok {
		m = meta.NewError(he.Code).AppendMessage(1000, strings.ToLower(he.Message.(string)))
	} else {
		m = meta.MetaErrorInternalServer.AppendError(1000, err)
	}

	c.JSON(m.HttpStatus, dto.BaseErrorResponse{
		BaseResponse: dto.BaseResponse{
			Code: m.Code,
		},
		Message: m.Message,
	})
}
