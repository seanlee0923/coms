package protocol

type LoginReq struct {
	Id string `json:"id"`
}

type HeartBeatReq struct {
}

type StatusReq struct {
	CpuUsage   float64 `json:"cpu_usage"`
	CpuCores   int     `json:"cpu_cores"`
	MemTotal   float64 `json:"mem_total"`
	MemUsage   float64 `json:"mem_usage"`
	DiskTotal  float64 `json:"disk_total"`
	DiskUsage  float64 `json:"disk_usage"`
	ServerTime string  `json:"server_time"`
}
