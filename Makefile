
go-normalize:
	@go fmt ./...
	@go vet ./...

test: go-normalize
	@go test -p 1 ./...

proto:
	@protoc --go_out=. ./messages/hagallpb/hagall.proto
	@protoc --go_out=. ./messages//vikjapb/vikja.proto
	@protoc --go_out=. ./messages//odalpb/odal.proto
	@protoc --go_out=. ./messages//dagazpb/dagaz.proto

tag: check-version test
	@echo "\033[94m\n• Tagging ${VERSION}\033[00m"
	@git tag ${VERSION}
	@git push origin ${VERSION}

check-version:
	@echo "\033[94m\n• Checking Version\033[00m"
ifdef VERSION
	@echo "version set to $(VERSION)"
else
	@echo "\033[91mVERSION is not defined\033[00m"
	@echo "~> make VERSION=\033[90mv0.0.x\033[00m command"
	@exit 1
endif
