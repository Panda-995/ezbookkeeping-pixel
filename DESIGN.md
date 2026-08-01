# Pixel Ledger Design System

## Direction

Pixel Ledger is a calm, self-hosted financial workspace with restrained pixel
details. It should feel like a precise digital ledger, not an arcade game and
not a generic admin template.

The interface uses clean surfaces, compact navigation, tabular financial
figures and small stepped corners. Pixel character appears in brand marks,
status pips, selected navigation and data decoration—not as heavy borders or
hard shadows on every component.

## Principles

1. **Money first.** Amounts, dates and category labels are the strongest visual
   anchors. Decorative imagery must never compete with financial data.
2. **Dense on desktop, effortless on mobile.** Desktop supports scanning and
   comparison; mobile keeps large touch targets, swipe actions, quick entry,
   bottom navigation and safe-area handling.
3. **Calm confidence.** Use graphite, paper white and ledger green. Expense,
   income and warning colors are semantic and always paired with text or icons.
4. **One component language.** Cards, dialogs, tables, forms, sheets and lists
   share the same spacing, border, radius and state tokens.
5. **Accessible by default.** WCAG AA contrast, visible focus, 44px mobile
   targets, reduced motion and no color-only meaning.

## Foundation

- UI font: system sans-serif with Chinese-first fallbacks.
- Data font: system monospace for amounts, timestamps and compact labels.
- Canvas: cool neutral gray; surfaces: white / deep graphite in dark mode.
- Primary: ledger green. Accent: warm amber. Expense: brick red. Income: green.
- Radius: 10–14px for major surfaces, 8–10px for controls.
- Border: 1px quiet neutral. Shadows: soft and rare.
- Pixel detail: 2px corner steps, square status pips and tabular numerals.

## Responsive behavior

- Desktop uses a 248px navigation rail and a compact 64px utility bar.
- Tablet collapses navigation into an overlay without changing page hierarchy.
- Mobile retains Framework7 navigation, swipe-back, pull-to-refresh, sheets,
  swipe actions, bottom tabs and the central quick-add control.
- Authentication uses one focused column on phones and a ledger preview plus
  form panel on wide screens.

## Registration flow

Registration is a single, uninterrupted task rather than a setup wizard. The
user enters only a username, email address, password, and password
confirmation. The username also becomes the initial display name; language,
currency, and first-day-of-week defaults follow the active locale; and the
concise common category tree is created automatically. Successful internal
registration goes directly to the overview. Email-verification deployments
show the verification result and return to sign-in without exposing unrelated
preferences.

Advanced choices remain available in Settings after the account exists. This
keeps first use fast without removing control from experienced users.

## Component and interaction states

- **Buttons:** one filled primary action per region; secondary actions use
  outline or text treatment. Hover lifts by at most 1px, press scales slightly,
  disabled states remain legible, and icon-only buttons have accessible names.
- **Inputs:** persistent visible labels, quiet neutral borders, a 3px focus
  halo, inline error text, and no competing placeholder/label copy. Password
  fields use password-manager-compatible autocomplete values.
- **Checkboxes and radios:** 44px interactive rows on mobile, primary outline
  on focus, and label text owns the click target. Selection is never
  communicated by color alone.
- **Switches:** the track changes color and the thumb changes position; the
  label states what the setting controls. Destructive or privacy-sensitive
  switches are not auto-saved without feedback.
- **Sliders:** enlarged thumb, visible active track, keyboard support, and a
  numeric value or text label adjacent to the control.
- **Selects, menus, sheets, and dialogs:** the trigger keeps a visible label;
  overlays use one clear title, explicit confirm/cancel actions, Escape/back
  dismissal, and a dimmed backdrop. Mobile sheets remain swipe-to-close.
- **Tables and lists:** headers are visually quieter than amounts; rows have
  hover/focus states; edit and overflow actions stay explicit. Bulk actions
  appear only after selection.
- **Drag and reorder:** handles have a 44px hit area, `grab`/`grabbing`
  cursors, a raised active item, and a distinct drop target. Reordering is an
  enhancement, not the only way to understand or edit an item.
- **Loading and feedback:** skeletons preserve layout; buttons expose loading
  state; validation uses alert/status regions; success and errors are concise
  and do not shift the primary action unexpectedly.

## Page architecture

- The global desktop add action is the single entry point for a new
  transaction; page-level duplicates are removed.
- Overview prioritizes net assets, monthly cash flow, period summaries, and
  trends in that order.
- Transaction pages keep view mode, date range, type, page size, search, and
  row actions in stable locations.
- Accounts, categories, tags, templates, and settings inherit the same field,
  dialog, list, empty-state, and drag-state language.
- Mobile preserves native-feeling navigation and puts quick entry, accounts,
  and transaction details within thumb reach.

## Anti-patterns

- No full-screen grids, checkerboards, neon arcade styling or ornamental noise.
- No hard offset shadow on every card or button.
- No illustration carried over from the original authentication template.
- No hidden editing affordances, color-only states or small mobile controls.
