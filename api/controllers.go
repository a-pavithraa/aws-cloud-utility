package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/apavithraa/aws-cloud-utility/service"
	"github.com/apavithraa/aws-cloud-utility/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func NewApiServer(settings util.AppSettings) *Server {
	router := gin.Default()
	server := &Server{
		Router: router,
	}
	router.Use(contextMiddleware(settings.Regions))
	router.GET("/api/ec2/details", server.EC2InstanceDetailsHandler)
	router.GET("/api/ec2/start", server.EC2StartInstanceHandler)
	router.GET("/api/ec2/stop", server.EC2StopInstanceHandler)
	router.GET("/api/eip/details", server.EIPDetailsHandler)
	router.GET("/api/eip/release", server.EIPReleaseHandler)
	router.GET("/api/s3", server.S3Handler)
	router.GET("/api/lb/details", server.LBDetailsHandler)
	router.GET("/api/lb/delete", server.LBDeleteHandler)
	router.GET("/api/dd/details", server.DynamoDBHandler)
	return server

}
func contextMiddleware(regions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("regions", regions)
		c.Next()

	}
}
func (server *Server) Start(addr string) error {
	return server.Router.Run(addr)
}

func commonHandler(resourceType service.ResourceType, c *gin.Context) {
	var (
		response map[string][]service.RestResponse
		err      error
	)
	regions, _ := c.Get("regions")

	response, err = service.GetResourcesForAllRegions(resourceType, regions.([]string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

func (server *Server) EC2InstanceDetailsHandler(c *gin.Context) {

	commonHandler(service.LB, c)
}
func (server *Server) DynamoDBHandler(c *gin.Context) {

	commonHandler(service.DYNAMODB, c)
}
func (server *Server) EIPDetailsHandler(c *gin.Context) {
	commonHandler(service.EIP, c)
}

func (server *Server) LBDetailsHandler(c *gin.Context) {

	commonHandler(service.LB, c)
}

func (server *Server) EIPReleaseHandler(c *gin.Context) {

	allocationId := c.Query("allocationId")
	region := c.Query("region")
	fmt.Println("allocationId===", allocationId)
	if CheckEmptyString(allocationId) || CheckEmptyString(region) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both Allocation Id and Region must be non-empty",
		})
		return

	}

	err := service.ReleaseEIP(region, allocationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Released EIP",
	})
}

func (server *Server) LBDeleteHandler(c *gin.Context) {

	lbArn := c.Query("lbArn")
	region := c.Query("region")
	fmt.Println("lbArn===", lbArn)
	if CheckEmptyString(lbArn) || CheckEmptyString(region) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both Load Balancer ARN and Region must be non-empty",
		})
		return

	}

	err := service.DeleteLoadBalancer(region, lbArn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Deleted LB",
	})
}

func (server *Server) EC2StartInstanceHandler(c *gin.Context) {

	instanceId := c.Query("instanceId")
	region := c.Query("region")
	fmt.Println("instanceId===", instanceId)
	if CheckEmptyString(instanceId) || CheckEmptyString(region) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both Instance Id and Region must be non-empty",
		})
		return

	}
	err := service.StartInstance(region, instanceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Instance Started",
	})
}

func (server *Server) EC2StopInstanceHandler(c *gin.Context) {

	instanceId := c.Query("instanceId")
	region := c.Query("region")
	if CheckEmptyString(instanceId) || CheckEmptyString(region) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both Instance Id and Region must be non-empty",
		})
		return

	}
	fmt.Println("instanceId===", instanceId)
	err := service.StopInstance(region, instanceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Instance Stopped",
	})
}

func (server *Server) S3Handler(c *gin.Context) {
	regions, _ := c.Get("regions")
	fmt.Println("regions===", regions)
	response, err := service.ListAllBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func CheckEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
