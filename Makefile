.PHONY: build
build:
	install -m 0755 -d build
	go build -o build/convert-md-to-jira ./convert-md-to-jira

.PHONY: clean
clean:
	go clean
	-rm -rf build/convert-md-to-jira

.PHONY: install
install:
	install -m 0755 build/convert-md-to-jira ~/bin
