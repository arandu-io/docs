# Arquitetura

## O caminho obrigatório

```
Request
  ↓  Recover        panic → página de debug (dev) / 500 mudo (prod)
  ↓  Observe        request_id, logger de escopo, Collector
  ↓  Headers        CSP, HSTS, nosniff
  ↓  RateLimit      por IP e por subject
  ↓  CSRF           header (HTMX) ou campo de form
  ↓  Handler        extrai entrada, delega. Não fala com banco.
  ↓  Request        Validate() → Errors por campo
  ↓  Service        regra de negócio
  ↓  Policy         Authorize() → emite Grant
  ↓  Repository     exige Grant. Sem Grant não compila.
  ↓  DB             instrumentado no Collector
```

A diferença em relação ao Laravel está em uma linha: **o Grant**.

No Laravel, `Gate::authorize()` é uma chamada que você pode esquecer de fazer, e
nada te avisa. Aqui, `security.Grant` tem apenas campos não exportados e nenhum
construtor público além de `security.Authorize`. Como toda assinatura de
repositório exige um `Grant`, não existe caminho do handler até o banco que não
passe por uma policy. Não é convenção nem lint: é o sistema de tipos.

## Módulo vertical, não camada horizontal

O Laravel organiza por camada (`app/Http/Controllers`, `app/Models`,
`app/Policies`), o que espalha uma feature por seis diretórios. Aqui um módulo é
um diretório:

```
modules/cliente/
  module.go            registro e rotas
  cliente.entity.go
  cliente.policy.go
  cliente.repo.go
  cliente.service.go
  cliente.request.go
  queries.sql          entrada do sqlc
  views/*.templ
  cliente_test.go
```

Consequência prática: um módulo é publicável. Um diretório é uma unidade de
distribuição; trinta arquivos espalhados não são.

## O que foi deliberadamente deixado de fora

| Laravel | Decisão aqui | Por quê |
|---|---|---|
| Service container | wiring explícito gerado | reflection em Go custa e esconde |
| Facades | não existe | exigiria estado global mutável |
| Eloquent | sqlc + repository | SQL previsível, injection impossível |
| `$fillable` | não existe | request struct explícita elimina a classe de bug |
| Blade | templ | tipagem em tempo de compilação |

As três últimas linhas são o ponto: metade da superfície de segurança do Laravel
existe para remediar perigos de linguagem dinâmica. Portar a mitigação de um bug
que Go não permite é complexidade sem retorno.
