---
name: ezbookkeeping
description: Query and fully manage authenticated ezBookkeeping financial data: accounts, balances, transactions, categories, tags and tag groups. Supports safe MCP mutations and exact REST JSON bodies for advanced edits and batch operations.
---

# ezBookkeeping API Tools

Use this skill when the user wants to inspect or change their ezBookkeeping data. It can create, read, update, hide, reconcile, move, and delete financial records while preserving ezBookkeeping's ownership, validation, transfer-pair, and audit rules.

## Operating rules

1. Read the target first and use the returned ID for every update or delete.
2. Show the proposed values before destructive or broad batch operations. Use MCP `dry_run` first whenever the tool exposes it.
3. Change an account balance with `adjust_account_balance`; this creates an auditable balance-modification transaction. Never rewrite a stored balance directly.
4. Treat `update_transaction.tags` as a replacement set: an empty list clears all tags.
5. Never bypass authentication, resource ownership, transaction edit scope, currency compatibility, or in-use deletion protection.
6. For transfers, update or delete the transfer-out record; ezBookkeeping keeps the related transfer-in record consistent.

## MCP tools

The embedded MCP server exposes authenticated financial mutations:

- Accounts: `query_all_accounts`, `create_account`, `update_account`, `adjust_account_balance`, `delete_account`
- Transactions: `query_transactions`, `add_transaction`, `update_transaction`, `delete_transaction`
- Categories: `query_all_transaction_categories`, `create_transaction_category`, `update_transaction_category`, `delete_transaction_category`
- Tags and groups: `query_all_transaction_tags`, create/update/delete tools for transaction tags and tag groups

Mutation tools return stable record IDs. Writes support `dry_run`; use it for user-visible previews before committing material changes.

## Usage

### List all supported commands

Linux / macOS

```bash
sh scripts/ebktools.sh list
```

Windows

```powershell
scripts\ebktools.ps1 list
```

### Show help for a specific command

Linux / macOS

```bash
sh scripts/ebktools.sh help <command>
```

Windows

```powershell
scripts\ebktools.ps1 help <command>
```

### Call API

Linux / macOS

```bash
sh scripts/ebktools.sh [global-options] <command> [command-options]
```

Windows

```powershell
scripts\ebktools.ps1 [global-options] <command> [command-options]
```

For complex updates, nested sub-accounts, picture IDs, geographic data, or batch operations, provide the server's exact request JSON in a file:

Linux / macOS

```bash
sh scripts/ebktools.sh --body-file ./request.json transactions-modify
```

Windows

```powershell
scripts\ebktools.ps1 -bodyFile .\request.json transactions-modify
```

Run `list` to see the full authenticated command set. It includes account modification/reconciliation/deletion, full transaction modification and batch editing, category/tag/tag-group CRUD, querying, and exchange rates.

Typical exact-body commands include:

- `accounts-modify`, `accounts-reconcile`, `accounts-hide`, `accounts-delete`
- `transactions-modify`, `transactions-batch-category`, `transactions-batch-account`
- `transactions-batch-tag-add`, `transactions-batch-tag-remove`, `transactions-batch-tag-clear`
- `transactions-move-all`, `transactions-delete`, `transactions-batch-delete`
- `transaction-categories-modify|hide|delete`
- `transaction-tags-modify|hide|delete`
- `transaction-tag-groups-add|modify|delete`

Use the corresponding `*-list` or `*-get` command first. Exact bodies deliberately pass through the official REST models, which retains support for every current field without lossy command-line conversion.

## Troubleshooting

If the script reports that the environment variable `EBKTOOL_SERVER_BASEURL` or `EBKTOOL_TOKEN` is not set, user can define them as system environment variables, or create a `.env` file in the user home directory that contains these two variables and place it there.

The meanings of these environment variables are as follows:

| Variable | Required | Description |
| --- | --- | --- |
| `EBKTOOL_SERVER_BASEURL` | Required | ezBookkeeping server base URL (e.g., `http://localhost:8080`) |
| `EBKTOOL_TOKEN` | Required | ezBookkeeping API token |

## Reference

ezBookkeeping: [https://ezbookkeeping.mayswind.net](https://ezbookkeeping.mayswind.net)
