.PHONY: mockgen

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
	#-not -path "./network_planning/*" | sort

clean-deepcopy:
	@# 清除自动生成的deepcopy
	find . -type f -name "*zz_generated.deepcopy.go" \
	-not -path "./network_planning/*" | xargs rm -f

re-generate-deepcopy:
	@# 重新生成的deepcopy
	make clean-deepcopy
	make generate-deepcopy