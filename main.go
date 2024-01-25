package main

import (
	"fmt"
	"os"
	"time"
	"os/signal"
	"syscall"
	"path/filepath"
	utils "github.com/0187773933/Blogger/v1/utils"
	server "github.com/0187773933/Blogger/v1/server"
	bolt_api "github.com/boltdb/bolt"
	logger "github.com/0187773933/Blogger/v1/logger"
)

var s server.Server

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "Shutting Down Blogger Server" )
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func main() {

	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )
	// fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )

	logger.Init()
	logger.Log.Printf( "Loaded Config File From : %s" , config_file_path )

	db , _ := bolt_api.Open( config.BoltDBPath , 0600 , &bolt_api.Options{ Timeout: ( 3 * time.Second ) } )
	fmt.Println( db )

	SetupCloseHandler()

	// utils.GenerateNewKeys()
	s = server.New( config )
	s.Start()

}
