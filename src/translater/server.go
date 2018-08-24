package main

import (
	"log"
	"net/http"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"response"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing"
	"translate"
	"errors"
	"strings"
)

var (
	cfg *Config
)

type Config struct {
	Debug      bool
	ListenPort string
	PID        string
	SecretKey  string
	Languages  map[string]string
	Accounts   []translate.Account
}

type InvalidConfig struct {
	file   string
	config string
}

func (e *InvalidConfig) Error() string {
	return fmt.Sprintf("%v", e.file)
}

// 载入配置文件
func loadConfig() (*Config, error) {
	cfg := &Config{
		Debug:      true,
		ListenPort: "80",
	}
	filePath := "src/config/conf.json"
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, &InvalidConfig{file: filePath}
	}

	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return nil, &InvalidConfig{file: filePath, config: string(jsonFile)}
	}

	return cfg, nil
}

func init() {
	if c, err := loadConfig(); err != nil {
		ae := err.(*InvalidConfig)
		panic("Config file read error:\nfile = " + ae.file + "\nconfig = " + ae.config)
	} else {
		cfg = c
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
		}
		toLang := c.Query("to", "zh-CHS")
		checkLanguages := []string{
			toLang,
		}
		if fromLang != "auto" {
			checkLanguages = append(checkLanguages, fromLang)
		}
		for _, v := range checkLanguages {
			if _, exists := cfg.Languages[v]; !exists {
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
			Debug:    cfg.Debug,
			From:     fromLang,
			To:       toLang,
			Accounts: cfg.Accounts,
		}
		translate := &translate.SogoTranslate{
			Translate: t,
		}
		translate.SetRawContent(text).Parse()

		if cfg.Debug {
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
	addr := cfg.ListenPort
	if len(addr) == 0 {
		addr = "8080"
	}
	http.ListenAndServe(":"+addr, nil)
}
