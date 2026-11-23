# Templates Kyse

- [Templates se tornam Go comum](#template-source-becomes-ordinary-go)
- [Construa um formulário tipado com um único fluxo de CSRF](#build-a-typed-form-with-one-csrf-path)
- [Use deliberadamente o conjunto fechado de diretivas](#use-the-closed-directive-set-deliberately)
- [Trate HTML sem escape como uma decisão local de segurança](#treat-raw-html-as-a-local-security-decision)
- [Compile os fontes e inspecione o comportamento](#compile-source-and-inspect-behavior)
- [Interfaces publicadas de templates e componentes](#released-template-and-component-surface)
- [Conecte template, requisição e binário](#connect-template-request-and-binary)

<a name="template-source-becomes-ordinary-go"></a>
## Templates se tornam Go comum

Kyse é a linguagem fechada de templates compilada por `aru view:build` e o módulo de componentes tipados importados pelas views. Um `.kyse.go` não é Go válido antes da compilação, por escolha de projeto. O Go gerado verifica tipos de dados, escapa a interpolação comum e relaciona erros do compilador à linha do fonte.

<a name="build-a-typed-form-with-one-csrf-path"></a>
## Construa um formulário tipado com um único fluxo de CSRF

Use `@csrf` em formulários e chame funções de componentes tipados nas posições de saída confiável sem escape. Funções de componentes e structs de props são símbolos Go exportados; um nome ou campo incorreto vira erro do compilador na linha da view.

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
## Use deliberadamente o conjunto fechado de diretivas

A interpolação escapada é o padrão, e o escape segue a posição: um valor dentro de ypy-vals ou ypy-headers é escapado como JSON, então não fecha uma string nem acrescenta chaves ao que a requisição envia. O conjunto fechado de diretivas fornece controle de fluxo, um layout com seções nomeadas, um campo CSRF e declarações `@go`. Registrar diretivas em runtime, localizar serviços nos templates e referenciar componentes por nomes em strings não são pontos de extensão.

- Use `{{ expression }}` para texto e atributos escapados.
- Use `{!! expression !!}` apenas para um valor `template.HTML`; qualquer outro tipo falha o build na linha da view.
- Use `@if` e `@foreach` para controle de fluxo compilado.
- Use `@extends`, `@section` e `@yield` para um único contrato de layout.
- Use `@go` para tipos de dados da página e declarações locais.

<a name="treat-raw-html-as-a-local-security-decision"></a>
## Trate HTML sem escape como uma decisão local de segurança

A saída crua pula o escape, então nunca converta texto de requisição, banco de dados, Markdown ou tradução em `template.HTML` só para renderizar marcação. Layout e dados da página precisam concordar, todo formulário de navegador precisa de CSRF, e o Go gerado nunca pode virar a fonte editada. Um atributo cujo nome é interpolado reprova o build, porque o escape de um valor é escolhido pelo nome do atributo e um nome que é dado poderia ser onclick ou href. Um nome que a camada de views recusa, escrito em Attrs ou `Parts[...]`.Attrs de um componente, como `ypy-*`, `on*`, style, class, role, `aria-*` ou um endereço, entra em panic com um `*components.RefusedAttribute` que nomeia a parte, o atributo e o motivo, em vez de desenhar um componente vazio.

> [!WARNING]
> Um componente tipado pode retornar HTML confiável; conteúdo arbitrário de quem o chama, não. A saída sem escape aceita apenas `template.HTML`, então uma string vinda de uma requisição, do banco de dados ou de um catálogo interrompe o build em vez de chegar à página como marcação, e um componente cujo valor é recusado na renderização não renderiza nada, em vez de um elemento quebrado. Revise cada expressão de saída sem escape como um limite de segurança e preserve a interpolação escapada comum para strings controladas pelo usuário.

<a name="compile-source-and-inspect-behavior"></a>
## Compile os fontes e inspecione o comportamento

Execute o compilador e os testes Go para verificar texto escapado, o campo exato `_token`, marcação dos erros de campo, rótulos acessíveis, conteúdo do layout e uma fixture com campo ou prop inválido de propósito cujo diagnóstico aponte a linha do `.kyse.go`.

- `aru view:build` termina com sucesso.
- A saída escapada não permite injetar marcação.
- O campo oculto de CSRF está presente.
- Props incorretos de componentes falham na linha do fonte.

<a name="released-template-and-component-surface"></a>
## Interfaces publicadas de templates e componentes

`aru view:build` analisa Kyse e gera Go. `components.Field`, Button, seus tipos Props, `kyse/icons`, `view.RenderInto` e `RegisterLayout` são as interfaces publicadas de compilação e runtime gerado.

| Superfície | Contrato publicado |
| --- | --- |
| `aru view:build` | analisar Kyse e gerar Go |
| `components.Field` / Button | controles tipados |
| `components.*Props` | contratos verificados na compilação |
| `kyse/icons` | ícones acessíveis que o linker pode remover quando não usados |

<a name="connect-template-request-and-binary"></a>
## Conecte template, requisição e binário

As views registram o código compilado, os assets o acompanham no binário, a validação fornece o estado dos campos e o CSRF vincula cada alteração feita pelo navegador à sessão.

- [Views](/docs/{{version}}/frontend/views)
- [Assets](/docs/{{version}}/frontend/assets)
- [Validação](/docs/{{version}}/http/validation)
- [CSRF](/docs/{{version}}/http/csrf)
