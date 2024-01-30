package server

import (
	"fmt"
	json "encoding/json"
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/Blogger/v1/types"
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
	context_body := context.Body()
	var p types.Page
	json.Unmarshal( context_body , &p )
	log.Debug( fmt.Sprintf( "Setting URL : %s" , p.URL ) )
	s.Set( "pages" , p.URL ,  p.HTMLB64 )
	return context.JSON( fiber.Map{
		"route": "/page/add" ,
		"result": true ,
	})
}

// so like https://github.com/quilljs/quill
// then we save the quill.root.innerHTML ? idk
// we have to do 2 lookups then , theres no way.
// first lookup == confirm static page exists , render html parent
// second lookup == GET JSON request sent by html parent for actual content
// or we could just render parent template as catch all , and then do 1 lookup for if any content exists ?
func ( s *Server ) PageHandler( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageHandler()" )
	sent_path := context.Path()
	sent_queries := context.Queries()
	page_html := s.Get( "pages" , sent_path )
	if page_html == "" {
		return context.JSON( fiber.Map{
			"route": "/*" ,
			"sent_path": sent_path ,
			"sent_queries": sent_queries ,
			"page_html": page_html ,
			"result": false ,
		})
	}
	return context.SendFile( "./v1/server/html/page.html" )
}