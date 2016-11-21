// Package post is a default implementation of blog post related interfaces declared in
// "github.com/ottemo/foundation/app/models/blog/post" package
package post

import (
	"time"

	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstBlogPostCollectionName = "blog_post"

	ConstErrorModule = "blog"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultBlogPost is a default implementer of InterfaceBlogPost
type DefaultBlogPost struct {
	id string

	identifier string

	title   string
	excerpt string
	content string

	tags          []interface{}
	featuredImage string

	createdAt time.Time
	updatedAt time.Time
	published bool
}
