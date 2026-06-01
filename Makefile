.PHONY: mockgen clean-deepcopy re-generate-deepcopy build-push deploy

all:

# 配合go:generate, 自动生成mock类代码
mockgen:
	@echo "Generating..."
	@go generate -run="mockgen" -x ./...

check-deepcopy-gen:
	@if ! command -v deepcopy-gen > /dev/null 2>&1; then \
		echo "deepcopy-gen does not exist, installing..."; \
		go install k8s.io/code-generator/cmd/deepcopy-gen@v0.32.3; \
	fi

generate-deepcopy:
	@echo "Generating DeepCopy methods..."
	@# 检查 deepcopy-gen
	make check-deepcopy-gen

	deepcopy-gen \
	 -v 10 \
	 --output-file zz_generated.deepcopy.go \
	 --go-header-file ./boilerplate.txt \
	 ./models/... \
	|| { echo "deepcopy-gen failed"; exit 1; }

	@echo "DeepCopy generation completed!"

	@# 显示生成的文件
	@echo "Generated files:"
	find . -type f -name "*zz_generated.deepcopy.go" | sort
	@#-not -path "./xxx/*" | sort

clean-deepcopy:
	@# 清除自动生成的deepcopy
	find . -type f -name "*zz_generated.deepcopy.go" \
	-not -path "./network_planning/*" | xargs rm -f

re-generate-deepcopy:
	@# 重新生成的deepcopy
	make clean-deepcopy
	make generate-deepcopy

# Docker 部署到NAS
# 构建并推送镜像到Registry
build-push:
	@echo "构建并推送镜像..."
	chmod +x build-and-push.sh
	./build-and-push.sh

# 部署到NAS：上传脚本并提示执行（镜像已存在）
deploy:
	@echo "上传部署脚本到NAS..."
	scp deploy-on-nas.sh mellolo@192.168.5.178:~/deploy-on-nas.sh
	@echo ""
	@echo "======================================"
	@echo "部署脚本已上传"
	@echo "======================================"
	@echo ""
	@echo "请在NAS上执行部署："
	@echo "  ssh mellolo@192.168.5.178"
	@echo "  sudo ./deploy-on-nas.sh"
	@echo ""
	@echo "完成后访问: http://192.168.5.178:18080"
	@echo ""