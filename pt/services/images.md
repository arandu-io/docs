# Manipulação de imagens

- [Decodifique só após conhecer limites de bytes e pixels](#decode-only-after-byte-and-pixel-limits-are-known)
- [Inspecione, decodifique, transforme e codifique uma imagem](#inspect-decode-transform-and-encode-one-image)
- [Escolha transformações determinísticas e metadados de saída](#choose-deterministic-transforms-and-output-metadata)
- [Recuse bombas de descompressão e formatos não suportados](#reject-decompression-bombs-and-unsupported-formats)
- [Use fixtures reais em cada limite configurado](#use-real-fixtures-at-every-configured-ceiling)
- [Fronteira do componente de imagens publicado](#released-image-component-boundary)
- [Armazene os bytes gerados por um contrato de filesystem](#store-generated-bytes-through-a-filesystem-contract)

> [!WARNING]
> **Limitações atuais.** Limites de decodificação e transformações com contexto estão publicados, mas nenhum fluxo de upload de primeira parte consome o pipeline. Valide tipo de mídia e bytes antes de decodificar e armazene resultados apenas por disco com escopo de tenant.

<a name="decode-only-after-byte-and-pixel-limits-are-known"></a>
## Decodifique só após conhecer limites de bytes e pixels

O componente de imagens publicado decodifica, transforma e codifica formatos comuns por opções tipadas. Nenhum consumidor atual em tag o conecta; este guia limitado documenta a fronteira do pacote e exige que cada aplicação comprove juntos upload, limites, armazenamento e entrega.

<a name="inspect-decode-transform-and-encode-one-image"></a>
## Inspecione, decodifique, transforme e codifique uma imagem

Limite bytes do upload na aplicação, reduza o teto de pixels do manager e então oriente, recorte, escolha um encoder publicado e solicite os bytes finais com contexto. A biblioteca limita pixels decodificados; a aplicação ainda define limites de bytes e de saída solicitada.

```go example
func Thumbnail(ctx context.Context, contents []byte, maxPixels int) ([]byte, error) {
	images := image.NewImageManager().MaxPixels(maxPixels)
	img := images.FromBytes(contents)
	encoded, err := img.Orient().Cover(640, 640).ToJpg().Quality(82).ToBytes(ctx)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}
```

<a name="choose-deterministic-transforms-and-output-metadata"></a>
## Escolha transformações determinísticas e metadados de saída

Normalize orientação uma vez, remova metadados desnecessários, escolha recorte ou ajuste explicitamente, gere chaves de armazenamento independentes de nomes de arquivo e escreva pelo adapter de filesystem apenas após sucesso de toda a transformação.

- Limite bytes de entrada antes de `FromBytes`.
- Defina teto de pixels decodificados.
- Valide a geometria de saída solicitada no código da aplicação.
- Use JPEG, PNG ou GIF com o driver padrão.

<a name="reject-decompression-bombs-and-unsupported-formats"></a>
## Recuse bombas de descompressão e formatos não suportados

Headers malformados, formato não suportado, pixels excessivos, entrada truncada, falha de encoder e de armazenamento são resultados distintos. Nunca sirva o original sem validação quando a transformação o recusar.

> [!WARNING]
> Tamanho comprimido não é limite seguro de memória. Aplique o teto publicado de pixels decodificados antes de alocar buffers completos e ensaie os limites escolhidos sob uploads concorrentes.

<a name="use-real-fixtures-at-every-configured-ceiling"></a>
## Use fixtures reais em cada limite configurado

Use fixtures válidos pequenos, dimensões no limite, um pixel acima, dados truncados e políglotas, animação quando não suportada, metadados de orientação, transparência, resize determinístico, falha de encoder e perfil de memória concorrente.

- Excesso de pixels declarados falha antes da alocação do canvas.
- Dimensões e tipo de mídia da saída cumprem o contrato.
- Limites de bytes de entrada e saída solicitada são testados pela aplicação.
- Uma transformação falha não grava objeto no armazenamento.

<a name="released-image-component-boundary"></a>
## Fronteira do componente de imagens publicado

`NewImageManager`, `MaxPixels`, FromBytes/FromStream, Orient, Cover, Contain, Resize, Scale, `ToJpg/ToPng/ToGif` e `ToBytes` compõem o driver padrão publicado. Nomes de formatos como `WebP` são reconhecidos, mas não têm encoder padrão e devem falhar, sem gerar bytes com rótulo incorreto.

| Superfície | Contrato publicado |
| --- | --- |
| `MaxPixels` / `FromBytes` | decodificação com teto de pixels |
| Cover / Contain / Resize | geometria determinística |
| `ToJpg` / `ToPng` / `ToGif` | encoders padrão publicados |
| `ToBytes` / Dimensions | saída e metadados vinculados ao contexto |

<a name="store-generated-bytes-through-a-filesystem-contract"></a>
## Armazene os bytes gerados por um contrato de filesystem

Requisições HTTP limitam uploads, filesystems armazenam o objeto final, autorização protege a propriedade e implantação dimensiona a memória do processo.

- [Requisições HTTP](/docs/{{version}}/http/requests)
- [Sistema de arquivos](/docs/{{version}}/services/filesystem)
- [Autorização](/docs/{{version}}/security/authorization)
- [Implantação](/docs/{{version}}/operations/deployment)
