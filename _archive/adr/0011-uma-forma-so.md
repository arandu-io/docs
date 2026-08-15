# ADR 0011-uma-forma-so — Uma forma só de fazer cada coisa

> Renumerado de 0005 para 0011: o número 0005 já estava publicado com o ADR de
> idioma do código, referenciado pelo `CLAUDE.md` e pelo `00-indice.md`.

**Status:** Aceita.

## Contexto

O Laravel acumulou três camadas de frontend (Blade, Livewire, Inertia), três
runtimes no Octane (Swoole, RoadRunner, FrankenPHP) e dependência de Node dentro
do core — inclusive um build npm para a própria página de erro. Ver
`docs/09-uma-forma-so.md` para as evidências no tree de referência.

Cada adição foi resposta legítima a um pedido real. O custo aparece somado: o
desenvolvedor novo escolhe três vezes antes da primeira linha, e o mantenedor
paga matriz de compatibilidade entre todas as combinações.

## Decisão

Uma única forma oficial para cada dimensão: view, runtime, autorização, dados,
build e config. Sem alternativa, sem modo, sem "opcional para quem precisa".

Teste de admissão de qualquer proposta: **cria uma segunda maneira de fazer algo
que já tem uma?** Se sim, é rejeitada.

## Consequência

Perdemos flexibilidade e perderemos usuários que queriam a outra opção. Em troca,
a curva de aprendizado tem um caminho só, a matriz de teste não explode, e o
gerador de código tem um alvo único — que é a condição para a Fase 4 existir.

Revisitar apenas com evidência de que a ausência da alternativa está custando
adoção medida, nunca com pedido isolado.
