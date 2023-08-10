# Example bomshell programs

This directory contains a number of example programs that can be used to 
understand how bomshell works. Each program is simple enough to demo a
single feature of bomshell.

The following list has a summary of the files, open each `.cel` file to 
read the full documentation of the example and instructions on how to run it.

## Example List

Example SBOMs used to run these examples are also found in this directory.

| File | Description | 
| --- | --- |
| [compose.cel](compose.cel) | Example of SBOM composition using `RelateNodeListAtID()` |
| [files.cel](files.cel) | Generate a new SBOM containing only the files found in an SBOM. |
| [packages.cel](packages.cel) | Generate a new SBOM containing only the packages found in an SBOM. |
| [loadsbom.cel](loadsbom.cel) | Demo of SBOM loading directly from bomshell. |
| [nodesbypurltype.cel](nodesbypurltype.cel) | Example showing how to extract all nodes of a certain purl type. |
| [getnodebyid.cel](getnodebyid.cel) | Demo querying an SBOM for a specific node. |

If you'd like to see more examples here file an issue, we'de be happy to create
more!
