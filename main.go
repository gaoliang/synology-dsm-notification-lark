package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

var LarkWebhookUrl string

func init() {
	LarkWebhookUrl = os.Getenv("LARK_WEBHOOK_URL")

	_, err := url.ParseRequestURI(LarkWebhookUrl)
	if err != nil {
		log.Fatalf("LARK_WEBHOOK_URL '%s' is not a valid url", LarkWebhookUrl)
	}
	log.Printf("init lark webhook url with : %s", LarkWebhookUrl)
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

		larkRequest := LarkRequest{
			MsgType: "text",
			Content: LarkContent{
				Text: synologyRequest.Content,
			},
		}
		jsonBody, _ := json.Marshal(larkRequest)

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
