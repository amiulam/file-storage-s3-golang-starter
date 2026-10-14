package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20

	r.ParseMultipartForm(maxMemory)
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse multipart form", err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find the video", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the video", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't upload a thumbnail for this video", err)
		return
	}

	contentType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid content type", err)
		return
	}

	if contentType != "image/jpeg" && contentType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Only JPEG and PNG thumbnails are supported", nil)
		return
	}

	randomByte := make([]byte, 32)
	_, err = rand.Read(randomByte)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate filename", nil)
		return
	}

	randomFileName := base64.RawURLEncoding.EncodeToString(randomByte)

	assetWithExt := getAssetWithExt(randomFileName, contentType)
	assetDiskPath := cfg.getAssetDiskPath(assetWithExt)

	dst, err := os.Create(assetDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create the file on server", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save the file", err)
		return
	}

	thumbnailUrl := cfg.getAssetURL(assetWithExt)
	video.ThumbnailURL = &thumbnailUrl

	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update the video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
