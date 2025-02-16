# Format can be one of "rss", "atom" or "json"
# Token can also be passed via basic auth
GET {{base_url}}/gitea.com/notifications.{{format}}
[QueryStringParams]
token: {{gitea_token}}
