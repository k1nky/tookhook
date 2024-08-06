# Build

## Prepare
To install dependencies run:
```
make prepare
```

## Build binary
```
make build
```

## Build docker
```
make docker
```

## Build docker with plugins
Add your plugins as git submodule:
```
git submodule add <your_plugin_project_url> plugins/<your_plugin_name>
make submodules
make plugin
make docker
```
