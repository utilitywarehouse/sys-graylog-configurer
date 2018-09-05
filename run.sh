#!/bin/sh
until curl -s "http://127.0.0.1:9000/api/system/lbstatus"; do
  echo 'graylog not ready, sleeping for 3 seconds'
  sleep 3
done

graylog_api="http://admin:${ADMIN_PASSWORD}@127.0.0.1:9000/api"

if [ "$AWS_ENABLED" = "true" ]; then
  printf "\nSetup aws configuration\n"
  aws_plugin_config='{"lookup_regions":"eu-west-1", "access_key": "'"${AWS_ACCESS_KEY_ID}"'", "secret_key": "'"${AWS_SECRET_ACCESS_KEY}"'"}'
  curl -s -X PUT -H "Content-Type: application/json" -d "${aws_plugin_config}" "${graylog_api}/system/cluster_config/org.graylog.aws.config.AWSPluginConfiguration"

  printf "\nSetup aws input\n"
  aws_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "aws_cloudtrail_input") | .id')
  aws_input='{"title":"aws_cloudtrail_input","type":"org.graylog.aws.inputs.cloudtrail.CloudTrailInput","configuration":{"aws_sqs_region":"eu-west-1","aws_s3_region":"eu-west-1","aws_sqs_queue_name":"'"${AWS_SQS_QUEUE_NAME}"'"},"global":true}'
  if [ ! "$aws_id" ]; then
    curl -s -X POST -H "Content-Type: application/json" -d "${aws_input}" "${graylog_api}/system/inputs"
    printf "\nAWS input created\n"
  else
    curl -s -X PUT -H "Content-Type: application/json" -d "${aws_input}" "${graylog_api}/system/inputs/${aws_id}"
    echo "\nAWS input updated\n"
  fi
fi

printf "\nSetup udp input\n"
udp_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "gelf_udp_input") | .id')
udp_input='{"title":"gelf_udp_input","type":"org.graylog2.inputs.gelf.udp.GELFUDPInput","configuration":{"port":12201,"bind_address":"0.0.0.0"},"global":true}'
if [ ! "$udp_id" ]; then
  curl -s -X POST -H "Content-Type: application/json" -d "${udp_input}" "${graylog_api}/system/inputs"
  printf "\nUDP input created\n"
else
  curl -s -X PUT -H "Content-Type: application/json" -d "${udp_input}" "${graylog_api}/system/inputs/${udp_id}"
  printf "\nUDP input updated\n"
fi

printf "\nSetup tcp input\n"
tcp_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "gelf_tcp_input") | .id')
tcp_input='{"title":"gelf_tcp_input","type":"org.graylog2.inputs.gelf.tcp.GELFTCPInput","configuration":{"port":12201,"bind_address":"0.0.0.0"},"global":true}'
if [ ! "$tcp_id" ]; then
  curl -s -X POST -H "Content-Type: application/json" -d "${tcp_input}" "${graylog_api}/system/inputs"
  printf "\nUDP input created\n"
else
  curl -s -X PUT -H "Content-Type: application/json" -d "${tcp_input}" "${graylog_api}/system/inputs/${tcp_id}"
  printf "\nUDP input updated\n"
fi

printf "\nSetup kubernetes extractor\n"
udp_id=$(curl -s -XGET "${graylog_api}/system/inputs" | jq -r '.inputs[] | select(.title == "gelf_udp_input") | .id')
k8s_extractor_id=$(curl -s -XGET "${graylog_api}/system/inputs/${udp_id}/extractors" | jq -r '.extractors[] | select(.title == "kubernetes") | .id')
k8s_extractor='{"title":"kubernetes","cut_or_copy":"copy","source_field":"kubernetes","extractor_type":"json","target_field":"","extractor_config":{"key_prefix":"k8s_"},"converters":{},"condition_type":"none","condition_value":""}'

if [ ! "$k8s_extractor_id" ]; then
  curl -s -X POST -H "Content-Type: application/json" -d "${k8s_extractor}" "${graylog_api}/system/inputs/${udp_id}/extractors"
  printf "\nK8s extractor created\n"
else
  curl -s -X PUT -H "Content-Type: application/json" -d "${k8s_extractor}" "${graylog_api}/system/inputs/${udp_id}/extractors/${k8s_extractor_id}"
  printf "\nK8s extractor updated\n"
fi

printf "\nSetup SSO plugin\n"
sso_plugin_config='{"username_header":"X-Forwarded-User","email_header":"X-Forwarded-Email","default_group":"Admin","auto_create_user":true,"require_trusted_proxies":true}'
curl -s -X PUT -H "Content-Type: application/json" -d "${sso_plugin_config}" "${graylog_api}/plugins/org.graylog.plugins.auth.sso/config"

printf "\nSetup custom elastic search template\n"
until curl -s "http://${ELASTICSEARCH_AUTHORITY}/_cluster/health"; do
  echo 'elasticsearch not ready, sleeping for 3 seconds'
  sleep 3
done
custom_mapping='{"template":"graylog_*","mappings":{"message":{"properties":{"id":{"type":"keyword"},"status":{"type":"keyword"}}}}}'
curl -s -X PUT -H "Content-Type: application/json" -d "${custom_mapping}" "http://${ELASTICSEARCH_AUTHORITY}/_template/graylog-custom-mapping"
# Rotate the active index to activate the new template
curl -s -X POST "${graylog_api}/system/deflector/cycle"

printf "\nGoing to sleep...\n"
sleep infinity
