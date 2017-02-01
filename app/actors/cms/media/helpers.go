package media

import (
	"github.com/ottemo/foundation/media"
)

// correctMediaType returns only supported media type according to srcMediaType specified
func correctMediaType(srcMediaType string) string {
	var mediaType = srcMediaType

	if len(srcMediaType) == 0 {
		mediaType = media.ConstMediaTypeImage
	} else if mediaType != media.ConstMediaTypeImage {
		mediaType = media.ConstMediaTypeDocument
	}

	return mediaType
}
