# 💣🐚 bomshell

An SBOM query language and associated utilities to work with data in any format.

`bomshell` is a runtime environment designed to evaluate expressions, called 
_recipes_, that operate on the SBOM graph. bomshell recipes can extract,
rearrange and remix data from SBOMs in any format, making SBOM composition a
reality. 

### __⚠️ Experimental Notice ⚠️__

`bomshell` is evolving rapidly but it should still be considered pre-release software. The language
is still incomplete and changing constantly.

## SBOM Querying and Remixing Examples

In essence, a bomshell invocation parses a set of SBOMs and executes a recipe.
At runrime, the preloaded SBOMs are accesible to the running program from the
bomshell environment. For more details be sure to check out the 
[`bomshell` tutorial](tutorial/) and the 
[examples directory](examples/).

### Extract Files and Packages from an SBOM

This example reads an SBOM, extracts its files and returns a new document 
with no packages, only those files:

```
bomshell -e 'sbom.files().ToDocument()' mysbom.spdx.json
```

This recipe the same but with nodes that are package data:

```
bomshell -e 'sbom.packages().ToDocument()' mysbom.spdx.json
```

### Multiformat Support

`bomshell` can read any SBOM format (that `protobom` supports). By default,
output is written as SPDX 2.3 but it can also be rendered to any format:

```
bomshell --document-format="application/vnd.cyclonedx+json;version=1.4" \
         --execute 'sbom.packages().ToDocument()' mysbom.spdx.json
```

Reading an SBOM into bomshell and writing it to another format essentially 
converts it into another format:

```
bomshell --document-format="application/vnd.cyclonedx+json;version=1.4" \
         --execute 'sbom' mysbom.spdx.json
```

### Quering SBOM Data

bomshell is still very young 👶🏽 but it already offers a few functions and methods
to query SBOM data. The following example extracts all go packages from an SBOM:

```
bomshell -e 'sbom.NodesByPurlType("golang")' mysbom.spdx.json 
```

Specific nodes can be looked up by ID too:

```
bomshell -e 'sbom.NodeByID("com.github.kubernetes-kubectl")' mysbom.spdx.json
```

### SBOM Composition

Loaded SBOMs are accessible through the `sbom[]` array. Nodes in
a document can be augmented or replaced. New graph sections can 
be remixed into a point in a document graph.

The following recipe extracts the npm packages from one SBOM and 
remixes them as dependencies of a binary in the other:

```
bomshell -e 'sbom[0].RelateNodeListAtID(sbom[1].NodesByPurlType("npm"), "my-binary", "DEPENDS_ON)' \
         --sbom=sbom1.spdx.json \
         --sbom=sbom2.cdx.json 
```

Note in the previous example that each SBOM is in a different format. Remixing
from different makes `bomshell` a powerful tool to work with any SBOM, tools can specialize in what they do best and bomshell
can compose documents assembled from multiple sources of
data.

## The `bomshell` Core

bomshell recipes are written in CEL 
([Common Expression Language](https://github.com/google/cel-spec))
making the runtime small and embeddable in other applications.

The backing library of Bomshell is 
[`protobom` the universal Software Bill of Materials I/O library ](https://github.com/bom-squad/protobom).
The bomshell runtime reads SBOMs and exposes the protobom
data graph to the CEL environment, emulating some methods and adding
some of its own.

Just as its core components, bomshell is open source, released under the
Apache 2.0 license.
