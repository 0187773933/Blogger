package server

import (
	fiber "github.com/gofiber/fiber/v2"
	// bolt "github.com/boltdb/bolt"
)

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

// Adds { key: HTML_STRING-b64 } to static-routes
func ( s *Server ) PageAdd( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageAdd()" )
	return context.JSON( fiber.Map{
		"route": "/page/add" ,
		"result": true ,
	})
}

func ( s *Server ) StaticHandler( context *fiber.Ctx ) ( error ) {
	log.Debug( "StaticHandler()" )
	sent_path := context.Path()
	sent_queries := context.Queries()

	saved_path := s.Get( "static-routes" , sent_path )

	return context.JSON( fiber.Map{
		"route": "/*" ,
		"sent_path": sent_path ,
		"sent_queries": sent_queries ,
		"saved_path": saved_path ,
		"result": true ,
	})
}