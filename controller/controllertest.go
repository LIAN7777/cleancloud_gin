package controller

import (
	"GinProject/dao"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Controllertest struct {
}

func New() Controllertest {
	return Controllertest{}
}

func (con Controllertest) Test(c *gin.Context) {
	c.String(http.StatusOK, "start success !")
}

func (con Controllertest) GetCity(c *gin.Context) {
	cityname := c.Param("city")
	city, err := dao.GetCityByName(cityname)
	if err != nil {
		c.String(http.StatusOK, "get error")
	} else {
		c.JSON(http.StatusOK, city)
	}
}

func (con Controllertest) GetUserById(c *gin.Context) {
	userId := c.Param("userId")
	userIdInt, err := strconv.ParseInt(userId, 10, 64) // 将 userId 转换为 int64 类型
	if err != nil {
		c.String(http.StatusOK, "Invalid user ID")
		return
	}
	user, err := dao.GetUserById(userIdInt)
	if err != nil {
		c.String(http.StatusOK, "get user error")
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func (con Controllertest) JwtAuthTest(c *gin.Context) {
	tokenStr, _ := c.Get("newTokenString")
	c.JSON(http.StatusOK, gin.H{
		"status":   "pass token auth !",
		"newToken": tokenStr,
	})
}

func (con Controllertest) FileTranTest(c *gin.Context) {
	filePath := c.PostForm("path")
	fileName := "file_to_download"
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusOK, "file open err")
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			c.String(http.StatusOK, "file close err")
			return
		}
	}(file)
	stat, err := file.Stat()
	if err != nil {
		c.String(http.StatusOK, "get file info err")
		return
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	c.Writer.Flush()
	var offset int64 = 0
	var bufsize int64 = 1024 * 1024 // 1MB
	buf := make([]byte, bufsize)
	for {
		n, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			log.Println("read file error", err)
			break
		}
		if n == 0 {
			break
		}
		_, err = c.Writer.Write(buf[:n])
		if err != nil {
			log.Println("write file error", err)
			break
		}
		offset += int64(n)
	}
	c.Writer.Flush()
}
