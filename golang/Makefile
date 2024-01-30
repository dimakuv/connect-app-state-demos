# Copyright (C) 2024 Intel Labs
# SPDX-License-Identifier: BSD-3-Clause

ifeq ($(DEBUG),1)
GRAMINE_LOG_LEVEL = all
else
GRAMINE_LOG_LEVEL = error
endif

.PHONY: all
all: main main.manifest
ifeq ($(SGX),1)
all: main.manifest.sgx main.sig
endif

main: main.go
	go build $<

main.manifest: main.manifest.template
	gramine-manifest \
		-Dlog_level=$(GRAMINE_LOG_LEVEL) \
		$< $@

main.sig main.manifest.sgx: sgx_sign
	@:

.INTERMEDIATE: sgx_sign
sgx_sign: main.manifest main
	gramine-sgx-sign \
		--manifest $< \
		--output $<.sgx

.PHONY: clean
clean:
	$(RM) *.token *.sig *.manifest.sgx *.manifest *.dat main

.PHONY: distclean
distclean: clean
