package main

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ── New models ──

type FileVersion struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID      string             `json:"file_id" bson:"file_id"`
	VersionNum  int                `json:"version_num" bson:"version_num"`
	StorageKey  string             `json:"storage_key" bson:"storage_key"`
	Size        int64              `json:"size" bson:"size"`
	Checksum    string             `json:"checksum" bson:"checksum"`
	UploadedBy  string             `json:"uploaded_by" bson:"uploaded_by"`
	Comment     string             `json:"comment" bson:"comment"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type FileComment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Content   string             `json:"content" bson:"content"`
	ParentID  string             `json:"parent_id" bson:"parent_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type FileTag struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	Tag       string             `json:"tag" bson:"tag"`
	AddedBy   string             `json:"added_by" bson:"added_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FileFavorite struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FileCollection struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	OwnerID     string             `json:"owner_id" bson:"owner_id"`
	WorkspaceID string             `json:"workspace_id" bson:"workspace_id"`
	FileIDs     []string           `json:"file_ids" bson:"file_ids"`
	IsPublic    bool               `json:"is_public" bson:"is_public"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type FilePreview struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID       string             `json:"file_id" bson:"file_id"`
	PreviewURL   string             `json:"preview_url" bson:"preview_url"`
	ThumbnailURL string             `json:"thumbnail_url" bson:"thumbnail_url"`
	PreviewType  string             `json:"preview_type" bson:"preview_type"`
	Width        int                `json:"width" bson:"width"`
	Height       int                `json:"height" bson:"height"`
	GeneratedAt  time.Time          `json:"generated_at" bson:"generated_at"`
}

type FileActivity struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Action    string             `json:"action" bson:"action"`
	Details   string             `json:"details" bson:"details"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FilePermission struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID     string             `json:"file_id" bson:"file_id"`
	UserID     string             `json:"user_id" bson:"user_id"`
	Permission string             `json:"permission" bson:"permission"`
	GrantedBy  string             `json:"granted_by" bson:"granted_by"`
	ExpiresAt  *time.Time         `json:"expires_at" bson:"expires_at"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type FileLink struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	Token     string             `json:"token" bson:"token"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	ExpiresAt *time.Time         `json:"expires_at" bson:"expires_at"`
	MaxViews  int                `json:"max_views" bson:"max_views"`
	Views     int                `json:"views" bson:"views"`
	Password  string             `json:"password" bson:"password"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FileScanResult struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID     string             `json:"file_id" bson:"file_id"`
	ScanType   string             `json:"scan_type" bson:"scan_type"`
	Status     string             `json:"status" bson:"status"`
	Findings   []string           `json:"findings" bson:"findings"`
	ScannedAt  time.Time          `json:"scanned_at" bson:"scanned_at"`
}

// ── Collections ──

func versionsCol() *mongo.Collection   { return mongoDB.Collection("file_versions") }
func commentsCol() *mongo.Collection   { return mongoDB.Collection("file_comments") }
func tagsCol() *mongo.Collection       { return mongoDB.Collection("file_tags") }
func favoritesCol() *mongo.Collection  { return mongoDB.Collection("file_favorites") }
func collectionsCol() *mongo.Collection { return mongoDB.Collection("file_collections") }
func previewsCol() *mongo.Collection   { return mongoDB.Collection("file_previews") }
func activityCol() *mongo.Collection   { return mongoDB.Collection("file_activity") }
func permissionsCol() *mongo.Collection { return mongoDB.Collection("file_permissions") }
func linksCol() *mongo.Collection      { return mongoDB.Collection("file_links") }
func scansCol() *mongo.Collection      { return mongoDB.Collection("file_scans") }

// ── Register extended routes ──

func registerExtendedRoutes(api *gin.RouterGroup) {
	// Versions
	api.GET("/:id/versions", listVersions)
	api.POST("/:id/versions", createVersion)
	api.GET("/:id/versions/:versionId", getVersion)
	api.DELETE("/:id/versions/:versionId", deleteVersion)
	api.POST("/:id/versions/:versionId/restore", restoreVersion)

	// Comments
	api.GET("/:id/comments", listComments)
	api.POST("/:id/comments", createComment)
	api.PUT("/:id/comments/:commentId", updateComment)
	api.DELETE("/:id/comments/:commentId", deleteComment)

	// Tags
	api.GET("/:id/tags", listTags)
	api.POST("/:id/tags", addTag)
	api.DELETE("/:id/tags/:tag", removeTag)
	api.GET("/tags/search", searchByTag)

	// Favorites
	api.POST("/:id/favorite", addFavorite)
	api.DELETE("/:id/favorite", removeFavorite)
	api.GET("/favorites", listFavorites)

	// Collections
	api.POST("/collections", createCollection)
	api.GET("/collections", listCollections)
	api.GET("/collections/:collectionId", getCollection)
	api.PUT("/collections/:collectionId", updateCollection)
	api.DELETE("/collections/:collectionId", deleteCollection)
	api.POST("/collections/:collectionId/files", addToCollection)
	api.DELETE("/collections/:collectionId/files/:fileId", removeFromCollection)

	// Previews
	api.GET("/:id/preview", getPreview)

	// Activity
	api.GET("/:id/activity", listActivity)
	api.GET("/activity/user/:userId", listUserActivity)

	// Permissions
	api.GET("/:id/permissions", listPermissions)
	api.POST("/:id/permissions", grantPermission)
	api.DELETE("/:id/permissions/:permissionId", revokePermission)

	// Shared links
	api.POST("/:id/links", createShareLink)
	api.GET("/:id/links", listShareLinks)
	api.DELETE("/:id/links/:linkId", deleteShareLink)
	api.GET("/shared/:token", accessSharedFile)

	// Scans
	api.GET("/:id/scan", getScanResult)
	api.POST("/:id/scan", triggerScan)

	// Bulk operations
	api.POST("/bulk/move", bulkMoveFiles)
	api.POST("/bulk/copy", bulkCopyFiles)
	api.POST("/bulk/tags", bulkAddTags)

	// Search
	api.GET("/search", searchFiles)

	// Duplicate detection
	api.GET("/:id/duplicates", findDuplicates)

	// Recent
	api.GET("/recent", listRecentFiles)

	// Trash
	api.GET("/trash", listTrash)
	api.POST("/trash/:id/restore", restoreFromTrash)

	// Storage quota
	api.GET("/quota/:workspaceId", getStorageQuota)
}

// ── Version handlers ──

func listVersions(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	cursor, err := versionsCol().Find(ctx, bson.M{"file_id": fileID}, options.Find().SetSort(bson.D{{Key: "version_num", Value: -1}}))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	var versions []FileVersion
	cursor.All(ctx, &versions)
	c.JSON(200, gin.H{"success": true, "data": versions})
}

func createVersion(c *gin.Context) {
	fileID := c.Param("id")
	var req struct {
		Comment string `json:"comment"`
	}
	c.ShouldBindJSON(&req)
	ctx := c.Request.Context()

	// Get current max version
	var maxVer FileVersion
	err := versionsCol().FindOne(ctx, bson.M{"file_id": fileID}, options.FindOne().SetSort(bson.D{{Key: "version_num", Value: -1}})).Decode(&maxVer)
	nextVer := 1
	if err == nil {
		nextVer = maxVer.VersionNum + 1
	}

	version := FileVersion{
		FileID:     fileID,
		VersionNum: nextVer,
		UploadedBy: c.Query("user_id"),
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}
	result, err := versionsCol().InsertOne(ctx, version)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	version.ID = result.InsertedID.(primitive.ObjectID)
	logFileActivity(ctx, fileID, c.Query("user_id"), "version_created", req.Comment)
	c.JSON(201, gin.H{"success": true, "data": version})
}

func getVersion(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("versionId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid version ID"})
		return
	}
	var version FileVersion
	if err := versionsCol().FindOne(c.Request.Context(), bson.M{"_id": versionID}).Decode(&version); err != nil {
		c.JSON(404, gin.H{"error": "version not found"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": version})
}

func deleteVersion(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("versionId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid version ID"})
		return
	}
	versionsCol().DeleteOne(c.Request.Context(), bson.M{"_id": versionID})
	c.JSON(200, gin.H{"success": true})
}

func restoreVersion(c *gin.Context) {
	versionID, err := primitive.ObjectIDFromHex(c.Param("versionId"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid version ID"})
		return
	}
	var version FileVersion
	if err := versionsCol().FindOne(c.Request.Context(), bson.M{"_id": versionID}).Decode(&version); err != nil {
		c.JSON(404, gin.H{"error": "version not found"})
		return
	}
	filesCol.UpdateOne(c.Request.Context(), bson.M{"file_id": version.FileID}, bson.M{
		"$set": bson.M{"storage_key": version.StorageKey, "size": version.Size, "checksum": version.Checksum, "updated_at": time.Now()},
	})
	logFileActivity(c.Request.Context(), version.FileID, c.Query("user_id"), "version_restored", strconv.Itoa(version.VersionNum))
	c.JSON(200, gin.H{"success": true, "message": "version restored"})
}

// ── Comment handlers ──

func listComments(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	cursor, err := commentsCol().Find(ctx, bson.M{"file_id": fileID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	var comments []FileComment
	cursor.All(ctx, &comments)
	c.JSON(200, gin.H{"success": true, "data": comments})
}

func createComment(c *gin.Context) {
	fileID := c.Param("id")
	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID string `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	comment := FileComment{
		FileID: fileID, UserID: c.Query("user_id"), Content: req.Content,
		ParentID: req.ParentID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	result, _ := commentsCol().InsertOne(c.Request.Context(), comment)
	comment.ID = result.InsertedID.(primitive.ObjectID)
	logFileActivity(c.Request.Context(), fileID, c.Query("user_id"), "comment_added", "")
	c.JSON(201, gin.H{"success": true, "data": comment})
}

func updateComment(c *gin.Context) {
	commentID, _ := primitive.ObjectIDFromHex(c.Param("commentId"))
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	commentsCol().UpdateOne(c.Request.Context(), bson.M{"_id": commentID}, bson.M{"$set": bson.M{"content": req.Content, "updated_at": time.Now()}})
	c.JSON(200, gin.H{"success": true})
}

func deleteComment(c *gin.Context) {
	commentID, _ := primitive.ObjectIDFromHex(c.Param("commentId"))
	commentsCol().DeleteOne(c.Request.Context(), bson.M{"_id": commentID})
	c.JSON(200, gin.H{"success": true})
}

// ── Tag handlers ──

func listTags(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	cursor, _ := tagsCol().Find(ctx, bson.M{"file_id": fileID})
	var tags []FileTag
	cursor.All(ctx, &tags)
	c.JSON(200, gin.H{"success": true, "data": tags})
}

func addTag(c *gin.Context) {
	fileID := c.Param("id")
	var req struct {
		Tag string `json:"tag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tag := FileTag{FileID: fileID, Tag: req.Tag, AddedBy: c.Query("user_id"), CreatedAt: time.Now()}
	tagsCol().InsertOne(c.Request.Context(), tag)
	c.JSON(201, gin.H{"success": true})
}

func removeTag(c *gin.Context) {
	fileID := c.Param("id")
	tag := c.Param("tag")
	tagsCol().DeleteOne(c.Request.Context(), bson.M{"file_id": fileID, "tag": tag})
	c.JSON(200, gin.H{"success": true})
}

func searchByTag(c *gin.Context) {
	tag := c.Query("tag")
	ctx := c.Request.Context()
	cursor, _ := tagsCol().Find(ctx, bson.M{"tag": tag})
	var tags []FileTag
	cursor.All(ctx, &tags)
	// Get file IDs
	fileIDs := []string{}
	for _, t := range tags {
		fileIDs = append(fileIDs, t.FileID)
	}
	if len(fileIDs) == 0 {
		c.JSON(200, gin.H{"success": true, "data": []File{}})
		return
	}
	cursor2, _ := filesCol.Find(ctx, bson.M{"file_id": bson.M{"$in": fileIDs}})
	var files []File
	cursor2.All(ctx, &files)
	c.JSON(200, gin.H{"success": true, "data": files})
}

// ── Favorite handlers ──

func addFavorite(c *gin.Context) {
	fileID := c.Param("id")
	userID := c.Query("user_id")
	fav := FileFavorite{FileID: fileID, UserID: userID, CreatedAt: time.Now()}
	favoritesCol().InsertOne(c.Request.Context(), fav)
	c.JSON(201, gin.H{"success": true})
}

func removeFavorite(c *gin.Context) {
	fileID := c.Param("id")
	userID := c.Query("user_id")
	favoritesCol().DeleteOne(c.Request.Context(), bson.M{"file_id": fileID, "user_id": userID})
	c.JSON(200, gin.H{"success": true})
}

func listFavorites(c *gin.Context) {
	userID := c.Query("user_id")
	ctx := c.Request.Context()
	cursor, _ := favoritesCol().Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	var favs []FileFavorite
	cursor.All(ctx, &favs)
	c.JSON(200, gin.H{"success": true, "data": favs})
}

// ── Collection handlers ──

func createCollection(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		WorkspaceID string `json:"workspace_id"`
		IsPublic    bool   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	col := FileCollection{
		Name: req.Name, Description: req.Description, OwnerID: c.Query("user_id"),
		WorkspaceID: req.WorkspaceID, IsPublic: req.IsPublic, FileIDs: []string{},
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	result, _ := collectionsCol().InsertOne(c.Request.Context(), col)
	col.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": col})
}

func listCollections(c *gin.Context) {
	userID := c.Query("user_id")
	workspaceID := c.Query("workspace_id")
	filter := bson.M{}
	if userID != "" {
		filter["owner_id"] = userID
	}
	if workspaceID != "" {
		filter["workspace_id"] = workspaceID
	}
	ctx := c.Request.Context()
	cursor, _ := collectionsCol().Find(ctx, filter)
	var cols []FileCollection
	cursor.All(ctx, &cols)
	c.JSON(200, gin.H{"success": true, "data": cols})
}

func getCollection(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("collectionId"))
	var col FileCollection
	if err := collectionsCol().FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&col); err != nil {
		c.JSON(404, gin.H{"error": "collection not found"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": col})
}

func updateCollection(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("collectionId"))
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPublic    *bool  `json:"is_public"`
	}
	c.ShouldBindJSON(&req)
	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.IsPublic != nil {
		update["is_public"] = *req.IsPublic
	}
	collectionsCol().UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$set": update})
	c.JSON(200, gin.H{"success": true})
}

func deleteCollection(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("collectionId"))
	collectionsCol().DeleteOne(c.Request.Context(), bson.M{"_id": id})
	c.JSON(200, gin.H{"success": true})
}

func addToCollection(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("collectionId"))
	var req struct {
		FileID string `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	collectionsCol().UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"file_ids": req.FileID}})
	c.JSON(200, gin.H{"success": true})
}

func removeFromCollection(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("collectionId"))
	fileID := c.Param("fileId")
	collectionsCol().UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$pull": bson.M{"file_ids": fileID}})
	c.JSON(200, gin.H{"success": true})
}

// ── Preview handlers ──

func getPreview(c *gin.Context) {
	fileID := c.Param("id")
	var preview FilePreview
	if err := previewsCol().FindOne(c.Request.Context(), bson.M{"file_id": fileID}).Decode(&preview); err != nil {
		c.JSON(404, gin.H{"error": "preview not found"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": preview})
}

// ── Activity handlers ──

func listActivity(c *gin.Context) {
	fileID := c.Param("id")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	ctx := c.Request.Context()
	cursor, _ := activityCol().Find(ctx, bson.M{"file_id": fileID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	var activities []FileActivity
	cursor.All(ctx, &activities)
	c.JSON(200, gin.H{"success": true, "data": activities})
}

func listUserActivity(c *gin.Context) {
	userID := c.Param("userId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	ctx := c.Request.Context()
	cursor, _ := activityCol().Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit))
	var activities []FileActivity
	cursor.All(ctx, &activities)
	c.JSON(200, gin.H{"success": true, "data": activities})
}

func logFileActivity(ctx context.Context, fileID, userID, action, details string) {
	activityCol().InsertOne(ctx, FileActivity{
		FileID: fileID, UserID: userID, Action: action, Details: details, CreatedAt: time.Now(),
	})
}

// ── Permission handlers ──

func listPermissions(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	cursor, _ := permissionsCol().Find(ctx, bson.M{"file_id": fileID})
	var perms []FilePermission
	cursor.All(ctx, &perms)
	c.JSON(200, gin.H{"success": true, "data": perms})
}

func grantPermission(c *gin.Context) {
	fileID := c.Param("id")
	var req struct {
		UserID     string `json:"user_id" binding:"required"`
		Permission string `json:"permission" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	perm := FilePermission{
		FileID: fileID, UserID: req.UserID, Permission: req.Permission,
		GrantedBy: c.Query("user_id"), CreatedAt: time.Now(),
	}
	result, _ := permissionsCol().InsertOne(c.Request.Context(), perm)
	perm.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": perm})
}

func revokePermission(c *gin.Context) {
	permID, _ := primitive.ObjectIDFromHex(c.Param("permissionId"))
	permissionsCol().DeleteOne(c.Request.Context(), bson.M{"_id": permID})
	c.JSON(200, gin.H{"success": true})
}

// ── Share link handlers ──

func createShareLink(c *gin.Context) {
	fileID := c.Param("id")
	var req struct {
		MaxViews int    `json:"max_views"`
		Password string `json:"password"`
	}
	c.ShouldBindJSON(&req)
	link := FileLink{
		FileID: fileID, Token: primitive.NewObjectID().Hex(),
		CreatedBy: c.Query("user_id"), MaxViews: req.MaxViews, Password: req.Password,
		IsActive: true, CreatedAt: time.Now(),
	}
	result, _ := linksCol().InsertOne(c.Request.Context(), link)
	link.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": link})
}

func listShareLinks(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	cursor, _ := linksCol().Find(ctx, bson.M{"file_id": fileID})
	var links []FileLink
	cursor.All(ctx, &links)
	c.JSON(200, gin.H{"success": true, "data": links})
}

func deleteShareLink(c *gin.Context) {
	linkID, _ := primitive.ObjectIDFromHex(c.Param("linkId"))
	linksCol().DeleteOne(c.Request.Context(), bson.M{"_id": linkID})
	c.JSON(200, gin.H{"success": true})
}

func accessSharedFile(c *gin.Context) {
	token := c.Param("token")
	var link FileLink
	if err := linksCol().FindOne(c.Request.Context(), bson.M{"token": token, "is_active": true}).Decode(&link); err != nil {
		c.JSON(404, gin.H{"error": "link not found or expired"})
		return
	}
	// Increment views
	linksCol().UpdateOne(c.Request.Context(), bson.M{"_id": link.ID}, bson.M{"$inc": bson.M{"views": 1}})
	// Get file
	var file File
	if err := filesCol.FindOne(c.Request.Context(), bson.M{"file_id": link.FileID}).Decode(&file); err != nil {
		c.JSON(404, gin.H{"error": "file not found"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": file})
}

// ── Scan handlers ──

func getScanResult(c *gin.Context) {
	fileID := c.Param("id")
	var scan FileScanResult
	if err := scansCol().FindOne(c.Request.Context(), bson.M{"file_id": fileID}, options.FindOne().SetSort(bson.D{{Key: "scanned_at", Value: -1}})).Decode(&scan); err != nil {
		c.JSON(404, gin.H{"error": "no scan results"})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": scan})
}

func triggerScan(c *gin.Context) {
	fileID := c.Param("id")
	scan := FileScanResult{
		FileID: fileID, ScanType: "antivirus", Status: "pending", Findings: []string{}, ScannedAt: time.Now(),
	}
	result, _ := scansCol().InsertOne(c.Request.Context(), scan)
	scan.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(202, gin.H{"success": true, "data": scan})
}

// ── Bulk handlers ──

func bulkMoveFiles(c *gin.Context) {
	var req struct {
		FileIDs       []string `json:"file_ids" binding:"required"`
		TargetChannel string   `json:"target_channel" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	for _, fid := range req.FileIDs {
		filesCol.UpdateOne(ctx, bson.M{"file_id": fid}, bson.M{"$set": bson.M{"channel_id": req.TargetChannel, "updated_at": time.Now()}})
	}
	c.JSON(200, gin.H{"success": true, "moved": len(req.FileIDs)})
}

func bulkCopyFiles(c *gin.Context) {
	var req struct {
		FileIDs       []string `json:"file_ids" binding:"required"`
		TargetChannel string   `json:"target_channel" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "copied": len(req.FileIDs)})
}

func bulkAddTags(c *gin.Context) {
	var req struct {
		FileIDs []string `json:"file_ids" binding:"required"`
		Tags    []string `json:"tags" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	for _, fid := range req.FileIDs {
		for _, tag := range req.Tags {
			tagsCol().InsertOne(ctx, FileTag{FileID: fid, Tag: tag, AddedBy: c.Query("user_id"), CreatedAt: time.Now()})
		}
	}
	c.JSON(200, gin.H{"success": true})
}

// ── Search ──

func searchFiles(c *gin.Context) {
	query := c.Query("q")
	workspaceID := c.Query("workspace_id")
	fileType := c.Query("file_type")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	filter := bson.M{"deleted_at": nil}
	if query != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"original_name": bson.M{"$regex": query, "$options": "i"}},
		}
	}
	if workspaceID != "" {
		filter["workspace_id"] = workspaceID
	}
	if fileType != "" {
		filter["file_type"] = fileType
	}

	ctx := c.Request.Context()
	cursor, _ := filesCol.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}).SetLimit(limit))
	var files []File
	cursor.All(ctx, &files)
	c.JSON(200, gin.H{"success": true, "data": files})
}

// ── Duplicate detection ──

func findDuplicates(c *gin.Context) {
	fileID := c.Param("id")
	ctx := c.Request.Context()
	var file File
	if err := filesCol.FindOne(ctx, bson.M{"file_id": fileID}).Decode(&file); err != nil {
		c.JSON(404, gin.H{"error": "file not found"})
		return
	}
	// Find by checksum
	cursor, _ := filesCol.Find(ctx, bson.M{"checksum": file.Checksum, "file_id": bson.M{"$ne": fileID}})
	var dupes []File
	cursor.All(ctx, &dupes)
	c.JSON(200, gin.H{"success": true, "data": dupes})
}

// ── Recent files ──

func listRecentFiles(c *gin.Context) {
	userID := c.Query("user_id")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	ctx := c.Request.Context()
	filter := bson.M{"deleted_at": nil}
	if userID != "" {
		filter["uploaded_by"] = userID
	}
	cursor, _ := filesCol.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}).SetLimit(limit))
	var files []File
	cursor.All(ctx, &files)
	c.JSON(200, gin.H{"success": true, "data": files})
}

// ── Trash ──

func listTrash(c *gin.Context) {
	userID := c.Query("user_id")
	ctx := c.Request.Context()
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if userID != "" {
		filter["uploaded_by"] = userID
	}
	cursor, _ := filesCol.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "deleted_at", Value: -1}}))
	var files []File
	cursor.All(ctx, &files)
	c.JSON(200, gin.H{"success": true, "data": files})
}

func restoreFromTrash(c *gin.Context) {
	fileID := c.Param("id")
	filesCol.UpdateOne(c.Request.Context(), bson.M{"file_id": fileID}, bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}})
	c.JSON(200, gin.H{"success": true, "message": "file restored"})
}

// ── Storage quota ──

func getStorageQuota(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	ctx := c.Request.Context()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"workspace_id": workspaceID, "deleted_at": nil}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total_size": bson.M{"$sum": "$size"}, "file_count": bson.M{"$sum": 1}}}},
	}
	cursor, _ := filesCol.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	var result []bson.M
	cursor.All(ctx, &result)

	totalSize := int64(0)
	fileCount := 0
	if len(result) > 0 {
		if v, ok := result[0]["total_size"].(int64); ok {
			totalSize = v
		}
		if v, ok := result[0]["file_count"].(int32); ok {
			fileCount = int(v)
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"workspace_id": workspaceID,
			"total_size":   totalSize,
			"file_count":   fileCount,
			"quota_limit":  10 * 1024 * 1024 * 1024, // 10GB default
		},
	})
}
