# git-issue


git-issue harness the power of existing standard to allow Human First collaboration, priorization, and monitoring
- `Git` for storage, auth, sign, automation (via hook), right management (via branch protection), offline first (decentralised collaboration)
- `JSON` for it universally supported syntax that any modern programming language can use natively
- `JSON Schema` for it zero-code datamodeling/validation, state management ability, that user can adapt to it usecase

## The right tool for the right job

CLI, TUI, WebUI, Script, Bot

Example Tools also makes use of
- `Markdown`: for rich description in WebUI
- `JQ`: for CLI filtering and WebUI filtering

## CLI

```sh
git issue list
git issue create --title first --status OPEN
git issue create '{"title":"first","status":"OPEN"}' # raw style
git issue update 1 --status CLOSE
git issue show   1 # {"title":"first","status":"CLOSE"}
git issue list | jq .title==first 
```

## WebApp


## TUI

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