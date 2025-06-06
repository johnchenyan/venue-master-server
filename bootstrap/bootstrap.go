package bootstrap

import (
	"fmt"
	"github.com/deatil/lakego-doak/lakego/kernel"
)

// 添加服务提供者
func AddProvider(f func() any) {
	kernel.AddProvider(f)
}

// 执行
func Execute() {
	// 服务提供者
	providers := kernel.GetAllProvider()

	fmt.Printf("providers = %#v\n", len(providers))

	// 运行
	kernel.New().
		LoadDefaultServiceProvider().
		WithServiceProviders(providers).
		Terminate()
}
