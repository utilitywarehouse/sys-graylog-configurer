# sys-graylog-configurer
Idempotent configurer for graylog

## Configuration through environmental variables
* `ADMIN_PASSWORD`: graylog api admin password
* `ADMINS`: space separated list of users to become admins. The name must match the SSO utilitywarehouse username. Eg: `jdoe ppig`
* `EXTRACTORS`: space separated list of fields to be json parsed and have their child fields indexed. Child fields will be indexed as `parent_field.child_field`. Eg: `"kubernetes message my_field"`
* `ELASTICSEARCH_AUTHORITY`: elasticserach host and port. Eg: `elasticsearch:9200`
* `ELASTICSEARCH_CUSTOM_TEMPLATE`: json template to complement graylog default template. Eg:
```json
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
* `SSO_PERMISSIONS`: json array with the list of permissions for the users created via SSO. Defaults to admin if not present. Permission details in http://docs.graylog.org/en/2.4/pages/users_and_roles/permission_system.html. Eg:
```json
[
  "system:read",
  "savedsearches:*",
  "journal:*"
]
```

### Temporary cloudtrail configuration
For configuration of https://github.com/Graylog2/graylog-plugin-aws. All variables have a `DEV` counterpart.
* `AWS_CLOUDTRAIL_ACCOUNTS`: space separated list of accounts to setup ingestion. Eg: `"dev prod"`
* `AWS_ID_PROD`: aws user id for `prod` account
* `AWS_SECRET_PROD`: aws user secret for `prod` account

### Graylog default index configuration
Basic configuration of the graylog index. It is assumed that the rotation strategy is "size" , and the retention strategy is "delete" (https://docs.graylog.org/en/latest/pages/configuration/server.conf.html#rotation)
* `GRAYLOG_INDEX_SHARDS`: number of shards of each index (for graylog workload, a good number is 1 per node)
* `GRAYLOG_INDEX_REPLICAS`: number of replicas for each index (1 is generally enough for logs)
* `GRAYLOG_INDEX_MAX_SIZE`: maximum size in bytes of an index before a new index is created
* `GRAYLOG_MAX_INDICES`: how many indices to keep
