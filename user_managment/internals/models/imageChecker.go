package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (m *Models) GenerateFileName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes) + ext, nil
}

func (m *Models) SaveFile(file multipart.File, fileName string) (string, error) {
	path := filepath.Join("ui/static/images", fileName)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(out, file)
	return path, err
}

func (m *Models) ValidateImage(file multipart.File) (bool, error) {
	_, _, err := image.Decode(file)
	if err != nil {
		return false, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Models) InsertImage(ctx context.Context, userID int, picNumber int, fileURI string) error {
	stmt := `INSERT INTO user_images (user_id, image_number, image_url)
             VALUES ($1, $2, $3)
             ON CONFLICT (user_id, image_number)
             DO UPDATE SET image_url = EXCLUDED.image_url`

	_, err := m.DB.Exec(ctx, stmt, userID, picNumber, fileURI)
	if err != nil {
		return err
	}
	return nil
}
