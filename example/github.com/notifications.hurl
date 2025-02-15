# Format can be one of "rss", "atom" or "json"
GET {{base_url}}/github.com/notifications.{{format}}
[QueryStringParams]
token: {{github_token}}
