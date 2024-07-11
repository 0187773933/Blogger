package server

import (
	// "fmt"
	// "time"
	// "strings"
	"encoding/json"
	bolt "github.com/boltdb/bolt"
	// fiber "github.com/gofiber/fiber/v2"
)

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