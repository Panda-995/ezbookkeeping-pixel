---
name: ezbookkeeping
description: "Query and fully manage authenticated ezBookkeeping financial data: accounts, balances, transactions, categories, tags and tag groups. Supports safe MCP mutations and exact REST JSON bodies for advanced edits and batch operations."
---

# ezBookkeeping API Tools

Use this skill when the user wants to inspect or change their ezBookkeeping data. It can create, read, update, hide, reconcile, move, and delete financial records while preserving ezBookkeeping's ownership, validation, transfer-pair, and audit rules.

## Operating rules

1. Read the target first and use the returned ID for every update or delete.
2. Show the proposed values before destructive or broad batch operations. Use MCP `dry_run` first whenever the tool exposes it.
3. Use `update_account.balance` for a direct correction of the account's current balance. Use `adjust_account_balance` only when the user explicitly wants an auditable balance-modification transaction.
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

## Partial update contracts

Treat omitted fields as unchanged. Send only the fields the user asked to change, plus the stable record ID. An explicit empty string or list clears that value where supported.

### Edit an account

Use `update_account` for every mutable account field:

- Identity: `id` (preferred) or the current exact `account_name`
- General: `name`, `category`, `icon`, `color`, `comment`, `hidden`
- Money: `currency`, `balance`
- Credit card and reconciliation: `credit_card_statement_date`, `last_reconciled_time`, `clear_last_reconciled_time`
- Safety: `dry_run`

`balance` is the user-facing displayed balance; liability signs are normalized by the tool. Currency and balance apply to single accounts, including sub-accounts. Account `type` and parent/sub-account hierarchy are structural and are not changed by `update_account`.

Preview and then apply a full account correction:

```json
{"id":"123","name":"Daily Card","category":"credit_card","icon":1,"color":"176b5b","currency":"CNY","balance":"2688.50","comment":"Primary card","hidden":false,"credit_card_statement_date":8,"dry_run":true}
```

### Edit a transaction

Use `update_transaction` for every mutable transaction field:

- Identity and kind: `id`, `type`
- Time and classification: `time`, `category_name`
- Money and accounts: `account_name`, `amount`, `destination_account_name`, `destination_amount`
- Details: `tags`, `picture_ids`, `comment`, `hide_amount`, `geo_location`, `clear_geo_location`
- Safety: `dry_run`

Valid `type` values are `income`, `expense`, `transfer`, and `balance_modification`. `time` uses RFC 3339. `tags` and `picture_ids` are replacement lists; pass `[]` to clear them. For a transfer, update the transfer-out record and provide destination fields when changing the destination side.

Preview a multi-field transaction edit:

```json
{"id":"456","time":"2026-08-02T14:30:00+08:00","category_name":"Dining","account_name":"Daily Card","amount":"88.60","tags":["Work","Reimbursable"],"picture_ids":[],"comment":"Client lunch","hide_amount":false,"geo_location":{"latitude":31.2304,"longitude":121.4737},"dry_run":true}
```

After verifying a preview, repeat the same call with `dry_run: false` or omit `dry_run`.

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
