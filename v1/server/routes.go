package server

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	types "github.com/0187773933/Blogger/v1/types"
	utils "github.com/0187773933/Blogger/v1/utils"
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

func ( s *Server ) Upload( context *fiber.Ctx ) ( error ) {
	form , err := context.MultipartForm()
	if err != nil {
		return context.Status( fiber.StatusBadRequest ).SendString( "Error Parsing Form" )
	}
	var file_data types.FileData
	bytes , exists := form.File[ "bytes" ]
	if exists && len( bytes ) > 0 {
		file , err := bytes[ 0 ].Open()
		defer file.Close()
		if err == nil {
			data , err := ioutil.ReadAll( file )
			if err == nil {
				file_data.FileName = bytes[ 0 ].Filename
				file_data.Data = data
			}
		}
	}
	if file_data.Data != nil {
		fmt.Println( "Uploaded Bytes from :" , file_data.FileName )
		fmt.Println( file_data.Data )
	}
	return context.JSON( fiber.Map{
		"url": "/upload" ,
		"method": "POST" ,
		"file": file_data ,
		"result": true ,
	})
}

func ( s *Server ) Post( context *fiber.Ctx ) ( error ) {
	form , err := context.MultipartForm()
	if err != nil {
		return context.Status( fiber.StatusBadRequest ).SendString( "Error Parsing Form" )
	}
	json_data , exists := form.Value[ "json" ]
	if !exists || len( json_data ) == 0 {
		return context.Status( fiber.StatusBadRequest ).SendString( "Missing JSON data" )
	}
	var p types.Post
	err = json.Unmarshal( []byte( json_data[ 0 ] ) , &p )
	if err != nil {
	    return context.Status(fiber.StatusBadRequest).SendString( "Error Parsing JSON Data: " + err.Error() )
	}
	utils.PrettyPrint( p )
	return context.JSON( fiber.Map{
		"url": "/post" ,
		"method": "POST" ,
		"post": p ,
		"result": true ,
	})
}

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )
	s.FiberApp.Get( "/login" , public_limiter , s.LoginGet )
	s.FiberApp.Post( "/login" , public_limiter , s.LoginPost )
	s.FiberApp.Post( "/post" , private_limiter , validate_session_mw , s.Post )
	s.FiberApp.Post( "/upload" , private_limiter , validate_session_mw , s.Upload )
}