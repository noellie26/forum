package pages

import (
	"database/sql"
	"forum/database"
	"forum/handlers"
	"forum/structs"
	"log"
	"net/http"
	"sort"
	"time"
)

func EditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := handlers.IsLoggedIn(r, db).User
		type forPage struct {
			structs.User
			Post    structs.Post
			Comment structs.Comment
			Tags    []structs.Tag
			Error   structs.ErrorMessage
		}

		if r.Method == http.MethodGet {
			id := r.URL.Query().Get("id")
			editType := r.URL.Query().Get("type")

			var post structs.Post
			var comment structs.Comment
			if editType == "post" {
				post, _ = GetPost(w, db, id)
				if post.Username != user.Username {
					handlers.ErrorHandler(w, http.StatusUnauthorized, "You are not allowed to edit other users posts")
					return
				}
			} else if editType == "comment" {
				comment, _ = GetComment(w, db, id)
				if comment.Username != user.Username {
					handlers.ErrorHandler(w, http.StatusUnauthorized, "You are not allowed to edit other users comments")
					return
				}
			}

			// Create a data struct to pass to the template
			data := forPage{
				User:    user,
				Post:    post,
				Comment: comment,
				Tags:    database.Tags,
			}

			handlers.RenderTemplates("edit", data, w, r)

		} else if r.Method == http.MethodPost {
			log.Println("[edit] post")
			id := r.FormValue("id")
			editType := r.FormValue("type")

			if editType == "post" {
				title, description, selectedTags, errStr := CheckPostCreation(r)
				if errStr != "" {
					handlers.ErrorHandler(w, http.StatusInternalServerError, errStr)
					return
				}

				post, _ := GetPost(w, db, id)
				if post.Username != user.Username {
					handlers.ErrorHandler(w, http.StatusUnauthorized, "Vous n'etes pas autorisée a modifier les posts des autres")
					return
				}

				// Check if an image file was uploaded
				imageFilename, errStr := ProcessImageUpload(r)
				if errStr != "" {
					data := forPage{
						User: user,
						Post: post,
						Tags: database.Tags,
					}
					data.Error = structs.ErrorMessage{
						Error:   true,
						Message: errStr,
						Field1:  title,
						Field2:  description,
						Field3:  selectedTags,
						Image:   imageFilename,
					}
					handlers.RenderTemplates("edit?type=post&id="+id, data, w, r)
					return
				}

				vidFilename, errStr := ProcessVideoUpload(r)
				if errStr != "" {
					data := forPage{
						User: user,
						Post: post,
						Tags: database.Tags,
					}
					data.Error = structs.ErrorMessage{
						Error:   true,
						Message: errStr,
						Field1:  title,
						Field2:  description,
						Field3:  selectedTags,
						Video: vidFilename,
					}
					handlers.RenderTemplates("edit?type=post&id="+id, data, w, r)
					return
				}

				err := editPost(db, id, title, description, imageFilename,vidFilename, selectedTags)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}

				http.Redirect(w, r, "/viewPost?id="+id, http.StatusFound)

			} else if editType == "comment" {
				content := r.FormValue("content")

				err := editComment(w, db, id, content)
				if err != nil {
					log.Println(err)
					http.Error(w, "Impossible de modifier le commentaire", http.StatusInternalServerError)
				}
				http.Redirect(w, r, "/activity", http.StatusFound)
			}
		}
	}
}

// GetPost from database
func GetPost(w http.ResponseWriter, db *sql.DB, postID string) (structs.Post, error) {
	postQuery := `
			SELECT postID, title, description, imageFileName, vidFilename, creationDate, username, likes, dislikes
			FROM posts
			WHERE postID = ?
		`

	postRow := db.QueryRow(postQuery, postID)

	var post structs.Post
	var imageFileName sql.NullString
	var vidFileName sql.NullString
	err := postRow.Scan(&post.ID, &post.Title, &post.Description, &imageFileName, &vidFileName, &post.CreationDate, &post.Username, &post.Likes, &post.Dislikes)
	if err != nil {
		http.Error(w, "Echec", http.StatusInternalServerError)
		return structs.Post{}, err
	}

	if imageFileName.Valid {
		post.ImageFileName = imageFileName.String
	} else {
		post.ImageFileName = "" // Set a default value for imageFileName when it is NULL
	}
    
	if vidFileName.Valid {
		post.VidFileName = vidFileName.String
	} else {
		post.VidFileName = "" // Set a default value for VidFileName when it is NULL
	}
	post.Tags, _ = database.GetPostTags(db, postID)

	return post, nil
}

func editPost(db *sql.DB, postID string, title, description, imageFilename string, vidFileName string, tags []string) error {
	// Retrieve the existing post from the database
	oldTags, err := database.GetPostTags(db, postID)
	if err != nil {
		return err
	}

	// Update the post's title and description
	updatePostQuery := `
		UPDATE posts
		SET title = ?, description = ?, edited = true, timeEdited = CURRENT_TIMESTAMP
		WHERE postID = ?
	`
	_, err = db.Exec(updatePostQuery, title, description, postID)
	if err != nil {
		return err
	}

	// Update the image filename if provided
	if imageFilename != "" {
		updateImageQuery := `
			UPDATE posts
			SET imageFilename = ?
			WHERE postID = ?
		`
		_, err = db.Exec(updateImageQuery, imageFilename, postID)
		if err != nil {
			return err
		}
	}

	// Update the video filename if provided
	if vidFileName != "" {
		updateVideoQuery := `
			UPDATE posts
			SET vidFilename = ?
			WHERE postID = ?
		`
		_, err = db.Exec(updateVideoQuery, vidFileName, postID)
		if err != nil {
			return err
		}
	}

	// Update the post's tags if they have changed
	if !equalStringSlices(oldTags, tags) {
		// Delete existing tags associated with the post
		deleteTagsQuery := `
			DELETE FROM post_tags
			WHERE postID = ?
		`
		_, err = db.Exec(deleteTagsQuery, postID)
		if err != nil {
			return err
		}

		// Insert the selected tags into the post_tags table
		for _, tagID := range tags {
			_, err = db.Exec("INSERT INTO post_tags (postID, tagID) VALUES (?, ?)", postID, tagID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetComment(w http.ResponseWriter, db *sql.DB, commentID string) (structs.Comment, error) {
	commentQuery := `
			SELECT commentID, content, creationDate, username, likes, dislikes, edited, timeEdited
			FROM comments
			WHERE commentID = ?
		`

	commentRow := db.QueryRow(commentQuery, commentID)

	var comment structs.Comment
	var timeEdited sql.NullTime
	err := commentRow.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &comment.Username, &comment.Likes, &comment.Dislikes, &comment.Edited, &timeEdited)
	if err != nil {
		http.Error(w, "Failed to retrieve comment", http.StatusInternalServerError)
		return structs.Comment{}, err
	}

	return comment, nil
}

func editComment(w http.ResponseWriter, db *sql.DB, commentID string, comment string) error {
	editQuery := `
		SELECT commentID, postID, content, username, creationDate, likes, dislikes, edited, timeEdited
		FROM comments
		WHERE commentID = ?
	`

	row := db.QueryRow(editQuery, commentID)

	var existingComment structs.Comment
	var timeEdited sql.NullTime
	err := row.Scan(&existingComment.ID, &existingComment.PostID, &existingComment.Content, &existingComment.Username, &existingComment.CreationDate, &existingComment.Likes, &existingComment.Dislikes, &existingComment.Edited, &timeEdited)
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return err
	}

	// Update the comment details
	existingComment.Content = comment
	existingComment.Edited = true
	existingComment.TimeEdited = time.Now().String()

	updateQuery := `
		UPDATE comments
		SET content = ?, edited = ?, timeEdited = ?
		WHERE commentID = ?
	`

	_, err = db.Exec(updateQuery, existingComment.Content, existingComment.Edited, existingComment.TimeEdited, existingComment.ID)
	if err != nil {
		http.Error(w, "Impossible de modifier le commentaire", http.StatusInternalServerError)
		return err
	}

	return nil
}

func equalStringSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	// Sort the slices to ensure consistent ordering
	sort.Strings(slice1)
	sort.Strings(slice2)

	// Compare each element in the slices
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
