# Plantillas Kyse

- [Las plantillas se convierten en Go común](#template-source-becomes-ordinary-go)
- [Construye un formulario tipado con un único flujo de CSRF](#build-a-typed-form-with-one-csrf-path)
- [Usa deliberadamente el conjunto cerrado de directivas](#use-the-closed-directive-set-deliberately)
- [Trata el HTML sin escape como una decisión local de seguridad](#treat-raw-html-as-a-local-security-decision)
- [Compila los fuentes e inspecciona el comportamiento](#compile-source-and-inspect-behavior)
- [Interfaces publicadas de plantillas y componentes](#released-template-and-component-surface)
- [Conecta plantilla, solicitud y binario](#connect-template-request-and-binary)

<a name="template-source-becomes-ordinary-go"></a>
## Las plantillas se convierten en Go común

Kyse es el lenguaje cerrado de plantillas compilado por `aru view:build` y el módulo de componentes tipados importados por las vistas. Por diseño, un `.kyse.go` no es Go válido antes de compilarse. El Go generado comprueba tipos de datos, escapa la interpolación común y relaciona errores del compilador con la línea del fuente.

<a name="build-a-typed-form-with-one-csrf-path"></a>
## Construye un formulario tipado con un único flujo de CSRF

Usa `@csrf` en formularios y llama a funciones de componentes tipados en las posiciones de salida confiable sin escape. Las funciones de componentes y las structs de props son símbolos Go exportados; un nombre o campo incorrecto se convierte en un error del compilador en la línea de la vista.

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
## Usa deliberadamente el conjunto cerrado de directivas

La interpolación escapada es el valor predeterminado, y el escape sigue la posición: un valor dentro de ypy-vals o ypy-headers se escapa como JSON, así que no puede cerrar un string ni añadir claves a lo que envía la solicitud. El conjunto cerrado de directivas proporciona control de flujo, un layout con secciones nombradas, un campo CSRF y declaraciones `@go`. Registrar directivas en runtime, localizar servicios en las plantillas y referirse a componentes por nombres en strings no son puntos de extensión.

- Usa `{{ expression }}` para texto y atributos escapados.
- Usa `{!! expression !!}` solo para un valor `template.HTML`, una función que devuelve uno o un literal de string escrito en la vista; un campo o variable string hace fallar el build en la línea de la vista.
- Usa `@if` y `@foreach` para control de flujo compilado.
- Usa `@extends`, `@section` y `@yield` para un único contrato de layout.
- Usa `@go` para tipos de datos de página y declaraciones locales.

<a name="treat-raw-html-as-a-local-security-decision"></a>
## Trata el HTML sin escape como una decisión local de seguridad

La salida cruda omite el escape, así que nunca convierta texto de una petición, la base de datos, Markdown o una traducción en `template.HTML` solo para renderizar marcado. El layout y los datos de la página deben coincidir, todo formulario de navegador necesita CSRF y el Go generado nunca debe convertirse en la fuente editada. Un atributo cuyo nombre se interpola hace fallar el build, porque el escape de un valor se elige por el nombre del atributo y un nombre que es dato podría ser onclick o href. Un nombre que la capa de vistas rechaza, escrito en Attrs o `Parts[...]`.Attrs de un componente, como `ypy-*`, `on*`, style, class, role, `aria-*` o una dirección, entra en panic con un `*components.RefusedAttribute` que nombra la parte, el atributo y el motivo, en lugar de dibujar un componente vacío.

> [!WARNING]
> Un componente tipado puede devolver HTML confiable; el contenido arbitrario de quien lo llama, no. La salida sin escape acepta un valor `template.HTML`, una función que devuelve uno o un literal de string escrito en la vista, así que un campo o variable string procedente de una solicitud, de la base de datos o de un catálogo detiene el build en lugar de llegar a la página como marcado. Convertirlo en `template.HTML` para silenciar el error convierte el texto del visitante en marcado. Un componente cuyo valor se rechaza al renderizar no renderiza nada, en lugar de un elemento roto. Revisa cada expresión de salida sin escape como un límite de seguridad y conserva la interpolación escapada común para strings controlados por el usuario.

<a name="compile-source-and-inspect-behavior"></a>
## Compila los fuentes e inspecciona el comportamiento

Ejecuta el compilador y las pruebas Go para verificar texto escapado, el campo exacto `_token`, marcado de errores de campo, etiquetas accesibles, contenido del layout y una fixture con campo o prop deliberadamente inválido cuyo diagnóstico indique la línea del `.kyse.go`.

- `aru view:build` termina correctamente.
- La salida escapada no permite inyectar marcado.
- El campo oculto de CSRF está presente.
- Los props incorrectos de componentes fallan en la línea del fuente.

<a name="released-template-and-component-surface"></a>
## Interfaces publicadas de plantillas y componentes

`aru view:build` analiza Kyse y genera Go. `components.Field`, Button, sus tipos Props, `kyse/icons`, `view.RenderInto` y `RegisterLayout` son las interfaces publicadas de compilación y runtime generado.

| Superficie | Contrato publicado |
| --- | --- |
| `aru view:build` | analizar Kyse y generar Go |
| `components.Field` / Button | controles tipados |
| `components.*Props` | contratos comprobados durante la compilación |
| `kyse/icons` | iconos accesibles que el linker puede eliminar cuando no se usan |

<a name="connect-template-request-and-binary"></a>
## Conecta plantilla, solicitud y binario

Las vistas registran el código compilado, los assets lo incluyen en el binario, la validación proporciona el estado de los campos y CSRF vincula cada modificación del navegador con la sesión.

- [Vistas](/docs/{{version}}/frontend/views)
- [Assets](/docs/{{version}}/frontend/assets)
- [Validación](/docs/{{version}}/http/validation)
- [CSRF](/docs/{{version}}/http/csrf)
