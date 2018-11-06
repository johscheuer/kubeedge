
# make edge_core
.PHONY: default edge_core
edge_core:
	go build cmd/edge_core.go

# unit tests
.PHONY: test
ifeq ($(WHAT),)
       TEST_DIR="./pkg/"
else
       TEST_DIR=${WHAT}	
endif

test:
	find ${TEST_DIR} -name "*_test.go"|xargs -i dirname {}|uniq|xargs -i go test ${T} {}
