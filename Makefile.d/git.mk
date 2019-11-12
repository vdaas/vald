.PHONY: git/config/init
## add git configs required for development
git/config/init:
	git config commit.template ".commit_template"

.PHONY: git/hooks/init
## add configs for registering pre-defined git hooks
git/hooks/init:
	ln -sf ../../hack/git/hooks/pre-commit .git/hooks/pre-commit
