package server

import (


	fiber "github.com/gofiber/fiber/v2"
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

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )

	// Auth
	s.FiberApp.Get( "/login" , public_limiter , s.LoginGet )
	s.FiberApp.Post( "/login" , public_limiter , s.LoginPost )

	// Posts
	s.FiberApp.Post( "/post" , private_limiter , validate_session_mw , s.Post )
	// s.FiberApp.Get( "/post/:seq_id" , private_limiter , validate_session_mw , s.PostGetViaSeqID )
	s.FiberApp.Get( "/post/get/all" , private_limiter , validate_session_mw , s.PostGetAll )
	s.FiberApp.Get( "/post/get/:ulid" , private_limiter , validate_session_mw , s.PostGetViaULID )
	s.FiberApp.Get( "/post/get/range/unix/:start/:stop" , private_limiter , validate_session_mw , s.PostGetRangeViaUNIX )
	s.FiberApp.Get( "/post/get/range/ulid/:start/:stop" , private_limiter , validate_session_mw , s.PostGetRangeViaULID )
	// s.FiberApp.Get( "/post/:uuid" , private_limiter , validate_session_mw , s.PostGetViaUUID )

	// Uploads
	s.FiberApp.Post( "/upload" , private_limiter , validate_session_mw , s.Upload )
}