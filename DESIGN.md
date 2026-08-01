# Ledger Studio Design System

## Direction

Ledger Studio is an editorial, self-hosted financial workspace. It combines a
warm paper canvas, a deep-ink navigation island, coral actions and asymmetrical
bento composition. It must feel authored and tactile rather than like a themed
admin template.

Pixel character appears as tiny registration marks, tabular numerals, square
status pips and chart rulers. The visual identity comes from layout, scale and
contrast—not from covering every component in pixel borders.

## Principles

1. **Money first.** Amounts, dates and category labels are the strongest visual
   anchors. Decorative imagery must never compete with financial data.
2. **Editorial hierarchy.** A few intentionally oversized amounts and headings
   create orientation; secondary metadata stays quiet and compact.
3. **Spacious on desktop, thumb-ready on mobile.** Desktop supports scanning
   and comparison; mobile keeps large touch targets, swipe actions, quick entry,
   bottom navigation and safe-area handling.
4. **Warm confidence.** Paper, ink, coral and muted mint form the product
   palette. Expense, income and warning colors remain semantic and are always
   paired with text or icons.
5. **One component language.** Cards, dialogs, tables, forms, sheets and lists
   share the same spacing, border, radius and state tokens.
6. **Accessible by default.** WCAG AA contrast, visible focus, 44px mobile
   targets, reduced motion and no color-only meaning.

## Foundation

- UI font: system sans-serif with Chinese-first fallbacks.
- Data font: system monospace for amounts, timestamps and compact labels.
- Canvas: warm paper `#f3eee4`; surface: `#fffdf8`; deep ink: `#14251e`.
- Primary action: coral `#e56845`; supporting accents: amber `#efc65f` and
  muted mint `#a9d7c5`.
- Radius: 20–28px for signature surfaces, 12–16px for controls. A single
  squared corner gives cards a recognizable editorial cut.
- Border: 1px warm neutral. Shadows communicate elevation only on navigation,
  overlays and signature cards.
- Pixel detail: square pips, registration marks, small rulers and tabular
  numerals. Body copy always uses a readable Chinese-first sans-serif stack.

## Responsive behavior

- Desktop uses a floating 272px navigation island, a translucent 76px context
  bar and an asymmetric three-part asset bento.
- Tablet collapses navigation into an overlay without changing page hierarchy.
- Mobile retains Framework7 navigation, swipe-back, pull-to-refresh, sheets,
  swipe actions, bottom tabs and the central quick-add control.
- Authentication uses an editorial manifesto and taped ledger card on desktop.
  Phones use a dedicated pocket-ledger flow rather than shrinking the desktop
  composition.

## Motion language

- Page changes use a 170ms opacity/8px vertical transition and never animate on
  first paint.
- Buttons compress to `0.97` on press. Hover lift is limited to pointer devices
  and at most 3px.
- Menus and dialogs use origin-aware opacity/scale transitions. Switches and
  sliders animate position and color rather than layout properties.
- Motion is interruptible and uses only `transform`, `opacity` and color.
  `prefers-reduced-motion` removes travel while retaining state feedback.

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
  outline or text treatment. Hover lifts by at most 3px, press scales slightly,
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

- No arcade checkerboards, neon styling or ornamental motion. The faint paper
  grid is structural and must never compete with the content.
- No hard offset shadow on every card or button.
- No illustration carried over from the original authentication template.
- No hidden editing affordances, color-only states or small mobile controls.
