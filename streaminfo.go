package mediastreaminfo

import "gopkg.in/vansante/go-ffprobe.v2"

type StreamInfo struct {
	Name    string            `json:"name"`
	Ffprobe ffprobe.ProbeData `json:"ffprobe"`
}
