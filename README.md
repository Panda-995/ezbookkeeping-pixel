# ezBookkeeping Pixel

[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CI and Multi-Arch Container](https://github.com/Panda-995/ezbookkeeping-pixel/actions/workflows/publish-container.yml/badge.svg)](https://github.com/Panda-995/ezbookkeeping-pixel/actions/workflows/publish-container.yml)
[![Architectures](https://img.shields.io/badge/image-amd64%20%7C%20arm64-176b5b)](https://github.com/Panda-995/ezbookkeeping-pixel/pkgs/container/ezbookkeeping-pixel)

[简体中文](README.zh_CN.md) · [Docker deployment](DOCKER_DEPLOYMENT.md) · [MCP & Agent guide](MCP_AGENT_GUIDE.md) · [Roadmap](FUTURE_FEATURES.md)

> This is an independent community refactor based on
> [mayswind/ezbookkeeping](https://github.com/mayswind/ezbookkeeping).
> The original project and documentation are authored by MaysWind. This
> edition adds a pixel-style interface, clearer editing flows, complete
> authenticated financial mutations for MCP/Agent clients, and a self-contained
> Docker/GHCR release workflow. It is not an official upstream release.

## Refactor highlights

- Pixel-inspired desktop and mobile design system with dark mode, visible focus states, reduced-motion support, and accessible amount colors.
- Explicit account, balance, and transaction editing entrances on desktop and mobile.
- Authenticated MCP CRUD for accounts, transactions, categories, tags, and tag groups, including `dry_run`, stable IDs, pictures, locations, and auditable balance changes.
- Agent skill scripts for exact REST bodies and batch financial operations.
- Docker Compose with SQLite persistence, health checks, complete environment-variable documentation, and automatic `linux/amd64` + `linux/arm64` GHCR publishing.
- The existing transaction import flow accepts original ezBookkeeping CSV/TSV exports from `v0.1.0` onward.

Public image:

```bash
docker pull ghcr.io/panda-995/ezbookkeeping-pixel:latest
```

## Import data exported by original ezBookkeeping

Use the existing desktop **Transaction Details → Import** flow and select
**Original ezBookkeeping Data Export File**. No migration plugin or separate
import tool is required.

- Exports from `v0.1.0–v0.4.1` are normalized from minute-precision timestamps
  and the legacy `Comment` column.
- Exports from `v0.5.0` and newer retain descriptions, locations, time zones,
  account currencies, tags, balance modifications, and cross-currency
  transfers.
- The same compatibility parser is used by the web UI, REST API, and
  `transaction-import` CLI command.

Original CSV/TSV exports contain transactions and the accounts, categories,
and tags referenced by those transactions. They do not contain passwords,
application settings, or transaction image files.

## Introduction
ezBookkeeping is a lightweight, self-hosted personal finance app with a user-friendly interface and powerful bookkeeping features. It helps you record daily transactions, import data from various sources, and quickly search and filter your bills. You can analyze historical data using built-in charts or perform custom queries with your own chart dimensions to better understand spending patterns and financial trends. ezBookkeeping is easy to deploy, and you can start it with just one single Docker command. Designed to be resource-efficient, it runs smoothly on devices such as Raspberry Pi, NAS, and MicroServers.

ezBookkeeping offers tailored interfaces for both mobile and desktop devices. With support for PWA (Progressive Web Apps), you can even [add it to your mobile home screen](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/add_to_home_screen.gif) and use it like a native app.

Live Demo: [https://ezbookkeeping-demo.mayswind.net](https://ezbookkeeping-demo.mayswind.net)

## Features
- **Open Source & Self-Hosted**
    - Built for privacy and control
- **Lightweight & Fast**
    - Minimal resource usage, runs smoothly even on low-resource devices
- **Easy Installation**
    - Docker support
    - Supports SQLite, MySQL, PostgreSQL
    - Cross-platform (Windows, macOS, Linux)
    - Works on x86, amd64, ARM architectures
- **User-Friendly Interface**
    - UI optimized for both mobile and desktop
    - PWA support for native-like mobile experience
    - Dark mode
- **AI-Powered Features**
    - Receipt image recognition
    - MCP (Model Context Protocol) support for AI integration
    - Agent Skill and API command-line script tools support for AI integration
- **Powerful Bookkeeping**
    - Two-level accounts and categories
    - Image attachments for transactions
    - Location tracking with maps
    - Scheduled transactions
    - Advanced filtering, search, visualization and analysis
- **Localization & Internationalization**
    - Multi-language and multi-currency support
    - Multiple exchange rate sources with automatic updates
    - Multi-timezone support
    - Custom formats for dates, numbers and currencies
- **Security**
    - Two-factor authentication (2FA)
    - OIDC external authentication
    - Login rate limiting
    - Application lock (PIN code / WebAuthn)
- **Data Import & Export**
    - Supports CSV, OFX, QFX, QIF, IIF, Camt.052, Camt.053, MT940, GnuCash, Firefly III, Beancount and more

For a full list of features, visit the [Full Feature List](https://ezbookkeeping.mayswind.net/features/).

## Installation
### Run with Docker
The public multi-architecture image is published to GHCR:

    $ docker pull ghcr.io/panda-995/ezbookkeeping-pixel:latest

For persistent volumes, secrets, health checks, reverse proxy settings, and
the complete environment-variable reference, use
[Docker deployment and environment variables](DOCKER_DEPLOYMENT.md).

### Install from Binary
The upstream binary releases remain available from
[mayswind/ezbookkeeping](https://github.com/mayswind/ezbookkeeping/releases);
they do not include this refactor. Build this edition from source or use its
GHCR image.

**Linux / macOS**

    $ ./ezbookkeeping server run

**Windows**

    > .\ezbookkeeping.exe server run

By default, ezBookkeeping listens on port 8080. You can then visit `http://{YOUR_HOST_ADDRESS}:8080/` .

### Build from Source
Make sure you have [Golang](https://golang.org/), [GCC](https://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps:

**Linux / macOS**

    $ ./build.sh package -o ezbookkeeping.tar.gz

All the files will be packaged in `ezbookkeeping.tar.gz`.

**Windows**

    > .\build.bat package -o ezbookkeeping.zip

or

    PS > .\build.ps1 package -Output ezbookkeeping.zip

All the files will be packaged in `ezbookkeeping.zip`.

You can also build a Docker image. Make sure you have [Docker](https://www.docker.com/) installed, then follow these steps:

**Linux**

    $ ./build.sh docker

## Contributing
We welcome contributions of all kinds.

If you find a bug, please [submit an issue](https://github.com/mayswind/ezbookkeeping/issues) on GitHub.

If you would like to contribute code, you can fork the repository and open a pull request.

Improvements to documentation, feature suggestions, and other forms of feedback are also appreciated.

You can view existing contributors on the [Contributor Graph](https://github.com/mayswind/ezbookkeeping/graphs/contributors).

## Translating
Help make ezBookkeeping accessible to users around the world. We welcome help to improve existing translations or add new ones. If you would like to contribute a translation, please refer to the [translation guide](https://ezbookkeeping.mayswind.net/translating).

Currently available translations:

| Tag | Language | Progress | Contributors |
| --- | --- | --- | --- |
| de | Deutsch | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fde.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/de.json) | [@chrgm](https://github.com/chrgm), [@1270o1](https://github.com/1270o1), [@martinschilliger](https://github.com/martinschilliger) |
| en | English | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fen.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/en.json) | / |
| es | Español | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fes.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/es.json) | [@Miguelonlonlon](https://github.com/Miguelonlonlon), [@abrugues](https://github.com/abrugues), [@AndresTeller](https://github.com/AndresTeller), [@diegofercri](https://github.com/diegofercri) |
| fr | Français | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Ffr.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/fr.json) | [@brieucdlf](https://github.com/brieucdlf) |
| it | Italiano | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fit.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/it.json) | [@waron97](https://github.com/waron97) |
| ja | 日本語 | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fja.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/ja.json) | [@tkymmm](https://github.com/tkymmm), [@Mink16](https://github.com/Mink16), [@x0x0b](https://github.com/x0x0b) |
| kn | ಕನ್ನಡ | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fkn.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/kn.json) | [@Darshanbm05](https://github.com/Darshanbm05) |
| ko | 한국어 | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fko.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/ko.json) | [@overworks](https://github.com/overworks) |
| nl | Nederlands | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fnl.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/nl.json) | [@automagics](https://github.com/automagics) |
| pt-BR | Português (Brasil) | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fpt-BR.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/pt-BR.json) | [@thecodergus](https://github.com/thecodergus), [@balaios](https://github.com/balaios) |
| ro | Română | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fro.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/ro.json) | [@gg64nou](https://github.com/gg64nou) |
| ru | Русский | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fru.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/ru.json) | [@artegoser](https://github.com/artegoser), [@dshemin](https://github.com/dshemin) |
| sl | Slovenščina | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fsl.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/sl.json) | [@thehijacker](https://github.com/thehijacker) |
| ta | தமிழ் | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fta.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/ta.json) | [@hhharsha36](https://github.com/hhharsha36) |
| th | ไทย | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fth.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/th.json) | [@natthavat28](https://github.com/natthavat28) |
| tr | Türkçe | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Ftr.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/tr.json) | [@aydnykn](https://github.com/aydnykn), [@snizamaddinov](https://github.com/snizamaddinov) |
| uk | Українська | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fuk.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/uk.json) | [@nktlitvinenko](https://github.com/nktlitvinenko), [@grid-pilot](https://github.com/grid-pilot), [@infinit1ve](https://github.com/infinit1ve) |
| vi | Tiếng Việt | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fvi.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/vi.json) | [@f97](https://github.com/f97) |
| zh-Hans | 中文 (简体) | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fzh-Hans.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/zh-Hans.json) | / |
| zh-Hant | 中文 (繁體) | [![Translation Progress](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fmayswind%2FezBookkeeping-i18n-badge%2Fmain%2Fbadges%2Fzh-Hant.json)](https://github.com/mayswind/ezBookkeeping-i18n-badge/blob/main/untranslated/zh-Hant.json) | / |

## Documentation
1. [English](https://ezbookkeeping.mayswind.net)
1. [中文 (简体)](https://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
