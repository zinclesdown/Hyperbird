package main

import (
	"fmt"
	directAccess "hyperbird/core/fileaccess/direct-access"
	fakes3Access "hyperbird/core/fileaccess/fakes3-access"
	"hyperbird/servers"
)

const TEST_ON_LAUNCH = true

// 单元测试全部
func testAll() {
	if TEST_ON_LAUNCH {
		directAccess.Test()
		fakes3Access.Test()

		fmt.Println("运行服务器...")
		servers.PreTestServer()
		servers.Runservers()
		servers.TestServers()
	}
}
