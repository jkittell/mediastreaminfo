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

func getContentInfo(content *content) {
	results := array.New[StreamInfo]()
	// verify ffprobe is available
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		log.Println("ffprobe is not available", err)
		content.Status = "error"
		content.EndTime = time.Now().UTC()
		return
	} else if strings.HasSuffix(content.URL, ".mpd") {
		log.Println("skip ffprobe for dash", nil)
		content.Status = "skipped"
		content.EndTime = time.Now().UTC()
		return
	} else {
		content.Status = "processing"
		// specify directory to download segments
		dir := path.Join("/tmp", uuid.New().String())
		// download segments to mp4
		streams := downloader.Run(dir, content.URL)
		if streams.Length() > 0 {
			// run ffprobe on each stream
			for i := 0; i < streams.Length(); i++ {
				str := streams.Lookup(i)
				info := probe(str.File)
				//Name: str.Name,
				//Info: *data,
				strInfo := StreamInfo{
					Name:    str.Name,
					Ffprobe: *info,
				}
				results.Push(strInfo)
			}

			content.StreamInfo = results
			content.Status = "completed"
		} else {
			log.Println("no streams found", nil)
			content.Status = "error"
		}
		err = os.RemoveAll(dir)
		if err != nil {
			log.Println("unable to remove downloaded segments")
		}
		content.EndTime = time.Now().UTC()
	}
}
