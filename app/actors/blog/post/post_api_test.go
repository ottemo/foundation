package post_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/actors/blog/post"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"
)

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceSession test implementation
//--------------------------------------------------------------------------------------------------------------

type testSession struct {
	_test_data_ map[string]interface{}
}

func (it *testSession) Close() error {
	return nil
}
func (it *testSession) Get(key string) interface{} {
	return it._test_data_[key]
}
func (it *testSession) GetID() string {
	return "ApplicationSession GetID"
}
func (it *testSession) IsEmpty() bool {
	return true
}
func (it *testSession) Set(key string, value interface{}) {
	it._test_data_[key] = value
}
func (it *testSession) Touch() error {
	return nil
}

//--------------------------------------------------------------------------------------------------------------
// api.InterfaceApplicationContext test implementation
//--------------------------------------------------------------------------------------------------------------

type testContext struct {
	//ResponseWriter    http.ResponseWriter
	//Request           *http.Request
	Request string
	//RequestParameters map[string]string
	RequestArguments map[string]string
	RequestContent   interface{}
	//RequestFiles      map[string]io.Reader

	Session       api.InterfaceSession
	ContextValues map[string]interface{}
	//Result        interface{}
}

func (it *testContext) GetRequestArguments() map[string]string {
	return it.RequestArguments
}
func (it *testContext) GetContextValues() map[string]interface{} {
	return it.ContextValues
}
func (it *testContext) GetContextValue(key string) interface{} {
	return it.ContextValues[key]
}
func (it *testContext) GetRequest() interface{} {
	return it.Request
}
func (it *testContext) GetRequestArgument(name string) string {
	return it.RequestArguments[name]
}
func (it *testContext) GetRequestContent() interface{} {
	return it.RequestContent
}
func (it *testContext) GetRequestContentType() string {
	return "request content type"
}
func (it *testContext) GetRequestFile(name string) io.Reader {
	return nil
}
func (it *testContext) GetRequestFiles() map[string]io.Reader {
	return nil
}
func (it *testContext) GetRequestSettings() map[string]interface{} {
	return map[string]interface{}{}
}
func (it *testContext) GetRequestSetting(name string) interface{} {
	return "request setting"
}
func (it *testContext) GetResponse() interface{} {
	return "response"
}
func (it *testContext) GetResponseContentType() string {
	return "response content type"
}
func (it *testContext) GetResponseResult() interface{} {
	return "response result"
}
func (it *testContext) GetResponseSetting(name string) interface{} {
	return "response setting"
}
func (it *testContext) GetResponseWriter() io.Writer {
	return nil
}
func (it *testContext) GetSession() api.InterfaceSession {
	return it.Session
}
func (it *testContext) SetContextValue(key string, value interface{}) {
	//return it.Session
}
func (it *testContext) SetResponseContentType(mimeType string) error {
	return nil
}
func (it *testContext) SetResponseResult(value interface{}) error {
	return nil
}
func (it *testContext) SetResponseSetting(name string, value interface{}) error {
	return nil
}
func (it *testContext) SetResponseStatus(code int) {
	//return nil
}
func (it *testContext) SetResponseStatusBadRequest()          {}
func (it *testContext) SetResponseStatusForbidden()           {}
func (it *testContext) SetResponseStatusNotFound()            {}
func (it *testContext) SetResponseStatusInternalServerError() {}
func (it *testContext) SetSession(session api.InterfaceSession) error {
	it.Session = session
	return nil
}

//--------------------------------------------------------------------------------------------------------------
// test functions
//--------------------------------------------------------------------------------------------------------------

func TestBlogPostAPI(t *testing.T) {

	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	// api redeclaration to mimic priveleged calls
	apiListPosts := post.APIListPosts
	apiPostByID := post.APIPostByID
	apiCreateBlogPost := api.IsAdmin(post.APICreateBlogPost)
	apiUpdateByID := api.IsAdmin(post.APIUpdateByID)
	apiDeleteByID := api.IsAdmin(post.APIDeleteByID)

	// init session
	session := new(testSession)
	session._test_data_ = map[string]interface{}{}

	// init context
	context := new(testContext)
	context.SetSession(session)

	numberOfPosts := 15
	numberOfPublished := 10
	numberOfPostsWithTag := 5
	numberOfPostsWithSpecialTags := 0

	//--------------------------------------------------------------------------------------------------------------
	// cleanup: delete all existing posts
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"limit": "0,1000",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, true)

	result, err := apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	for _, item := range utils.InterfaceToArray(result) {
		itemHashMap := utils.InterfaceToMap(item)
		context.RequestArguments = map[string]string{
			"id": utils.InterfaceToString(itemHashMap["_id"]),
		}
		_, err := apiDeleteByID(context)
		if err != nil {
			t.Error(err)
		}
	}

	//--------------------------------------------------------------------------------------------------------------
	// create post as user
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}
	context.RequestContent = map[string]interface{}{
		"identifier": "identifier",
	}
	session.Set(api.ConstSessionKeyAdminRights, false)

	result, createError := apiCreateBlogPost(context)
	if createError == nil {
		t.Error("blog post created by user")
	}

	//--------------------------------------------------------------------------------------------------------------
	// create few posts with delays to get different update time
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, true)

	tagFormat := 0
	for i := 0; i < numberOfPosts; i++ {
		published := "false"
		if i < numberOfPublished {
			published = "true"
		}

		tags := []interface{}{}
		if i < numberOfPostsWithTag {
			switch tagFormat {
			case 0:
				{
					tagFormat = 1
					tags = append(tags, "tag", "tag0")
					numberOfPostsWithSpecialTags++
				}
			case 1:
				{
					tagFormat = 2
					tags = append(tags, "tag1", "tag")
					numberOfPostsWithSpecialTags++
				}
			case 2:
				{
					tagFormat = 3
					tags = append(tags, "tag")
				}
			case 3:
				{
					tagFormat = 0
					tags = append(tags, "any", "tag", "tag3")
				}
			}
		}

		context.RequestContent = map[string]interface{}{
			"identifier": "identifier" + fmt.Sprintf("%03d", i),
			"published":  published,
			"tags":       tags,
		}

		fmt.Println("Created [" + utils.InterfaceToString(i+1) + "] posts of [" + utils.InterfaceToString(numberOfPosts) + "] scheduled. Wait 1 second.")
		time.Sleep(time.Second)
		_, err := apiCreateBlogPost(context)
		if err != nil {
			t.Error(err)
		}
	}

	//--------------------------------------------------------------------------------------------------------------
	// check count
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, true)

	result, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	if result != numberOfPosts {
		t.Error("Not all posts have been created [" + utils.InterfaceToString(result) + "] != [" + utils.InterfaceToString(numberOfPosts) + "].")
	}

	//--------------------------------------------------------------------------------------------------------------
	// check count by user (published)
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	result, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	if result != numberOfPublished {
		t.Error("Number of published is incorrect [" + utils.InterfaceToString(result) + "] != [" + utils.InterfaceToString(numberOfPublished) + "]")
	}

	//--------------------------------------------------------------------------------------------------------------
	// check prev / next posts by user
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	postsResult, err := apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	posts := utils.InterfaceToArray(postsResult)

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	idStack := [2]string{"", ""}
	var prevPostMap map[string]interface{}

	for i, currentPost := range posts {
		postMap := utils.InterfaceToMap(currentPost)
		context.RequestArguments = map[string]string{
			"id": utils.InterfaceToString(postMap["_id"]),
		}

		resultByID, err := apiPostByID(context)
		if err != nil {
			t.Error(err)
		}
		result := utils.InterfaceToMap(resultByID)

		if i == 0 {
			extra := utils.InterfaceToMap(result["extra"])
			if extra["prev"] != nil {
				t.Error("previous post exists for first record")
			}
		} else if i == len(posts)-1 {
			extra := utils.InterfaceToMap(result["extra"])
			if extra["next"] != nil {
				t.Error("next post exists for last record")
			}
		} else if i >= 2 {
			extra := utils.InterfaceToMap(prevPostMap["extra"])
			extraPrev := utils.InterfaceToMap(extra["prev"])
			extraNext := utils.InterfaceToMap(extra["next"])
			if extraPrev["_id"] != idStack[0] {
				t.Error("previous post incorrect")
			}
			if extraNext["_id"] != result["_id"] {
				t.Error("next post incorrect")
			}
		}
		prevPostMap = result
		idStack[0] = idStack[1]
		idStack[1] = utils.InterfaceToString(result["_id"])
	}

	//--------------------------------------------------------------------------------------------------------------
	// update post
	//--------------------------------------------------------------------------------------------------------------

	// get first post
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"limit": "0,1",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, true)

	postsResult, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	firstPost := utils.InterfaceToArray(postsResult)[0]

	postHashMap := utils.InterfaceToMap(firstPost)

	postHashMap["identifier"] = utils.InterfaceToString(postHashMap["identifier"]) + "_updated"
	postHashMap["published"] = utils.InterfaceToString(!utils.InterfaceToBool(postHashMap["identifier"]))
	postHashMap["title"] = utils.InterfaceToString(postHashMap["identifier"]) + "_updated"
	postHashMap["excerpt"] = utils.InterfaceToString(postHashMap["excerpt"]) + "_updated"
	postHashMap["content"] = utils.InterfaceToString(postHashMap["content"]) + "_updated"
	postHashMap["featured_image"] = "updated_" + utils.InterfaceToString(postHashMap["featured_image"])

	tags := utils.InterfaceToArray(postHashMap["tags"])
	tags = append(tags, "updated")
	postHashMap["tags"] = tags

	// update post
	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"id": utils.InterfaceToString(postHashMap["_id"]),
	}
	context.RequestContent = map[string]interface{}{
		"identifier":     postHashMap["identifier"],
		"published":      postHashMap["published"],
		"title":          postHashMap["title"],
		"excerpt":        postHashMap["excerpt"],
		"content":        postHashMap["content"],
		"tags":           postHashMap["tags"],
		"featured_image": postHashMap["featured_image"],
	}
	session.Set(api.ConstSessionKeyAdminRights, true)

	result, err = apiUpdateByID(context)
	if err != nil {
		t.Error(err)
	}
	resultHashMap := utils.InterfaceToMap(result)

	for key, value := range utils.InterfaceToMap(context.RequestContent) {
		resultStr := utils.InterfaceToString(resultHashMap[key])
		valueStr := utils.InterfaceToString(value)
		if resultStr != valueStr {
			msg := "updated [" + key + "] [" + resultStr + "] != [" + valueStr + "]"
			t.Error(msg)
		}
	}

	//--------------------------------------------------------------------------------------------------------------
	// check count by user (published)
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
		"limit":  "0,1000",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	result, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	if result != numberOfPublished {
		t.Error("Number of published is incorrect [" + utils.InterfaceToString(result) + "] != [" + utils.InterfaceToString(numberOfPublished) + "]")
	}

	//--------------------------------------------------------------------------------------------------------------
	// check count by user (tag)
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
		"limit":  "0,1000",
		"tags":   "tag1,tag0",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	result, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	if result != numberOfPostsWithSpecialTags {
		t.Error("Number of posts with tag [tag] is incorrect [" + utils.InterfaceToString(result) + "] != [" + utils.InterfaceToString(numberOfPostsWithTag) + "]")
	}

	//--------------------------------------------------------------------------------------------------------------
	// check count by user (tag)
	//--------------------------------------------------------------------------------------------------------------

	context.ContextValues = map[string]interface{}{}
	context.RequestArguments = map[string]string{
		"action": "count",
		"limit":  "0,1000",
		"tags":   "tag",
	}
	context.RequestContent = map[string]interface{}{}
	session.Set(api.ConstSessionKeyAdminRights, false)

	result, err = apiListPosts(context)
	if err != nil {
		t.Error(err)
	}

	if result != numberOfPostsWithTag {
		t.Error("Number of posts with tag [tag] is incorrect [" + utils.InterfaceToString(result) + "] != [" + utils.InterfaceToString(numberOfPostsWithTag) + "]")
	}

	//--------------------------------------------------------------------------------------------------------------
	//
	//--------------------------------------------------------------------------------------------------------------

}
