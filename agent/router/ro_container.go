package router

import (
	v2 "github.com/1Panel-dev/1Panel/agent/app/api/v2"
	"github.com/gin-gonic/gin"
)

type ContainerRouter struct{}

func (s *ContainerRouter) InitRouter(Router *gin.RouterGroup) {
	byRouter := Router.Group("containers")
	baseApi := v2.ApiGroupApp.BaseApi
	{
		byRouter.GET("/exec", baseApi.ContainerWsSSH)
		byRouter.GET("/stats/:id", baseApi.ContainerStats)

		byRouter.POST("", baseApi.ContainerCreate)
		byRouter.POST("/update", baseApi.ContainerUpdate)
		byRouter.POST("/upgrade", baseApi.ContainerUpgrade)
		byRouter.POST("/info", baseApi.ContainerInfo)
		byRouter.POST("/search", baseApi.SearchContainer)
		byRouter.POST("/list", baseApi.ListContainer)
		byRouter.POST("/list/byimage", baseApi.ListContainerByImage)
		byRouter.GET("/status", baseApi.LoadContainerStatus)
		byRouter.GET("/list/stats", baseApi.ContainerListStats)
		byRouter.POST("/item/stats", baseApi.ContainerItemStats)
		byRouter.GET("/search/log", baseApi.ContainerStreamLogs)
		byRouter.POST("/download/log", baseApi.DownloadContainerLogs)
		byRouter.GET("/limit", baseApi.LoadResourceLimit)
		byRouter.POST("/clean/log", baseApi.CleanContainerLog)
		byRouter.POST("/inspect", baseApi.Inspect)
		byRouter.POST("/rename", baseApi.ContainerRename)
		byRouter.POST("/commit", baseApi.ContainerCommit)
		byRouter.POST("/operate", baseApi.ContainerOperation)
		byRouter.POST("/prune", baseApi.ContainerPrune)

		byRouter.POST("/users", baseApi.LoadContainerUsers)
		byRouter.POST("/files/search", baseApi.ListContainerFiles)
		byRouter.POST("/files/upload", baseApi.UploadContainerFile)
		byRouter.POST("/files/content", baseApi.GetContainerFileContent)
		byRouter.POST("/files/size", baseApi.GetContainerFileSize)
		byRouter.POST("/files/del", baseApi.DeleteContainerFile)
		byRouter.POST("/files/download", baseApi.DownloadContainerFile)

		byRouter.GET("/repo", baseApi.ListRepo)
		byRouter.POST("/repo/status", baseApi.CheckRepoStatus)
		byRouter.POST("/repo/search", baseApi.SearchRepo)
		byRouter.POST("/repo/update", baseApi.UpdateRepo)
		byRouter.POST("/repo", baseApi.CreateRepo)
		byRouter.POST("/repo/del", baseApi.DeleteRepo)

		byRouter.POST("/compose/search", baseApi.SearchCompose)
		byRouter.POST("/compose", baseApi.CreateCompose)
		byRouter.POST("/compose/env", baseApi.LoadComposeEnv)
		byRouter.POST("/compose/test", baseApi.TestCompose)
		byRouter.POST("/compose/operate", baseApi.OperatorCompose)
		byRouter.POST("/compose/clean/log", baseApi.CleanComposeLog)
		byRouter.POST("/compose/update", baseApi.ComposeUpdate)

		byRouter.GET("/template", baseApi.ListComposeTemplate)
		byRouter.POST("/template/search", baseApi.SearchComposeTemplate)
		byRouter.POST("/template/update", baseApi.UpdateComposeTemplate)
		byRouter.POST("/template/batch", baseApi.BatchComposeTemplate)
		byRouter.POST("/template", baseApi.CreateComposeTemplate)
		byRouter.POST("/template/del", baseApi.DeleteComposeTemplate)

		byRouter.GET("/image", baseApi.ListImage)
		byRouter.GET("/image/all", baseApi.ListAllImage)
		byRouter.POST("/image/search", baseApi.SearchImage)
		byRouter.POST("/image/pull", baseApi.ImagePull)
		byRouter.POST("/image/push", baseApi.ImagePush)
		byRouter.POST("/image/save", baseApi.ImageSave)
		byRouter.POST("/image/load", baseApi.ImageLoad)
		byRouter.POST("/image/remove", baseApi.ImageRemove)
		byRouter.POST("/image/tag", baseApi.ImageTag)
		byRouter.POST("/image/build", baseApi.ImageBuild)

		byRouter.GET("/network", baseApi.ListNetwork)
		byRouter.POST("/network/del", baseApi.DeleteNetwork)
		byRouter.POST("/network/search", baseApi.SearchNetwork)
		byRouter.POST("/network", baseApi.CreateNetwork)
		byRouter.GET("/volume", baseApi.ListVolume)
		byRouter.POST("/volume/del", baseApi.DeleteVolume)
		byRouter.POST("/volume/search", baseApi.SearchVolume)
		byRouter.POST("/volume", baseApi.CreateVolume)

		byRouter.GET("/daemonjson", baseApi.LoadDaemonJson)
		byRouter.GET("/daemonjson/file", baseApi.LoadDaemonJsonFile)
		byRouter.GET("/docker/status", baseApi.LoadDockerStatus)
		byRouter.POST("/docker/operate", baseApi.OperateDocker)
		byRouter.POST("/daemonjson/update", baseApi.UpdateDaemonJson)
		byRouter.POST("/logoption/update", baseApi.UpdateLogOption)
		byRouter.POST("/ipv6option/update", baseApi.UpdateIpv6Option)
		byRouter.POST("/daemonjson/update/byfile", baseApi.UpdateDaemonJsonByFile)
	}
}
