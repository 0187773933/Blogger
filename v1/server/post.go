package server

import (
	"bytes"
	"encoding/json"
	bolt "github.com/boltdb/bolt"
	types "github.com/0187773933/Blogger/v1/types"
)

func ( s *Server ) PostGetRange( min string , max string ) ( result []types.Post ) {
	s.DB.View( func( tx *bolt.Tx ) error {
		c := tx.Bucket( []byte( "posts" ) ).Cursor()
		for k , v := c.Seek( []byte( min ) ); k != nil && bytes.Compare( k , []byte( max ) ) <= 0; k , v = c.Next() {
			var p types.Post
			json.Unmarshal( v , &p )
			result = append( result , p )
		}
		return nil
	})
	return
}