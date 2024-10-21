package Upload

import (
	"fmt"
	"net/http"
	db "videoback/model"
)
const fipath ="/backend/cmd/upload/"


func get_upload( w http.ResponseWriter,r *http.Response ){
	w.Header().Set("Content-Type", "application/json")
   queryParams:=r.URL.query()
   fielname:=queryParams.Get("file")
   if filename==""{
	fmt.Println("nikmha got ")
	
   }
   var videos []db.Video
   result:=db.DB.Find("VideoName=?",filename).First(&videos)
   if result.Error!=nil{
   fmt.Println("nikmha ")
   }
   path:=videos[0].HLSPath
   fmt.Println(path)
   

   


   

 



