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
	"os"
	"slog"
)

var (
	logger        slog.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	config        *translate.Config
	v             *viper.Viper
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
	v.AddConfigPath("./config/")
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
	filename := "./runtime/logs/log.log"
	logFile := &os.File{}
	defer logFile.Close()
	exists := false
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsExist(err) {
			exists = true
		}
	} else {
		exists = true
	}
	if exists {
		logFile, err = os.OpenFile(filename, os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatalln(filename + " open failed.")
		}
	} else {
		logFile, err = os.Create(filename)
		if err != nil {
			log.Fatalln(filename + " create failed.")
		}
	}

	flag := log.LstdFlags | log.Lshortfile
	logger = slog.Logger{
		InfoLogger:    log.New(logFile, "[INFO] ", flag),
		WarningLogger: log.New(logFile, "[WARNING] ", flag),
		ErrorLogger:   log.New(logFile, "[ERROR] ", flag),
	}
	logger.InfoLogger.Println("Start server ...")
	router := routing.New()
	router.Use(
		access.Logger(logger.InfoLogger.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(logger.InfoLogger.Printf),
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
		fromLang := c.Query("from", "")
		toLang := c.Query("to", "")
		c.Request.ParseForm()
		text := strings.TrimSpace(c.Request.PostFormValue("text"))
		if len(text) == 0 {
			errors.New("`text` param is not allow empty.")
		}
		t := translate.Translate{
			Logger:   logger,
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
			t.Logger.InfoLogger.Println(text)
		}

		doc, err := translate.Do()
		if err == nil {
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
	log.Println("Begin start serve.")
	err = http.ListenAndServe(":"+addr, nil)
	if err != nil {
		logger.ErrorLogger.Panic(err)
	}
}
