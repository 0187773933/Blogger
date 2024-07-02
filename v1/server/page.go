package server

import (
	"fmt"
	json "encoding/json"
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/Blogger/v1/types"
	uuid "github.com/satori/go.uuid"
	ulid "github.com/oklog/ulid/v2"
	bolt "github.com/boltdb/bolt"
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
	p.UUID = uuid.NewV4().String()
	p.ULID = ulid.Make().String()
	p.Created = s.GetFormattedTimeString()
	p.SortedOrder = -1
	log.Debug( fmt.Sprintf( "Storing Content for URL : %s" , p.URL ) )
	// s.Set( "pages" , p.URL ,  p.HTMLB64 )
	s.SetOBJ( "pages" , p.URL , p )
	return context.JSON( fiber.Map{
		"route": "/page/add" ,
		"page": p ,
		"result": true ,
	})
}

func ( s *Server ) PageDelete( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageDelete()" )
	x_url := context.Query( "url" )
	s.DB.Update( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( "pages" ) )
		if b == nil { return nil }
		b.Delete( []byte( x_url ) )
		return nil
	})
	return context.JSON( fiber.Map{
		"route": "/page/delete" ,
		"url": x_url ,
		"result": true ,
	})
}

func ( s *Server ) PageAddGetWYSIWYG( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/page_add_wysiwyg.html" )
}

func ( s *Server ) PageAddGetWYSIWYG2( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/page_add_wysiwyg_2.html" )
}

func ( s *Server ) PageAddGetWYSIWYG3( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/page_add_wysiwyg_3.html" )
}

func ( s *Server ) PageAddGetWYSIWYG3Edit( context *fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/page_add_wysiwyg_3_edit.html" )
}

func ( s *Server ) PageGet( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageGet()" )
	// x_url := context.Params( "url" )
	x_url := context.Query( "url" )
	// page_html_b64 := s.Get( "pages" , x_url )
	var p types.Page
	s.DB.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( "pages" ) )
		if b == nil { return nil }
		v := b.Get( []byte( x_url ) )
		if v == nil { return nil }
		err := json.Unmarshal( v , &p )
		if err != nil {
			log.Debug( err )
			return nil
		}
		return nil
	})
	return context.JSON( fiber.Map{
		"route": "/page/get/:url" ,
		"page": p ,
		"result": true ,
	})
}

func ( s *Server ) PageGetAll( context *fiber.Ctx ) ( error ) {
	var pages []types.Page
	s.DB.View( func( tx *bolt.Tx ) error {
		c := tx.Bucket( []byte( "pages" ) ).Cursor()
		for k , v := c.First(); k != nil; k , v = c.Next() {
			var p types.Page
			json.Unmarshal( v , &p )
			if p.UUID == "" { continue }
			pages = append( pages , p )
		}
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post/get/all" ,
		"method": "GET" ,
		"pages": pages ,
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
	sent_path := context.Path()
	log.Debug( fmt.Sprintf( "PageHandler( %s )" , sent_path ) )
	// sent_queries := context.Queries()
	// page_html := s.Get( "pages" , sent_path )
	// if page_html == "" {
	// 	return context.JSON( fiber.Map{
	// 		"route": "/*" ,
	// 		"sent_path": sent_path ,
	// 		"sent_queries": sent_queries ,
	// 		"page_html": page_html ,
	// 		"result": false ,
	// 	})
	// }
	return context.SendFile( "./v1/server/html/page.html" )
}