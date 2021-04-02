package post

import (
	"encoding/json"
	"github.com/user/stories/business/post"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var paramsUrl map[string]interface{}

func HandlerPostController(r *http.Request) map[string]interface{} {
	var validController = regexp.MustCompile("(post)(?:\\/[a-zA-Z0-9]+|)")
	m := validController.FindStringSubmatch(r.URL.Path)
	//fmt.Println("m", m)
	dataRoute := strings.Split(m[0], "/")
	//fmt.Println("dataRoute", dataRoute)

	l, _ := url.Parse(r.URL.RawQuery)
	ss := strings.Split(l.String(), "&")
	mm := map[string]interface{}{}

	if len(ss) > 0 {
		for _, pair := range ss {
			z := strings.Split(pair, "=")
			if len(z) == 1 && z[0] != "" {
				mm[z[0]] = ""
			}

			if len(z) == 2 {
				mm[z[0]] = z[1]
			}
		}
		paramsUrl = mm
	}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.

	var data map[string]interface{}
	switch r.Method {
	case "GET":
		if len(dataRoute) == 2 {
			data = get(dataRoute[1])
		} else {
			data = getList()
		}
	case "POST":
		var paramsPost struct {
			Title       string
			Description string
			Content     string
			CateId      int
			AuthId      int
			CreateDate  string
		}
		err := json.NewDecoder(r.Body).Decode(&paramsPost)
		if err != nil {
			panic(err)
		}
		data = create(map[string]interface{}{
			"title":       paramsPost.Title,
			"description": paramsPost.Description,
			"content":     paramsPost.Content,
			"cateId":      paramsPost.CateId,
			"authId":      paramsPost.AuthId,
			"createDate":  paramsPost.CreateDate,
		})
	case "PUT":
		var paramsPost struct {
			Title       string
			Description string
			Content     string
			CateId      int
			AuthId      int
			//Image       int
			//Status      string
			//TotalView   string
		}
		err := json.NewDecoder(r.Body).Decode(&paramsPost)
		if err != nil {
			panic(err)
			//http.Error(w, err.Error(), http.StatusBadRequest)
			//return
		}
		data = update(dataRoute[1], map[string]interface{}{
			"title":       paramsPost.Title,
			"description": paramsPost.Description,
			"content":     paramsPost.Content,
			"cateId":      paramsPost.CateId,
			"authId":      paramsPost.AuthId,
		})
	case "DELETE":
		if len(dataRoute) < 2 {
			data = map[string]interface{}{
				"status":  false,
				"message": "Invalid",
			}
		} else {
			data = deletePost(dataRoute[1])
		}
	}
	return data
}

func getList() map[string]interface{} {
	return post.GetAllPost(paramsUrl)
}

func get(id string) map[string]interface{} {
	return post.GetDetail(id)
}

func create(params map[string]interface{}) map[string]interface{} {
	return post.CreatePost(params)
}

func update(id string, params map[string]interface{}) map[string]interface{} {
	return post.UpdatePost(id, params)
}

func deletePost(id string) map[string]interface{} {
	return post.DeletePost(id)
}
