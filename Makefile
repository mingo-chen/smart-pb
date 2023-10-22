.PHONY: test # 运行所有单元测试
test:
	go test -run ^Test github.com/mingo-chen/smart-pb -v

