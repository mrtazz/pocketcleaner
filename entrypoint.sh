#!/bin/bash
export POCKETCLEANER_CONFIG=/root/.pocketcleaner.ini
echo "consumer_key = ${POCKETCLEANER_CONSUMER_SECRET}" >> ${POCKETCLEANER_CONFIG}
echo "access_token = ${POCKETCLEANER_ACCESS_TOKEN}" >> ${POCKETCLEANER_CONFIG}
echo "keep_count   = ${POCKETCLEANER_KEEP_COUNT}" >> ${POCKETCLEANER_CONFIG}

pocketcleaner --debug
