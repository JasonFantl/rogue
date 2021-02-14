# Rogue
## TODO
* ECS
    * Components
        * Display
        * Position
        * Controller
        * Blocking
        * Movable
        * Handshake
    * Systems
        * Display
        * Movement
        * Handshake
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
* Shop keeper
* More complex monster AI (Fight each other, have inventory, field of view, etc)

### Dependencies

* go (version 1.15)
    
### Installing

* git pull this and try to run the dependencies below

* might have to run  ``` go build ```
* if that doesnt work, try ``` go get ```
* and finally if that doesnt work, try ``` go get github.com/gdamore/tcell/v2 ```

### Executing program
From the root directory
```
go run *.go
```

### handy resources
* https://www.reddit.com/r/roguelikedev/
* https://www.fatoldyeti.com/
* http://gameprogrammingpatterns.com/contents.html
* https://github.com/wadsworj/rogue
* https://austinmorlan.com/posts/entity_component_system/
