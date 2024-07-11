package server

import (
	"fmt"
	"os"
	// "time"
	// "strings"
	"path/filepath"
	"encoding/json"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
)

func FileExists( name string ) ( result bool ) {
	result = false
	_ , err := os.Stat( name )
	if os.IsNotExist( err ) {
		return
	}
	if err != nil {
		return
	}
	result = true
	return
}

func ( s *Server ) Set( bucket_name string , key string , value string ) {
	s.DB.Update( func( tx *bolt.Tx ) error {
		b , err := tx.CreateBucketIfNotExists( []byte( bucket_name ) )
		if err != nil { log.Debug( err ); return nil }
		err = b.Put( []byte( key ) , []byte( value ) )
		if err != nil { log.Debug( err ); return nil }
		return nil
	})
	return
}

func ( s *Server ) Get( bucket_name string , key string ) ( result string ) {
	s.DB.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( bucket_name ) )
		if b == nil { return nil }
		v := b.Get( []byte( key ) )
		if v == nil { return nil }
		result = string( v )
		return nil
	})
	return
}

func ( s *Server ) SetOBJ( bucket_name string , key string , obj interface{} ) {
	obj_json , err := json.Marshal( obj )
	if err != nil {
		log.Debug( err )
		return
	}
	s.DB.Update( func( tx *bolt.Tx ) error {
		b , err := tx.CreateBucketIfNotExists( []byte( bucket_name ) )
		if err != nil { log.Debug( err ); return nil }
		err = b.Put( []byte( key ) , obj_json )
		if err != nil { log.Debug( err ); return nil }
		return nil
	})
	return
}

func ( s *Server ) GetOBJ( bucket_name string , key string ) ( result interface{} ) {
	s.DB.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( bucket_name ) )
		if b == nil { return nil }
		v := b.Get( []byte( key ) )
		if v == nil { return nil }
		err := json.Unmarshal( v , &result )
		if err != nil {
			log.Debug( err )
			return nil
		}
		return nil
	})
	return
}


func ( s *Server ) ServeImages( context *fiber.Ctx ) ( error ) {
	uuid := context.Params( "uuid" )
	ext := context.Params( "ext" )
	x_path := filepath.Join( s.Config.ImagesSavePath , fmt.Sprintf( "%s.%s" , uuid , ext ) )
	fmt.Println( "ServeImages() -->" , x_path )
	if FileExists( x_path ) == false {
		return context.Status( fiber.StatusInternalServerError ).SendString( "File Doesn't Exist" )
	}
	return context.SendFile( x_path , false )
}