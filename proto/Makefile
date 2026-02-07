all: generate-all

# 生成所有proto文件
generate-all:
	trpc create \
		-p llm/llm.proto \
		-o llm \
		--rpconly \
		--nogomod \
		--mock=false

# 清理生成的代码
clean:
	rm -rf */*.pb.go
	rm -rf */*.trpc.go

.PHONY: all generate-all clean
