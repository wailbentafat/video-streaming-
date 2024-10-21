package Upload

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	db "videoback/model"
)

const uploadPath = "./upload/"

type ErrorResponse struct {
	Message string `json:"message"`
}



func Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		RespondWithError(w, http.StatusMethodNotAllowed, "wrong method")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	customFilename := r.FormValue("filename")
	if customFilename == "" {
		RespondWithError(w, http.StatusBadRequest, "filename is required")
		return
	}

	file, handler, err := r.FormFile("video")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %s, File Size: %d, File Header: %v\n", handler.Filename, handler.Size, handler.Header)

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to create upload directory")
		return
	}

filePath := uploadPath + customFilename + ".mp4"

	dst, err := os.Create(filePath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	fmt.Printf("File saved successfully at: %s\n", filePath)

	var hlspath string
	err, hlspath = Savesegmentations(filePath, customFilename)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("HLS path generated: %s\n", hlspath)

	video := db.Video{VideoName: customFilename, VideoPath: filePath, HLSPath: hlspath, User_id: 1}
	if err := db.DB.Create(&video).Error; err != nil {
		fmt.Printf("Error saving video to database: %v\n", err)
		RespondWithError(w, http.StatusInternalServerError, "failed to save video to database")
		return
	}

	fmt.Println("Video record created in database")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "File uploaded successfully: %s"}`, filePath)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func Savesegmentations(filepath string, csfilenm string) (error, string) {
    hlspath := uploadPath + csfilenm + "/"
    if err := os.MkdirAll(hlspath, os.ModePerm); err != nil {
        return err, ""
    }

    hlsplaylist := hlspath + csfilenm + ".m3u8"
    fmt.Println("HLS Playlist: ", hlsplaylist)
    segmentPattern := hlspath + csfilenm + ".%d.ts"

    cmd := exec.Command(
        "ffmpeg",
        "-i", filepath,
        "-c", "copy",
        "-bsf:v", "h264_mp4toannexb",
        "-f", "segment",
        "-segment_list", hlsplaylist,
        "-segment_time", "10",
        "-segment_format", "ts",
        segmentPattern,
    )

    
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Printf("Error running FFmpeg command: %v\n", err)
        return err, ""
    }

    return nil, hlsplaylist
}
