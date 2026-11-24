# File storage

- [A Disk owns tenancy before an adapter sees a path](#a-disk-owns-tenancy-before-an-adapter-sees-a-path)
- [Create a disk and stream an object](#create-a-disk-and-stream-an-object)
- [List, move, download, and sign temporary access](#list-move-download-and-sign-temporary-access)
- [Reject unsafe identity instead of rewriting it](#reject-unsafe-identity-instead-of-rewriting-it)
- [Use the disk assertions against the configured adapter](#use-the-disk-assertions-against-the-configured-adapter)
- [Disk operations and adapter capability](#disk-operations-and-adapter-capability)
- [Storage is not the application's plain filesystem](#storage-is-not-the-application-s-plain-filesystem)

<a name="a-disk-owns-tenancy-before-an-adapter-sees-a-path"></a>
## A Disk owns tenancy before an adapter sees a path

Disk is the application boundary over a storage Adapter. Every read and write takes a Grant, derives the tenant prefix, validates the logical key, and only then hands a finished path to the adapter. That order prevents one driver from forgetting the isolation another driver enforces.

<a name="create-a-disk-and-stream-an-object"></a>
## Create a disk and stream an object

`foundation.NewDrivers` builds `drivers.Disks` from `FILESYSTEM_DISK`: local writes under `FILESYSTEM_ROOT`, and r2 and s3 are built over the S3 protocol from the bucket, endpoint, region, and keys, with a missing variable refused at boot. `Disk("")` answers the default disk. Put accepts an `io.Reader` so an upload need not be buffered in memory; Get returns file information and content through the same tenant boundary. The application key remains logical and never includes a tenant prefix.

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
## List, move, download, and sign temporary access

Files and Directories enumerate one tenant's prefix; Copy and Move keep both keys inside it; Response and Download write HTTP responses; `TemporaryURL` delegates expiring access only when the adapter provides it. `PutFile` and `PutFileAs` accept upload values while preserving the same path policy.

`filesystem.Serve` writes a stored file to a response from a route whose Policy produced the Grant. The content type comes from the stored metadata, X-Content-Type-Options is nosniff, and the Content-Security-Policy is default-src 'none'; sandbox, so a rendered file runs no script. SVG, HTML, and XML can carry script and go out as an attachment unless `ServeOptions.InlineActiveContent` asks to show them inline, still sandboxed.

<a name="reject-unsafe-identity-instead-of-rewriting-it"></a>
## Reject unsafe identity instead of rewriting it

An empty Grant tenant, separator-bearing tenant, escaping key, missing object, or unsupported temporary URL returns an error. Do not normalize an unsafe tenant into a different identifier: two original identifiers could collapse into one directory. Close `File.Body` returned by Get and every stream returned by `ReadStream`, and propagate partial write failures. A write under a read-only Grant, issued while an operator only views an account, is refused with `auth.ErrReadOnly`.

<a name="use-the-disk-assertions-against-the-configured-adapter"></a>
## Use the disk assertions against the configured adapter

`AssertExists`, `AssertMissing`, `AssertCount`, and `AssertDirectoryEmpty` check through the public Disk boundary. Add cross-tenant tests with identical logical keys, path traversal cases, stream close behavior, and temporary URL expiry. Run the adapter contract for each production driver rather than trusting a local-only test.

<a name="disk-operations-and-adapter-capability"></a>
## Disk operations and adapter capability

Put, Get, `ReadStream`, `WriteStream`, Delete, Copy, Move, Stat, checksum, visibility, and directory methods make up the storage surface. `ProvidesTemporaryURLs` and `ProvidesTemporaryUploadURLs` reveal optional adapter behavior before use. `GetAdapter` exists for integration boundaries, but ordinary application operations stay on Disk so Grant enforcement is not bypassed.

<a name="storage-is-not-the-application-s-plain-filesystem"></a>
## Storage is not the application's plain filesystem

Tenant-owned objects use Disk. Process configuration, embedded assets, generated views, and other application files use the ordinary filesystem because they are not customer data. A database record may own an object key, but storing the absolute adapter path in that record leaks deployment detail and makes driver changes harder.
