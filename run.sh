#!/bin/sh
until curl -s "http://127.0.0.1:9000/api/system/lbstatus"; do
  echo 'graylog not ready, sleeping for 3 seconds'
  sleep 3
done
sleep 10
printf "\n"

graylog_api="http://admin:${ADMIN_PASSWORD}@127.0.0.1:9000/api"

set_input()
{
  input_name="$1"
  input_data="$2"
  origin="$3"

  printf "\nSetup ${input_name}\n"

  input_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "'"${input_name}"'") | .id')
  if [ ! "${input_id}" ]; then
    curl -s -X POST -H "Content-Type: application/json" -d "${input_data}" "${graylog_api}/system/inputs"
    sleep 5
    printf "\n${input_name} created\n"
  else
    curl -s -X PUT -H "Content-Type: application/json" -d "${input_data}" "${graylog_api}/system/inputs/${input_id}"
    sleep 5
    printf "\n${input_name} updated\n"
  fi

  if [ "${origin}" ]; then
    input_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "'"${input_name}"'") | .id')
    curl -s -X POST -H "Content-Type: application/json" -d '{"key":"origin","value":"'"${origin}"'"}' "${graylog_api}/system/inputs/${input_id}/staticfields"
    sleep 5
    printf "${origin} origin added to ${input_name}\n"
  fi
}

if [ "${AWS_CLOUDTRAIL_PROD_ENABLED}" = "true" ]; then
  input_name=aws_cloudtrail_input_prod
  input_data='{"title":"'"${input_name}"'","type":"org.graylog.aws.inputs.cloudtrail.CloudTrailInput","configuration":{"aws_sqs_region":"eu-west-1","aws_s3_region":"eu-west-1","aws_sqs_queue_name":"'"${AWS_SQS_QUEUE_PROD}"'","aws_access_key":"'"${AWS_ID_PROD}"'","aws_secret_key":"'"${AWS_SECRET_PROD}"'"},"global":true}'
  set_input ${input_name} ${input_data} cloudtrail-prod
fi

if [ "${AWS_CLOUDTRAIL_DEV_ENABLED}" = "true" ]; then
  input_name=aws_cloudtrail_input_dev
  input_data='{"title":"'"${input_name}"'","type":"org.graylog.aws.inputs.cloudtrail.CloudTrailInput","configuration":{"aws_sqs_region":"eu-west-1","aws_s3_region":"eu-west-1","aws_sqs_queue_name":"'"${AWS_SQS_QUEUE_DEV}"'","aws_access_key":"'"${AWS_ID_DEV}"'","aws_secret_key":"'"${AWS_SECRET_DEV}"'"},"global":true}'
  set_input ${input_name} ${input_data} cloudtrail-dev
fi

input_data='{"title":"gelf_tcp_input","type":"org.graylog2.inputs.gelf.tcp.GELFTCPInput","configuration":{"port":12202,"bind_address":"0.0.0.0"},"global":true}'
set_input gelf_tcp_input ${input_data}

printf "\nSetup kubernetes extractor\n"
tcp_input_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "gelf_tcp_input") | .id')
k8s_extractor_id=$(curl -s -XGET "${graylog_api}/system/inputs/${tcp_input_id}/extractors" | jq -r '.extractors[] | select(.title == "kubernetes") | .id')
k8s_extractor='{"title":"kubernetes","cut_or_copy":"copy","source_field":"kubernetes","extractor_type":"json","target_field":"","extractor_config":{"key_prefix":"k8s_"},"converters":{},"condition_type":"none","condition_value":""}'

if [ ! "${k8s_extractor_id}" ]; then
  curl -s -X POST -H "Content-Type: application/json" -d "${k8s_extractor}" "${graylog_api}/system/inputs/${tcp_input_id}/extractors"
  printf "\nk8s extractor created\n"
else
  curl -s -X PUT -H "Content-Type: application/json" -d "${k8s_extractor}" "${graylog_api}/system/inputs/${tcp_input_id}/extractors/${k8s_extractor_id}"
  printf "\nk8s extractor updated\n"
fi

printf "\nSetup SSO plugin\n"
sso_plugin_config='{"username_header":"X-Forwarded-User","email_header":"X-Forwarded-Email","default_group":"Admin","auto_create_user":true,"require_trusted_proxies":true}'
curl -s -X PUT -H "Content-Type: application/json" -d "${sso_plugin_config}" "${graylog_api}/plugins/org.graylog.plugins.auth.sso/config"

printf "\n\nSetup custom elastic search template\n"
until curl -s "http://${ELASTICSEARCH_AUTHORITY}/_cluster/health"; do
  echo 'elasticsearch not ready, sleeping for 3 seconds'
  sleep 3
done
custom_template='{"template":"graylog_*","settings":{"refresh_interval":"30s"},"mappings":{"message":{"properties":{"id":{"type":"keyword"},"level":{"type":"keyword"},"date":{"type":"keyword"},"status":{"type":"keyword"}}}}}'
curl -s -X PUT -H "Content-Type: application/json" -d "${custom_template}" "http://${ELASTICSEARCH_AUTHORITY}/_template/graylog-custom-template"
# Rotate the active index to activate the new template
curl -s -X POST "${graylog_api}/system/deflector/cycle"

printf "\n\nGoing to sleep...\n"
sleep infinity
