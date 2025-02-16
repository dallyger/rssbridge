# Format can be one of "rss", "atom" or "json"
# Token can also be passed via basic auth
GET {{base_url}}/codeberg.org/notifications.{{format}}
[QueryStringParams]
token: {{codeberg_token}}
