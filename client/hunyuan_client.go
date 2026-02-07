package client

import (
	"os"
	"sync"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

var (
	hunyuanClient     *hunyuan.Client
	hunyuanClientOnce sync.Once
)

// GetHunyuanClient 获取混元客户端单例
func GetHunyuanClient() *hunyuan.Client {
	hunyuanClientOnce.Do(func() {
		// 从环境变量读取密钥
		credential := common.NewCredential(
			os.Getenv("TENCENTCLOUD_SECRET_ID"),
			os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		)

		// 实例化一个client选项
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "hunyuan.tencentcloudapi.com"

		// 实例化要请求产品的client对象
		client, err := hunyuan.NewClient(credential, "", cpf)
		if err != nil {
			panic(err)
		}

		hunyuanClient = client
	})

	return hunyuanClient
}
