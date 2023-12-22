APP       := mall
TARGET    := serverHealthy
MFLAGS    :=
DFLAGS    :=
# CONFIG    := client
STRIP_FLAG:= N
J2GO_FLAG:= 

libpath=${subst :, ,$(GOPATH)}
$(foreach path,$(libpath),$(eval -include $(path)/src/github.com/TarsCloud/TarsGo/hack/scripts/makefile.tars))
