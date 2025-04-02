# Plugins

## Built-in plugins

### ~log
The incoming request with body will be logged.

```
hooks:
  - income: log_it
    handlers:
      - type: ~log
        comment: Just log it
```

### ~exec
Executes the specified command. The request body is passed via environment variables.

```
hooks:
  - income: exec_it
    handlers:
      - type: ~exec
        comment: Run command
        # If the Shell value is empty, /bin/sh will be used.
        shell: /bin/bash
        args: /path/to/myscript -arg1 1
```

``` /path/to/myscript
# where PLUGIN_EXEC_DATA is incoming body
/bin/sh -c "echo $PLUGIN_EXEC_DATA > /tmp/test"
```

### ~http
```
hooks:
  - income: post_it
    handlers:
      - type: ~http
        comment: Post a message
        method: POST
        uri: https://my.supersite.org/
```
