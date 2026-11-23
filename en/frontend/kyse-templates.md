# Kyse templates

- [Template source becomes ordinary Go](#template-source-becomes-ordinary-go)
- [Build a typed form with one CSRF path](#build-a-typed-form-with-one-csrf-path)
- [Use the closed directive set deliberately](#use-the-closed-directive-set-deliberately)
- [Treat raw HTML as a local security decision](#treat-raw-html-as-a-local-security-decision)
- [Compile source and inspect behavior](#compile-source-and-inspect-behavior)
- [Released template and component surface](#released-template-and-component-surface)
- [Connect template, request, and binary](#connect-template-request-and-binary)

<a name="template-source-becomes-ordinary-go"></a>
## Template source becomes ordinary Go

Kyse is the closed template language compiled by `aru view:build` and the module of typed components imported by views. A `.kyse.go` file intentionally is not valid Go before compilation; generated Go asserts data types, escapes ordinary interpolation, and maps compiler errors back to the source line.

<a name="build-a-typed-form-with-one-csrf-path"></a>
## Build a typed form with one CSRF path

Use `@csrf` in a form and call typed component functions from trusted raw-output positions. Component functions and prop structs are exported Go symbols, so a wrong name or field becomes a compiler error at the view line.

```kyse
@extends('layouts.app')
@section('content')
<form method="post" action="{{ .SubmitURL }}">
    @csrf
    {!! components.Field(components.FieldProps{Name: "title", Label: "Title", Value: .Post.Title, Page: ., Required: true}) !!}
    {!! components.Button(components.ButtonProps{Label: "Publish", Type: "submit"}) !!}
</form>
@endsection
```

<a name="use-the-closed-directive-set-deliberately"></a>
## Use the closed directive set deliberately

Escaped interpolation is the default, and the escape follows the position: a value inside ypy-vals or ypy-headers is escaped as JSON, so it cannot close a string or add keys to what the request sends. The closed directive set supplies control flow, one layout and named sections, a CSRF field, and `@go` declarations. Runtime directive registration, service lookup inside templates, and string-named components are not extension points.

- Use `{{ expression }}` for escaped text and attributes.
- Use `{`!! expression `!!}` only for a `template.HTML` value; any other type fails the build at the view line.
- Use `@if` and `@foreach` for compiled control flow.
- Use `@extends`, `@section`, and `@yield` for one layout contract.
- Use `@go` for page data types and local declarations.

<a name="treat-raw-html-as-a-local-security-decision"></a>
## Treat raw HTML as a local security decision

Raw output bypasses escaping, so never convert request, database, Markdown, or translated text to `template.HTML` merely to render markup. Layout and page data must agree, every browser form needs CSRF, and generated Go must never become the edited source. An attribute whose name is interpolated fails the build, because the escape of a value is chosen from the attribute's name and a name that is data could be onclick or href. A name the view layer refuses, written in a component's Attrs or `Parts[...]`.Attrs, such as `ypy-*`, `on*`, style, class, role, `aria-*`, or an address, panics with a `*components.RefusedAttribute` naming the part, the attribute, and the reason, instead of drawing an empty component.

> [!WARNING]
> A typed component may return trusted HTML; arbitrary caller content does not. Raw output accepts only `template.HTML`, so a string from a request, the database, or a catalogue stops the build instead of reaching the page as markup, and a component whose value is refused at render time renders nothing rather than a broken element. Review every raw-output expression as a security boundary and preserve ordinary escaped interpolation for user-controlled strings.

<a name="compile-source-and-inspect-behavior"></a>
## Compile source and inspect behavior

Run the compiler and Go tests for escaped text, the exact `_token` field, field-error markup, accessible labels, layout content, and a deliberate bad field or prop fixture whose diagnostic names the `.kyse.go` line.

- `aru view:build` succeeds.
- Escaped output cannot inject markup.
- The hidden CSRF field is present.
- Bad component props fail at the source line.

<a name="released-template-and-component-surface"></a>
## Released template and component surface

`aru view:build` parses Kyse and generates Go. `components.Field`, Button, their Props types, `kyse/icons`, `view.RenderInto`, and `RegisterLayout` are the released compile-time and generated runtime surfaces.

| Surface | Released contract |
| --- | --- |
| `aru view:build` | parse Kyse and generate Go |
| `components.Field` / Button | typed controls |
| `components.*Props` | compile-time contracts |
| `kyse/icons` | accessible linker-prunable icons |

<a name="connect-template-request-and-binary"></a>
## Connect template, request, and binary

Views register compiled output, assets ship it in the binary, validation supplies field state, and CSRF binds each browser mutation to the session boundary.

- [Views](/docs/{{version}}/frontend/views)
- [Assets](/docs/{{version}}/frontend/assets)
- [Validation](/docs/{{version}}/http/validation)
- [CSRF](/docs/{{version}}/http/csrf)
