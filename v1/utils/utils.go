package utils

import (
	"fmt"
	"os"
	"net"
	"unicode"
	"strings"
	"encoding/json"
	"strconv"
	"time"
	"encoding/base64"
	filepath "path/filepath"
	yaml "gopkg.in/yaml.v3"
	ioutil "io/ioutil"
	runtime "runtime"
	ulid "github.com/oklog/ulid/v2"
	types "github.com/0187773933/Blogger/v1/types"
	// fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

func SetupStackTraceReport() {
	if r := recover(); r != nil {
		stacktrace := make( []byte , 1024 )
		runtime.Stack( stacktrace , true )
		fmt.Printf( "%s\n" , stacktrace )
	}
}

func GenerateNewKeys() {
	admin_login_url := encryption.GenerateRandomString( 16 )
	admin_prefix := encryption.GenerateRandomString( 6 )
	login_url := encryption.GenerateRandomString( 16 )
	prefix := encryption.GenerateRandomString( 6 )
	// https://github.com/gofiber/fiber/blob/main/middleware/encryptcookie/utils.go#L91
	// https://github.com/0187773933/encryption/blob/master/v1/encryption/encryption.go#L46
	// cookie_secret := fiber_cookie.GenerateKey()
	cookie_secret_bytes := encryption.GenerateRandomBytes( 32 )
	cookie_secret := base64.StdEncoding.EncodeToString( cookie_secret_bytes )
	cookie_secret_message := encryption.GenerateRandomString( 16 )
	admin_cookie_secret_message := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	api_key := encryption.GenerateRandomString( 16 )
	encryption_key := encryption.GenerateRandomString( 32 )
	bolt_prefix := encryption.GenerateRandomString( 6 )
	redis_prefix := encryption.GenerateRandomString( 6 )
	log_key := encryption.GenerateRandomString( 6 )
	log_encryption_key := encryption.GenerateRandomString( 32 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tURL - Admin Login === %s\n" , admin_login_url )
	fmt.Printf( "\tURL - Admin Prefix === %s\n" , admin_prefix )
	fmt.Printf( "\tURL - Login === %s\n" , login_url )
	fmt.Printf( "\tURL - Prefix === %s\n" , prefix )
	fmt.Printf( "\tCOOKIE - Secret === %s\n" , cookie_secret )
	fmt.Printf( "\tCOOKIE - USER - Message === %s\n" , cookie_secret_message )
	fmt.Printf( "\tCOOKIE - ADMIN - Message === %s\n" , admin_cookie_secret_message )
	fmt.Printf( "\tCREDS - Admin Username === %s\n" , admin_username )
	fmt.Printf( "\tCREDS - Admin Password === %s\n" , admin_password )
	fmt.Printf( "\tCREDS - API Key === %s\n" , api_key )
	fmt.Printf( "\tCREDS - Encryption Key === %s\n" , encryption_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n" , admin_password )
	fmt.Printf( "\tLOG - Log Key === %s\n" , log_key )
	fmt.Printf( "\tLOG - Encryption Key === %s\n" , log_encryption_key )
	fmt.Printf( "\tBOLT - Prefix === %s\n" , bolt_prefix )
	fmt.Printf( "\tREDIS - Prefix === %s\n" , redis_prefix )
	panic( "Exiting" )
}

func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	encountered := make( map[ string ]bool )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip := ipv4.String()
			if !encountered[ ip ] {
				encountered[ ip ] = true
				ip_addresses = append( ip_addresses , ip )
			}
		}
	}
	return
}

func ParseConfig( file_path string ) ( result types.Config ) {
	config_file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( config_file , &result )
	if error != nil { panic( error ) }
	return
}

func GetConfig() ( result types.Config ) {
	var config_file_path string
	if len( os.Args ) > 1 {
		config_file_path , _ = filepath.Abs( os.Args[ 1 ] )
	} else {
		config_file_path , _ = filepath.Abs( "./SAVE_FILES/config.yaml" )
		if _ , err := os.Stat( config_file_path ); os.IsNotExist( err ) {
			panic( "Config File Not Found" )
		}
	}
	result = ParseConfig( config_file_path )
	return
}

func PrettyPrint( input interface{} ) {
	jd , _ := json.MarshalIndent( input , "" , "  " )
	fmt.Println( string( jd ) )
}

func RemoveNonASCII( input string ) ( result string ) {
	for _ , i := range input {
		if i > unicode.MaxASCII { continue }
		result += string( i )
	}
	return
}
const SanitizedStringSizeLimit = 100
func SanitizeInputString( input string ) ( result string ) {
	trimmed := strings.TrimSpace( input )
	if len( trimmed ) > SanitizedStringSizeLimit { trimmed = strings.TrimSpace( trimmed[ 0 : SanitizedStringSizeLimit ] ) }
	result = RemoveNonASCII( trimmed )
	return
}

func WriteJSON( filePath string , data interface{} ) {
	file, _ := json.MarshalIndent( data , "" , " " )
	_ = ioutil.WriteFile( filePath , file , 0644 )
}

func UnixToULID( input_string string ) ( result string ) {
	x_input_in64 , _ := strconv.ParseInt( input_string , 10 , 64 )
	x_unix_time := time.Unix( x_input_in64 , 0 )
	x_ulid_time_stamp := ulid.Timestamp( x_unix_time )
	x_ulid , _ := ulid.New( x_ulid_time_stamp , nil )
	result = x_ulid.String()
	return
}