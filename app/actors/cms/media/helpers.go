package media

import "github.com/ottemo/foundation/media/fsmedia"

// correctMediaType returns only supported media type according to srcMediaType specified
func correctMediaType(srcMediaType string) string {
	var mediaType = srcMediaType

	if len(srcMediaType) == 0 {
		mediaType = fsmedia.ConstMediaTypeImage
	} else if mediaType != fsmedia.ConstMediaTypeImage {
		mediaType = fsmedia.ConstMediaTypeDocument
	}

	return mediaType
}
