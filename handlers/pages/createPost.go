package pages

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var ApprovalNeeded = false // if new posts need moderator approval

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		forPage := structs.ForPage{}
		forPage.User = handlers.IsLoggedIn(r, db).User
		forPage.Error.Error = false
		forPage.Tags = database.Tags

		// Check if the user is logged in
		if !forPage.User.LoggedIn {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		// Handle GET request to render the create post page
		if r.Method == http.MethodGet {
			handlers.RenderTemplates("createPost", forPage, w, r)
			return
		}

		// Handle POST request
		if r.Method != http.MethodPost {
			return
		}

		title, description, selectedTags, errStr := CheckPostCreation(r)
		if errStr != "" {
			handlers.ErrorHandler(w, http.StatusInternalServerError, errStr)
			return
		}

		// Check if an image file was uploaded
		imageFilename, errStr := ProcessImageUpload(r)
		if errStr != "" {
			forPage.Error.Error = true
			forPage.Error.Message = errStr
			forPage.Error.Field1 = title
			forPage.Error.Field2 = description
			forPage.Error.Field3 = selectedTags
			handlers.RenderTemplates("createPost", forPage, w, r)
			return
		}

		// Voir si La video a ete telechargée
		vidFilename, errStr := ProcessVideoUpload(r)
		if errStr != "" {
			forPage.Error.Error = true
			forPage.Error.Message = errStr
			forPage.Error.Field1 = title
			forPage.Error.Field2 = description
			forPage.Error.Field3 = selectedTags
			handlers.RenderTemplates("createPost", forPage, w, r)
			return
		}

		// Prepare the SQL statement for inserting post data
		stmt := "INSERT INTO posts (username, title, description, imageFilename, vidFilename, approved) VALUES (?, ?, ?, ?, ?, ?)"

		var approval2 bool
		if forPage.User.TypeInt > 0 {
			approval2 = false
		} else {
			approval2 = ApprovalNeeded
		}
		// Execute the SQL statement to insert post data into the database
		result, err := db.Exec(stmt, forPage.User.Username, title, description, imageFilename,vidFilename, !approval2)
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "N’a pas réussi à insérer les données de poste dans la base de données.")
			return
		}

		// Get the ID of the newly created post
		postID, err := result.LastInsertId()
		if err != nil {
			handlers.ErrorHandler(w, http.StatusInternalServerError, "Impossible d’obtenir l’ID de post.")
			return
		}

		// Insert the selected tags into the post_tags table
		for _, tagID := range selectedTags {
			_, err = db.Exec("INSERT INTO post_tags (postID, tagID) VALUES (?, ?)", postID, tagID)
			if err != nil {
				handlers.ErrorHandler(w, http.StatusInternalServerError, "Impossible d’insérer la balise dans la table post_tags")
				return
			}
		}

		// Redirect to the homepage or display a success message
		http.Redirect(w, r, "/viewPost?id="+strconv.FormatInt(postID, 10), http.StatusFound)
	}
}

func CheckPostCreation(r *http.Request) (string, string, []string, string) {
	// Parse the form data from the request
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", "", []string{}, "Échec de l’analyse des données de formulaire"
	}

	// Extract the post data
	title := r.FormValue("title")
	description := r.FormValue("description")
	selectedTags := r.Form["tags"]

	// Check for empty fields
	if title == "" || description == "" || len(selectedTags) == 0 {
		return "", "", []string{}, "Champs vides interdits"
	}

	return title, description, selectedTags, ""
}

func ProcessImageUpload(r *http.Request) (string, string) {
	maxFileSize := 20 * 1024 * 1024 // 20MB
	file, header, err := r.FormFile("image")
	if err != nil {
		// No image file uploaded
		return "", ""
	}
	defer file.Close()

	// Validate the image size
	if header.Size > int64(maxFileSize) {
		return header.Filename, "Champs vides interdits La taille de l’image dépasse la limite maximale de 20 Mo"
	}

	// Validate the image format
	fileExt := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[fileExt] {
		return "", "Format d’image non valide. Seuls les formats JPEG, PNG et GIF sont autorisés"
	}

	// Generate a unique filename for the image
	imageFilename := generateUniqueFileName(header.Filename)

	// Save the image file to the server
	imagePath := filepath.Join("./assets/uploads", imageFilename)
	err = saveImageToFile(file, imagePath)
	if err != nil {
		return header.Filename, "Erreur."
	}

	return imageFilename, ""
}

func saveImageToFile(file multipart.File, imagePath string) error {
	// Create a new file at the specified path
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func ProcessVideoUpload(r *http.Request) (string, string) {
    maxFileSize := 100 * 1024 * 1024 // 100MB pour une vidéo
    video, header, err := r.FormFile("video")
    if err != nil {
        // Aucune vidéo n'a été uploadée
        return "", ""
    }
    defer video.Close()

    // Valider la taille de la vidéo
    if header.Size > int64(maxFileSize) {
        return header.Filename, "La taille de la vidéo doit être inférieure à 100 Mo"
    }

    // Valider le format de la vidéo
    videoExt := filepath.Ext(header.Filename)
    allowedVideoExts := map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".mkv": true, ".webm": true}
    if !allowedVideoExts[videoExt] {
        return "", "Format vidéo non valide. Seuls les formats WEPM, MP4, MOV, AVI et MKV sont autorisés"
    }

    // Générer un nom de fichier unique pour la vidéo
    videoFilename := generateUniqueFileName(header.Filename)

    // Enregistrer le fichier vidéo sur le serveur
    videoPath := filepath.Join("./assets/uploads/videos", videoFilename)
    err = saveVideoToFile(video, videoPath)
    if err != nil {
        return header.Filename, "Impossible d'enregistrer la vidéo"
    }

    return videoFilename, ""
}


func saveVideoToFile(file multipart.File, videoPath string) error {
	// Create a new file at the specified path
	f, err := os.Create(videoPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func generateUniqueFileName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	return filename
}
