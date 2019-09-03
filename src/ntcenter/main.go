package main

import (
	"net/http"
	"time"

	"common"
	"github.com/beinan/fastid"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.JWT([]byte(SECRET_STRING)))

	// Routes
	//e.GET("/", hello)
	e.POST("/getToken", buildJWT)
	e.GET("/getID",buildID)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e:=err.(*echo.HTTPError)
		c.String(e.Code,e.Message.(string))
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func buildID(c echo.Context) error  {
	id := fastid.CommonConfig.GenInt64ID()
	return c.JSON(http.StatusOK,map[string]int64{
		"id": id,
	})
	return nil
}

func buildJWT(c echo.Context) error {
	username := c.FormValue("name")
	password := c.FormValue("password")
	_ = password
	if true /*username == "jon" && password == "jon"*/ {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["iss"] = "ntcenter"
		claims["sub"] = "use esb api"
		//claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour).Unix()

		t, err := token.SignedString([]byte(common.JWTSecretKey))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}
