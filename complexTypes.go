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

type ElementStats struct {
	InputAudioLatency float64
	InputVideoLatency float64

	AudioE2ELatency []MediaLatencyStat
	VideoE2ELatency []MediaLatencyStat
	Timestamp       uint64
	Type            MediaType
	Id              string
	// RTCPeerConnection Stats
	DataChannelsClosed int
	DataChannelsOpened int
	// RTCCertificate stats
	Base64Certificate    string
	Fingerprint          string
	FingerprintAlgorithm string
	IssuerCertificateId  string
	// RTCCodec
	Channels    int
	Clockrate   uint64
	Codec       string
	Parametesrs string
	PayloadType string
	// RTC DatachannelStats
	BytesReceived    int
	BytesSent        int
	DatachannelId    int
	Label            string
	MessagesReceived int // Represents the total number of API 'message' events received.
	MessagesSent     int // Represents the total number of API 'message' events sent.
	Protocol         string
	// State RTCDataChannelState Not sure how to do enums

	// RTCIceCandidateAttributes
	AddressSourceUrl string // The URL of the TURN or STUN server indicated in the RTCIceServers that translated this IP address.
	//CandidateType RTCStatsIceCandidateType //The enumeration RTCStatsIceCandidateType is based on the cand-type defined in [RFC5245] section 15.1.
	IpAddress  string
	PortNumber int
	Priority   uint64
	Transport  string

	//RTCIceCandidatePairStats
	AvailableIncomingBitrate float64 //Measured in Bits per second, and is implementation dependent.
	AvailableOutgoingBitrate float64 //Measured in Bits per second, and is implementation dependent.
	// Duplicate BytesReceived            uint64  // Represents the total number of payload bytes received on this candidate pair, i.e., not including headers or padding.
	// Duplicate BytesSent                uint64  //Represents the total number of payload bytes sent on this candidate pair, i.e., not including headers or padding.
	LocalCandidateId         string  //It is a unique identifier that is associated to the object that was inspected to produce the RTCIceCandidateAttributes for the local candidate associated with this candidate pair.
	Nominated                bool    //Related to updating the nominated flag described in Section 7.1.3.2.4 of [RFC5245].
	// Duplicate Priority                 uint64  //Calculated from candidate priorities as defined in [RFC5245] section 5.7.2.
	Readable                 bool    //Has gotten a valid incoming ICE request.
	RemoteCandidateId        string  //It is a unique identifier that is associated to the object that was inspected to produce the RTCIceCandidateAttributes for the remote candidate associated with this candidate pair.
	RoundTripTime            float64 // Represents the RTT computed by the STUN connectivity checks
	//State RTCStatsIceCandidatePairState //Represents the state of the checklist for the local and remote candidates in a pair.
	TransportId string //It is a unique identifier that is associated to the object that was inspected to produce the RTCTransportStats associated with this candidate pair.
	Writable    bool   // Has gotten ACK to an ICE request.

	//RTCMediaStreamStats
	StreamIdentifier string //Stream identifier.
	TrackIds         []string

	//RTCMediaStreamTrackStats
	AudioLevel                   float64  //Only valid for audio, and the value is between 0..1 (linear), where 1.0 represents 0 dBov.
	getEchoReturnLoss            float64  //Only present on audio tracks sourced from a microphone where echo cancellation is applied.
	getEchoReturnLossEnhancement float64  //Only present on audio tracks sourced from a microphone where echo cancellation is applied.
	getFrameHeight               int      //Only makes sense for video media streams and represents the height of the video frame for this SSRC.
	FramesCorrupted              uint64   // Only valid for video.
	FramesDecoded                uint64   // Only valid for video.
	FramesDropped                uint64   // Only valid for video.
	FramesPerSecond              float64  // Only valid for video.
	FramesReceived               uint64   // Only valid for video and when remoteSource is set to true.
	FramesSent                   uint64   // Only valid for video.
	FrameWidth                   uint64   // Only makes sense for video media streams and represents the width of the video frame for this SSRC.
	RemoteSource                 bool     // true indicates that this is a remote source.
	SsrcIds                      []string // Synchronized sources.
	TrackIdentifier              string   //Represents the track.id property.

	//RTCRTPStreamStats
	AssociateStatsId string  //The associateStatsId is used for looking up the corresponding (local/remote) RTCStats object for a given SSRC.
	CodecId          string  // The codec identifier
	FirCount         uint64  // Count the total number of Full Intra Request (FIR) packets received by the sender.
	FractionLost     float64 // The fraction packet loss reported for this SSRC.
	IsRemote         bool    // false indicates that the statistics are measured locally, while true indicates that the measurements were done at the remote endpoint and reported in an RTCP RR/XR.
	MediaTrackId     string  // Track identifier.
	NackCount        uint64  // Count the total number of Negative ACKnowledgement (NACK) packets received by the sender and is sent by receiver.
	PacketsLost      uint64  // Total number of RTP packets lost for this SSRC.
	PliCount         uint64  // Count the total number of Packet Loss Indication (PLI) packets received by the sender and is sent by receiver.
	Remb             uint64  // The Receiver Estimated Maximum Bitrate (REMB).
	SliCount         uint64  // Count the total number of Slice Loss Indication (SLI) packets received by the sender.
	Ssrc             string  // The synchronized source SSRC
	// Duplicate TransportId      string  // It is a unique identifier that is associated to the object that was inspected to produce the RTCTransportStats associated with this RTP stream.
	// 
	// // RTCInboundRTPStreamStats
	Jitter float64 // Packet Jitter measured in seconds for this SSRC.
	PacketsReceived uint64 //Total number of RTP packets received for this SSRC.
	// Duplicate BytesReceived uint64 //Total number of bytes received for this SSRC.
	// 
	// // RTCOutboundRTPStreamStats
	// Duplicate Jitter float64 // Packet Jitter measured in seconds for this SSRC.
	PacketsSent uint64 //Total number of RTP packets received for this SSRC.
	// Duplicate BytesSent uint64 //Total number of bytes received for this SSRC.

	//RTCTransportStats
	ActiveConnection        bool   //Set to true when transport is active.
	// Duplicate BytesReceived           uint64 // Represents the total number of bytes received on this PeerConnection, i.e., not including headers or padding.
	// Duplicate BytesSent               uint64 // Represents the total number of payload bytes sent on this PeerConnection, i.e., not including headers or padding.
	LocalCertificateId      string // For components where DTLS is negotiated, give local certificate.
	RemoteCertificateId     string // For components where DTLS is negotiated, give remote certificate.
	RtcpTransportStatsId    string // If RTP and RTCP are not multiplexed, this is the id of the transport that gives stats for the RTCP component, and this record has only the RTP component stats.
	SelectedCandidatePairId string // It is a unique identifier that is associated to the object that was inspected to produce the RTCIceCandidatePairStats associated with this transport.
}

type MediaLatencyStat struct {
	Avg  float64
	Name string
	Type MediaType
}
