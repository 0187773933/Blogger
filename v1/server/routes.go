package server

import (
	"fmt"
	"time"
	"strconv"
	"io/ioutil"
	"encoding/json"
	ulid "github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
	bolt "github.com/boltdb/bolt"
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
	context_body := context.Body()
	var p types.Post
	json.Unmarshal( context_body , &p )
	p.Date = utils.GetFormattedTimeString()
	p.UUID = uuid.NewV4().String()
	p.ULID = ulid.Make().String()
	post_json_bytes , _ := json.Marshal( p )
	s.DB.Update( func( tx *bolt.Tx ) error {
		posts_bucket , _ := tx.CreateBucketIfNotExists( []byte( "posts" ) )
		posts_bucket.Put( []byte( p.ULID ) , post_json_bytes )
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post" ,
		"method": "POST" ,
		"post": p ,
		"result": true ,
	})
}

func ( s *Server ) PostGetViaSeqID( context *fiber.Ctx ) ( error ) {
	// var p types.Post
	// json.Unmarshal( context_body , &p )
	seq_id := context.Params( "seq_id" )
	seq_id_int , _ := strconv.Atoi( seq_id )
	// seq_id_index_int := ( seq_id_int - 1 )
	// seq_id_index_string := strconv.Itoa( seq_id_index_int )
	// var post_string string
	var t_post types.Post
	fmt.Println( "Seq ID ===" , seq_id )
	s.DB.View( func( tx *bolt.Tx ) error {
		// c := tx.Bucket( []byte( "posts" ) ).Cursor()
		// k , v := c.Seek( []byte( seq_id ) )
		// fmt.Println( k , v )
		b := tx.Bucket( []byte( "posts" ) )
		c := b.Cursor()
		_ , v := c.Seek( []byte( seq_id ) )
		fmt.Println( string( v ) )
		if v == nil {
			_ , v := c.Prev() // you have to do it this way. there is no c.Curr(). and c.Seek always goes +1 somehow
			json.Unmarshal( v , &t_post )
		}

		// post_string = string( v )
		// fmt.Printf( "key=%s, value=%s\n" , k , v )
		// for k, v := c.First(); k != nil; k, v = c.Next() {
		// 	fmt.Printf("key=%s, value=%s\n", k, v)
		// }
		return nil
	})
	if t_post.SeqID != seq_id_int {
		return context.JSON( fiber.Map{
			"url": "/post/:seq_id" ,
			"method": "GET" ,
			"result": false ,
		})
	}
	return context.JSON( fiber.Map{
		"url": "/post/:seq_id" ,
		"method": "GET" ,
		"post": t_post ,
		"result": true ,
	})
}

// func ( s *Server ) PostGetViaUUID( context *fiber.Ctx ) ( error ) {
// 	var p types.Post
// 	json.Unmarshal( context_body , &p )
// 	s.DB.View( func( tx *bolt.Tx ) error {

// 		c := tx.Bucket( []byte( "posts" ) ).Cursor()
// 		k , v := c.Seek( min )

// 		posts_bucket , _ := tx.CreateBucketIfNotExists( []byte( "posts" ) )
// 		post_id , _ := posts_bucket.NextSequence()
// 		posts_bucket.Put( utils.IToB( post_id ) , context_body )
// 		return nil
// 	})
// 	return context.JSON( fiber.Map{
// 		"url": "/post/:uuid" ,
// 		"method": "GET" ,
// 		"result": true ,
// 	})
// }

func ( s *Server ) PostGetAll( context *fiber.Ctx ) ( error ) {

	s.DB.View( func( tx *bolt.Tx ) error {
		c := tx.Bucket( []byte( "posts" ) ).Cursor()
		for k , v := c.First(); k != nil; k , v = c.Next() {
			var p types.Post
			json.Unmarshal( v , &p )
			t , _ := ulid.Parse( p.ULID )
			fmt.Println( t.Time() )
		}
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post/get/all" ,
		"method": "GET" ,
		"result": true ,
	})
}

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )

	// Auth
	s.FiberApp.Get( "/login" , public_limiter , s.LoginGet )
	s.FiberApp.Post( "/login" , public_limiter , s.LoginPost )

	// Posts
	s.FiberApp.Post( "/post" , private_limiter , validate_session_mw , s.Post )
	s.FiberApp.Get( "/post/:seq_id" , private_limiter , validate_session_mw , s.PostGetViaSeqID )
	s.FiberApp.Get( "/post/get/all" , private_limiter , validate_session_mw , s.PostGetAll )
	// s.FiberApp.Get( "/post/:uuid" , private_limiter , validate_session_mw , s.PostGetViaUUID )

	// Uploads
	s.FiberApp.Post( "/upload" , private_limiter , validate_session_mw , s.Upload )
}