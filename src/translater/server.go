package main

import (
	"log"
	"net/http"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"fmt"
	"response"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing"
	"translate"
	"errors"
	"strings"
	"github.com/spf13/viper"
)

var (
	config *translate.Config
	v      *viper.Viper
)

type InvalidConfig struct {
	file   string
	config string
}

func (e *InvalidConfig) Error() string {
	return fmt.Sprintf("%v", e.file)
}

func init() {
	v = viper.New()
	v.AddConfigPath("./src/config/")
	v.SetConfigName("conf")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	err = v.Unmarshal(&config)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	router := routing.New()
	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	api := router.Group("/api")
	api.Use(
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	// GET /api/ping
	api.Get("/ping", func(c *routing.Context) error {
		return c.Write("OK")
	})

	// 翻译文本
	// POST /api/translate
	api.Post("/translate", func(c *routing.Context) error {
		success := true
		errorMessage := ""
		fromLang := c.Query("from", "auto")
		if len(fromLang) == 0 {
			fromLang = "auto"
		} else {
			fromLang = strings.ToLower(fromLang)
		}
		toLang := c.Query("to", "zh-CHS")
		//toLang = strings.ToLower(toLang)
		checkLanguages := []string{
			strings.ToLower(toLang),
		}
		if fromLang != "auto" {
			checkLanguages = append(checkLanguages, strings.ToLower(fromLang))
		}
		for _, v := range checkLanguages {
			if _, exists := config.Languages[v]; !exists {
				success = false
				errorMessage = fmt.Sprintf("Not Support `%v` language.", v)
				//log.Panic(fmt.Sprintf("Not Support `%v` language.", v))
				break
			}
		}
		if !success {
			return c.Write(&response.FailResponse{
				Success: false,
				Error: response.Error{
					Message: errorMessage,
				},
			})
		}

		c.Request.ParseForm()
		text := c.Request.PostFormValue("text")
		if len(strings.TrimSpace(text)) == 0 {
			errors.New("`text` param is not allow empty.")
		}
		t := translate.Translate{
			Viper:    v,
			Config:   config,
			From:     fromLang,
			To:       toLang,
			Accounts: config.Accounts,
		}
		translate := &translate.SogoTranslate{
			Translate: t,
		}
		translate.SetRawContent(text).Parse()

		if config.Debug {
			log.Println(text)
		}

		doc, err := translate.Do()
		if err == nil {
			success = true
			resp := &response.SuccessResponse{
				Success: true,
				Data: response.SuccessData{
					RawContent: translate.GetRawContent(),
					Content:    doc.Render(),
				},
			}

			return c.Write(resp)
		} else {
			error := &response.Error{
				Message: fmt.Errorf("%v", err).Error(),
			}
			resp := &response.FailResponse{
				Success: false,
				Error:   *error,
			}

			return c.Write(resp)
		}
	})

	http.Handle("/", router)
	addr := config.ListenPort
	if len(addr) == 0 {
		addr = "8080"
	}
	http.ListenAndServe(":"+addr, nil)
}
