package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ── Additional Models ──

type FileWatcher struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	NotifyOn  string             `json:"notify_on" bson:"notify_on"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FilePin struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	ChannelID string             `json:"channel_id" bson:"channel_id"`
	PinnedBy  string             `json:"pinned_by" bson:"pinned_by"`
	PinnedAt  time.Time          `json:"pinned_at" bson:"pinned_at"`
}

type FileReaction struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Emoji     string             `json:"emoji" bson:"emoji"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FileDownloadLog struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID     string             `json:"file_id" bson:"file_id"`
	UserID     string             `json:"user_id" bson:"user_id"`
	IP         string             `json:"ip" bson:"ip"`
	UserAgent  string             `json:"user_agent" bson:"user_agent"`
	DownloadAt time.Time          `json:"download_at" bson:"download_at"`
}

type FileAccessRequest struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID     string             `json:"file_id" bson:"file_id"`
	RequesterID string            `json:"requester_id" bson:"requester_id"`
	Reason     string             `json:"reason" bson:"reason"`
	Status     string             `json:"status" bson:"status"` // pending, approved, denied
	ReviewedBy string             `json:"reviewed_by" bson:"reviewed_by"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	ReviewedAt *time.Time         `json:"reviewed_at" bson:"reviewed_at"`
}

type FileTemplate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	MimeType    string             `json:"mime_type" bson:"mime_type"`
	StorageKey  string             `json:"storage_key" bson:"storage_key"`
	Category    string             `json:"category" bson:"category"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	UseCount    int                `json:"use_count" bson:"use_count"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type FileLabel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID    string             `json:"file_id" bson:"file_id"`
	Label     string             `json:"label" bson:"label"`
	Color     string             `json:"color" bson:"color"`
	AddedBy   string             `json:"added_by" bson:"added_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type FileExport struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileIDs     []string           `json:"file_ids" bson:"file_ids"`
	Format      string             `json:"format" bson:"format"` // zip, tar
	Status      string             `json:"status" bson:"status"`
	URL         string             `json:"url" bson:"url"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	CompletedAt *time.Time         `json:"completed_at" bson:"completed_at"`
}

type FileNotificationPref struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileID   string             `json:"file_id" bson:"file_id"`
	UserID   string             `json:"user_id" bson:"user_id"`
	Muted    bool               `json:"muted" bson:"muted"`
	Desktop  bool               `json:"desktop" bson:"desktop"`
	Mobile   bool               `json:"mobile" bson:"mobile"`
}

// ── Collection accessors ──

func fileWatchersCol() *mongo.Collection   { return mongoDB.Collection("file_watchers") }
func filePinsCol() *mongo.Collection       { return mongoDB.Collection("file_pins") }
func fileReactionsCol() *mongo.Collection  { return mongoDB.Collection("file_reactions") }
func fileDownloadsCol() *mongo.Collection  { return mongoDB.Collection("file_download_logs") }
func fileAccessReqsCol() *mongo.Collection { return mongoDB.Collection("file_access_requests") }
func fileTemplatesCol() *mongo.Collection  { return mongoDB.Collection("file_templates") }
func fileLabelsCol() *mongo.Collection     { return mongoDB.Collection("file_labels") }
func fileExportsCol() *mongo.Collection    { return mongoDB.Collection("file_exports") }
func fileNotifPrefsCol() *mongo.Collection { return mongoDB.Collection("file_notif_prefs") }

// ── Route registration ──

func registerExtendedRoutes2(api *gin.RouterGroup) {
	// Watchers
	api.POST("/:id/watchers", addFileWatcher)
	api.DELETE("/:id/watchers", removeFileWatcher)
	api.GET("/:id/watchers", listFileWatchers)
	api.GET("/:id/watching", isWatchingFile)

	// Pins
	api.POST("/:id/pin", pinFile)
	api.DELETE("/:id/pin", unpinFile)
	api.GET("/channel/:channelId/pins", listChannelPins)
	api.GET("/:id/pinned", isFilePinned)

	// Reactions
	api.POST("/:id/reactions", addFileReaction)
	api.DELETE("/:id/reactions", removeFileReaction)
	api.GET("/:id/reactions", listFileReactions)
	api.GET("/:id/reactions/summary", getFileReactionSummary)

	// Download logs
	api.GET("/:id/downloads", listFileDownloads)
	api.GET("/:id/downloads/count", countFileDownloads)
	api.GET("/downloads/recent", listRecentDownloads)

	// Access requests
	api.POST("/:id/access-requests", createAccessRequest)
	api.GET("/:id/access-requests", listAccessRequests)
	api.PUT("/:id/access-requests/:requestId", reviewAccessRequest)
	api.GET("/access-requests/pending", listPendingAccessRequests)

	// Templates
	api.POST("/templates", createFileTemplate)
	api.GET("/templates", listFileTemplates)
	api.GET("/templates/:templateId", getFileTemplate)
	api.PUT("/templates/:templateId", updateFileTemplate)
	api.DELETE("/templates/:templateId", deleteFileTemplate)
	api.POST("/templates/:templateId/use", useFileTemplate)

	// Labels
	api.POST("/:id/labels", addFileLabel)
	api.DELETE("/:id/labels/:label", removeFileLabel)
	api.GET("/:id/labels", listFileLabels)
	api.GET("/labels/search", searchByLabel)

	// Notification prefs
	api.PUT("/:id/notifications", setFileNotifPref)
	api.GET("/:id/notifications", getFileNotifPref)

	// Export
	api.POST("/export", createFileExport)
	api.GET("/export/:exportId", getFileExportStatus)

	// Bulk additional
	api.POST("/bulk/delete", bulkDeleteFiles)
	api.POST("/bulk/favorite", bulkFavoriteFiles)
	api.POST("/bulk/label", bulkLabelFiles)

	// File stats
	api.GET("/:id/stats", getFileStats)
	api.GET("/workspace/:workspaceId/stats", getWorkspaceFileStats)
	api.GET("/channel/:channelId/stats", getChannelFileStats)
	api.GET("/user/:userId/stats", getFileUserStats)

	// Rename & copy
	api.PUT("/:id/rename", renameFile)
	api.POST("/:id/copy", copyFile)
	api.PUT("/:id/move", moveFile)
}

// ── Handlers ──

func addFileWatcher(c *gin.Context) {
	w := FileWatcher{
		FileID: c.Param("id"), UserID: c.GetHeader("X-User-ID"), NotifyOn: "all", CreatedAt: time.Now(),
	}
	res, err := fileWatchersCol().InsertOne(context.TODO(), w)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	w.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": w})
}

func removeFileWatcher(c *gin.Context) {
	_, err := fileWatchersCol().DeleteOne(context.TODO(), bson.M{"file_id": c.Param("id"), "user_id": c.GetHeader("X-User-ID")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func listFileWatchers(c *gin.Context) {
	cur, err := fileWatchersCol().Find(context.TODO(), bson.M{"file_id": c.Param("id")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var watchers []FileWatcher
	_ = cur.All(context.TODO(), &watchers)
	c.JSON(200, gin.H{"success": true, "data": watchers})
}

func isWatchingFile(c *gin.Context) {
	count, _ := fileWatchersCol().CountDocuments(context.TODO(), bson.M{"file_id": c.Param("id"), "user_id": c.GetHeader("X-User-ID")})
	c.JSON(200, gin.H{"success": true, "watching": count > 0})
}

func pinFile(c *gin.Context) {
	var req struct{ ChannelID string `json:"channel_id"` }
	_ = c.ShouldBindJSON(&req)
	p := FilePin{FileID: c.Param("id"), ChannelID: req.ChannelID, PinnedBy: c.GetHeader("X-User-ID"), PinnedAt: time.Now()}
	res, err := filePinsCol().InsertOne(context.TODO(), p)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	p.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": p})
}

func unpinFile(c *gin.Context) {
	_, err := filePinsCol().DeleteOne(context.TODO(), bson.M{"file_id": c.Param("id")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func listChannelPins(c *gin.Context) {
	cur, err := filePinsCol().Find(context.TODO(), bson.M{"channel_id": c.Param("channelId")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var pins []FilePin
	_ = cur.All(context.TODO(), &pins)
	c.JSON(200, gin.H{"success": true, "data": pins})
}

func isFilePinned(c *gin.Context) {
	count, _ := filePinsCol().CountDocuments(context.TODO(), bson.M{"file_id": c.Param("id")})
	c.JSON(200, gin.H{"success": true, "pinned": count > 0})
}

func addFileReaction(c *gin.Context) {
	var req struct{ Emoji string `json:"emoji"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	r := FileReaction{FileID: c.Param("id"), UserID: c.GetHeader("X-User-ID"), Emoji: req.Emoji, CreatedAt: time.Now()}
	res, err := fileReactionsCol().InsertOne(context.TODO(), r)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	r.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": r})
}

func removeFileReaction(c *gin.Context) {
	_, err := fileReactionsCol().DeleteOne(context.TODO(), bson.M{"file_id": c.Param("id"), "user_id": c.GetHeader("X-User-ID"), "emoji": c.Query("emoji")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func listFileReactions(c *gin.Context) {
	cur, err := fileReactionsCol().Find(context.TODO(), bson.M{"file_id": c.Param("id")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var reactions []FileReaction
	_ = cur.All(context.TODO(), &reactions)
	c.JSON(200, gin.H{"success": true, "data": reactions})
}

func getFileReactionSummary(c *gin.Context) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"file_id": c.Param("id")}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$emoji"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
	}
	cur, err := fileReactionsCol().Aggregate(context.TODO(), pipeline)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var results []bson.M
	_ = cur.All(context.TODO(), &results)
	c.JSON(200, gin.H{"success": true, "data": results})
}

func listFileDownloads(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "download_at", Value: -1}}).SetLimit(50)
	cur, err := fileDownloadsCol().Find(context.TODO(), bson.M{"file_id": c.Param("id")}, opts)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var logs []FileDownloadLog
	_ = cur.All(context.TODO(), &logs)
	c.JSON(200, gin.H{"success": true, "data": logs})
}

func countFileDownloads(c *gin.Context) {
	count, _ := fileDownloadsCol().CountDocuments(context.TODO(), bson.M{"file_id": c.Param("id")})
	c.JSON(200, gin.H{"success": true, "count": count})
}

func listRecentDownloads(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "download_at", Value: -1}}).SetLimit(50)
	cur, err := fileDownloadsCol().Find(context.TODO(), bson.M{"user_id": c.GetHeader("X-User-ID")}, opts)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var logs []FileDownloadLog
	_ = cur.All(context.TODO(), &logs)
	c.JSON(200, gin.H{"success": true, "data": logs})
}

func createAccessRequest(c *gin.Context) {
	var req struct{ Reason string `json:"reason"` }
	_ = c.ShouldBindJSON(&req)
	ar := FileAccessRequest{FileID: c.Param("id"), RequesterID: c.GetHeader("X-User-ID"), Reason: req.Reason, Status: "pending", CreatedAt: time.Now()}
	res, err := fileAccessReqsCol().InsertOne(context.TODO(), ar)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	ar.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": ar})
}

func listAccessRequests(c *gin.Context) {
	cur, err := fileAccessReqsCol().Find(context.TODO(), bson.M{"file_id": c.Param("id")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var reqs []FileAccessRequest
	_ = cur.All(context.TODO(), &reqs)
	c.JSON(200, gin.H{"success": true, "data": reqs})
}

func reviewAccessRequest(c *gin.Context) {
	var req struct{ Status string `json:"status"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	objID, _ := primitive.ObjectIDFromHex(c.Param("requestId"))
	now := time.Now()
	_, err := fileAccessReqsCol().UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": req.Status, "reviewed_by": c.GetHeader("X-User-ID"), "reviewed_at": now}})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func listPendingAccessRequests(c *gin.Context) {
	cur, err := fileAccessReqsCol().Find(context.TODO(), bson.M{"status": "pending"})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var reqs []FileAccessRequest
	_ = cur.All(context.TODO(), &reqs)
	c.JSON(200, gin.H{"success": true, "data": reqs})
}

func createFileTemplate(c *gin.Context) {
	var t FileTemplate
	if err := c.ShouldBindJSON(&t); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	t.CreatedBy = c.GetHeader("X-User-ID")
	t.CreatedAt = time.Now()
	res, err := fileTemplatesCol().InsertOne(context.TODO(), t)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	t.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": t})
}

func listFileTemplates(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "use_count", Value: -1}})
	cur, err := fileTemplatesCol().Find(context.TODO(), bson.M{}, opts)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var templates []FileTemplate
	_ = cur.All(context.TODO(), &templates)
	c.JSON(200, gin.H{"success": true, "data": templates})
}

func getFileTemplate(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("templateId"))
	var t FileTemplate
	err := fileTemplatesCol().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&t)
	if err != nil { c.JSON(404, gin.H{"error": "template not found"}); return }
	c.JSON(200, gin.H{"success": true, "data": t})
}

func updateFileTemplate(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("templateId"))
	var req bson.M
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	_, err := fileTemplatesCol().UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": req})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func deleteFileTemplate(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("templateId"))
	_, err := fileTemplatesCol().DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func useFileTemplate(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("templateId"))
	_, _ = fileTemplatesCol().UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$inc": bson.M{"use_count": 1}})
	var t FileTemplate
	_ = fileTemplatesCol().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&t)
	c.JSON(200, gin.H{"success": true, "data": t})
}

func addFileLabel(c *gin.Context) {
	var req struct{ Label string `json:"label"`; Color string `json:"color"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	l := FileLabel{FileID: c.Param("id"), Label: req.Label, Color: req.Color, AddedBy: c.GetHeader("X-User-ID"), CreatedAt: time.Now()}
	res, err := fileLabelsCol().InsertOne(context.TODO(), l)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	l.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(201, gin.H{"success": true, "data": l})
}

func removeFileLabel(c *gin.Context) {
	_, err := fileLabelsCol().DeleteOne(context.TODO(), bson.M{"file_id": c.Param("id"), "label": c.Param("label")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func listFileLabels(c *gin.Context) {
	cur, err := fileLabelsCol().Find(context.TODO(), bson.M{"file_id": c.Param("id")})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var labels []FileLabel
	_ = cur.All(context.TODO(), &labels)
	c.JSON(200, gin.H{"success": true, "data": labels})
}

func searchByLabel(c *gin.Context) {
	label := c.Query("label")
	cur, err := fileLabelsCol().Find(context.TODO(), bson.M{"label": label}, options.Find().SetLimit(50))
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	defer cur.Close(context.TODO())
	var labels []FileLabel
	_ = cur.All(context.TODO(), &labels)
	c.JSON(200, gin.H{"success": true, "data": labels})
}

func setFileNotifPref(c *gin.Context) {
	var req FileNotificationPref
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	req.FileID = c.Param("id")
	req.UserID = c.GetHeader("X-User-ID")
	opts := options.Update().SetUpsert(true)
	_, err := fileNotifPrefsCol().UpdateOne(context.TODO(), bson.M{"file_id": req.FileID, "user_id": req.UserID}, bson.M{"$set": req}, opts)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func getFileNotifPref(c *gin.Context) {
	var pref FileNotificationPref
	err := fileNotifPrefsCol().FindOne(context.TODO(), bson.M{"file_id": c.Param("id"), "user_id": c.GetHeader("X-User-ID")}).Decode(&pref)
	if err != nil { c.JSON(200, gin.H{"success": true, "data": nil}); return }
	c.JSON(200, gin.H{"success": true, "data": pref})
}

func createFileExport(c *gin.Context) {
	var req struct{ FileIDs []string `json:"file_ids"`; Format string `json:"format"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	if req.Format == "" { req.Format = "zip" }
	export := FileExport{FileIDs: req.FileIDs, Format: req.Format, Status: "pending", CreatedBy: c.GetHeader("X-User-ID"), CreatedAt: time.Now()}
	res, err := fileExportsCol().InsertOne(context.TODO(), export)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	export.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(202, gin.H{"success": true, "data": export})
}

func getFileExportStatus(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("exportId"))
	var export FileExport
	err := fileExportsCol().FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&export)
	if err != nil { c.JSON(404, gin.H{"error": "export not found"}); return }
	c.JSON(200, gin.H{"success": true, "data": export})
}

func bulkDeleteFiles(c *gin.Context) {
	var req struct{ IDs []string `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	var objIDs []primitive.ObjectID
	for _, id := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil { objIDs = append(objIDs, objID) }
	}
	_, _ = filesCol.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in": objIDs}}, bson.M{"$set": bson.M{"status": "deleted", "deleted_at": time.Now()}})
	c.JSON(200, gin.H{"success": true, "deleted": len(req.IDs)})
}

func bulkFavoriteFiles(c *gin.Context) {
	var req struct{ IDs []string `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	for _, id := range req.IDs {
		f := FileFavorite{FileID: id, UserID: c.GetHeader("X-User-ID"), CreatedAt: time.Now()}
		_, _ = favoritesCol().InsertOne(context.TODO(), f)
	}
	c.JSON(200, gin.H{"success": true, "favorited": len(req.IDs)})
}

func bulkLabelFiles(c *gin.Context) {
	var req struct{ IDs []string `json:"ids"`; Label string `json:"label"`; Color string `json:"color"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	for _, id := range req.IDs {
		l := FileLabel{FileID: id, Label: req.Label, Color: req.Color, AddedBy: c.GetHeader("X-User-ID"), CreatedAt: time.Now()}
		_, _ = fileLabelsCol().InsertOne(context.TODO(), l)
	}
	c.JSON(200, gin.H{"success": true, "labeled": len(req.IDs)})
}

func getFileStats(c *gin.Context) {
	fileID := c.Param("id")
	downloads, _ := fileDownloadsCol().CountDocuments(context.TODO(), bson.M{"file_id": fileID})
	reactions, _ := fileReactionsCol().CountDocuments(context.TODO(), bson.M{"file_id": fileID})
	comments, _ := commentsCol().CountDocuments(context.TODO(), bson.M{"file_id": fileID})
	versions, _ := versionsCol().CountDocuments(context.TODO(), bson.M{"file_id": fileID})
	c.JSON(200, gin.H{"success": true, "data": gin.H{"downloads": downloads, "reactions": reactions, "comments": comments, "versions": versions}})
}

func getWorkspaceFileStats(c *gin.Context) {
	wsID := c.Param("workspaceId")
	total, _ := filesCol.CountDocuments(context.TODO(), bson.M{"workspace_id": wsID, "status": bson.M{"$ne": "deleted"}})
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"workspace_id": wsID, "status": bson.M{"$ne": "deleted"}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$metadata.size"}}}}}},
	}
	cur, _ := filesCol.Aggregate(context.TODO(), pipeline)
	var result struct{ TotalSize int64 `bson:"total_size"` }
	if cur != nil && cur.Next(context.TODO()) { _ = cur.Decode(&result) }
	c.JSON(200, gin.H{"success": true, "data": gin.H{"total_files": total, "total_size": result.TotalSize}})
}

func getChannelFileStats(c *gin.Context) {
	chID := c.Param("channelId")
	total, _ := filesCol.CountDocuments(context.TODO(), bson.M{"channel_id": chID, "status": bson.M{"$ne": "deleted"}})
	c.JSON(200, gin.H{"success": true, "data": gin.H{"total_files": total}})
}

func getFileUserStats(c *gin.Context) {
	userID := c.Param("userId")
	total, _ := filesCol.CountDocuments(context.TODO(), bson.M{"uploaded_by": userID, "status": bson.M{"$ne": "deleted"}})
	downloads, _ := fileDownloadsCol().CountDocuments(context.TODO(), bson.M{"user_id": userID})
	c.JSON(200, gin.H{"success": true, "data": gin.H{"total_files": total, "total_downloads": downloads}})
}

func renameFile(c *gin.Context) {
	var req struct{ Name string `json:"name"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	_, err := filesCol.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"original_name": req.Name, "updated_at": time.Now()}})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}

func copyFile(c *gin.Context) {
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var file bson.M
	err := filesCol.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&file)
	if err != nil { c.JSON(404, gin.H{"error": "file not found"}); return }
	file["_id"] = primitive.NewObjectID()
	file["created_at"] = time.Now()
	file["updated_at"] = time.Now()
	res, err := filesCol.InsertOne(context.TODO(), file)
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(201, gin.H{"success": true, "new_id": res.InsertedID})
}

func moveFile(c *gin.Context) {
	var req struct{ ChannelID string `json:"channel_id"`; WorkspaceID string `json:"workspace_id"` }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	update := bson.M{"updated_at": time.Now()}
	if req.ChannelID != "" { update["channel_id"] = req.ChannelID }
	if req.WorkspaceID != "" { update["workspace_id"] = req.WorkspaceID }
	_, err := filesCol.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
	c.JSON(200, gin.H{"success": true})
}
