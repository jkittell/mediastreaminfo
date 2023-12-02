package mediastreaminfo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"github.com/jkittell/mediastreamdownloader/downloader"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
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

func getStreamInfo(content StreamInfo) (StreamInfo, error) {
	results := StreamInfo{
		Id:            content.Id,
		URL:           content.URL,
		ABRStreamInfo: array.New[ABRStreamInfo](),
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
				strInfo := ABRStreamInfo{
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

func getContentInfo(content StreamInfo) {
	debugMsg("%s start getting stream info", content.Id)
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
	debugMsg("%s done getting stream info", content.Id)
	debugMsg("%s elapsed time %s", content.Id, end.Sub(start).String())
}
