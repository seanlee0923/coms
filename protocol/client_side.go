package protocol

type LoginReq struct {
	Id string `json:"id"`
}

type HeartBeatReq struct {
}

type StatusReq struct {
	CpuUsage   float64 `json:"cpu_usage"`
	MemUsage   float64 `json:"mem_usage"`
	DiskUsage  float64 `json:"disk_usage"`
	ServerTime string  `json:"server_time"`
}
