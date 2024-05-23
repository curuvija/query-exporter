# Docker Build Tools Changelog

{{#releases}}
## [{{name}}](https://github.com/curuvija/query-exporter/releases/{{name}}) ({{date}})

{{#sections}}
### {{name}}

{{#commits}}
* [{{#short5}}{{sha}}{{/short5}}](https://github.com/curuvija/query-exporter/commit/{{sha}}) {{message.shortMessage}} ({{authorAction.identity.name}})

{{/commits}}
{{^commits}}
No changes.
{{/commits}}
{{/sections}}
{{^sections}}
No changes.
{{/sections}}
{{/releases}}
{{^releases}}
No releases.
{{/releases}}

Also check out the [release notes](https://github.com/curuvija/query-exporter/blob/master/RELEASE_NOTES.md#{{name}}) for this release.
