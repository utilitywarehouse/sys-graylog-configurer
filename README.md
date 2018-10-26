# sys-graylog-configurer
Idempotent configurer for graylog

## Configuration through environmental variables
* `ADMIN_PASSWORD`: graylog api admin password
* `EXTRACTORS`: space separated list of fields to be json parsed and have their child fields indexed. Child fields will be indexed as `parent_field.child_field`. Eg: `"kubernetes message my_field"`
* `ELASTICSEARCH_AUTHORITY`: elasticserach host and port. Eg: `elasticsearch:9200`
* `ELASTICSEARCH_CUSTOM_TEMPLATE`: json template to complement graylog default template. Eg:
```
{
  "template": "graylog_*",
  "settings": {
    "refresh_interval": "30s"
  },
  "mappings": {
    "message": {
      "properties": {
        "id": {
          "type": "keyword"
        },
        "level": {
          "type": "keyword"
        },
        "log_level": {
          "type": "keyword"
        },
        "log_status": {
          "type": "keyword"
        },
        "log_timestamp": {
          "type": "keyword"
        },
        "date": {
          "type": "keyword"
        },
        "status": {
          "type": "keyword"
        }
      }
    }
  }
}
```

### Temporary cloudtrail configuration
For configuration of https://github.com/Graylog2/graylog-plugin-aws. All variables have a `DEV` counterpart.
* `AWS_CLOUDTRAIL_PROD_ENABLED`: boolean
* `AWS_SQS_QUEUE_PROD`: sqs queue name
* `AWS_ID_PROD`: aws user id
* `AWS_SECRET_PROD`: aws user secret
