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

## Estrutura: a do Laravel, diretório a diretório

O esqueleto é o do Laravel, com os mesmos nomes e nos mesmos lugares —
`app/Http/Controllers`, `app/Models`, `app/Policies`, `app/Providers`,
`bootstrap/`, `config/`, `database/`, `resources/views`, `routes/`, `storage/`,
`public/`. Ver `28-estrutura-do-projeto.md`.

O objetivo é único e vale mais que qualquer elegância de organização: **o dev
que vem do Laravel abre o projeto e reconhece tudo.**

Somam-se três diretórios que o Laravel não tem, e eles são a diferença do
produto:

```
app/Services/       ← regra de negócio, obrigatório
app/Repositories/   ← acesso a dado, exige Grant
app/Policies/       ← existe no Laravel, aqui é obrigatório
```

No Laravel, Service e Repository são convenção de time organizado. Aqui são
esqueleto, e o `doctor` cobra. É a frase de onboarding virada em diretório.

## O que foi deliberadamente deixado de fora

| Laravel | Decisão aqui | Por quê |
|---|---|---|
| Service container | wiring explícito gerado | reflection em Go custa e esconde |
| Facades | não existe | exigiria estado global mutável |
| Eloquent | SQL parametrizado + repository | SQL previsível, injection impossível |
| `$fillable` | não existe | request struct explícita elimina a classe de bug |
| Blade | kyse | tipagem em tempo de compilação |

As três últimas linhas são o ponto: metade da superfície de segurança do Laravel
existe para remediar perigos de linguagem dinâmica. Portar a mitigação de um bug
que Go não permite é complexidade sem retorno.
