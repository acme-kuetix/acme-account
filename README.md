# std-account

A Kuetix package providing the **Chart of Accounts** domain primitive for ERP systems. Each account has a code, name, type (asset/liability/equity/revenue/expense), optional parent code (for hierarchical charts), and a balance direction (debit or credit) that determines which side increases the account's balance.

## What's included

- **Transition module** (`modules/account/transitions/account.go`) — in-memory CRUD + `Tree` (nested chart builder). Enforces:
  - Type → direction mapping: `asset`/`expense` → `debit`, `liability`/`equity`/`revenue` → `credit`
  - Type-consistency: a child must match its parent's type
  - Cycle prevention: `parentCode` cannot create a cycle
  - Leaf-only deletion: deleting an account with children fails (delete leaves first)
- **Endpoint workflows** (`workflows/solutions/account/`) — ready-to-wire HTTP endpoint workflows:
  - `list.wsl` — `GET /accounts?type=asset`
  - `create.swsl` — `POST /accounts`
  - `get.wsl` — `GET /accounts/{code}`
  - `update.wsl` — `PUT /accounts/{code}`
  - `delete.wsl` — `DELETE /accounts/{code}`
  - `tree.wsl` — `GET /accounts/tree` (nested chart)

## Account fields

| Field | Type | Notes |
|---|---|---|
| `code` | string | unique, required (e.g. `1000`) |
| `name` | string | required |
| `type` | string | `asset`, `liability`, `equity`, `revenue`, `expense` — required |
| `parentCode` | string | optional, must reference existing account of same type |
| `debitCredit` | string | `debit` or `credit` — defaults from type |
| `active` | bool | default true |

## Installation

```bash
go get github.com/acme-kuetix/acme-account
```

Enable in your app's `modules/modules.go`:

```go
import (
	_ "github.com/acme-kuetix/acme-account/modules/account/transitions"
	stdAccountModules "github.com/acme-kuetix/acme-account/modules"
)

func Enable() {
	stdAccountModules.Enable()
}
```

## Usage in WSL

```wsl
import account/account

state Setup {
  action account/account.Create(code: "1000", name: "Assets", typ: "asset") as A1
  on success -> Child
}

state Child {
  action account/account.Create(code: "1100", name: "Cash", typ: "asset", parentCode: "1000") as A2
  // $A2.response.debitCredit → "debit" (derived from type=asset)
  end ok
}
```

## Wiring into an app's routes table

```wsl
const {
  routes: {
    "/accounts": [
      { method: "GET",    "workflow": "@solutions/account/list",   description: "List accounts" },
      { method: "POST",   "workflow": "@solutions/account/create", description: "Create account", require: { json: ["code", "name", "type"] } },
    ],
    "/accounts/tree": [
      { method: "GET",    "workflow": "@solutions/account/tree", description: "Account tree" },
    ],
    "/accounts/{code}": [
      { method: "GET",    "workflow": "@solutions/account/get",    description: "Get account",    require: { url: ["code"] } },
      { method: "PUT",    "workflow": "@solutions/account/update", description: "Update account", require: { url: ["code"] } },
      { method: "DELETE", "workflow": "@solutions/account/delete", description: "Delete account", require: { url: ["code"] } },
    ],
  }
}
```

**Note:** Route `/accounts/tree` must be registered **before** `/accounts/{code}` — otherwise `tree` will be captured as the `{code}` parameter. Put specific routes ahead of parameterized ones.

## Development

```bash
./runner.sh validate       # validate all workflows
./runner.sh build           # build runtime/bin/std-account
GOWORK=off go build -o runtime/bin/std-account ./cmd/cli   # clean build
```

After editing transition files, run `kue update` to regenerate `modules/di.go`, `modules/meta.go`, `modules/modules.json`.
