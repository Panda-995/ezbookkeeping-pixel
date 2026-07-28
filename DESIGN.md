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

## Anti-patterns

- No full-screen grids, checkerboards, neon arcade styling or ornamental noise.
- No hard offset shadow on every card or button.
- No illustration carried over from the original authentication template.
- No hidden editing affordances, color-only states or small mobile controls.
