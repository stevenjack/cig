VERISON ?= ""

PHONY: bump
SILENT: bump

current_version:
	$(eval CURRENT_VERSION:=$(shell cat VERSION))

bump: current_version
	echo "=> Bumping release to v${VERSION}"
	for file in "cig.go" "README.md" "install.sh" ; do \
		sed -i '' "s/${CURRENT_VERSION}/${VERSION}/g" $$file; \
	done
