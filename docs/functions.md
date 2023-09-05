# Functions Inventory

This is an initial inventory of the functions planned for each of the three
objects currently exposed to the bomshell runtime environment:

| method | Return Value | Description | SBOM | NodeList | Node |
| --- | --- | --- | --- | --- | --- |
| files() | NodeList | Returns all nodes that are files | ✔️ | TBD | TBD |
| packages() | NodeList | Returns all nodes that are files | ✔️ | TBD | TBD |
| <td colspan="6">__Node Querying Functions__</td> |
| nodeByID() | Node | Returns the node with the matching identifier | ✔️ | TBD | TBD |
| nodesByName() | NodeList | Returns all elements whose name matches | TBD | TBD | TBD |
| nodesByPurl() | NodeList | Returns all elements with a matching purl | TBD | TBD | TBD |
| nodesByPurlType() | NodeList | Returns all elements whose purl is of a certain type | ✔️ | TBD | TBD |
| nodesByDepth() | NodeList | Returns nodes at X degrees of separation from the root  | TBD | TBD | TBD |
| <td colspan="6">__Graph Fragment Querying Functions__</td> |
| graphByID() | NodeList | Returns the graph fragment of a Node that matches | TBD | TBD | TBD |
| graphByName() | NodeList | Returns the graph fragment of elements whose name match the query | TBD | TBD | TBD |
| graphByPurl() | NodeList | Returns the graph of all elements with a matching purl | TBD | TBD | TBD |
| graphByPurlType() | NodeList | Returns all elements whose purl is of a certain type | TBD | TBD | TBD |
| graphByDepth() | NodeList | Returns graph fragments starting at X degrees of separation from the root  | TBD | TBD | TBD |
|<td colspan="6">__Element Transformation__</td>|
| toNodeList() | NodeList | Returns a NodeList from the object | TBD | TBD | ✔️ |
|<td colspan="6">__Composition Functions__</td> |
| add() | NodeList | Combines nodelists or nodes into a single nodelist | ✔️ | TBD | TBD |
| union() | NodeList | Returns a new nodelist with elements in common | ✔️ | TBD | TBD |
| union() | NodeList | Returns the nodes from a nodelist not present in the second | ✔️ | TBD | TBD |
| relateAt() | NodeList | Inserts a nodelist or node at a point | TBD | TBD | N/A |
