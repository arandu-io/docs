# ADR 0003-modulos-verticais — Módulos verticais em vez de camadas

**Status:** Superseded pela ADR 0019. O `modules/<nome>/` vertical caiu em
agosto/2026: a arvore passou a ser a do Laravel, diretorio a diretorio.
Mantida como registro do porque a decisao original fazia sentido.

## Contexto

Ver docs/01-arquitetura.md.

## Decisão

Uma feature vive em um diretório. Habilita distribuição de módulos e reduz o custo de entender uma feature de seis arquivos espalhados para um diretório.

## Consequência

Revisitar apenas com dado de uso real, não com preferência de estilo.
