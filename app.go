package mediastreaminfo

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"net/http"
	"sync"
	"time"
)

// ABRStreamInfo represents data about an ABR variant stream.
type ABRStreamInfo struct {
	Name    string            `json:"name"`
	Ffprobe ffprobe.ProbeData `json:"ffprobe"`
}

// StreamInfo represents data about a video/audio stream.
type StreamInfo struct {
	Id            string                      `json:"id"`
	URL           string                      `json:"url"`
	ABRStreamInfo *array.Array[ABRStreamInfo] `json:"abr_stream_info"`
	Status        string                      `json:"status"`
	StartTime     time.Time                   `json:"start_time"`
	EndTime       time.Time                   `json:"end_time"`
}

var lock = &sync.Mutex{}

type db struct {
	Database *array.Array[*StreamInfo]
}

// contents array to store contents data.
var contents db

func getInstance() db {
	if contents.Database == nil {
		lock.Lock()
		defer lock.Unlock()
		if contents.Database == nil {
			log.Println("creating new contents database")
			contents.Database = array.New[*StreamInfo]()
		}
	}
	return contents
}

func Start() {
	router := gin.Default()
	router.GET("/contents", getContents)
	router.GET("/contents/:id", getContentByID)
	router.POST("/contents", postContents)

	err := router.Run(":3000")
	if err != nil {
		log.Println(err)
	}
}

// getContents responds with the list of all contents as JSON.
func getContents(c *gin.Context) {
	getInstance()
	c.JSON(http.StatusOK, contents.Database)
}

// postContents adds content from JSON received in the request body.
func postContents(c *gin.Context) {
	newStreamInfo := &StreamInfo{
		Id:            uuid.New().String(),
		URL:           "",
		ABRStreamInfo: nil,
		Status:        "queued",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}

	// Call BindJSON to bind the received JSON to
	// newStreamInfo.
	if err := c.BindJSON(&newStreamInfo); err != nil {
		log.Println(err)
		return
	}

	// Add the new content to the array.
	getInstance()
	contents.Database.Push(newStreamInfo)
	c.JSON(http.StatusCreated, newStreamInfo)
	go getContentInfo(newStreamInfo)
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
