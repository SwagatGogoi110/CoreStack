package model

type TranscriptionJob struct {
	TranscriptionJobName   string      `json:"TranscriptionJobName"`
	TranscriptionJobStatus string      `json:"TranscriptionJobStatus"` // COMPLETED, FAILED
	LanguageCode           string      `json:"LanguageCode"`
	MediaFormat            string      `json:"MediaFormat"`
	Transcript             *Transcript `json:"Transcript,omitempty"`
	CreationTime           int64       `json:"CreationTime"`
}

type Transcript struct {
	TranscriptFileUri string `json:"TranscriptFileUri"`
}
