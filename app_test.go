package mediastreaminfo

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_JSON(t *testing.T) {
	// var info StreamInfo
	newStreamInfo := &StreamInfo{
		Id:            uuid.New().String(),
		URL:           "http://foo.com/index.m3u8",
		ABRStreamInfo: nil,
		Status:        "queued",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}
	
	data, err := json.Marshal(newStreamInfo)

	err = json.Unmarshal(data, &StreamInfo{})
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
}
