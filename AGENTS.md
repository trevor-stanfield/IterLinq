# Agent Configuration

## Style Preference: Swift-Like Readability (All Languages)

Apply Swift-style clarity conventions to naming and API design across languages in this repository.

- Prefer descriptive, intention-revealing names over abbreviations.
- Avoid short ambiguous identifiers (for example `cmp`, `src`, `dst`) when a clearer name is practical.
- Use boolean names that read as predicates: `is...`, `has...`, `should...`, `can...`.
- Favor API names that read naturally at call sites (verb/noun phrasing).
- Keep loop indices short only when conventional and localized (`i`, `j` are acceptable in tight loops).
- When public names are improved, preserve compatibility aliases/wrappers when feasible.
- Prioritize readability and maintainability over minimal character count.
