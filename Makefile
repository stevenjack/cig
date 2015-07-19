VERISON ?= ""
TYPE ?= ""

PHONY: bump
SILENT: bump

current_version:
	$(eval CURRENT_VERSION:=$(shell cat VERSION))

update_files: current_version
	echo "=> Bumping release to v${VERSION}"
	for file in "cig.go" "README.md" "install.sh" "VERSION" ; do \
		sed -i '' "s/${CURRENT_VERSION}/${VERSION}/g" $$file; \
	done

tag:
	echo "=> Tagging latest release"
	git tag v${VERSION}

commit:
	echo "=> Commiting version bump"
	git add cig.go README.md install.sh VERSION
	git commit -m "${TYPE} release ${VERSION}"

push:
	echo "=> pushing changes up"
	git push origin master
	git push origin --tags

bump: update_files commit tag
