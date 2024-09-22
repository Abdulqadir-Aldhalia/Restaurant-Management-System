package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"server-side/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

var (
	Domain          = os.Getenv("DOMAIN")
	DOMAIN          = os.Getenv("Domain")
	ImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
)

func SendJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func HandelError(w http.ResponseWriter, status int, message string) {
	SendJsonResponse(w, status, map[string]string{
		"error": message,
	})
}

func GetRootpath(dir string) string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return filepath.Join(cwd, dir)
}

func HandleError(w http.ResponseWriter, status int, message string) {
	SendJsonResponse(w, status, map[string]string{
		"error": message,
	})
}

// SaveImageFile saves the uploaded image file to a specified directory with a new name
func SaveImageFile(file io.Reader, table string, filename string) (string, error) {
	// Create directory structure if it doesn't exist
	fullPath := filepath.Join("uploads", table)
	fmt.Printf(" filePath of the saving image  => %s \n", fullPath)
	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return "", err
	}

	// Generate new filename
	randomNumber := rand.Intn(1000)
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	if !validImageFile(ext) {
		return "", fmt.Errorf("Not a valid file type")
	}

	newFileName := fmt.Sprintf("%s_%d_%d%s", filepath.Base(table), timestamp, randomNumber, ext)
	newFilePath := filepath.Join(fullPath, newFileName)

	// Save the file
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", err
	}

	// Return the full path including directory
	return newFilePath, nil
}

func DeleteImageFile(filePath string) error {
	if err := os.Remove(filePath[1:]); err != nil {
		fmt.Printf("Failed to delete file: %s, error: %v\n", filePath, err)
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File successfully deleted")
	} else {
		fmt.Println("File still exists after deletion attempt")
	}
	return nil
}

func validImageFile(ext string) bool {
	for _, validExt := range ImageExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidUUID(userId string) bool {
	if userId == "" {
		return false
	}
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(userId)
}

func ValidRole(roleId string) bool {
	roles := map[string]bool{
		"1": true,
		"2": true,
		"3": true,
	}

	return roles[roleId]
}

func ValidatePhoneNumber(phone string) bool {
	regex := `^\+?[1-9]\d{1,14}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(phone)
}

func ValidateUser(user model.User) error {
	if len(user.Email) > 30 || len(user.Email) < 10 {
		return errors.New("username must be between 10 and 30 characters")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(user.Name) < 3 {
		return errors.New("name must be greater than 2characters")
	}

	if !ValidatePhoneNumber(user.Phone) {
		return errors.New("invalid phone number")
	}

	return nil
}

func isEmptyOrNil(values ...interface{}) bool {
	for _, value := range values {
		switch v := value.(type) {
		case string:
			if v == "" {
				return true
			}
		default:
			if v == nil {
				return true
			}
		}
	}
	return false
}
