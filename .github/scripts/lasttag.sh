#!/bin/sh

last_tag() {
  LAST_TAG=$(git tag | sort --version-sort --filed-separator=. | tail -n 1 2> /dev/null)

  if test -z "$LAST_TAG"; then
    LAST_TAG="v0.0.0"
  fi

  echo $LAST_TAG
}

last_tag