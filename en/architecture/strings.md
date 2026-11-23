# Strings

- [Standard library first](#standard-library-first)
- [Shared naming conventions](#shared-naming-conventions)
- [Boundaries and empty slugs](#boundaries-and-empty-slugs)

<a name="standard-library-first"></a>
## Standard library first

Use the standard library for ordinary string work. strings provides searching, splitting, replacement, trimming, comparison, and case conversion without adding an application convention. strconv owns parsing and formatting primitive values, while unicode and utf8 own character and encoding questions. Keeping those calls visible makes familiar Go behavior easy to recognize.

Reach for str when several packages must agree on a naming convention that the standard library does not define: converting a field name into a label, deriving a Go-style identifier, inflecting an English noun, folding supported Latin text to ASCII, or producing an address slug. A transformation used by one feature stays with that feature instead of becoming a global helper.

<a name="shared-naming-conventions"></a>
## Shared naming conventions

Studly and Camel build identifiers from dash-, underscore-, or space-separated words. Snake and Kebab make the opposite boundary explicit. Headline turns a field-like value into a label; framework :frameworkVersion uses that exact function when validation humanizes a field name, so Password Confirmation is the shared output rather than wording chosen separately by each form.

Slug takes its separator explicitly because Go has no default arguments. Package functions are the clearest form for one operation. `str.Of` is available when a sequence of transformations describes one value, and `ToString` returns the final string; both paths call the same package behavior.

```go example
func NamingExamples() (string, string, string) {
	return str.Studly("invoice_item"), // InvoiceItem
		str.Headline("password_confirmation"), // Password Confirmation
		str.Slug("Quarterly Report 2026", "-") // quarterly-report-2026
}
```

> [!NOTE]
> Keep the original value as the source of truth. A derived identifier can be produced again; storing several case variants gives them separate opportunities to drift.

<a name="boundaries-and-empty-slugs"></a>
## Boundaries and empty slugs

ASCII is a documented folding table, not a general transliterator. It covers Latin-1 Supplement and Latin Extended-A; unsupported writing systems are dropped rather than guessed. Slug first uses that fold, so a title written only in an unsupported script can produce an empty slug. Check the result and use a separately issued stable key when empty output is possible.

A slug is presentation, not identity, authorization, or a database allowlist. Tenant identifiers, storage keys, route parameters, and column choices still need validation by the subsystem that consumes them. Human-language sentences belong in translation catalogues: changing word boundaries or case does not translate meaning.
