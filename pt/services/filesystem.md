# Armazenamento de arquivos

- [Disk controla o tenant antes de o adapter receber um caminho](#a-disk-owns-tenancy-before-an-adapter-sees-a-path)
- [Crie um disco e transmita um objeto](#create-a-disk-and-stream-an-object)
- [Liste, mova, baixe e assine acesso temporário](#list-move-download-and-sign-temporary-access)
- [Recuse identidade insegura em vez de reescrevê-la](#reject-unsafe-identity-instead-of-rewriting-it)
- [Use as asserções de disco no adapter configurado](#use-the-disk-assertions-against-the-configured-adapter)
- [Operações de Disk e capacidade do adapter](#disk-operations-and-adapter-capability)
- [Storage não é o filesystem comum da aplicação](#storage-is-not-the-application-s-plain-filesystem)

<a name="a-disk-owns-tenancy-before-an-adapter-sees-a-path"></a>
## Disk controla o tenant antes de o adapter receber um caminho

Disk é a fronteira da aplicação sobre Adapter. Toda leitura e escrita recebe Grant, deriva o prefixo do tenant, valida a chave lógica e só então entrega o caminho final ao adapter. Essa ordem impede que um driver esqueça o isolamento garantido por outro.

<a name="create-a-disk-and-stream-an-object"></a>
## Crie um disco e transmita um objeto

`foundation.NewDrivers` constrói `drivers.Disks` a partir de `FILESYSTEM_DISK`: local escreve sob `FILESYSTEM_ROOT`, e r2 e s3 são construídos sobre o protocolo S3 a partir do bucket, endpoint, região e chaves, com uma variável ausente recusada no boot. `Disk("")` responde o disco padrão. Put aceita um `io.Reader` para que um upload não precise ficar em memória; Get devolve informações e conteúdo do arquivo pela mesma fronteira de tenant. A chave da aplicação continua lógica e nunca inclui um prefixo de tenant.

```go example
func StoreReport(ctx context.Context, drivers *foundation.Drivers, g security.Grant, body io.Reader) error {
	disk, err := drivers.Disks.Disk("")
	if err != nil {
		return err
	}
	if err := disk.Put(ctx, g, "reports/q1.pdf", body, "application/pdf"); err != nil {
		return err
	}
	file, err := disk.Get(ctx, g, "reports/q1.pdf")
	if err != nil {
		return err
	}
	defer file.Body.Close()
	return nil
}
```

<a name="list-move-download-and-sign-temporary-access"></a>
## Liste, mova, baixe e assine acesso temporário

Files e Directories enumeram o prefixo de um tenant; Copy e Move mantêm ambas as chaves nele; Response e Download escrevem respostas HTTP; `TemporaryURL` delega acesso com expiração apenas quando o adapter o fornece. `PutFile` e `PutFileAs` recebem uploads preservando a mesma política de caminhos.

`filesystem.Serve` escreve um arquivo guardado numa resposta a partir de uma rota cuja Policy produziu o Grant. O tipo de conteúdo vem dos metadados guardados, X-Content-Type-Options é nosniff e a Content-Security-Policy é default-src 'none'; sandbox, de modo que um arquivo renderizado não executa script. SVG, HTML e XML podem levar script e saem como anexo, a menos que `ServeOptions.InlineActiveContent` peça para mostrá-los inline, ainda em sandbox.

<a name="reject-unsafe-identity-instead-of-rewriting-it"></a>
## Recuse identidade insegura em vez de reescrevê-la

Um tenant de Grant vazio, um tenant com separador, uma chave que escapa, um objeto ausente ou uma URL temporária sem suporte retornam erro. Não normalize um tenant inseguro em outro identificador: dois identificadores originais poderiam se fundir num mesmo diretório. Feche o `File.Body` retornado por Get e todo stream retornado por `ReadStream`, e propague falhas parciais de escrita. Uma escrita sob um Grant somente leitura, emitido enquanto um operador apenas visualiza uma conta, é recusada com `auth.ErrReadOnly`.

<a name="use-the-disk-assertions-against-the-configured-adapter"></a>
## Use as asserções de disco no adapter configurado

`AssertExists`, `AssertMissing`, `AssertCount` e `AssertDirectoryEmpty` verificam pela fronteira pública Disk. Adicione testes entre tenants com chaves lógicas iguais, travessia de caminhos, fechamento de streams e expiração de URLs temporárias. Execute o contrato de cada driver de produção, sem confiar só no teste local.

<a name="disk-operations-and-adapter-capability"></a>
## Operações de Disk e capacidade do adapter

Put, Get, `ReadStream`, `WriteStream`, Delete, Copy, Move, Stat, checksum, visibilidade e métodos de diretório compõem o armazenamento. `ProvidesTemporaryURLs` e `ProvidesTemporaryUploadURLs` revelam capacidades opcionais antes do uso. `GetAdapter` existe para integração; operações normais ficam em Disk para não contornar Grant.

<a name="storage-is-not-the-application-s-plain-filesystem"></a>
## Storage não é o filesystem comum da aplicação

Objetos do tenant usam Disk. Configuração do processo, assets embutidos, views geradas e outros arquivos da aplicação usam o filesystem comum, pois não são dados do cliente. Um registro de banco pode possuir a chave do objeto, mas guardar o caminho absoluto do adapter expõe detalhes de implantação e dificulta trocar drivers.
