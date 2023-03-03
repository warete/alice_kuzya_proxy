package aqara

type AqaraPayload struct {
	Intent string      `json:"intent"`
	Data   interface{} `json:"data"`
}

type ResourceHistoryItem struct {
	ResourceId string `json:"resourceId"`
	SubjectId  string `json:"subjectId"`
	TimeStamp  int64  `json:"timeStamp"`
	Value      string `json:"value"`
}

type ResourceHistoryResult struct {
	Data []ResourceHistoryItem `json:"data"`
}

type ResourceHistoryResponse struct {
	Result ResourceHistoryResult `json:"result"`
}

type KuzyaPayload struct {
	Value      string `json:"value"`
	SceneIdOn  string `json:"sceneIdOn"`
	SceneIdOff string `json:"sceneIdOff"`
	DeviceId   string `json:"deviceId"`
	ResourceId string `json:"resourceId"`
}
