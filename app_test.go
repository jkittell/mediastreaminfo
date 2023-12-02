package mediastreaminfo

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jkittell/array"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// var info StreamInfo
	newStreamInfo := StreamInfo{
		Id:            uuid.New().String(),
		URL:           "http://foo.com/index.m3u8",
		ABRStreamInfo: nil,
		Status:        "queued",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}

	data, err := json.Marshal(&newStreamInfo)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	info := StreamInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Printf("%+v\n", info)
}

func mockGet(info StreamInfo) StreamInfo {
	updatedInfo := StreamInfo{
		Id:            info.Id,
		URL:           info.URL,
		ABRStreamInfo: array.New[ABRStreamInfo](),
		Status:        "completed",
		StartTime:     info.StartTime,
		EndTime:       time.Now().UTC(),
	}

	return updatedInfo
}

func Test2(t *testing.T) {
	newStreamInfo := StreamInfo{
		Id:            uuid.New().String(),
		URL:           "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8",
		ABRStreamInfo: nil,
		Status:        "started",
		StartTime:     time.Now().UTC(),
		EndTime:       time.Time{},
	}

	info := mockGet(newStreamInfo)

	b, err := json.Marshal(&info)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(string(b))
}
