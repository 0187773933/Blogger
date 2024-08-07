const uuid_v4_regex = /^[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i
function is_uuid( str ) { return uuid_v4_regex.test( str ); }
const barcode_regex = /^\d+$/;
function is_barcode( str ) { return barcode_regex.test( str ); }
function sleep( ms ) { return new Promise( resolve => setTimeout( resolve , ms ) ); }

function convert_milliseconds_to_time_string( milliseconds ) {
	let seconds = Math.floor( milliseconds / 1000 );
	let minutes = Math.floor( seconds / 60 );
	let hours = Math.floor( minutes / 60 );
	let days = Math.floor( hours / 24 );
	hours %= 24;
	minutes %= 60;
	seconds %= 60;

	let time_string = `${days} days , ${hours} hours , ${minutes} minutes , and ${seconds} seconds`;
	return time_string;
}

function set_nested_property( obj , keys , value ) {
	if ( keys.length === 1 ) {
		obj[ keys[ 0 ] ] = value;
	} else {
		const key = keys.shift();
		obj[ key ] = obj[ key ] || {};
		set_nested_property( obj[ key ] , keys , value );
	}
}

function add_qr_code( text , element_id ) {
	let x_element = document.getElementById( element_id );
	x_element.innerHTML = "";
	let user_qrcode = new QRCode( x_element , {
		text: text ,
		width: 256 ,
		height: 256 ,
		colorDark : "#000000" ,
		colorLight : "#ffffff" ,
		correctLevel : QRCode.CorrectLevel.H
	});
}

function set_url( new_url ) {
	// no page reload ?
	console.log( `Changing URL , FROM = ${window.location.href} || TO = ${new_url}` );
	window.history.pushState( null , null , new_url );

	// Update the query parameters
	// url.searchParams.set("q", "example");

	// Update the URL with a full page reload
	// window.location.href = url.toString();
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