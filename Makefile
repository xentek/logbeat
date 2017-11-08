GLIDE := $(shell command -v glide 2> /dev/null)
GOVERALLS := $(shell command -v goveralls 2> /dev/null)
GOLINT := $(shell command -v golint 2> /dev/null)
GEM := $(shell command -v gem 2> /dev/null)
GCG := $(shell command -v github_changelog_generator 2> /dev/null)

.PHONY: init deps test ci release

# setup project
init:
ifndef GLIDE
	go get -u github.com/Masterminds/glide
else
	@echo "Already Installed: glide. Skipping."
endif
ifndef GOVERALLS
	go get -u github.com/mattn/goveralls
else
	@echo "Already Installed: goveralls. Skipping."
endif
ifndef GOLINT
	go get -u github.com/golang/lint/golint
else
	@echo "Already Installed: golint. Skipping."
endif
ifdef GEM
ifndef GCG
	gem install github_changelog_generator --no-ri --no-rdoc
else
	@echo "Already Installed: github_changelog_generator. Skipping."
endif
endif
	glide install

# install dependencies
deps:
	glide up

# test package
test:
	go test $(glide nv)

ci: test
	goveralls -service=travis-ci

# release package
#
# Example:
#
# 	make V=v0.0.0 M='my release description' release
release:
ifndef V
	@echo  "Missing Required Argument: V"
	exit 1
endif
ifndef M
	@echo  "Missing Required Argument: M"
	exit 1
endif

	rm -f .git/RELEASE_EDITMSG
	touch .git/RELEASE_EDITMSG
	echo "$(V)\n\n$(M)" >> .git/RELEASE_EDITMSG
	git tag -s $(V) -F .git/RELEASE_EDITMSG
	git push --tags
	github_changelog_generator --issue-line-labels="ALL" --release-url="https://github.com/xentek/logbeat/releases/tag/%s"
	git commit CHANGELOG.md -m "updates changelog for $(V)"
	git push origin master
	hub release create -f .git/RELEASE_EDITMSG $(V)
