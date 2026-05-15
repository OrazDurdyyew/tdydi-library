# Electronic Library Redesign System
Version: 1.0
Target: Nuxt 4 Academic Electronic Library Platform

---

# IMPORTANT RULES

This project already has:
- working Nuxt 4 architecture
- backend
- routing
- API
- business logic
- composables
- stores

DO NOT:
- refactor architecture
- rewrite backend logic
- change API structure
- change routing logic
- rewrite composables
- modify authentication logic
- change data models

ONLY:
- redesign UI/UX
- improve visual consistency
- improve responsiveness
- improve accessibility
- improve typography
- improve spacing/layout
- modernize components

This is a production redesign task, NOT a prototype.

---

# DESIGN DIRECTION

Create a premium academic electronic library interface.

The design should feel like a combination of:
- JSTOR
- Elsevier
- arXiv
- Medium
- university research portals

The UI must:
- feel professional
- feel calm
- prioritize readability
- prioritize search and research workflows
- avoid startup/SaaS aesthetics
- avoid flashy effects

The design language should be:
- academic
- modern
- minimal
- typography-focused
- structured
- premium

---

# DO NOT USE

Avoid:
- glassmorphism
- neon colors
- huge gradients
- floating blobs
- oversized shadows
- giant border radius
- excessive animations
- parallax effects
- startup landing page style
- crypto/Web3 aesthetics
- overly colorful cards
- random accent colors

This is an academic platform, not a marketing website.

---

# DESIGN PRINCIPLES

Priority order:
1. readability
2. content hierarchy
3. search UX
4. accessibility
5. responsiveness
6. consistency
7. subtle aesthetics

Everything must feel:
- clean
- balanced
- spacious
- readable
- organized

---

# COLOR SYSTEM

Use this palette consistently.

## Primary Colors

Primary:
#1E3A5F

Secondary:
#2C5F8D

Accent:
#C8A96B

## Neutral Colors

Background:
#F8FAFC

Surface:
#FFFFFF

Muted Surface:
#F1F5F9

Border:
#E5E7EB

## Typography

Primary Text:
#111827

Secondary Text:
#374151

Muted Text:
#6B7280

Disabled:
#9CA3AF

---

# DARK MODE

Dark mode must exist.

DO NOT use pure black.

Dark palette:

Background:
#0F172A

Surface:
#111827

Border:
#1F2937

Primary Text:
#F9FAFB

Secondary Text:
#D1D5DB

Muted Text:
#9CA3AF

Dark mode should feel soft and readable.

---

# TYPOGRAPHY SYSTEM

## Fonts

Headings/UI:
Inter

Article Content:
Source Serif 4

Fallback:
sans-serif / serif

---

# TYPOGRAPHY RULES

## Body Text

Minimum:
16px

Preferred:
16-18px

Line height:
1.6 - 1.8

---

## Headings

H1:
36-48px

H2:
28-36px

H3:
22-28px

H4:
18-22px

Typography must feel spacious and readable.

---

# SPACING SYSTEM

Use consistent 8px spacing system.

Allowed spacing:
4
8
12
16
24
32
40
48
64
80

Use generous whitespace.

The interface should breathe.

---

# BORDER RADIUS

Preferred:
rounded-lg
rounded-xl

Avoid:
rounded-3xl
rounded-full everywhere

Academic UI should not feel playful.

---

# SHADOW SYSTEM

Use minimal shadows.

Preferred:
shadow-sm
shadow-md

Most cards should rely on:
- borders
- spacing
- typography

NOT heavy shadows.

---

# LAYOUT RULES

## Max Width

Use:
max-w-7xl

Reading areas:
65ch - 75ch

Avoid extremely wide text blocks.

---

# HEADER

Header must:
- be sticky
- contain subtle blur
- have bottom border
- height 64-72px

Structure:

Left:
- logo
- platform name

Center:
- global search

Right:
- profile
- language switcher
- dark mode toggle

Header should feel clean and premium.

---

# SEARCH UX

Search is the MOST IMPORTANT feature.

Requirements:
- large centered search bar
- excellent mobile experience
- filters
- tags
- advanced search support
- keyboard accessibility

Search input should feel modern and fast.

---

# HOMEPAGE STRUCTURE

Homepage should contain:

1. Hero Search Section
2. Research Categories
3. Latest Publications
4. Featured Journals
5. Recently Added
6. Statistics
7. Collections
8. Footer

---

# HERO SECTION

DO NOT create oversized marketing hero sections.

Hero should contain:
- search bar
- short academic message
- quick statistics
- categories

Keep it compact and useful.

---

# ARTICLE CARDS

Cards should contain:
- title
- authors
- journal
- year
- DOI
- abstract preview
- tags
- PDF action

Cards must:
- prioritize readability
- have clean spacing
- use subtle hover states

Avoid:
- giant thumbnails
- unnecessary gradients
- excessive buttons

---

# ARTICLE PAGE

This is the most important page.

Layout:

Left Sidebar:
- metadata
- DOI
- journal
- citation
- PDF download
- keywords

Main Content:
- title
- abstract
- article content
- references

Right Sidebar:
- related articles
- metrics
- recommendations

---

# READING EXPERIENCE

Article pages must:
- feel comfortable for long reading
- use serif typography
- have excellent line spacing
- avoid visual noise

Reading width:
65-75ch

---

# TABLES

Tables should:
- be readable
- support overflow on mobile
- use sticky headers if needed
- use subtle borders

Avoid:
- dense spreadsheet look
- tiny text

---

# MOBILE DESIGN

Mobile-first responsive design is REQUIRED.

---

# MOBILE RULES

## Navigation

Desktop sidebar → mobile drawer.

---

## Filters

Desktop filters → bottom sheet or drawer on mobile.

---

## Cards

Desktop:
3-4 columns

Tablet:
2 columns

Mobile:
1 column

---

## Touch Targets

Minimum:
44px

Buttons must be touch-friendly.

---

# PERFORMANCE

Prioritize:
- fast loading
- lazy loading
- optimized images
- SSR compatibility
- accessibility
- semantic HTML

---

# ACCESSIBILITY

Must support:
- keyboard navigation
- screen readers
- visible focus states
- sufficient contrast
- semantic HTML
- aria labels

Accessibility is mandatory.

---

# COMPONENT STYLE

All components should:
- share consistent spacing
- use consistent radius
- use consistent typography
- use subtle interactions

The platform should feel unified.

---

# INTERACTIONS

Animations should be subtle.

Preferred:
- 150-250ms
- opacity
- translateY
- soft hover transitions

Avoid:
- bouncing animations
- large scaling
- flashy motion

---

# BUTTONS

Primary buttons:
- solid primary color

Secondary buttons:
- border style

Danger:
- minimal red usage

Buttons should feel professional and restrained.

---

# FORMS

Forms must:
- have proper spacing
- large readable inputs
- visible focus states
- proper validation UI

Inputs should:
- be clean
- modern
- readable

---

# SIDEBARS

Sidebars should:
- use subtle borders
- avoid heavy backgrounds
- remain readable
- collapse properly on mobile

---

# EMPTY STATES

Create premium empty states.

Should include:
- icon
- explanation
- action suggestion

Avoid generic placeholders.

---

# LOADING STATES

Use:
- skeleton loaders
- progressive loading
- subtle shimmer

Avoid:
- spinner-only loading everywhere

---

# FOOTER

Footer should contain:
- ISSN
- indexing systems
- policies
- copyright
- contact info
- journal links

Footer should feel academic and structured.

---

# UI QUALITY TARGET

The final design should feel:
- university-grade
- production-ready
- premium
- minimal
- elegant
- readable
- scalable

The redesign must improve:
- readability
- navigation
- research workflows
- mobile UX
- accessibility
- visual consistency

---

# FINAL IMPORTANT RULE

Do NOT redesign for visual trends.

Design for:
- researchers
- students
- professors
- long reading sessions
- academic workflows

The platform should feel trustworthy, stable, and professional.