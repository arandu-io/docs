# ADR 0022-espelhar-o-scaffolding-do-laravel-ui — A estrutura do `laravel/ui`, com a nossa stack

**Status:** Aceita. Substitui integralmente a versão anterior desta ADR, que
dizia o contrário e nasceu de um mal-entendido registrado abaixo.

## Contexto

O `docs/28-estrutura-do-projeto.md` fixou a árvore do Laravel diretório a
diretório, e ficou faltando decidir o que preenche `resources/views/` num
projeto novo. O Laravel tem uma resposta pronta e conhecida:

```bash
php artisan ui bootstrap --auth
```

Ele gera nove views, e elas são o alvo:

```
resources/views/layouts/app.blade.php
resources/views/home.blade.php
resources/views/welcome.blade.php
resources/views/auth/login.blade.php
resources/views/auth/register.blade.php
resources/views/auth/verify.blade.php
resources/views/auth/passwords/confirm.blade.php
resources/views/auth/passwords/email.blade.php
resources/views/auth/passwords/reset.blade.php
```

São 665 linhas, e o desenvolvedor Laravel conhece cada uma: o layout com a
navbar e o dropdown do usuário, o card de login, o formulário de registro, os
três passos de recuperação de senha.

## Decisão

**O `aru make:auth` gera as mesmas nove views, nos mesmos caminhos, com a mesma
função.** O que muda é o que está dentro delas:

| Camada | `laravel/ui` | Arandu |
|---|---|---|
| Estrutura de arquivos | nove views, os caminhos acima | **idêntica** |
| Motor de template | Blade | `kyse` (ADR 0020) |
| CSS | Bootstrap 5 | **Tailwind v4**, binário standalone |
| Interatividade | jQuery + Bootstrap JS | **HTMX + Alpine** |
| Pipeline de assets | Vite + npm | `aru view:build`, zero Node (REGRA 13) |

O `laravel/ui` é a referência de **forma**, não de dependência. Ele é o
scaffolding mais próximo do nosso caso porque é server-rendered e sem build de
JavaScript no cliente — o mesmo desenho que HTMX e Alpine entregam. O Livewire é
a outra referência válida pelo mesmo motivo: componente renderizado no servidor,
reatividade sem SPA.

Os starter kits do Laravel 12 (Livewire, React, Vue com Inertia) resolvem o
problema no cliente. O `laravel/ui` resolve no servidor, e é onde estamos.

## O que isso NÃO decide

**Não troca o CSS.** Tailwind v4 continua, com o binário standalone que a `aru`
baixa e cacheia em `~/.arandu/bin`. A REGRA 13 segue como está: build é
permitido, ser Node não é.

## O erro que gerou a versão anterior desta ADR

O proprietário disse "se baseia no `laravel/ui` com Bootstrap, projeto básico".
Isso foi lido como *adotar Bootstrap no lugar do Tailwind*, e uma ADR foi escrita
decidindo exatamente isso — com o CSS do Bootstrap vendorizado, os componentes
reescritos e o binário do Tailwind removido do `aru view:build`.

Era a leitura errada. A frase era sobre **de onde copiar a estrutura**, e a
correção veio antes de qualquer linha de código mudar: o download do Bootstrap
foi interrompido, e nada do `porang` chegou a ser reescrito.

Fica registrado porque é o tipo de erro que se repete: um nome de tecnologia
numa frase sobre estrutura foi lido como decisão de stack. A stack estava
decidida havia semanas e não estava em discussão.

## Consequências

- `aru make:auth` passa a gerar nove views em `resources/views/`, e não quatro
  arquivos em `modules/authui/` (ver ADR 0012, que fica pendente de atualização)
- as classes são Tailwind, com a mesma semântica visual das do Bootstrap:
  navbar, card, form-control, btn-primary viram o equivalente em utilitário
- `aru view:build` continua com a etapa de CSS, e continua baixando o binário
  standalone do Tailwind com verificação de SHA-256
- o `laravel/ui` entra no `plans/REFERENCIAS.md` como o que olhar para o
  scaffolding de auth

## Referências

- ADR 0019 — a estrutura do Laravel, diretório a diretório
- ADR 0020 — `kyse` no lugar do templ
- ADR 0021 — `porang` dissolvido
- `docs/14-view-e-assets.md` — a stack de view, que esta ADR não altera
- `reference_laravel/laravel/resources/views/` — as nove views, para consulta
