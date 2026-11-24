# Manipulación de imágenes

- [Decodifique solo tras conocer límites de bytes y píxeles](#decode-only-after-byte-and-pixel-limits-are-known)
- [Inspeccione, decodifique, transforme y codifique una imagen](#inspect-decode-transform-and-encode-one-image)
- [Elija transformaciones deterministas y metadatos de salida](#choose-deterministic-transforms-and-output-metadata)
- [Rechace bombas de descompresión y formatos no compatibles](#reject-decompression-bombs-and-unsupported-formats)
- [Use fixtures reales en cada límite configurado](#use-real-fixtures-at-every-configured-ceiling)
- [Frontera del componente de imágenes publicado](#released-image-component-boundary)
- [Almacene los bytes generados mediante un contrato de filesystem](#store-generated-bytes-through-a-filesystem-contract)

<a name="decode-only-after-byte-and-pixel-limits-are-known"></a>
## Decodifique solo tras conocer límites de bytes y píxeles

El componente de imágenes publicado decodifica, transforma y codifica formatos comunes mediante opciones tipadas. El esqueleto no lo conecta; cada aplicación enlaza por su cuenta upload, límites, almacenamiento y entrega, y los demuestra juntos.

<a name="inspect-decode-transform-and-encode-one-image"></a>
## Inspeccione, decodifique, transforme y codifique una imagen

Limite bytes del upload en la aplicación, reduzca el techo de píxeles del manager y luego oriente, recorte, elija un encoder publicado y solicite los bytes finales con contexto. La biblioteca limita píxeles decodificados; la aplicación aún define límites de bytes y salida solicitada.

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
## Elija transformaciones deterministas y metadatos de salida

Normalice la orientación una vez, elimine metadatos innecesarios, elija recorte o ajuste explícitamente, genere claves de almacenamiento independientes de nombres de archivo y escriba mediante el adapter de filesystem solo tras completar toda la transformación.

- Limite bytes de entrada antes de `FromBytes`.
- Defina un límite de píxeles decodificados.
- Valide la geometría de salida solicitada en el código de aplicación.
- Use JPEG, PNG o GIF con el driver predeterminado.

<a name="reject-decompression-bombs-and-unsupported-formats"></a>
## Rechace bombas de descompresión y formatos no compatibles

Headers malformados, formato no compatible, píxeles excesivos, entrada truncada, fallo de encoder y de almacenamiento son resultados distintos. Nunca sirva el original sin validar cuando la transformación lo rechace.

> [!WARNING]
> El tamaño comprimido no es un límite seguro de memoria. Aplique el techo publicado de píxeles decodificados antes de asignar búferes completos y ensaye los límites elegidos con uploads concurrentes.

<a name="use-real-fixtures-at-every-configured-ceiling"></a>
## Use fixtures reales en cada límite configurado

Use fixtures válidos pequeños, dimensiones en el límite, un píxel por encima, datos truncados y políglotas, animación cuando no se admita, metadatos de orientación, transparencia, resize determinista, fallo de encoder y perfil de memoria concurrente.

- Un exceso de píxeles declarados falla antes de asignar el canvas.
- Las dimensiones y el tipo de medio de salida cumplen el contrato.
- Los límites de bytes de entrada y salida solicitada se prueban en la aplicación.
- Una transformación fallida no escribe un objeto en almacenamiento.

<a name="released-image-component-boundary"></a>
## Frontera del componente de imágenes publicado

`NewImageManager`, `MaxPixels`, FromBytes/FromStream, Orient, Cover, Contain, Resize, Scale, `ToJpg/ToPng/ToGif` y `ToBytes` forman el driver estándar publicado. Nombres de formatos como `WebP` se reconocen, pero no tienen encoder estándar y deben fallar, sin generar bytes con etiqueta incorrecta.

| Superficie | Contrato publicado |
| --- | --- |
| `MaxPixels` / `FromBytes` | decodificación con límite de píxeles |
| Cover / Contain / Resize | geometría determinista |
| `ToJpg` / `ToPng` / `ToGif` | encoders estándar publicados |
| `ToBytes` / Dimensions | salida y metadatos vinculados al contexto |

<a name="store-generated-bytes-through-a-filesystem-contract"></a>
## Almacene los bytes generados mediante un contrato de filesystem

Las peticiones HTTP limitan uploads, los filesystems almacenan el objeto final, la autorización protege la propiedad y el despliegue dimensiona la memoria del proceso.

- [Peticiones HTTP](/docs/{{version}}/http/requests)
- [Sistema de archivos](/docs/{{version}}/services/filesystem)
- [Autorización](/docs/{{version}}/security/authorization)
- [Despliegue](/docs/{{version}}/operations/deployment)
