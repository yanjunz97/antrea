#!/usr/bin/env bash

set -eu

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
CRON_LOG_FILE=$THIS_DIR/cron.log
ENV_VARIABLE_FILE=$THIS_DIR/env.sh
CRONTAB_FILE=$THIS_DIR/monitor-cron

touch $CRON_LOG_FILE
printenv | sed 's/^\(.*\)$/export \1/g' >> $ENV_VARIABLE_FILE
chmod +x $ENV_VARIABLE_FILE
crontab $CRONTAB_FILE
cron -f
