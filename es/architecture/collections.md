# Collections

- [Usa una collection cuando la operación exprese lo que necesitas](#use-a-collection-when-the-operation-is-the-vocabulary)
- [Construye, transforma y finaliza valores tipados](#build-transform-and-finish-typed-values)
- [La evaluación bajo demanda y los maps cambian el coste](#laziness-and-maps-change-the-cost-model)

> [!WARNING]
> **Limitaciones actuales.** Los auxiliares genéricos de collections y arrays están publicados y probados, pero ninguna aplicación oficial con tag los utiliza. Prefiere slices comunes salvo cuando la operación documentada elimine repetición en el código de la aplicación y verifica el auxiliar elegido en la aplicación consumidora.

<a name="use-a-collection-when-the-operation-is-the-vocabulary"></a>
## Usa una collection cuando la operación exprese lo que necesitas

`collections.Collection[T]` es un slice tipado con transformaciones nombradas. Es útil cuando una secuencia de filtros, agrupaciones o reducciones se lee mejor que bucles locales. La indexación común, un bucle corto y la ordenación de la biblioteca estándar siguen siendo la primera opción cuando resultan más claros.

<a name="build-transform-and-finish-typed-values"></a>
## Construye, transforma y finaliza valores tipados

Collect, Make, Empty, Range y Wrap construyen valores. Filter, Reject, First, Sole, Take, Skip y Slice los seleccionan. Map, `FlatMap`, `GroupBy`, `KeyBy` y `MapWithKeys` los transforman; Reduce, Percentage, Count, Join e Implode terminan un cálculo. La ordenación, unicidad, diferencia e intersección ocurren en memoria, sin garantías de base de datos ni autorización.

Las operaciones que cambian el tipo, como `Map[T,U]`, son funciones de paquete porque los métodos Go no pueden introducir un segundo parámetro de tipo. Las operaciones que modifican valores, como Push, Pop y Splice, usan receivers de puntero; una llamada de apariencia fluida no es automáticamente una copia independiente.

| Superficie | Contrato publicado |
| --- | --- |
| Construcción | Collect, Make, Empty, Range, Wrap, `FromJSON` |
| Selección | Filter, Reject, First, Sole, Take, Skip, Slice |
| Transformación | Map, `FlatMap`, `GroupBy`, `KeyBy`, `MapWithKeys` |
| Agregación | Reduce, Percentage, Count, Join, Implode |
| Trabajo bajo demanda | `RangeLazy`, `MapLazy`, Take, Remember, All |

```go example
func LargeAmounts() (collections.Collection[string], int) {
	amounts := collections.Collect([]int{125, 200, 75, 200})
	large := amounts.Filter(func(value, _ int) bool { return value >= 200 })
	labels := collections.Map(large, func(value, _ int) string {
		return strconv.Itoa(value)
	})
	total := collections.Reduce(amounts, func(sum, value, _ int) int {
		return sum + value
	}, 0)
	return labels, total
}
```

<a name="laziness-and-maps-change-the-cost-model"></a>
## La evaluación bajo demanda y los maps cambian el coste

Una collection bajo demanda difiere el trabajo hasta recorrerla; Remember memoriza lo obtenido. No añade cancelación de I/O, por lo que un productor bloqueante aún necesita código sensible al contexto. Los auxiliares de maps solo definen orden cuando lo establece su contrato de resultado.

collections/arr contiene auxiliares tipados de acceso y presentación en Hesape :hesapeVersion. Prefiere un bucle local cuando sea más corto y verifica el auxiliar elegido en la aplicación consumidora, porque el ecosistema etiquetado no tiene un flujo de aplicación oficial para estas collections genéricas.

- [API de collections](https://pkg.go.dev/github.com/arandu-io/hesape/collections)
- [Auxiliares de arrays](https://pkg.go.dev/github.com/arandu-io/hesape/collections/arr)
- [Concurrencia](/docs/{{version}}/architecture/concurrency)
