# git-issue

All contributors have differents needs, some would prefere the CLI, while others need a WebUI to get a global overview.

This is why git-issue harness the power of existing standard to allow Human First collaboration, priorization, and monitoring
- `Git` for storage, auth, sign, automation (via hook), right management (via branch protection), offline first (decentralised collaboration)
- `JSON` for it universally supported syntax that any modern programming language can use natively
- `JSON Schema` for it zero-code datamodeling/validation, state management ability, that user can adapt to it usecase

This obnoxius stack allow to build any desired interface and tooling on top

## CLI

the CLI require `JQ` for edit and filtering

```sh
git issue pull origin
git issue list
git issue create 42 # will use $EDITOR + $JSON_VALIDATOR
git issue show   42
git issue update 42 '.assignee="k@yne.fr"'
git issue push origin

git issue show | jq '.assignee=="k@yne.fr"' # filter example
```

## WebApp

The demo webapp make use of `Markdown`: for rich description and `JQ`: (via jq.js) for filtering

## TUI

current Bubbletea based TUI for issue browsing

<img width="500" height="250" alt="image" src="https://github.com/user-attachments/assets/ed69970e-3199-4b6c-b66e-7eb75f17ebe2" />


## State management

State-based fields are JSON Schema validated using the `:old` suffix before submiting the change.
- The previous value is kept as `field:old` alongside the new value.
- schema is validating both values 

```json
{
    "properties": { "status": { "enum": ["OPEN", "PENDING", "CLOSE"] } },
    "required": ["status"],

    "if": { "properties": { "status:old": { "type": "string" } }, "required": ["status:old"] },
    "then": { "allOf": [
      { "if":   { "properties": { "status:old": { "const": "OPEN" } } },
        "then": { "properties": { "status":     { "enum": ["OPEN", "PENDING"] } } } },
      { "if":   { "properties": { "status:old": { "const": "PENDING" } } },
        "then": { "properties": { "status":     { "enum": ["PENDING", "CLOSE"] } } } },
      { "if":   { "properties": { "status:old": { "const": "CLOSE" } } },
        "then": { "properties": { "status":     { "const": "CLOSE" } } } }
    ] }
}
```

## Similar projects
People did not waited 2025 to store project-related information in git, some of them also use the name git-issue
- aaiyer/bugseverywhere: store issues in filetree
- MichaelMure/git-bug: store issues in /bug/ git object namespace