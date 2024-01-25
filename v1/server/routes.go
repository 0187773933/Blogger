package server

import (
	"fmt"
	"time"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
)

var public_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 1 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		fmt.Println( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

var private_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 3 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c *fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		fmt.Println( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

func ( s *Server ) Home( context *fiber.Ctx ) ( error ) {
	session := validate_session( context )
	if session == false {
		return context.JSON( fiber.Map{
			"route": "/" ,
			"source": "https://github.com/0187773933/Blogger" ,
		})
	}
	log.Debug( "Logged In User , Sending Home Page" )
	return context.SendFile( "./v1/server/html/home.html" )
}

func ( s *Server ) Post( context *fiber.Ctx ) ( error ) {
	return context.JSON( fiber.Map{
		"url": "/post" ,
		"method": "POST" ,
		"result": true ,
	})
}

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )
	s.FiberApp.Get( "/login" , public_limiter , s.LoginGet )
	s.FiberApp.Post( "/login" , public_limiter , s.LoginPost )
	s.FiberApp.Post( "/post" , private_limiter , validate_session_mw , s.Post )
}