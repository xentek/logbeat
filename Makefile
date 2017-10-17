GLIDE := $(shell command -v glide 2> /dev/null)
GOVERALLS := $(shell command -v goveralls 2> /dev/null)
GOLINT := $(shell command -v golint 2> /dev/null)
GEM := $(shell command -v golint 2> /dev/null)
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
release:
	rm -f .git/RELEASE_EDITMSG
	touch .git/RELEASE_EDITMSG
	echo "${1}\n\n${@}" >> .git/RELEASE_EDITMSG
	git tag -s ${1} -F .git/RELEASE_EDITMSG
	git push --tags
	github_changelog_generator --issue-line-labels="ALL" --release-url="https://github.com/macandmia/logbeat/releases/tag/%s"
	git commit CHANGELOG.md -m "updates changelog for ${1}"
	git push origin master
	hub release create -f .git/RELEASE_EDITMSG ${1}
