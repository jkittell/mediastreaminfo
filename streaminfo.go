package mediastreaminfo

import (
	"encoding/json"
	"fmt"
	"github.com/jkittell/array"
	"github.com/jkittell/toolbox"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
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
