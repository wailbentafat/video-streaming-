package Upload

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	db "videoback/model"
)
const fipath ="/backend/cmd/upload/"
func countFilesInDirectory(directory string) (int, error) {
	
	files, err := os.ReadDir(directory)
	if err != nil {
		return 0, err
	}

	loop:=len(files)-2
	return loop, nil
}

func Get_upload( w http.ResponseWriter,r *http.Request ){
	w.Header().Set("Content-Type", "application/json")
   w.Header().Set("Access-Control-Allow-Origin", "*")
   w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

   name :=r.URL.Query().Get("name")
   filename:=name
   fmt.Println(filename)
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
   if _,err:=os.Stat(path);err!=nil{
   fmt.Println("nikmha ")
   }
   length,count:=countFilesInDirectory(path)
   fmt.Println(length)
   fmt.Println(count)

   for i:=1;i<=length;i++{
      http.ServeFile(w,r,fipath+filename+"/"+filename+"_"+strconv.Itoa(i)+".ts")
   }
   
}

  



   

   


   

 



