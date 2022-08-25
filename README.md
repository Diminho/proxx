# Proxx

###Game Proxx
Aim of the game is not to be trapped in black hole and find all cells avoiding black holes.

####Iternals:
Board is represented as graph data structure. 
So it is basically grid with connected field in the north, south, west and east directions.
Seeking and revealing of cells is made with breadth first search algorithm
since it gives ability to find neighbors more efficiently than depth first search, for instance.

####Details:
Added few unit tests. Some pieces of code is not covered since it was not in the challenge,
but I decided to add at least some.

In makefile there is `localbuild` tool to build the game for different OSes.

Game is started by `./proxx start`
