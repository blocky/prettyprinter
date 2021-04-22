GO=go
GOTEST=$(GO) test -count=1
GOMOD=$(GO) mod
GOTIDY=$(GOMOD) tidy

MOCK=mockery

default: test

mock: # autogenerate mocks for interface testing
	$(MOCK) --all

test:
	$(GOTEST) ./... && $(GOTIDY)

