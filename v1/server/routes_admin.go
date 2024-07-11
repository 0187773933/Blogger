package server

import (
	"fmt"
	"strconv"
	// net_url "net/url"
	// bolt_api "github.com/boltdb/bolt"
	// encryption "github.com/0187773933/encryption/v1/encryption"
	fiber "github.com/gofiber/fiber/v2"
	// logger "github.com/0187773933/Logger/v1/logger"
)

func ( s *Server ) GetLogMessages( c *fiber.Ctx ) ( error ) {
	count := c.Query( "count" )
	if count == "" {
		count = c.Query( "c" )
	}
	count_int , _ := strconv.Atoi( count )
	if count_int == 0 {
		count_int = -1
	}
	log.Debug( fmt.Sprintf( "Count === %d" , count_int ) )
	messages := log.GetMessages( count_int )
	return c.JSON( fiber.Map{
		"result": true ,
		"url": "/log/:count" ,
		"count": count ,
		"messages": messages ,
	})
}

func ( s *Server ) SetupAdminRoutes() {
	cdn_group := s.FiberApp.Group( "/cdn" )
	cdn_group.Use( CDNLimter )
	cdn_group.Use( s.ValidateAdminMW )
	s.FiberApp.Static( "/cdn" , "./v1/server/cdn" )
	admin := s.FiberApp.Group( fmt.Sprintf( "/%s" , s.Config.URLS.AdminPrefix ) )
	admin.Use( s.ValidateAdminMW )
	admin.Get( "/log/view" , s.GetLogMessages )
	admin.Get( "/page/add" , s.PageAddGetWYSIWYG )
	admin.Get( "/page/edit/*" , s.PageAddGetWYSIWYGEdit )
	admin.Get( "/page/delete" , s.PageDelete )
	admin.Post( "/page/add" , s.PageAdd )
	// admin.Get( "/page/get" , s.PageGet )
	admin.Get( "/pages/get/all" , s.PagesGetAll )
	admin.Post( "/pages/update/order" , s.PagesUpdateOrder )

}