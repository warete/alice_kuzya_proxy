package aqara

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/warete/alice_kuzya_proxy/pkg/serviceproxy"
)

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

type IAqara interface {
	serviceproxy.IServiceProxy
	sendRequest(method string, params url.Values, body AqaraPayload) (string, error)
	ExecScene(sceneId string) error
	GetResourceHistory(deviceId, resourceId string, startTime, endTime int64) ([]ResourceHistoryItem, error)
}

type AqaraImpl struct {
	config AqaraConfig
}

func NewAqara(config AqaraConfig) (IAqara, error) {
	return AqaraImpl{
		config: config,
	}, nil
}

func (a AqaraImpl) AddRoutes(r *gin.Engine) {
	r.GET("/aqara", func(c *gin.Context) {
		res, err := a.GetResourceHistory("lumi.54ef44100050f25e", "4.1.85", time.Now().Add(-24*time.Hour).Unix()*1000, time.Now().Unix()*1000)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, res[0].Value)
	})
}

func (a AqaraImpl) sendRequest(method string, params url.Values, body AqaraPayload) (string, error) {
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	request, err := http.NewRequest(method, a.config.ApiUrl+"?"+params.Encode(), bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", err
	}

	headerTime := strconv.Itoa(int(time.Now().Unix() * 1000))
	headerNonce := strconv.Itoa(int(time.Now().Unix()) * 1000)

	var sign string

	sParams := "Accesstoken=" + a.config.AccessToken + "&" + "Appid=" + a.config.AppId + "&" + "Keyid=" + a.config.KeyId + "&" + "Nonce=" + headerNonce + "&" + "Time=" + headerTime + a.config.AppKey
	hash := md5.Sum([]byte(strings.ToLower(sParams)))
	sign = hex.EncodeToString(hash[:])

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Appid", a.config.AppId)
	request.Header.Set("AccessToken", a.config.AccessToken)
	request.Header.Set("Keyid", a.config.KeyId)
	request.Header.Set("Time", headerTime)
	request.Header.Set("Nonce", headerNonce)
	request.Header.Set("Sign", sign)

	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", err
	}
	defer res.Body.Close()

	respData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(respData), nil
}

func (a AqaraImpl) ExecScene(sceneId string) error {
	_, err := a.sendRequest("POST", url.Values{}, AqaraPayload{
		Intent: "config.scene.run",
		Data: map[string]string{
			"sceneId": sceneId,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (a AqaraImpl) GetResourceHistory(deviceId, resourceId string, startTime, endTime int64) ([]ResourceHistoryItem, error) {
	res, err := a.sendRequest("POST", url.Values{}, AqaraPayload{
		Intent: "fetch.resource.history",
		Data: map[string]interface{}{
			"subjectId":   deviceId,
			"resourceIds": []string{resourceId},
			"startTime":   startTime,
			"endTime":     endTime,
		},
	})
	if err != nil {
		return nil, err
	}
	var resObj ResourceHistoryResponse
	if err = json.Unmarshal([]byte(res), &resObj); err != nil {
		return nil, err
	}
	return resObj.Result.Data, nil
}
