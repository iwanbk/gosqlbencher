PKGS = $(shell go list ./...)

lint:
	golint $(PKGS)
