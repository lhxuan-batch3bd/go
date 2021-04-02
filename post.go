//businnes
package post

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/user/stories/mysql/author"
	"github.com/user/stories/mysql/category"
	"github.com/user/stories/mysql/post"
	post2 "github.com/user/stories/redis/post"
	"strconv"
	"strings"
	"time"
)

func GetAllPost(params map[string]interface{}) map[string]interface{} {
	type data struct {
		Id          int
		Title       string
		Slug        string
		Description string
		CateId      int
		CateName    string
		AuthId      int
		AuthName    string
	}
	var filters = map[string]interface{}{}

	keyCache := ""

	if val, ok := params["az"]; ok == true {
		filters["az"] = val.(string)
		keyCache += "az" + val.(string)
	}
	if val, ok := params["sort"]; ok == true {
		filters["sort"] = val.(string)
		keyCache += "sort" + val.(string)
	}
	var lm int
	if val, ok := params["limit"]; ok == true {
		limit, err := strconv.Atoi(val.(string))
		if err == nil {
			lm = limit
			filters["limit"] = limit
			keyCache += "limit" + val.(string)
		}
	}
	if val, ok := params["page"]; ok {
		page, err := strconv.Atoi(val.(string))
		if err == nil {
			filters["offset"] = (page - 1) * lm
			keyCache += "page" + val.(string)
		}
	}
	if val, ok := params["search"]; ok == true {
		filters["search"] = val.(string)
		keyCache += "search" + val.(string)
	}

	dataCache := post2.GetAll(keyCache)
	total := post.GetTotal(filters)

	if len(dataCache) > 0 {
		return map[string]interface{}{
			"total": total,
			"data":  dataCache,
		}
	}

	listPost := post.Get(filters)

	if len(listPost) == 0 {
		return map[string]interface{}{
			"total": total,
			"data":  []map[string]interface{}{},
		}
	}

	var arrCateId []string
	var arrAuthorId []string

	for _, val := range listPost {
		cateId := val.CateId
		authorId := val.AuthId

		arrCateId = append(arrCateId, strconv.Itoa(cateId))
		arrAuthorId = append(arrAuthorId, strconv.Itoa(authorId))
	}

	listCate := category.Get(map[string]interface{}{
		"arrCateId": strings.Join(arrCateId, ","),
	})
	listAuth := author.Get(map[string]interface{}{
		"arrAuthorId": strings.Join(arrAuthorId, ","),
	})

	var arrData []data

	dataMappingCate := map[int]string{}
	for _, valCate := range listCate {
		dataMappingCate[valCate.Id] = valCate.Name
	}

	for _, val := range listPost {
		cateId := val.CateId
		authId := val.AuthId

		cateName := dataMappingCate[cateId]
		authName := ""
		for _, valAuth := range listAuth {
			if valAuth.Id == authId {
				authName = valAuth.Name
			}
		}
		p := data{
			Id:          val.Id,
			Title:       val.Title,
			Slug:        val.Slug,
			Description: val.Description,
			CateId:      val.CateId,
			CateName:    cateName,
			AuthId:      val.AuthId,
			AuthName:    authName,
		}

		arrData = append(arrData, p)
	}

	//set cache
	if len(arrData) > 0{
		b, _ := json.Marshal(arrData)
		post2.SetAll(keyCache, b)
	}

	return map[string]interface{}{
		"total": total,
		"data":  arrData,
	}

}

func CreatePost(params map[string]interface{}) map[string]interface{} {

	fmt.Println("params", params)

	title, ok := params["title"]
	if !ok {
		return map[string]interface{}{
			"status":  false,
			"message": "Params title required",
		}
	}
	description, ok := params["description"]
	if !ok {
		return map[string]interface{}{
			"status":  false,
			"message": "Params content required",
		}
	}
	content, ok := params["content"]
	if !ok {
		return map[string]interface{}{
			"status":  false,
			"message": "Params content required",
		}
	}
	cateId, ok := params["cateId"]
	if !ok {
		return map[string]interface{}{
			"status":  false,
			"message": "Params cateId required",
		}
	}
	authId, ok := params["authId"]
	if !ok {
		return map[string]interface{}{
			"status":  false,
			"message": "Params authId required",
		}
	}
	//Check
	postList := post.Get(map[string]interface{}{
		"Title":  params["title"],
		"CateId": params["cateId"],
		"AuthId": params["authId"],
	})

	if len(postList) != 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "The Post already exists",
		}
	}
	//Insert
	id, _ := post.Insert(post.Post{
		Title:       title.(string),
		Slug:        slug.Make(title.(string)),
		Description: description.(string),
		Content:     content.(string),
		CateId:      cateId.(int),
		AuthId:      authId.(int),
		CreateDate:  time.Now().Format("2006-01-02 15:04:05"),
	})
	//return
	return map[string]interface{}{
		"status":  true,
		"message": "Success",
		"postId":  id,
	}
}

func GetDetail(id string) map[string]interface{} {

	type data struct {
		Id          int
		Title       string
		Slug        string
		Description string
		Content     string
		CateId      int
		CateName    string
		AuthId      int
		AuthName    string
		CreateDate  string
	}

	dataCache := post2.GetDetail(id)

	if len(dataCache) > 0 {
		fmt.Println("dataCache", dataCache)
		pC := data{
			Id:          int(dataCache["Id"].(float64)),
			Title:       dataCache["Title"].(string),
			Slug:        dataCache["Slug"].(string),
			Description: dataCache["Description"].(string),
			Content:     dataCache["Content"].(string),
			CateId:      int(dataCache["CateId"].(float64)),
			CateName:    dataCache["CateName"].(string),
			AuthId:      int(dataCache["AuthId"].(float64)),
			AuthName:    dataCache["AuthName"].(string),
			CreateDate:  dataCache["CreateDate"].(string),
		}
		return map[string]interface{}{
			"postDetail": pC,
		}
	}

	listPost := post.Get(map[string]interface{}{
		"Id": id,
	})

	if len(listPost) == 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "Id does not exist ",
		}
	}

	postDetail := listPost[0]

	//cate
	cateName := ""
	cate := category.Get(map[string]interface{}{
		"Id": postDetail.CateId,
	})
	if len(cate) > 0 {
		cateDetail := cate[0]
		cateName = cateDetail.Name
	}

	//auth
	authName := ""
	auth := author.Get(map[string]interface{}{
		"Id": postDetail.AuthId,
	})
	if len(auth) > 0 {
		authDetail := auth[0]
		authName = authDetail.Name
	}

	//set cache
	post2.SetDetail(id, map[string]interface{}{
		"Id":          postDetail.Id,
		"Title":       postDetail.Title,
		"Slug":        postDetail.Slug,
		"Description": postDetail.Description,
		"Content":     postDetail.Content,
		"CateId":      postDetail.CateId,
		"CateName":    cateName,
		"AuthId":      postDetail.AuthId,
		"AuthName":    authName,
		"CreateDate":  postDetail.CreateDate,
	})

	p := data{
		Id:          postDetail.Id,
		Title:       postDetail.Title,
		Slug:        postDetail.Slug,
		Description: postDetail.Description,
		Content:     postDetail.Content,
		CateId:      postDetail.CateId,
		CateName:    cateName,
		AuthId:      postDetail.AuthId,
		AuthName:    authName,
		CreateDate:  postDetail.CreateDate,
	}

	return map[string]interface{}{
		"postDetail": p,
	}
}

func UpdatePost(id string, params map[string]interface{}) map[string]interface{} {

	listPost := post.Get(map[string]interface{}{
		"Id": id,
	})
	if len(listPost) == 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "Id does not exist ",
		}
	}
	if len(params) == 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "There is no params value",
		}
	}

	var filters = map[string]interface{}{}
	for _, val := range listPost {
		filters["Title"] = val.Title
		filters["Description"] = val.Description
		filters["Content"] = val.Content
		filters["CateId"] = val.CateId
		filters["AuthId"] = val.AuthId
	}

	title, _ := params["title"]
	if len(title.(string)) != 0 {
		filters["Title"] = title
	}

	description, _ := params["description"]
	if len(description.(string)) != 0 {
		filters["Description"] = description
	}

	content, _ := params["content"]
	if len(content.(string)) != 0 {
		filters["Content"] = content
	}

	cateId := category.Get(map[string]interface{}{
		"Id": params["cateId"].(int),
	})
	if len(cateId) == 0 {
		fmt.Println(params["cateId"])
		if params["cateId"] != 0 {
			return map[string]interface{}{
				"status":  false,
				"message": "The CateId does not exist",
			}
		}
	}
	if params["cateId"] != 0 {
		filters["CateId"] = params["cateId"]
	}

	authId := author.Get(map[string]interface{}{
		"Id": params["authId"],
	})
	if len(authId) == 0 {
		if params["authId"] != 0 {
			return map[string]interface{}{
				"status":  false,
				"message": "The AuthId does not exist",
			}
		}
	}
	if params["authId"] != 0 {
		filters["AuthId"] = params["authId"].(int)
	}

	checkPost := post.Get(map[string]interface{}{
		"Title":  filters["Title"],
		"CateId": filters["CateId"],
		"AuthId": filters["AuthId"],
	})

	if len(checkPost) != 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "The Post already exists",
		}
	}

	Id, _ := strconv.Atoi(id)
	_ = post.Update(filters, Id)

	return map[string]interface{}{
		"status":  true,
		"message": "The Post has been updated",
	}
}

func DeletePost(id string) map[string]interface{} {

	del := post.Get(map[string]interface{}{
		"Id": id,
	})
	if len(del) == 0 {
		return map[string]interface{}{
			"status":  false,
			"message": "Id does not exists",
		}
	}

	Id, _ := strconv.Atoi(id)
	_ = post.Delete(Id)

	return map[string]interface{}{
		"status":  true,
		"message": "Success",
	}
}
//redis
package post

import (
	"encoding/json"
	"github.com/user/stories/redis"
	"time"
)

const keyNameDetail = "Cache_Post_Detail"
const keyNameList = "Cache_Post_List"

func GetDetail(id string) map[string]interface{} {
	keyCache := keyNameDetail + ":" + id
	redisClient := redis.GetRedisClient()

	dataCache, _ := redisClient.Get(keyCache).Result()

	var dataDecode map[string]interface{}
	if dataCache != "" {
		json.Unmarshal([]byte(dataCache), &dataDecode)
	}

	return dataDecode
}

func SetDetail(id string, data map[string]interface{}) bool {
	keyCache := keyNameDetail + ":" + id
	redisClient := redis.GetRedisClient()
	dataByte, _ := json.Marshal(data)
	redisClient.Set(keyCache, dataByte, 60*time.Minute)
	return true
}

func GetAll(key string) []map[string]interface{} {
	keyCache := keyNameList + ":" + key
	redisClient := redis.GetRedisClient()

	dataCache, _ := redisClient.Get(keyCache).Result()

	var dataDecode []map[string]interface{}
	if dataCache != "" {
		json.Unmarshal([]byte(dataCache), &dataDecode)
	}

	return dataDecode
}

func SetAll(key string, dataByte []byte) bool {
	keyCache := keyNameList + ":" + key
	redisClient := redis.GetRedisClient()
	redisClient.Set(keyCache, dataByte, 60*time.Minute)
	return true
}
//api
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
//mysql
package post

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/user/stories/mysql"
	"log"
	"strconv"
	"strings"
	"time"

	//"time"
)

const table string = "posts"

const PrimaryKey string = "Id"

const StatusActive int = 1
const StatusInActive int = 2
const StatusRemove int = 0

const IsComplete int = 1

type Post struct {
	Id          int
	Title       string `json:"title"`
	Slug        string
	Description string `json:"description"`
	Content     string
	CateId      int
	AuthId      int
	Image       *string
	Status      int
	TotalView   int
	CreateDate  string
	UpdateDate  *string
}

func Insert(info Post) (id int, err error) {
	db := mysql.DbConnect()
	stmt, err := db.Prepare("insert into " + table + "(Title, Slug, Description, Content, CateId, AuthId, Image, Status, TotalView,CreateDate) value(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//Slug:=slug.Make(Title)
	res, err := stmt.Exec(info.Title, info.Slug, info.Description, info.Content, info.CateId, info.AuthId, info.Image, info.Status, info.TotalView, info.CreateDate)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return int(lastId), err
}

func Get(params map[string]interface{}) []Post {
	db := mysql.DbConnect()

	var where = buildWhere(params)

	var query = "SELECT * FROM " + table + " WHERE " + where

	if val, ok := params["sort"]; ok {
		query += " order by " + val.(string)
	} else {
		query += " order by " + PrimaryKey
	}
	if val, ok := params["az"]; ok {
		query += " " + val.(string)
	} else {
		query += " desc"
	}
	if val, ok := params["limit"]; ok {
		query += " limit " + strconv.Itoa(val.(int))
	} else {
		query += " limit 10"
	}
	if val, ok := params["offset"]; ok {
		query += " offset " + strconv.Itoa(val.(int))
	} else {
		query += " offset 0"
	}

	fmt.Println("query", query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var res []Post
	for rows.Next() {
		p := Post{}
		err = rows.Scan(&p.Id, &p.Title, &p.Slug, &p.Description, &p.Content, &p.CateId, &p.AuthId, &p.Image, &p.Status, &p.TotalView, &p.CreateDate, &p.UpdateDate)
		if err != nil {
			panic(err)
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
	defer db.Close()
	return res
}

func GetTotal(params map[string]interface{}) int {
	db := mysql.DbConnect()
	var where = buildWhere(params)
	var query = "SELECT COUNT(*) FROM " + table + " WHERE " + where

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	defer rows.Close()
	return count
}

func Update(data map[string]interface{}, id int) bool {
	if id <= 0 {
		return false
	}
	sql := ""
	var paramsBiding []interface{}
	if val, ok := data["Title"]; ok == true {
		sql += " Title=?,"
		paramsBiding = append(paramsBiding, val)
		sql += " Slug=?,"
		paramsBiding = append(paramsBiding,slug.Make(val.(string)))
	}
	if val, ok := data["Description"]; ok == true {
		sql += " Description=?,"
		paramsBiding = append(paramsBiding, val)
	}
	if val, ok := data["Content"]; ok == true {
		sql += " Content=?,"
		paramsBiding = append(paramsBiding, val)
	}
	if val, ok := data["CateId"]; ok == true {
		sql += " CateId=?,"
		paramsBiding = append(paramsBiding, val)
	}
	if val, ok := data["AuthId"]; ok == true {
		sql += " AuthId=?,"
		paramsBiding = append(paramsBiding, val)
	}
	if val, ok := data["Status"]; ok == true {
		sql += " Status=?,"
		paramsBiding = append(paramsBiding, val)
	}
	if val, ok := data["TotalView"]; ok == true {
		sql += " TotalView=?,"
		paramsBiding = append(paramsBiding, val)
	}
	sql += " UpdateDate=?,"
	paramsBiding = append(paramsBiding, time.Now().Format("2006-01-02 15:04:05"))

	if sql == "" {
		return false
	}
	sql = strings.TrimSuffix(sql, ",")
	sql = "update " + table + " set" + sql + " WHERE Id =?"
	paramsBiding = append(paramsBiding, id)
	db := mysql.DbConnect()
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	rs, err := stmt.Exec(paramsBiding...)
	if err != nil {
		log.Fatal(err)
	}
	_, err = rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return true
}

func buildWhere(params map[string]interface{}) string {
	var where = " 1=1"

	if val, ok := params["Id"]; ok {
		where += " AND Id = " + fmt.Sprint(val)
	}
	if val, ok := params["Title"]; ok {
		where += " AND Title = '" + val.(string) + "'"
	}
	if val, ok := params["Description"]; ok {
		where += " AND Description = '" + val.(string) + "'"
	}
	if val, ok := params["Content"]; ok {
		where += " AND Content = '" + val.(string) + "'"
	}
	if val, ok := params["CateId"]; ok {
		where += " AND CateId = " + fmt.Sprint(val)
	}
	if val, ok := params["AuthId"]; ok {
		where += " AND AuthId = " + fmt.Sprint(val)
	}
	if val, ok := params["TotalView"]; ok {
		where += " AND TotalView = " + fmt.Sprint(val)
	}
	if val, ok := params["Status"]; ok {
		where += " AND Status = " + fmt.Sprint(val)
	}

	if val, ok := params["search"]; ok {
		where += " AND Title LIKE '%" + val.(string) + "%'"
	}
	return where
}

func Delete(id int) bool {
	if id <= 0 {
		return false
	}
	db := mysql.DbConnect()

	sql := "DELETE FROM " + table + " WHERE id=?"
	res, err := db.Exec(sql, id)
	if err != nil {
		log.Fatal(err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return true
}
//cli
package main

import (
	"github.com/user/stories/config"
)

func main() {
	config.InitLoad()

	//	//Categories Table
	//	//arrData := []map[string]string{
	//	//	{
	//	//		"Name":     "Windows 10",
	//	//		"Slug":     "windows-10",
	//	//		"Position": "2",
	//	//	}, {
	//	//		"Name":     "Iphone",
	//	//		"Slug":     "iphone",
	//	//		"Position": "3",
	//	//	}, {
	//	//		"Name":     "Android",
	//	//		"Slug":     "android",
	//	//		"Position": "4",
	//	//	}, {
	//	//		"Name":     "Tan Cong Mang",
	//	//		"Slug":     "tan-cong-mang",
	//	//		"Position": "5",
	//	//	}, {
	//	//		"Name":     "Tien Ich",
	//	//		"Slug":     "tien-ich",
	//	//		"Position": "6",
	//	//	}, {
	//	//		"Name":     "Ung Dung",
	//	//		"Slug":     "ung-dung",
	//	//		"Position": "7",
	//	//	}, {
	//	//		"Name":     "Lap Trinh",
	//	//		"Slug":     "lap-trinh",
	//	//		"Position": "8",
	//	//	}, {
	//	//		"Name":     "Dien May",
	//	//		"Slug":     "dien-may",
	//	//		"Position": "9",
	//	//	}, {
	//	//		"Name":     "Video",
	//	//		"Slug":     "video",
	//	//		"Position": "10",
	//	//	},
	//	//}
	//	////
	//	//for _, v := range arrData {
	//	//	position, _ := strconv.Atoi(v["Position"])
	//	//	cateId, err := category.Insert(
	//	//		category.Cate{
	//	//			Name:     v["Name"],
	//	//			Slug:     v["Slug"],
	//	//			Position: position,
	//	//		})
	//	//
	//	//	println("cateId", cateId)
	//	//	println("err", err)
	//	//}

	//	//cateId, err := category.Insert(category.Cate{
	//	//	Name:     "Game",
	//	//	Slug:     "game",
	//	//	Position: 11,
	//	//})
	//	//
	//	//println("cateId", cateId)
	//	//println("err", err)

	//	//res := category.Get(map[string]interface{}{
	//	//	//"CateId" : 1,
	//	//	"limit" : 5,
	//	//	//"Id" : 1,
	//	//})

	//	//fmt.Println(res)
	//	//os.Exit(1)
	//	//Update
	//	//res:=category.Update(map[string]interface{}{
	//	//	"Name": "Điện Máy",
	//	//},9)
	//	//fmt.Println(res)
	//	//os.Exit(1)
	//	//Delete
	//	//res:=category.Delete(12)
	//	//fmt.Println(res)
	//	//os.Exit(1)

	//	//Posts Table
	//	//arrDataPost := []map[string]string{
	//	//	{
	//	//		"Title": "Những bí ẩn thú vị về các loại quả mà chúng ta ăn thường ngày",
	//	//		"Slug":  "Nhung-bi-an-thu-vi-ve-cac-loai-qua-ma-chung-ta-an-thuong-ngay",
	//	//		"CateId":      "1",
	//	//		"Description": "Những loại quả quen thuộc mà chúng ta vẫn ăn và nhìn thấy hàng ngày đều ẩn chứa những sự thật thú vị mà có lẽ bạn chưa bao giờ được biết.",
	//	//		"Status":      "1",
	//	//		//"CreateDate": "time.Now()",
	//	//	}, {
	//	//		"Title": "Cách thiết lập cấu hình SSH tùy chỉnh trong Windows Terminal",
	//	//		"Slug":  "Cach-thiet-lap-cau-hinh-ssh-tuy-chinh-trong-windows-terminal",
	//	//		"CateId":      "1",
	//	//		"Description": "Một trong những tính năng làm nên sự tuyệt vời của Windows Terminal chính là khả năng thiết lập cấu hình SSH tùy chỉnh,",
	//	//		"Status":      "1",
	//	//	}, {
	//	//		"Title": "Mã độc này được viết bằng một ngôn ngữ lập trình bất thường, khiến nó cực kỳ khó bị phát hiện",
	//	//		"Slug":  "Ma-doc-nay-duoc-viet-bang-mot-ngon-ngu-lap-trinh-bat-thuong,khiến nó cực kỳ khó bị phát hiện",
	//	//		"CateId":      "1",
	//	//		"Description": "Mã độc này có tên NimzaLoader và được viết bằng ngôn ngữ lập trình Nim",
	//	//		"Status":      "1",
	//	//
	//	//	},
	//	//}
	//	//
	//	//for _, v := range arrDataPost {
	//	//	cateid, _  := strconv.Atoi(v["CateId"])
	//	//	status, _ := strconv.Atoi(v["Status"])
	//	//
	//	//	PostId, err := post.Insert(
	//	//		post.Post{
	//	//			Title: v["Title"],
	//	//			Slug:  v["Slug"],
	//	//			CateId:      cateid,
	//	//			Description: v["Description"],
	//	//			Status:      status,
	//	//		})
	//	//
	//	//	println("PostId", PostId)
	//	//	println("err", err)
	//	//}

	//	//postId, err := post.Insert(post.Post{
	//	//	Title:       "Title",
	//	//	Slug:        "Slug",
	//	//	//CateId:      1,
	//	//	Description: "Description",
	//	//	Status:      1,
	//	//	CreateDate: time.Now(),
	//	//})
	//	//
	//	//println("postId", postId)
	//	//println("err", err)
	//
	//	//res := post.Get(map[string]interface{}{
	//	//	//"CateId" : 1,
	//	//	//"limit" : 3,
	//	//	//"Id" : 3,
	//	//	"Title" : "Title",
	//	//})
	//	//
	//	//fmt.Println(res)
	//	//os.Exit(1)
	//
	//res := post.Update(map[string]interface{}{
	//	//"Title": "Youth",
	//	"CateId": 9,
	//	//"AuthId" :4,
	//	//"TotalView":0,
	//	//"Description": "Description",
	//	//"Content": "Nội dung",
	//	"UpdateDate": time.Now().Format("2006-01-02 15:04:05"),
	//}, 1)
	//fmt.Println(res)
	//os.Exit(1)

	//	//res:=post.Delete(20)
	//	//fmt.Println(res)
	//	//os.Exit(1)

	//	//Author Table
	//	//arrDataAuthor := []map[string]string{
	//	//	{
	//	//		"Name": "Phạm Hải",
	//	//		"Slug":  "pham-hai",
	//	//		"Status":      "1",
	//	//	}, {
	//	//		"Name": "Phương Phùng",
	//	//		"Slug":  "phuong-phung",
	//	//		"Status":      "1",
	//	//	}, {
	//	//		"Name": "Nguyễn Nhật Minh",
	//	//		"Slug":  "nguyen-nhat-minh",
	//	//		"Status":      "1",
	//	//	},
	//	//}
	//	//
	//	//for _, v := range arrDataAuthor {
	//	//	status, _ := strconv.Atoi(v["Status"])
	//	//	authorId, err := author.Insert(
	//	//		author.Author{
	//	//			Name: v["Name"],
	//	//			Slug:  v["Slug"],
	//	//			Status:      status,
	//	//		})
	//	//
	//	//	println("authorId", authorId)
	//	//	println("err", err)
	//	//}
	//	//
	//	//authorId, err := author.Insert(author.Author{
	//	//	Name:       "Name",
	//	//	Slug:        "Slug",
	//	//	Status:      1,
	//	//	CreateDate: time.Now(),
	//	//})
	//	//
	//	//println("authorId", authorId)
	//	//println("err", err)
	//	//
	//	//res := author.Get(map[string]interface{}{
	//	//	//"limit" : 3,
	//	//	"Id" : 3,
	//	//})
	//	//
	//	//fmt.Println(res)
	//	//os.Exit(1)
	//
	//	res:=author.Update(map[string]interface{}{
	//		"Name": "Nguyễn Nhật Minh",
	//		"UpdateDate" : time.Now().Format("2006-01-02 15:04:05"),
	//	},4)
	//	fmt.Println(res)
	//	os.Exit(1)
	//	//
	//	//res:=author.Delete(1)
	//	//fmt.Println(res)
	//	//os.Exit(1)
	//
	//}

	// Users struct which contains
	// an array of users
	//type Posts struct {
	//	Posts []Post `json:"post"`
	//}

	// User struct which contains a name
	// a type and a list of social links
	//type Posts struct {
	//	Title       string `json:"title"`
	//	Author      string `json:"author"`
	//	Category    string `json:"category"`
	//	Description string `json:"description"`

}

//func main() {
//	config.InitLoad()
//// Open our jsonFile
//jsonFile, err := os.Open("post.json")
//// if we os.Open returns an error then handle it
//if err != nil {
//	fmt.Println(err)
//}
//
//fmt.Println("Successfully Opened post.json")
//// defer the closing of our jsonFile so that we can parse it later on
//defer jsonFile.Close()
//
//// read our opened xmlFile as a byte array.
//byteValue, _ := ioutil.ReadAll(jsonFile)
//
//// we initialize our Users array
//var posts []Posts
//// we unmarshal our byteArray which contains our
//// jsonFile's content into 'users' which we defined above
//json.Unmarshal(byteValue, &posts)
//// we iterate through every user within our users array and
//// print out the user Type, their name, and their facebook url
//// as just an example
////fmt.Println(posts)
//
//for _, val := range posts {
//	cateList := category.Get(map[string]interface{}{
//		"Name": val.Category,
//	})
//	if len(cateList) == 0 {
//		continue
//	}
//	cate := cateList[0]
//	authList := author.Get(map[string]interface{}{
//		"Name": val.Author,
//	})
//	if len(authList) == 0 {
//		newAuthor, _ := author.Insert(author.Author{
//			Name:       val.Author,
//			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
//		})
//		fmt.Println("new author: ", newAuthor)
//		continue
//	}
//	auth := authList[0]
//	postList := post.Get(map[string]interface{}{
//		"Title":  val.Title,
//		"CateId": cate.Id,
//		"AuthId": auth.Id,
//	})
//	if len(postList) != 0 {
//		continue
//	}
//
//	fmt.Println(val.Title)
//	fmt.Println(slug.Make(val.Title))
//	addId, _ := post.Insert(
//		post.Post{
//			Title:       val.Title,
//			Slug:        slug.Make(val.Title),
//			CateId:      cate.Id,
//			AuthId:      auth.Id,
//			Description: val.Description,
//			CreateDate:  time.Now().Format("2006-01-02 15:04:05"),
//		})
//	fmt.Println("Add Id: ", addId)
//	//	//os.Exit(1)
//}
//}

