ADDLICENSE ?= go run -modfile hack/tools/go.mod github.com/google/addlicense


.PHONY: copyright
copyright:
	$(ADDLICENSE) -f ./hack/copyright.txt .

check:
	hack/check/check-mdlint.sh
	hack/check/check-yaml.sh

test-e2e:
	./tests/e2e/01-test-basic.sh
