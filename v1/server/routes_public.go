package server

import (
	"fmt"
	"time"
	"os"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	// bcrypt "golang.org/x/crypto/bcrypt"
	// encryption "github.com/0187773933/encryption/v1/encryption"
	// try "github.com/manucorporat/try"
)

func CDNMaxedOut( c *fiber.Ctx ) error {
	ip_address := c.IP()
	log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
	log.Info( log_message )
	c.Set( "Content-Type" , "text/html" )
	return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
}

var CDNLimter = rate_limiter.New( rate_limiter.Config{
	Max: 20 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: CDNMaxedOut ,
	LimiterMiddleware: rate_limiter.SlidingWindow{} ,
})

func PublicMaxedOut( c *fiber.Ctx ) error {
	ip_address := c.IP()
	log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
	log.Info( log_message )
	c.Set( "Content-Type" , "text/html" )
	return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
}

var PublicLimter = rate_limiter.New( rate_limiter.Config{
	Max: 3 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: PublicMaxedOut ,
	LimiterMiddleware: rate_limiter.SlidingWindow{} ,
})

func ( s *Server ) RenderHomePage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	admin_logged_in := s.ValidateAdmin( context )
	if admin_logged_in == true {
		content , _ := os.ReadFile( "./v1/server/html/admin.html" )
		html := strings.Replace( string( content ) , "</body>" ,
			fmt.Sprintf( "<script>window.AP = '%s';</script></body>" , s.Config.URLS.AdminPrefix ) , 1 )
		// return context.SendFile( x_url )
		context.Set( "Content-Type" , "text/html" )
		return context.SendString( html )
	}
	return context.SendFile( "./v1/server/html/home.html" )
}

func ( s *Server ) SetupPublicRoutes() {
	home_url := "/"
	// login_url := "/login"
	// logout_url := "/logout"
	if s.Config.URLS.Prefix != "" {
		home_url = fmt.Sprintf( "/%s" , s.Config.URLS.Prefix )
		// login_url = fmt.Sprintf( "/%s/login" , s.Config.URLS.Prefix )
		// logout_url = fmt.Sprintf( "/%s/logout" , s.Config.URLS.Prefix )
	}
	var admin_login_url string
	if s.Config.URLS.AdminLogin != "" {
		admin_login_url = fmt.Sprintf( "/%s" , s.Config.URLS.AdminLogin )
	} else {
		admin_login_url = "/admin/login"
	}
	fmt.Println( "Admin Login URL ===" , admin_login_url )
	admin_logout_url := fmt.Sprintf( "/%s/logout" , s.Config.URLS.AdminPrefix )
	s.FiberApp.Get( home_url , PublicLimter , s.RenderHomePage )
	// s.FiberApp.Get( login_url , PublicLimter , s.LoginPage )
	// s.FiberApp.Post( login_url , PublicLimter , s.Login )
	// s.FiberApp.Get( logout_url , PublicLimter , s.Logout )
	s.FiberApp.Get( admin_login_url , PublicLimter , s.LoginPage )
	s.FiberApp.Post( admin_login_url , PublicLimter , s.AdminLogin )
	s.FiberApp.Get( admin_logout_url , PublicLimter , s.AdminLogout )

	images_match_url := fmt.Sprintf( "/images/:uuid.:ext" )
	s.FiberApp.Get( images_match_url , PublicLimter , s.ServeImages )

	s.FiberApp.Get( "/page/get" , PublicLimter , s.PageGet )
	s.FiberApp.Get( "/pages/get/all" , PublicLimter , s.PagesGetAll )



	s.FiberApp.Get( "/*" , PublicLimter , s.PageHandler )
}