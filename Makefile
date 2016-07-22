# Simple GNU Makefile.

define hdr
	@echo ""
	@echo "# ================================================================"
	@echo "# $1"
	@echo "# ================================================================"
endef

all: bin/runat

test: help test-after-sec test-before-sec test-after-ts test-before-ts

clean: ; rm -rf bin *~

bin/runat : runat.go
	$(call hdr,$@)
	GOBIN=$$(pwd)/bin go install $<

.PHONY: help
help: bin/runat
	$(call hdr,$@)
	bin/runat -h

# Test run is after the current time mark.
# Time specification is in seconds mark.
.PHONY: test-after-sec
test-after-sec: bin/runat
	$(call hdr,$@)
	@Sec=$$(date +'%S') ; \
	(( Sec += 0 )) ; \
	if (( Sec > 57 )) ; then \
		sleep 2 ; \
		Sec=$$(date +'%S') ; \
	fi ; \
	(( Sec += 0 )) ; \
	echo "current mark is $$Sec" ; \
	(( Sec += 2 )) ; \
	echo "running at the $$Sec second mark" ; \
	bin/runat -vv $$Sec /bin/bash -c "date && pwd"

# Test run is before the current time mark (wait a bit).
# Time specification is in seconds mark.
.PHONY: test-before-sec
test-before-sec: bin/runat
	$(call hdr,$@)
	@Sec=$$(date +'%S') ; \
	(( Sec += 0 )) ; \
	if (( Sec < 1 )) ; then \
		sleep 1 ; \
		Sec=$$(date +'%S'); \
	fi ; \
	(( Sec += 0 )) ; \
	echo "current mark is $$Sec" ; \
	(( Sec = 0  )) ; \
	echo "running at the $$Sec second mark" ; \
	bin/runat -vv $$Sec /bin/bash -c "date && pwd"

# Test run is after the current time mark.
# Time specification is a time stamp.
# Blast! The BSD date command does not have -s.
.PHONY: test-after-ts
test-after-ts: bin/runat
	$(call hdr,$@)
	@Sec=$$(date +'%S') ; \
	(( Sec += 0 )) ; \
	if (( Sec > 57 )) ; then \
		sleep 2 ; \
		Sec=$$(date +'%S') ; \
	fi ; \
	(( Sec += 0 )) ; \
	echo "current mark is $$Sec" ; \
	(( Sec += 2 )) ; \
	Hr=$$(date +'%H') ; \
	Min=$$(date +'%M') ; \
	Dts=$$(echo $$Hr $$Min $$Sec | awk '{printf("%02d:%02d:%02d",$$1,$$2,$$3)}') ; \
	echo "running at the $$Dts time stamp" ; \
	bin/runat -vv $$Dts /bin/bash -c "date && pwd"

# Test run is before the current time mark (wait a bit).
# Time specification is a time stamp.
# Blast! The BSD date command does not have -s.
.PHONY: test-before-ts
test-before-ts: bin/runat
	$(call hdr,$@)
	@Sec=$$(date +'%S') ; \
	(( Sec += 0 )) ; \
	if (( Sec < 1 )) ; then \
		sleep 1 ; \
		Sec=$$(date +'%S'); \
	fi ; \
	(( Sec += 0 )) ; \
	echo "current second mark is $$Sec" ; \
	(( Sec = 0  )) ; \
	Min=$$(date +'%M') ; \
	Hr=$$(date +'%H') ; \
	(( Min += 0 )) ; \
	(( Min++ )) ; \
	(( Hr += 0 )) ; \
	echo "current minute mark is $$Min" ; \
	if (( Min > 59 )) ; then \
		(( Min = 0 )) ; \
		(( Hr++ )) ; \
		if (( Hr > 23 )) ; then \
			(( Hr = 0 )) ; \
		fi ; \
	fi ; \
	Dts=$$(echo $$Hr $$Min $$Sec | awk '{printf("%02d:%02d:%02d",$$1,$$2,$$3)}') ; \
	echo "running at the $$Dts time stamp" ; \
	bin/runat -vv $$Dts /bin/bash -c "date && pwd"
