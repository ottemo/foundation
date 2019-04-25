// Copyright 2019 Ottemo. All rights reserved.

/*
Package fsmedia is a default implementation of InterfaceMediaStorage declared in "github.com/ottemo/commerce/media" package.

It using filesystem ot store media files in, and database to store existing media to object bindings. Following media types
special behaviour supposed:
	ConstMediaTypeImage ("image") - files stored in filesystem as set of specified image sizes
	ConstMediaTypeLink ("link") - external resource - not stored in filesystem
	ConstMediaTypeDocument ("document"), and others - file stored in filesystem but have not and image

Files are stored within filesystem using following pattern: [ConstMediaDefaultFolder]/[mediaType]/[objectModelName]/[obejctID]/[mediaFileName].

"image" type media re-sizes to [ConstConfigPathMediaImageSize] and set of [ConstConfigPathMediaImageSizes] sizes - these
sizes are specified by config values (stored supposedly in DB). On program first startup them are initialized to
[ConstDefaultImageSize] and [ConstDefaultImageSizes] (the same behaviour for config invalid values). Image resizing happens
with usage of white background if [ConstResizeOnBackground] is set to true.

Database record contain only "base image" record other sizes are named wit a following pattern:
	baseImage: [fileName].[fileExtension]
	sizedImage: [fileName]_[sizeName].[fileExtension]

Image sizes is a config value like "small: 75x75, thumb: 260x300, big: 560x650", where key means [sizeName] and value is
size [maxWidth]x[maxHeight] definition. If [sizeName] is not specified iw will equal to value (i.e. "75x75, thumb: 50x50"
equals to "75x75: 75x75, thumb: 50x50"). Image re-sizes to [maxWidth]x[maxHeight] bounding box, "0" dimension have a
special meaning - it means that in this direction image can be any size ("100x0" means that image height is not limited)
Image re-sizing happens with keeping of image aspect ratio.
*/
package fsmedia
