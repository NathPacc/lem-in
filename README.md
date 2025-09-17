# Lem-in

**Lem-in** is a project focused on finding the optimal movement strategy for ants within a predefined colony. The objective is to move all ants from a designated starting room to an ending room in the fewest possible turns. Each ant can move one room per turn, and no room (except the start and end) can contain more than one ant at a time.

---

### Algorithm Overview

The program follows a multi-step algorithm:

- Finds all possible paths from the start to the end room.
- Removes redundant paths. For example, if there's a path `start -> 1 -> end`, then a longer path like `start -> 1 -> 2 -> end` becomes unnecessary.
- Identifies all possible sets of independent paths (paths that only share the start and end rooms).
- Calculates how many turns are needed to move all ants optimally through each set of paths.
- Outputs the ant movements using the fastest set of paths.

---

### Project Structure

LEM-IN/
├── colony/               # Heart of the program
│   ├── algo.go           # Main algorithm
│   ├── prints.go         # Printing the different structs and the resolution
│   └── setup.go          # Initializing datas and creating the colony
│
├── datas/                # Dealing with the recovering and verification of the datas 
│   ├── datas.go          # Recovering the datas
│   └── errors.go         # Verifying the datas
│
├── files/                # Entry files describing the colony
│
├── modules/              
│   └── structColony.go   # Declaration of structs used by algo
│   └── modules.go        # Declaration of structs used for data recovering
│
├── go.mod                 
├── main.go               # Execution of the program
└── README.md             

---

### Usage

Prepare an input file (e.g., `yourfile.txt`) with the following format:

1. **Number of ants**  
   An integer greater than 0 representing how many ants will traverse the colony.

2. **Room definitions**  
   Each room is defined by a name and its coordinates:  
   `room_name x y`

3. **Link definitions**  
   Connections between rooms are defined by two room names separated by a dash:  
   `room1-room2`

4. **Special lines**  
   - `##start` indicates the next room is the starting room.  
   - `##end` indicates the next room is the ending room.  
   - Lines starting with `#` are comments and will be ignored.

**Example input:**

9                           //Number of ants
##start
start 1 24                 //Starting room's name is "start" and it's coordinates are (1,24)
0 2 3                      //Room "0" have (2,3) as coordinates
hello_world 0 0            //Room hello_world have (0,0) as coordinates
##end
queen 5 9                 //Ending room's name is "queen" and it's coordinates are (5,9)
start-0                   //Link beetween rooms start and 0
0-hello_world             //Link beetween rooms hello_world and 0
start-queen               //Link beetween rooms start and queen
queen-hello_world         //Link beetween rooms queen and hello_world


Once the file is correctly placed, type "./lem-in yourfile.txt in your terminal.

### Results

The output shows the movement of ants per turn. Each line represents one turn.  
Each movement is formatted as: `L<AntID>-<RoomName>`

Exemple :
L1-2
L1-3 L2-2
L1-end L2-3 L3-2
L2-end L3-3 L4-2 
L3-end L4-3
L4-end

### Author

Nathan PACCOUD - Program created during my formation in Zone01 Rouen.