# Simple GNU Makefile.

define hdr
	@echo ""
	@echo "# ================================================================"
	@echo "# $1"
	@echo "# ================================================================"
endef

all: bin/runat

test: help version test-after-sec test-before-sec test-after-ts test-before-ts

clean: ; rm -rf bin *~

bin/runat : runat.go
	$(call hdr,$@)
	GOBIN=$$(pwd)/bin go install $<

.PHONY: help
help: bin/runat
	$(call hdr,$@)
	bin/runat -h

.PHONY: version
version: bin/runat
	$(call hdr,$@)
	bin/runat --version

# Test run is after the current time mark.
# Time specification is in seconds mark.
.PHONY: test-after-sec
test-after-sec: bin/runat
	$(call hdr,$@)
	@ \
	Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	if (( Sec > 57 )) ; then \
		sleep 2 ; \
		Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	fi ; \
	echo "current mark is $$Sec" ; \
	(( Sec += 2 )) ; \
	echo "running at the $$Sec second mark in 2 seconds" ; \
	bin/runat -vv $$Sec /bin/bash -c "date && pwd"

# Test run is before the current time mark (wait a bit).
# Time specification is in seconds mark.
.PHONY: test-before-sec
test-before-sec: bin/runat
	$(call hdr,$@)
	@ \
	Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	if (( Sec < 1 )) ; then \
		sleep 1 ; \
		Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	fi ; \
	echo "current mark is $$Sec" ; \
	(( Diff = 60 - Sec )) ; \
	(( Sec = 0  )) ; \
	echo "running at the $$Sec second mark in $$Diff seconds" ; \
	bin/runat -vv $$Sec /bin/bash -c "date && pwd"

# Test run is after the current time mark.
# Time specification is a time stamp.
# Blast! The BSD date command does not have -s.
.PHONY: test-after-ts
test-after-ts: bin/runat
	$(call hdr,$@)
	@ \
	Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	if (( Sec > 57 )) ; then \
		sleep 2 ; \
		Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	fi ; \
	echo "current mark is $$Sec" ; \
	(( Sec += 2 )) ; \
	Min=$$(echo $$(date +'%M') | sed -e 's/^0//'); \
	Hr=$$(echo $$(date +'%H') | sed -e 's/^0//'); \
	Dts=$$(echo $$Hr $$Min $$Sec | awk '{printf("%02d:%02d:%02d",$$1,$$2,$$3)}') ; \
	echo "running at the $$Dts time stamp in 2 seconds" ; \
	bin/runat -vv $$Dts /bin/bash -c "date && pwd"

# Test run is before the current time mark (wait a bit).
# Time specification is a time stamp.
# Blast! The BSD date command does not have -s.
.PHONY: test-before-ts
test-before-ts: bin/runat
	$(call hdr,$@)
	@ \
	Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	if (( Sec < 1 )) ; then \
		sleep 1 ; \
		Sec=$$(echo $$(date +'%S') | sed -e 's/^0//'); \
	fi ; \
	echo "current second mark is $$Sec" ; \
	(( Diff = 60 - Sec )) ; \
	(( Sec = 0  )) ; \
	Min=$$(echo $$(date +'%M') | sed -e 's/^0//'); \
	Hr=$$(echo $$(date +'%H') | sed -e 's/^0//'); \
	(( Min++ )) ; \
	if (( Min > 59 )) ; then \
		(( Min = 0 )) ; \
		(( Hr++ )) ; \
		if (( Hr > 23 )) ; then \
			(( Hr = 0 )) ; \
		fi ; \
	fi ; \
	Dts=$$(echo $$Hr $$Min $$Sec | awk '{printf("%02d:%02d:%02d",$$1,$$2,$$3)}') ; \
	echo "running at the $$Dts time stamp in $$Diff seconds" ; \
	bin/runat -vv $$Dts /bin/bash -c "date && pwd"
