package api

type ChunkCheckRequest struct {
	MD5       string `json:"md5"`
	Chunk     int32  `json:"chunk"`
	ChunkSize int64  `json:"chunkSize"`
}

type ChunkUploadRequest struct {
	MD5      string `json:"md5"`
	Chunk    int32  `json:"chunk"`
	FileData []byte
}

type ChunkMergeRequest struct {
	MD5         string `json:"md5"`
	ChunkSize   int64  `json:"chunkSize"`
	ChunkCount  int32  `json:"chunkCount"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"desc"`
}

type FilePullRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"desc"`
	URL         string `json:"url"`
}
