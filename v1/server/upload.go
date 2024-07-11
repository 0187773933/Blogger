package server

import (
	"fmt"
	"io/ioutil"
	filepath "path/filepath"
	uuid "github.com/satori/go.uuid"
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/Blogger/v1/types"
)

func (s *Server) UploadImage(context *fiber.Ctx) error {
    form, err := context.MultipartForm()
    if err != nil {
        return context.Status(fiber.StatusBadRequest).SendString("Error Parsing Form")
    }
    var x_uuid string
    var fileData types.FileData
    var ext string
    var originalName string
    files, exists := form.File["file"]
    if exists && len(files) > 0 {
        file, err := files[0].Open()
        defer file.Close()
        if err == nil {
            data, err := ioutil.ReadAll(file)
            if err == nil {
                originalName = files[0].Filename
                ext = filepath.Ext(files[0].Filename)
                x_uuid = uuid.NewV4().String()
                newFileName := x_uuid + ext
                savePath := filepath.Join( s.Config.ImagesSavePath , newFileName )

                err = ioutil.WriteFile(savePath, data, 0644)
                if err != nil {
                    return context.Status(fiber.StatusInternalServerError).SendString("Error Saving File")
                }

                fileData.FileName = files[0].Filename

                fmt.Println("Uploaded image:", x_uuid, files[0].Filename)
                fmt.Println("Saved as:", savePath)
            }
        }
    }
    return context.JSON(fiber.Map{
        "success": 1,
        "file": fiber.Map{
            "url":      fmt.Sprintf("/images/%s%s", x_uuid, ext),
            "caption":    originalName,
            "extension": ext,
        },
    })
}

func (s *Server) Upload(context *fiber.Ctx) error {
	form, err := context.MultipartForm()
	if err != nil {
		return context.Status(fiber.StatusBadRequest).SendString("Error Parsing Form")
	}
	var x_uuid string
	var fileData types.FileData
	var ext string
	var original_name string
	files, exists := form.File["file"]
	if exists && len(files) > 0 {
		file, err := files[0].Open()
		defer file.Close()
		if err == nil {
			data, err := ioutil.ReadAll(file)
			if err == nil {
				original_name = files[0].Filename
				ext = filepath.Ext(files[0].Filename)
				x_uuid = uuid.NewV4().String()
				newFileName := x_uuid + ext
				savePath := filepath.Join( s.Config.UploadsSavePath , newFileName)

				err = ioutil.WriteFile(savePath, data, 0644)
				if err != nil {
					return context.Status(fiber.StatusInternalServerError).SendString("Error Saving File")
				}

				fileData.FileName = files[0].Filename

				fmt.Println("Uploaded file:", x_uuid , files[0].Filename)
				fmt.Println("Saved as:", savePath)
			}
		}
	}
	return context.JSON(fiber.Map{
		"success": 1 ,
		"file": fiber.Map{
			"url": fmt.Sprintf( "/files/%s%s" , x_uuid , ext ) ,
			"title": original_name ,
			"extension": ext ,
		} ,
	})
}

// func ( s *Server ) Upload( context *fiber.Ctx ) ( error ) {
// 	form , err := context.MultipartForm()
// 	if err != nil {
// 		return context.Status( fiber.StatusBadRequest ).SendString( "Error Parsing Form" )
// 	}
// 	var file_data types.FileData
// 	bytes , exists := form.File[ "bytes" ]
// 	if exists && len( bytes ) > 0 {
// 		file , err := bytes[ 0 ].Open()
// 		defer file.Close()
// 		if err == nil {
// 			data , err := ioutil.ReadAll( file )
// 			if err == nil {
// 				file_data.FileName = bytes[ 0 ].Filename
// 				file_data.Data = data
// 			}
// 		}
// 	}
// 	if file_data.Data == nil {
// 		return context.JSON(fiber.Map{
// 			"success": 0 ,
// 			"file": fiber.Map{
// 				"url": "" ,
// 			} ,
// 		})
// 	}
// 	fmt.Println( "Uploaded Bytes from :" , file_data.FileName )
// 	fmt.Println( file_data.Data )
// 	return context.JSON(fiber.Map{
// 		"success": 1 ,
// 		"file": fiber.Map{
// 			"url": "" ,
// 		} ,
// 	})
// }