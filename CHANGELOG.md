# Release Notes

## Version v0.0.5 - 2020-06-30

Features:

- Added declaration of CPU and memory usage limit in stack.yaml
- Added okteto binary installation
- Added new command `safira okteto login`
- Added new command `safira template pull`
- Added Dockerfile
- Added namespace flag in command `safira function deploy`
- Added .gitlab-ci.yml file creation at execution `safira function new`

Bug Fixes:

- Fix `safira remove function` for remove service and ingress

## Version v0.0.4 - 2020-06-22

Features:

- Added swagger-ui deploy
- Added swagger-editor deploy

Improvements:

- Improvements function deploy

Bug Fixes:

- Fix run `safira init` in user root

## Version v0.0.3 - 2020-06-10

Features:

- Added new flag `hostname` in `safira function deploy`
- Added config scale in a stack.yaml
- Added possibility to declare custom yamls in stack.yaml

Improvements:

- Added info UP-TO-DATE/AVAILABLE in `safira infra status`

## Version v0.0.2 - 2020-06-08

Features:

- Added new flag `kubeconfig` in `safira function deploy`
- Added the new command `sapira infra status`

Improvements:

- `safira function deploy` Add flag `update`, force the deploy to take a new image
- `safira infra secrets` Add Konga secrets
- `safira function new` Yaml project renamed to stack.yaml
- Changed the search for .env information from `safira function deploy` to yaml project
- Improved various error messages
- Removed kongplugin creation on deploy
- Added new `all-functions` flag in the deploy
- Separate the build and push command, maintaining an alternative to the two actions with the name of build-push
- Atualizado o namespace do vertigo ipaas

## Version v0.0.1-beta.2 - 2020-06-01

Features:

- Add flag verbose in all commands
- Add new command `safira init`

Improvements:

- Improvements in all errors messages

## Version v0.0.1-beta - 2020-05-26

Features:
    
- Add new command `safira infra secrets` to get user and password of applications
- `safira function new` Add a file .env

Improvements:
    
- `function new` Add a folder deploy in .gitignore

## Version v0.0.1-alpha - 2020-05-25

Initial Safira release
