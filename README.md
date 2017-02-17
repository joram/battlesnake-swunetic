# battlesnake-swunetic
Genetically training weight of snake heuristics

# Setup
run `swu services start` to get the game server running http://localhost:4000

run `swu run` to build and run the snake server at http://localhost:9000

*NOTE:* when you put your snake url in to the game server, it needs to be your public IP not localhost.

# Overview

## Comparing Snakes
This will, given a list of heuristics, run games with snakes of various weights against each other, and find a winner. For a list of snake configurations (weights for the heuristics), it will run 100 games, and return the count of wins for each of the configurations.

## Evolution
Given 100 sample games run with a set of snakes, a set of new snakes will be generated based on the winner(s) of the 100 games.
A certain amount of random varience will be introduced as well. This second generation of snakes, in theory, should be better than the first set of snakes.

## Heuristics
These are the core of this logic, the better the heuristics the better the snake can be. A heuristic is just a simple function that takes a game board and returns a direction.

### Suggested Heuristics (please add)
- move to food
- move to center
- move away from center (towards edge)
- move to increase board control
- move to other snake's head
- move away other snake's head
- move to other snake's tail
- move away other snake's tail

## Repetition
Iterating on these snakes should give us better and better snake... profit?

![Image to make Noah cringe](http://library.missouri.edu/exhibits/eugenics/exhibit_images/800px/eugenics_tree_1921.jpg)
