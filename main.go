package main

import (
	"time"

	"github.com/apavithraa/aws-cloud-utility/api"
	"github.com/apavithraa/aws-cloud-utility/util"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

func configureLogging(logLevel string) error {
	l, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(l)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableHTMLEscape: true,
	})
	return nil
}

func main() {
	settings := util.LoadAppConfig()
	configureLogging(settings.LogLevel)
	server := api.NewApiServer(settings)

	err := server.Start(settings.Port)
	if err != nil {
		log.Fatal(err)
	}

}
func contextMiddleware(regions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("regions", regions)
		c.Next()

	}
}
