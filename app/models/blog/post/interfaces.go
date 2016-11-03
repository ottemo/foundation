// Package post represents abstraction of business layer blog post object
package post

import (
	"time"

	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
)

// Package global constants
const (
	ConstModelNameBlogPost = "BlogPost"

	ConstErrorModule = "blog"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceBlogPost represents interface to access business layer implementation of blog post object
type InterfaceBlogPost interface {
	GetIdentifier() string
	SetIdentifier(string) error
	IsPublished() bool
	SetPublished(bool) error

	GetTitle() string
	SetTitle(string) error
	GetExcerpt() string
	SetExcerpt(string) error
	GetContent() string
	SetContent(string) error

	GetTags() []interface{}
	SetTags([]interface{}) error
	GetFeaturedImage() string
	SetFeaturedImage(string) error

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time) error
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
}
