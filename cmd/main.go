package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"github.com/jkittell/mediastreamdownloader/downloader"
	"github.com/jkittell/mediastreaminfo"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type db struct {
	Database *array.Array[mediastreaminfo.StreamInfo]
}

// contents array to store contents data.
var contents db

func getInstance() db {
	if contents.Database == nil {
		lock.Lock()
		defer lock.Unlock()
		if contents.Database == nil {
			contents.Database = array.New[mediastreaminfo.StreamInfo]()
		}
	}
	return contents
}

func main() {
	router := gin.Default()
	router.GET("/contents", getContents)
	router.GET("/contents/:id", getContentByID)
	router.POST("/contents", postContents)

	err := router.Run("127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
}

// probe mp4 file with ffprobe
func probe(path string) *ffprobe.ProbeData {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	fileReader, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}

	data, err := ffprobe.ProbeReader(ctx, fileReader)
	if err != nil {
		log.Printf("Error getting data: %v", err)
	}

	return data
}

func getStreamInfo(content mediastreaminfo.StreamInfo) (mediastreaminfo.StreamInfo, error) {
	results := mediastreaminfo.StreamInfo{
		Id:            content.Id,
		URL:           content.URL,
		ABRStreamInfo: array.New[mediastreaminfo.ABRStreamInfo](),
		Status:        content.Status,
		StartTime:     content.StartTime,
		EndTime:       content.EndTime,
	}

	// verify ffprobe is available
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		log.Println("ffprobe is not available", err)
		results.Status = "error"
		results.EndTime = time.Now().UTC()
		return results, err
	} else if strings.HasSuffix(content.URL, ".mpd") {
		log.Println("skip ffprobe for dash", nil)
		results.Status = "skipped"
		results.EndTime = time.Now().UTC()
		return results, nil
	} else {
		// specify directory to download segments
		dir := path.Join("/tmp", uuid.New().String())
		// download segments to mp4
		streams := downloader.Run(dir, content.URL)
		if streams.Length() > 0 {
			// run ffprobe on each stream
			for i := 0; i < streams.Length(); i++ {
				str := streams.Lookup(i)
				info := probe(str.File)
				strInfo := mediastreaminfo.ABRStreamInfo{
					Name:    str.Name,
					Ffprobe: *info,
				}
				results.ABRStreamInfo.Push(strInfo)
			}
			results.Status = "completed"
		} else {
			log.Println("no streams found", nil)
			results.Status = "error"
			return results, nil
		}
		err = os.RemoveAll(dir)
		if err != nil {
			log.Println("unable to remove downloaded segments")
		}
		results.EndTime = time.Now().UTC()
	}
	return results, nil
}

func getContentInfo(content mediastreaminfo.StreamInfo) {
	log.Printf("%s start getting stream info", content.Id)
	start := time.Now()
	streamInfo, err := getStreamInfo(content)
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < contents.Database.Length(); i++ {
		info := contents.Database.Lookup(i)
		if info.Id == streamInfo.Id {
			contents.Database.Set(i, streamInfo)
		}
	}

	end := time.Now()
	log.Printf("%s done getting stream info", content.Id)
	log.Printf("%s elapsed time %s", content.Id, end.Sub(start).String())
}

// getContents responds with the list of all contents as JSON.
func getContents(c *gin.Context) {
	getInstance()
	c.JSON(http.StatusOK, contents.Database)
}

// postContents adds content from JSON received in the request body.
func postContents(c *gin.Context) {
	streamInfo := mediastreaminfo.StreamInfo{
		Id:            uuid.New().String(),
		URL:           "",
		ABRStreamInfo: nil,
		Status:        "started",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}

	// Call BindJSON to bind the received JSON to stream info.
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
