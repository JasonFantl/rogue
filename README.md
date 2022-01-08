# Rogue

## Instructions

Every action will cause the game to take a step forward.

Use `WASD` to move the character and move the curser in menu. 

Some things you interact with by simply moving onto the objects tile, for example
  * Attacking
  * Unlocking doors

Use `E` to interact with the world, for example
  * Pick up items
  * Initiate trading with another entity
  * Confirm buttons in menu
  * Throw item (after choosing to throw it)

Use `Q` to open and close the menu.

Use `Esc` to exit the game.

## Fun Notes

* Monsters only pursue treasure, so they will only attack you if you are carrying any.
* If you find two monsters and treasure, they will fight each other to the death for it.
* Once you die, the perspective switches to a different creature, which brain you are now fighting with for control of its body.
* Since damage is dealt when a moving weapon hits a creature (this makes throwing useful); and trading works by: unequipping items, moving them to the other person, then they equip them; then when you trade a sword the other person takes damage since the game thinks you threw it at them.
## Planned features

* level exit
* Traps
* Dialogue with NPCs
* More terrain (water, jungle, mountains, paths, cities, ports, mines, generated islands)
* chunk loading

### Dependencies

* go (1.15 or above)
* sudo apt-get install libgl1-mesa-dev xorg-dev
### Installing

* Download this repo
* Compile with `go build`
* Run with `./rogue`

### handy resources
* https://www.reddit.com/r/roguelikedev/
* https://www.fatoldyeti.com/
* http://gameprogrammingpatterns.com/contents.html
* https://github.com/wadsworj/rogue
* https://austinmorlan.com/posts/entity_component_system/
