# Copyright (C) 2024 Intel Labs
# SPDX-License-Identifier: BSD-3-Clause

ARCH_LIBDIR ?= /lib/$(shell $(CC) -dumpmachine)

ifeq ($(DEBUG),1)
GRAMINE_LOG_LEVEL = debug
else
GRAMINE_LOG_LEVEL = error
endif

.PHONY: all
all: python3.manifest
ifeq ($(SGX),1)
all: python3.manifest.sgx python3.sig
endif

python3.manifest: python3.manifest.template
	gramine-manifest \
		-Dlog_level=$(GRAMINE_LOG_LEVEL) \
		-Darch_libdir=$(ARCH_LIBDIR) \
		-Dentrypoint=$(realpath $(shell sh -c "command -v python3")) \
		$< >$@

python3.manifest.sgx python3.sig: sgx_sign
	@:

.INTERMEDIATE: sgx_sign
sgx_sign: python3.manifest
	gramine-sgx-sign \
		--manifest $< \
		--output $<.sgx

.PHONY: clean
clean:
	$(RM) *.token *.sig *.manifest.sgx *.manifest *.dat
	$(RM) -r __pycache__

.PHONY: distclean
distclean: clean
