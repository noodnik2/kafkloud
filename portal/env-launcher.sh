#! /bin/sh

# functions

stderr() {
  echo "$*" 1>&2
}

fatal() {
  stderr "$@"
  exit 1
}

gen_envfile() {
  for var in $*; do
    eval val="\$$var"
    [ -z "$val" ] && fatal "environment variable '$var' was not set"
    echo "$var=$val"
  done
}

# parameters

TARGET_FILE=$1
[ -z "$TARGET_FILE" ] && fatal "no target file specified"
shift

ENV_VARS="$1"
[ -z "$ENV_VARS" ] && fatal "no environment variable list specified"
shift

[ $# -eq 0 ] && fatal "no target command was specified"

[ -s "$TARGET_FILE" ] && fatal "$TARGET_FILE is non-empty"

# execution

# do "practice run"
gen_envfile $ENV_VARS > /dev/null

# now the "real thing", outputting to the target file
gen_envfile $ENV_VARS > "$TARGET_FILE"

# invoke the target command
"$@"