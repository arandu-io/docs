# Pint

**Veredito:** NÃO EXISTE EM GO
**Nome Arandu:** —

## O que o Laravel faz

Formatador de código PHP com opinião, sobre o php-cs-fixer, com presets.

## Decisão

Não construir. Revisão em relação ao mapa anterior — havia sido nomeado `katu` por engano.

## Por quê

O Pint existe porque o PHP não tem formatador canônico: existem PSR-12, php-cs-fixer, phpcbf, e cada projeto escolhe e configura o seu.

Go tem `gofmt` embutido na toolchain desde o dia um. Não há configuração, não há debate de estilo, todo projeto Go do mundo é formatado igual. Um formatador nosso seria pior que o padrão e ninguém usaria.

O que sobra de útil — checar as **regras de arquitetura**, não o estilo — já é o `aru doctor`, que verifica se todo repositório tem policy, se handler não importa `data`, se não há SQL montado com Sprintf.

## O que conteria, se existir

Nada. `gofmt`, `go vet` e `aru doctor` cobrem o espaço.

## Gatilho para construir

Nenhum.

## Alternativa enquanto não existe

`gofmt -w .`, `go vet ./...`, `aru doctor`.
