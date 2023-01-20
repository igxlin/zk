# zk

A commandline tool to help organize markdown notes.

By default, the current working directory is taken as the root directory
of markdown notes.

Created for personal use. **USE IT AT YOUR OWN RISK.**

Install:

```
make install
```

Usage:

- `zk mv path/to/source path/to/target`: The same as `mv`. However,
  backlinks are updated together.

  - **The source and target paths cannot be directories.**
