package googlecast

// Ref: https://github.com/vishen/go-chromecast/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Payload interface {
	WithId(id int) Payload
}

type PayloadHeader struct {
	Type      string `json:"type"`
	RequestId int    `json:"requestId,omitempty"`
}

type SetVolumeRequest struct {
	PayloadHeader
	Volume `json:"volume"`
}

type LaunchAppRequest struct {
	PayloadHeader
	AppId string `json:"appId"`
}

type LoadMediaRequest struct {
	PayloadHeader
	Media       MediaItem `json:"media"`
	CurrentTime int       `json:"currentTime,omitempty"`
	Autoplay    bool      `json:"autoplay,omitempty"`
}

type LoadQueueRequest struct {
	PayloadHeader
	RepeatMode string          `json:"repeatMode"`
	Items      []LoadQueueItem `json:"items"`
}

type LoadQueueItem struct {
	Media            MediaItem `json:"media"`
	Autoplay         bool      `json:"autoplay"`
	PlaybackDuration uint      `json:"playbackDuration"`
}

type ReceiverStatusResponse struct {
	PayloadHeader
	Status struct {
		Apps   []App `json:"applications"`
		Volume `json:"volume"`
	} `json:"status"`
}

type MediaStatusResponse struct {
	PayloadHeader
	Status []Media `json:"status"`
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

func (this *PayloadHeader) WithId(id int) Payload {
	this.RequestId = id
	return this
}

func (this *SetVolumeRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *LaunchAppRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *LoadMediaRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}