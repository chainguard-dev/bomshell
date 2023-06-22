# Functions Inventory

This is an initial inventory of the functions available in each of the three
objects currently exposed to the CEL environment:

| method | Description | SBOM | NodeList | Node |
| --- | --- | --- | --- | --- |
| files() | Returns all nodes that are files | ✔️ | | |
| packages() | Returns all nodes that are files | ✔️ | | |
| add() | Combines nodelists or nodes into a single nodelist | ✔️ | | |
| toNodeList() | Converts a Node to a NodeList | N/A | N/A | ✔️ |
| nodeByID() | Converts a Node to a NodeList | ✔️ |  |  |
| nodesByName() | Returns all elements whose name matches | | | |
| nodesByPurlType() | Returns all elements whose purl is of a certain type | | | |
| addAt() | Inserts a nodelist or node at a point | | | |
| nodesByDepth() | | | N/A |
