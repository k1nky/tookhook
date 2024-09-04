# Build

## Prepare
To install dependencies run:
```
make prepare
```
You should run `prepare` before the first build.

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
git submodule update --init --recursive --remote
```
And run:
```
make docker
```
