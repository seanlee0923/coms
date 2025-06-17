package protocol

type StatusResp struct {
	Ok bool `json:"ok"`
}

type FileUploadReq struct {
	FileName  string `json:"file_name"`
	ChunkSize int    `json:"chunk_size"`
}
