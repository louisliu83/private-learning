package types

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

type DatasetMetaInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"desc"`
	ValidDays   int32  `json:"validDays"`
	BizContext  string `json:"bizCtx"`
}

type ChunkMergeRequest struct {
	MD5        string `json:"md5"`
	ChunkSize  int64  `json:"chunkSize"`
	ChunkCount int32  `json:"chunkCount"`
	DatasetMetaInfo
}

type FilePullRequest struct {
	URL string `json:"url"`
	DatasetMetaInfo
}

type FileLocalCopyRequest struct {
	FilePath string `json:"path"`
	DatasetMetaInfo
}
