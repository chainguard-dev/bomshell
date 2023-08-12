# Running `bomshell` Recipes

The plain bomshell runtime offers three ways of running recipes. In its default
mode, `bomshell` tries to understand what you are trying to run by checking
the values of the invocation. The default bomshell invocation is meant to be run
from the command line. For a more predictable, programatic way of running `bomshell`
programs see `bomshell run` below.

## The Plain `bomshell`` Command

Typing `bomshell` into your terminal invokes the _smart_ bomshell exec mode.
In this smart mode, bomshell will try to figure out what to run. It will try
to find a bomshell recipe in three ways:

1. __Recipe file.__ `bomshell` will look at the value of the first positional 
argument to check if it points to a file. If it does, it will load it and
avaluate its contents as a bomshell program.

1. __Inline bomshell code.__ If the first positional argument is not a file, 
bomshell will try to execute its contents as CEL bomshell code.

1. __`-e|--execute` Flag.__ If the `--execute` flag is defined, bomshell will
read its contents and use it as the recipe to run. All positional arguments will
be interpreted as paths of SBOMs to run into the runtime environment.

## `bomshell run`

When using bomshell in scripts or other automated environments, the recommended
subcomand is `bomshell run`. This subcomand will not try to interpret arguments
or flags. Positional arguments will always be interpreted in the same fashion:

* The first positional argument is the `bomshell` file to run.
* Any other arguments will be interpreted as SBOMs to preload and expose in the
runtime environment.

### Exmaple:

```
bomshell run myrecipe.cel sbom1.spdx.json sbom2.cdx.json
```

## Piping SBOMs into `bomshell`

The CLI supports three ways of preloading SBOMs, plus a fourth one at run time.

1. First by passing the paths of the documents in positional arguments:

```
bomshell 'sbom.files()' sbom1.cdx.json
```

2. Next, using the `--sbom` flag. The following command is equivalent to the previous
example:

```
bomshell 'sbom.files()' --sbom=sbom1.cdx.json
```

3. Finally, piping the SBOM to `bomshell`'s STDIN. Both the `run` subcommand and
the default `bomshell` mode support piping SBOMs to the runtime. Here is an example:

```
cat sbom.spdx.json | bomshell 'sbom.packages().ToDocument()'
```

### SBOM Ordering and Access at Runtime

ALl SBOMs preloaded at runtime are loaded and available at runtime in the global
`sboms` collection. This global variable is a map, integer based where all SBOMs
are stored in the order they were defined.

The order SBOMs are stored in the `sboms` array is the following:

1. SBOMs piped to `bomshell`'s STDIN. If a file was piped to the bomshell process,
it will always be `sboms[0]`.
1. SBOMs from positional arguments. The next entries in the SBOM collection will
be any SBOMs defined as positional arguments.
1. Files defined using `--sbom`. Finally any SBOMs defined using the `--sbom` flag
will be the last ones appended to the global SBOM collection.

For convenience, a global `sbom` variable is always available. It stores the 
first SBOM loaded into the runtime.
