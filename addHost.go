package gin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
)

func UpdateHost(cfg config.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string][]string
		cleaner := config.NewURIParser()
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		json.Unmarshal(body, &data)
		for _, j := range cfg.Endpoints {
			for k, i := range j.Backend {
				if k == 0 {
					for num, host := range i.Host {
						for _, zombie := range cleaner.CleanHosts(data["zombieIP"]) {
							if host == zombie {
								i.Host = append(i.Host[:num], i.Host[num+1:]...)
							}
						}
					}
					i.Host = append(i.Host, cleaner.CleanHosts(data["newIP"])...)
					fmt.Println(i.Host)
				}
			}

		}

	}

}
