package kurento

// Media Profile.
// Currently WEBM and MP4 are supported.
type MediaProfileSpecType string

// Implement fmt.Stringer interface
func (t MediaProfileSpecType) String() string {
	return string(t)
}

const (
	MEDIAPROFILESPECTYPE_WEBM            MediaProfileSpecType = "WEBM"
	MEDIAPROFILESPECTYPE_MP4             MediaProfileSpecType = "MP4"
	MEDIAPROFILESPECTYPE_WEBM_VIDEO_ONLY MediaProfileSpecType = "WEBM_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_WEBM_AUDIO_ONLY MediaProfileSpecType = "WEBM_AUDIO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_VIDEO_ONLY  MediaProfileSpecType = "MP4_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_AUDIO_ONLY  MediaProfileSpecType = "MP4_AUDIO_ONLY"
)

type IceCandidate struct {
	Candidate     string
	SdpMid        string
	SdpMLineIndex int
}

type ServerInfo struct {
	Version      string
	Modules      []ModuleInfo
	Type         ServerType
	Capabilities []string
}

// Indicates if the server is a real media server or a proxy
type ServerType string

// Implement fmt.Stringer interface
func (t ServerType) String() string {
	return string(t)
}

const (
	SERVERTYPE_KMS ServerType = "KMS"
	SERVERTYPE_KCS ServerType = "KCS"
)

type ModuleInfo struct {
	Version   string
	Name      string
	Factories []string
}

// Type of media stream to be exchanged.
// Can take the values AUDIO, DATA or VIDEO.
type MediaType string

// Implement fmt.Stringer interface
func (t MediaType) String() string {
	return string(t)
}

const (
	MEDIATYPE_AUDIO MediaType = "AUDIO"
	MEDIATYPE_DATA  MediaType = "DATA"
	MEDIATYPE_VIDEO MediaType = "VIDEO"
)

// Type of filter to be created.
// Can take the values AUDIO, VIDEO or AUTODETECT.
type FilterType string

// Implement fmt.Stringer interface
func (t FilterType) String() string {
	return string(t)
}

const (
	FILTERTYPE_AUDIO      FilterType = "AUDIO"
	FILTERTYPE_AUTODETECT FilterType = "AUTODETECT"
	FILTERTYPE_VIDEO      FilterType = "VIDEO"
)

// Codec used for transmission of video.
type VideoCodec string

// Implement fmt.Stringer interface
func (t VideoCodec) String() string {
	return string(t)
}

const (
	VIDEOCODEC_VP8  VideoCodec = "VP8"
	VIDEOCODEC_H264 VideoCodec = "H264"
	VIDEOCODEC_RAW  VideoCodec = "RAW"
)

// Codec used for transmission of audio.
type AudioCodec string

// Implement fmt.Stringer interface
func (t AudioCodec) String() string {
	return string(t)
}

const (
	AUDIOCODEC_OPUS AudioCodec = "OPUS"
	AUDIOCODEC_PCMU AudioCodec = "PCMU"
	AUDIOCODEC_RAW  AudioCodec = "RAW"
)

type Fraction struct {
	Numerator   int
	Denominator int
}

type AudioCaps struct {
	Codec   AudioCodec
	Bitrate int
}

type VideoCaps struct {
	Codec     VideoCodec
	Framerate Fraction
}

type ElementConnectionData struct {
	Source            MediaElement
	Sink              MediaElement
	Type              MediaType
	SourceDescription string
	SinkDescription   string
}
