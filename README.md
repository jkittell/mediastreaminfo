# mediastreaminfo

Uses ffprobe to get information about a stream. 

Containerized application
https://hub.docker.com/r/jpkitt/mediastreaminfo

Send the url to the web service.
```
curl --location 'http://127.0.0.1:3000/contents' \
--header 'Content-Type: text/plain' \
--data '{ "url" : "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8" }'
```

Web service returns an id for stream info.
```
{
    "id": "b60b1534-9635-4270-aefb-4f30364ec69e",
    "url": "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8",
    "info": [],
    "status": "queued",
    "start_time": "2023-11-30T17:24:23.413907301Z",
    "end_time": "0001-01-01T00:00:00Z"
}
```

Get the info about the stream.
```
curl --location 'http://127.0.0.1:3000/contents/b60b1534-9635-4270-aefb-4f30364ec69e'
```
