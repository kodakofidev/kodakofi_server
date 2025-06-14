package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleFileUpload menghapus file lama dan menyimpan file profil baru untuk user tertentu
func FileNameProfile(ctx *gin.Context, file *multipart.FileHeader, userID string) (filename, filePath string, err error) {
	// Hapus file lama user
	oldFiles, err := filepath.Glob(filepath.Join("public", "profile-images", "*_"+userID+"_profile*"))
	if err != nil {
		return "", "", fmt.Errorf("failed to check for existing files: %w", err)
	}

	for _, oldFile := range oldFiles {
		if err := os.Remove(oldFile); err != nil {
			log.Printf("Warning: failed to delete old file %s: %v", oldFile, err)
		}
	}

	// Generate nama file baru
	ext := filepath.Ext(file.Filename)
	filename = fmt.Sprintf("%d_%s_profile%s", time.Now().UnixNano(), userID, ext)
	filePath = filepath.Join("public", "profile-images", filename)

	// Simpan file baru
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, filePath, nil
}
