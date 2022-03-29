package todo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cod3rboy/go-todo-api/utils"
)

var service = http.NewServeMux()

const contentTypeJson = "application/json"

var methodContentType = map[string]string{
	http.MethodPost:  contentTypeJson,
	http.MethodPatch: contentTypeJson,
	http.MethodPut:   contentTypeJson,
}

func init() {
	service.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		path = strings.TrimSpace(path)
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimSuffix(path, "/")
		pathSegments := strings.Split(path, "/")
		segmentCount := len(pathSegments)
		match := !(segmentCount > 2)
		id := ""
		if segmentCount == 2 {
			id = pathSegments[1]
		}

		if !match {
			http.NotFound(rw, req)
			return
		}

		if contentType, found := methodContentType[req.Method]; found && req.Header.Get("content-type") != contentType {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
			rw.Write([]byte("unsupported content type"))
			return
		}

		switch req.Method {
		case http.MethodGet:
			get(id, rw, req)
		case http.MethodPost:
			if id != "" {
				http.NotFound(rw, req)
			} else {
				create(rw, req)
			}
		case http.MethodPut:
			put(id, rw, req)
		case http.MethodDelete:
			remove(id, rw, req)
		default:
			utils.MethodNotAllowed(rw, req)
		}
	})
}

func GetService() *http.ServeMux {
	return service
}

type timestamp = int64

type ToDoItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	Done        bool      `json:"done"`
	Due         timestamp `json:"due"`
}

var todoList = make(map[string]*ToDoItem)

func create(rw http.ResponseWriter, req *http.Request) {
	todoItem, err := unmarshalBody(req)
	if err != nil {
		utils.BadRequest(rw, req)
		return
	}
	todoItem.ID = utils.GenerateRandomID()
	todoList[todoItem.ID] = todoItem
	sendResponse(rw, todoItem, http.StatusCreated)
}

func get(id string, rw http.ResponseWriter, req *http.Request) {
	if id == "" {
		todoItems := make([]*ToDoItem, len(todoList))
		i := 0
		for _, item := range todoList {
			todoItems[i] = item
			i++
		}
		sendResponse(rw, todoItems, http.StatusOK)
	} else {
		item, found := todoList[id]
		if !found {
			http.NotFound(rw, req)
			return
		}
		sendResponse(rw, item, http.StatusOK)
	}
}

func put(id string, rw http.ResponseWriter, req *http.Request) {
	item, found := todoList[id]
	if !found {
		http.NotFound(rw, req)
		return
	}
	newItem, err := unmarshalBody(req)
	if err != nil {
		utils.BadRequest(rw, req)
		return
	}
	newItem.ID = item.ID
	todoList[id] = newItem
	sendResponse(rw, todoList[id], http.StatusCreated)
}

func remove(id string, rw http.ResponseWriter, req *http.Request) {
	if id == "" {
		for key := range todoList {
			delete(todoList, key)
		}
		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte{})
	} else {
		item, found := todoList[id]
		if !found {
			http.NotFound(rw, req)
			return
		}
		delete(todoList, item.ID)
		rw.WriteHeader(http.StatusNoContent)
		rw.Write([]byte{})
	}
}

func unmarshalBody(req *http.Request) (*ToDoItem, error) {
	defer req.Body.Close()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	item := &ToDoItem{}
	if err := json.Unmarshal(data, item); err != nil {
		return nil, err
	}
	return item, nil
}

func sendResponse(rw http.ResponseWriter, message any, statusCode int) {
	data, err := json.Marshal(message)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	} else {
		rw.WriteHeader(statusCode)
		rw.Write(data)
	}
}
