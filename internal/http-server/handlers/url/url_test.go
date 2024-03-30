package url_test

import (
	in_memory "github.com/DaDvoy/url-shortener-api.git/internal/storage/in-memory"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func getMap() map[string]string {
	return map[string]string{
		"http://google.com": "L7du_MycqX",
		"http://golang.org": "SoacZ1yCDu",
		"http://mysite.ru":  "cJQtkiv0zX",
	}
}

func fillStorage(storage *in_memory.Storage) {
	mp := getMap()
	for v := range mp {
		_ = storage.SaveURL(v, mp[v])
	}
}
