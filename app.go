package mediastreaminfo

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"github.com/jkittell/toolbox"
	"log"
	"net/http"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type db struct {
	Database *array.Array[StreamInfo]
}

// contents array to store contents data.
var contents db

func getInstance() db {
	if contents.Database == nil {
		lock.Lock()
		defer lock.Unlock()
		if contents.Database == nil {
			contents.Database = array.New[StreamInfo]()
		}
	}
	return contents
}

func StartService(port int) error {
	router := gin.Default()
	router.GET("/contents", getContents)
	router.GET("/contents/:id", getContentByID)
	router.POST("/contents", postContents)

	return router.Run(fmt.Sprintf(":%d", port))
}

func Get(id string) StreamInfo {
	var info StreamInfo
	apiURL := fmt.Sprintf("http://127.0.0.1:3000/contents/%s", id)
	status, res, err := toolbox.SendRequest(toolbox.GET, apiURL, "", nil)
	if err != nil {
		log.Println(err)
		return info
	}

	if status != 200 {
		log.Printf("get response code: %d", status)
	}

	err = json.Unmarshal(res, &info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func GetAll() *array.Array[StreamInfo] {
	infos := array.New[StreamInfo]()
	apiURL := "http://127.0.0.1:3000/contents"
	status, res, err := toolbox.SendRequest(toolbox.GET, apiURL, "", nil)
	if err != nil {
		log.Println(err)
		return infos
	}

	if status != 200 {
		log.Printf("get response code: %d", status)
	}

	err = json.Unmarshal(res, &infos)
	if err != nil {
		log.Println(err)
	}
	return infos
}

func Post(url string) StreamInfo {
	var info StreamInfo
	apiURL := "http://127.0.0.1:3000/contents"

	data, _ := json.Marshal(map[string]string{"url": url})
	status, res, err := toolbox.SendRequest(toolbox.POST, apiURL, string(data), nil)
	if err != nil {
		log.Println(err)
		return info
	}

	if status != 201 {
		log.Printf("post response code: %d\n", status)
		return info
	}

	err = json.Unmarshal(res, &info)
	if err != nil {
		log.Println(err)
	}
	return info
}

// getContents responds with the list of all contents as JSON.
func getContents(c *gin.Context) {
	getInstance()
	c.JSON(http.StatusOK, contents.Database)
}

// postContents adds content from JSON received in the request body.
func postContents(c *gin.Context) {
	streamInfo := StreamInfo{
		Id:            uuid.New().String(),
		URL:           "",
		ABRStreamInfo: nil,
		Status:        "started",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}

	// Call BindJSON to bind the received JSON to
	// newStreamInfo.
	if err := c.BindJSON(&streamInfo); err != nil {
		log.Println(err)
		return
	}

	// Add the new content to the array.
	getInstance()
	contents.Database.Push(streamInfo)
	c.JSON(http.StatusCreated, streamInfo)

	go getContentInfo(streamInfo)
}

// getContentByID locates the content whose ID value matches the id
// parameter sent by the client, then returns that content as a response.
func getContentByID(c *gin.Context) {
	getInstance()
	id := c.Param("id")

	// Loop through the list of contents, looking for
	// content whose ID value matches the parameter.
	for i := 0; i < contents.Database.Length(); i++ {
		j := contents.Database.Lookup(i)
		if j.Id == id {
			c.JSON(http.StatusOK, j)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "content not found"})
}
