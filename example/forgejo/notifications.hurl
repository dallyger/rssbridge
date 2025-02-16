# Forgejo domain example: codeberg.org
# Format can be one of "rss", "atom" or "json"
# Token can also be passed via basic auth
GET {{base_url}}/forgejo/{{forgejo_domain}}/notifications.{{format}}
[QueryStringParams]
token: {{forgejo_token}}
