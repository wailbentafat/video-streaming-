package Upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	db "videoback/model"
)

func Get_all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	type video1 struct {
		FrstSegments string `json:"frst_segments"`
		VideoName    string `json:"video_name"`
		M3u8         string `json:"m3u8"`
	}

	var videos []db.Video
	result := db.DB.Find(&videos)
	if result.Error != nil {
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		fmt.Println("Database query error:", result.Error)
		return
	}

	var videoList []video1 

	for _, video := range videos {
		hlsPath := video.HLSPath 
        s:= []string{video.VideoName,"m3u8"}
		m:=[]string{video.VideoName,"0.ts"}
		m3u8:=strings.Join(s,".")
		firstseg:=strings.Join(m,".")
		fmt.Println(firstseg)
		fmt.Println(m3u8)
		    
		firstSegmentPath := filepath.Join(hlsPath, firstseg)
		m3u8Path := filepath.Join(hlsPath, m3u8)
		
		if _, err := os.Stat(m3u8Path); os.IsNotExist(err) {
			fmt.Println("M3U8 file does not exist:", m3u8Path)
			continue
		}

		if _, err := os.Stat(firstSegmentPath); os.IsNotExist(err) {
			fmt.Println("First segment file does not exist:", firstSegmentPath)
			continue 
		}

		
		videoList = append(videoList, video1{
			FrstSegments: firstSegmentPath, 
			VideoName:    video.VideoName,
			M3u8:         m3u8Path, 
		})
	}

	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videoList)
}
