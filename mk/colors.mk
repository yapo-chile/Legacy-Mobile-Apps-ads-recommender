
export ENABLE_COLORS=true

# Echo color variables
## Header Color
export HC=\033[1;33m
## Title Color
export TC=\033[0;33m
## Error Color
export EC=\033[0;31m
## No Color
export NC=\033[0m

define BASH_FUNC_echoColor%%
() {
    if [ -z "$$2" ]
    then
        return
    fi

	local color="$$1"
    if [ -z "$$ENABLE_COLORS" ] || [ "$$ENABLE_COLORS" != "true"  ]
    then
        echo "$${2}"
    else
        echo -e "$${color}$${2}$${NC}"
    fi
}
endef

define BASH_FUNC_echoHeader%%
() {
	echoColor $${HC} "$$1"
}
endef

define BASH_FUNC_echoTitle%%
() {
	echoColor $${TC} "$$1"
}
endef

define BASH_FUNC_echoError%%
() {
	echoColor $${EC} "$$1"
}
endef

export BASH_FUNC_echoColor%%
export BASH_FUNC_echoHeader%%
export BASH_FUNC_echoTitle%%
export BASH_FUNC_echoError%%
