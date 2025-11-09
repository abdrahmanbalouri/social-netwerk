package helper

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Upload(PictureSend string) (string, error) {
	uploadDir := "../frontend/public/uploads"
	var imageFileName string
	// Créer le dossier s’il n’existe pas
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Nom unique pour l'image
	imageFileName = fmt.Sprintf("msg_%d.png", time.Now().UnixNano())
	imagePath := filepath.Join(uploadDir, imageFileName)

	// Supprimer le préfixe base64 si présent
	base64Data := PictureSend
	if strings.HasPrefix(base64Data, "data:image") {
		parts := strings.Split(base64Data, ",")
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// Décoder l'image base64
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Créer le fichier image
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %v", err)
	}
	defer file.Close()

	// Copier les bytes dans le fichier
	if _, err := io.Copy(file, strings.NewReader(string(imageBytes))); err != nil {
		return "", fmt.Errorf("failed to write image to file: %v", err)
	}
	return imageFileName, nil
}
