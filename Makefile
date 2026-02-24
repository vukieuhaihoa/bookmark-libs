.PHONY: clean
clean:
	go clean -testcache
	rm -rf $(COVERAGE_FOLDER)

COVERAGE_EXCLUDE=sqldb|redis|mocks|vendor|test|docs|main.go|config.go|client.go
COVERAGE_THRESHOLD = 80
COVERAGE_FOLDER=./coverage

#=========================== DEV TOOLS =========================== 
.PHONY: mock-gen
mock-gen:
	go generate ./...

.PHONY: test 
test: clean
	mkdir -p $(COVERAGE_FOLDER)
	go test ./... -coverprofile=$(COVERAGE_FOLDER)/coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -v -E "$(COVERAGE_EXCLUDE)" $(COVERAGE_FOLDER)/coverage.tmp > $(COVERAGE_FOLDER)/coverage.out
	go tool cover -html=$(COVERAGE_FOLDER)/coverage.out -o $(COVERAGE_FOLDER)/coverage.html
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi

.PHONY: docker-test
docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark_service:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi	

