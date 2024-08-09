# TOOKHOOK

---

## About

TookHook is a webhook server with pluggable handlers. The main goal of the prohect is to build custom integrations based on incoming webhook requests.

Possible use cases:
- sending messages to notification systems from service that does not support them out of the box;
- creating your own automation with services that support webhooks on events;
- collecting completion statuses of scheduled tasks (cron, etc).

The TookHook server has built-in plugin such as ~exec and ~log. You can use only the ~exec plugin, which executes specified program in the shell, but it is more efficient to create your own plugin for your tasks.. It is not difficult. Plugins are based on [go-plugin](https://github.com/hashicorp/go-plugin/tree/main) and use gRPC for communication. This allows you to build plugins not only on Go.
See more in [04-PLUGIN](https://github.com/k1nky/tookhook/blob/main/docs/04-PLUGIN.md)

## How it works

Let's look at a rule example:

```
hooks:
  - income: how_it_works_example
    handlers:
      - type: telegram
        pre:
          - template: |-
              **Jira**
              [{{ .issue.key }}] changed status to {{ .issue.fields.status.name }}.
              Project: {{ .issue.fields.project.name }}
              Assignee: {{ .issue.fields.assignee.name }}
        options:
          chat: "XXX"
          token: YYY
      - type: ~exec
        pre:
          - template: {{ .issue.key }}
        options:
          args:
            - echo $PLUGIN_EXEC_DATA > /tmp/test
        on: jira:issue_created
```

Now let's simulate an incoming webhook request:
```
curl -X POST \
    -H "content-type: application/json"\
    -d '{
  "timestamp": 1721128884783,
  "webhookEvent": "jira:issue_created",
  "user": {
    "name": "user_name",
  },
  "issue": {
    "key": "XX-YYY",
    "fields": {
      "project": {
        "name": "project_name",
      },
      "assignee": {
        "name": "user_name",
        "displayName": "User Name",
      },
      "status": {
        "name": "To Do",
      },
    }
  }
}' \
http://tookhook:8080/hook/how_it_works_example
```

What will happen:

1. The following message will be passed to the telegram plugin:
```
**Jira**
[XX-YYY] changed status to To Do.
Project: project_name
Assignee: user_name
```
2. The plugin will send this message to chat from `options.chat` with token `options.token`
3. Next the builtin plugin exec will be called. It will be called if the income message mathes to `jira:issue_created`. It is true in the example.
```
# where PLUGIN_EXEC_DATA=XX-YYY
/bin/sh -c "echo $PLUGIN_EXEC_DATA > /tmp/test"
```

You found more about the rules in [03-CONFIGURE](https://github.com/k1nky/tookhook/blob/main/docs/03-CONFIGURE.md).

## Documentation

* [01-BUILD](https://github.com/k1nky/tookhook/blob/main/docs/01-BUILD.md)
* [02-INSTALL](https://github.com/k1nky/tookhook/blob/main/docs/01-INSTALL.md)
* [03-CONFIGURE](https://github.com/k1nky/tookhook/blob/main/docs/03-CONFIGURE.md)
* [04-PLUGIN](https://github.com/k1nky/tookhook/blob/main/docs/04-PLUGIN.md)