package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SynologyWebHookRequest struct {
	Content string `json:"content"`
}

type LarkContent struct {
	Text string `json:"text"`
}

type LarkRequest struct {
	MsgType string      `json:"msg_type"`
	Content LarkContent `json:"content"`
}

type LarkSecretRequest struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	LarkRequest
}

var LarkWebhookUrl string
var LarkSecret string

func init() {
	LarkWebhookUrl = os.Getenv("LARK_WEBHOOK_URL")
	LarkSecret = os.Getenv("LARK_SECRET")

	_, err := url.ParseRequestURI(LarkWebhookUrl)
	if err != nil {
		log.Fatalf("LARK_WEBHOOK_URL '%s' is not a valid url", LarkWebhookUrl)
	}
	log.Printf("init lark webhook url with : %s", LarkWebhookUrl)
}

func GenSign(secret string, timestamp int64) (string, error) {
	//using timestamp + key to do sha256, then do base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func main() {
	r := gin.Default()
	r.POST("/lark", func(c *gin.Context) {
		synologyRequest := SynologyWebHookRequest{}
		data, err := c.GetRawData()
		jsonString := string(data)
		fmt.Printf("request data: %v\n", jsonString)
		// fix Synology test message fucking corrupted json format
		jsonString = strings.ReplaceAll(jsonString, "\n", " ")
		json.Unmarshal([]byte(jsonString), &synologyRequest)
		encrypted := len(LarkSecret) > 0

		larkRequest := LarkRequest{
			MsgType: "text",
			Content: LarkContent{
				Text: synologyRequest.Content,
			},
		}
		var larkSecretRequest LarkSecretRequest
		if encrypted {
			now := time.Now().Unix()
			sign, _ := GenSign(LarkSecret, now)
			larkSecretRequest = LarkSecretRequest{
				Timestamp:   strconv.FormatInt(now, 10),
				Sign:        sign,
				LarkRequest: larkRequest,
			}
		}
		var jsonBody []byte
		if encrypted {
			jsonBody, _ = json.Marshal(larkSecretRequest)

		} else {
			jsonBody, _ = json.Marshal(larkRequest)
		}

		resp, err := http.Post(LarkWebhookUrl,
			"application/json",
			bytes.NewBuffer(jsonBody))

		if resp.StatusCode != http.StatusOK {
			log.Panicf("faild to POST %s, status code is %d\n", LarkWebhookUrl, resp.StatusCode)
			c.JSON(resp.StatusCode, gin.H{
				"message": "faild to POST lark server\n",
			})
			return
		}

		if err != nil {
			log.Panicf("faild to POST %s, error is %s\n", LarkWebhookUrl, err)
			c.JSON(resp.StatusCode, gin.H{
				"message": "faild to POST lark server",
			})
			return
		}

		defer resp.Body.Close()
		larkBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("send content '%s' to lark server, response is %s\n", synologyRequest.Content, larkBody)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
