# Localization

- [One translator, one locale decision per request](#one-translator-one-locale-decision-per-request)
- [Register catalogues and resolve a key](#register-catalogues-and-resolve-a-key)
- [Replace values, choose plurals, and add namespaces](#replace-values-choose-plurals-and-add-namespaces)
- [An unresolved fallback returns the key](#an-unresolved-fallback-returns-the-key)
- [Test the locale source and the catalogue](#test-the-locale-source-and-the-catalogue)
- [Translator operations](#translator-operations)
- [Routes and prose have separate language contracts](#routes-and-prose-have-separate-language-contracts)

<a name="one-translator-one-locale-decision-per-request"></a>
## One translator, one locale decision per request

Translator loads keyed sentences from registered paths and caches groups safely for concurrent requests. The request locale belongs in context, chosen once by middleware or an explicit locale scope. Downstream code reads that decision instead of parsing the URL or language header again.

<a name="register-catalogues-and-resolve-a-key"></a>
## Register catalogues and resolve a key

Build the translator once from the catalogue directory with `translation.NewFromDirectory`, naming the locale to answer in and the fallback, and hand it to the builder with `WithTranslator`. Every request then carries the catalogue beside the application name, so `view.Page` answers Translator and Locale and a component draws its own words in the request's language. Get checks the configured fallback chain; when neither the requested nor fallback catalogue resolves the item, the missing-key handler runs and the unresolved key itself is returned. The skeleton already does this: its catalogue is `lang/en.json`, and `bootstrap/app.go` hands the translator to the builder with `WithTranslator`.

```go example
func WithTranslations(drivers *foundation.Drivers) (*foundation.Application, error) {
	translator, err := translation.NewFromDirectory(os.DirFS("."), "lang", "en", "en")
	if err != nil {
		return nil, err
	}

	return foundation.Configure(drivers).WithTranslator(translator).Create()
}
```

<a name="replace-values-choose-plurals-and-add-namespaces"></a>
## Replace values, choose plurals, and add namespaces

Get replaces named placeholders, Choice selects the plural form for a number, Has checks catalogue coverage, and `AddNamespace` associates package translations with their own hint. `InLocale` is useful for work outside a request; when route middleware already supplies locale, use that single source rather than nesting a competing choice.

<a name="an-unresolved-fallback-returns-the-key"></a>
## An unresolved fallback returns the key

Fallback can keep a page usable when the requested locale is incomplete, but it does not invent a sentence. If the item is absent after fallback and the missing-key handler supplies no replacement, Get returns the key, making the catalogue defect visible. Catalogue parity tests should compare key sets and reject copied source-language values unless a short, explained exception list permits them. An unsupported locale must never become a path used to read arbitrary files.

<a name="test-the-locale-source-and-the-catalogue"></a>
## Test the locale source and the catalogue

Send requests through every supported locale prefix and verify both translated chrome and stable route generation. Exercise replacement and plural boundaries, then run a catalogue check that finds missing, extra, and suspiciously identical values. A test must also assert it inspected a plausible number of keys so a broken loader cannot pass by checking nothing.

<a name="translator-operations"></a>
## Translator operations

`NewFromDirectory` builds the application's translator and Locales lists the locales its directory carries, which is the list the locale middleware negotiates against. New, Get, Choice, Has, and `HasForLocale` form the normal read path. `AddPath` and `AddJSONPath` register application catalogues; `AddNamespace` registers package catalogues; `DetermineLocalesUsing` customizes fallback order; `HandleMissingKeysUsing` observes unresolved keys. `SetLocale` changes the default and is not a substitute for request-scoped locale context.

<a name="routes-and-prose-have-separate-language-contracts"></a>
## Routes and prose have separate language contracts

Localized chrome, dates, numbers, and product sentences belong to translation. API identifiers, Go examples, and the manual body may deliberately remain in the source language when that is the published contract. URL helpers should receive the selected locale from routing context so links, middleware, and catalogues cannot disagree about the current request.
