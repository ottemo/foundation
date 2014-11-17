// Package fsmedia is a default implementation of I_MediaStorage declared in
// "github.com/ottemo/foundation/media" package
package fsmedia

// Package global constants
const (
	MEDIA_DB_COLLECTION  = "media"    // database collection name to store media assignment information into
	MEDIA_DEFAULT_FOLDER = "./media/" // filesystem folder path to store media files in there
)

// I_MediaStorage implementer class
type FilesystemMediaStorage struct {
	storageFolder string
	setupWaitCnt  int
}
