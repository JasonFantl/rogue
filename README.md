# Rogue
## Planned features
* Moveable Character
* Walls
* Inventory
* Items (only need to be picked up)
* Monsters
    * Movement
    * Attacking
* Health
* Field of vision
* Map memory
* Leveling
    * STR
    * AC
* Room generation
* level exit
* Inventory slots
* usable Items
* Traps
* Throwable items
* Trading
* Doors with keys
* AI Brains / economic related activities
  * example: bread
    * growing wheat
    * harvesting
    * processing
    * eating
  * example: clothing
    * sheer sheep
    * process wool
    * make clothing
  * example: 
* Dialoge with NPCs
* More terrain (water, jungle, mountains, paths, cities, ports, mines, generated islands)
* chunk loading
* More complex monster AI (Fight each other, have inventory, field of view, etc)

### Dependencies

* go (version 1.15)
    
### Installing

* must first have Go installed on computer
* download this repo
* sudo apt-get install -y libgl1-mesa-dev
* sudo apt-get install xorg-dev

### Executing program (Linux)
From the root directory
```
go build
```
then
```
./rogue
```
### handy resources
* https://www.reddit.com/r/roguelikedev/
* https://www.fatoldyeti.com/
* http://gameprogrammingpatterns.com/contents.html
* https://github.com/wadsworj/rogue
* https://austinmorlan.com/posts/entity_component_system/
