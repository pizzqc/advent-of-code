# --- Day 12: Hill Climbing Algorithm ---
You try contacting the Elves using your handheld device, but the river you're following must be too low to get a decent signal.

You ask the device for a heightmap of the surrounding area (your puzzle input). The heightmap shows the local area from above broken into a grid; the elevation of each square of the grid is given by a single lowercase letter, where a is the lowest elevation, b is the next-lowest, and so on up to the highest elevation, z.

Also included on the heightmap are marks for your current position (S) and the location that should get the best signal (E). Your current position (S) has elevation a, and the location that should get the best signal (E) has elevation z.

You'd like to reach E, but to save energy, you should do it in as few steps as possible. During each step, you can move exactly one square up, down, left, or right. To avoid needing to get out your climbing gear, the elevation of the destination square can be at most one higher than the elevation of your current square; that is, if your current elevation is m, you could step to elevation n, but not to elevation o. (This also means that the elevation of the destination square can be much lower than the elevation of your current square.)

For example:
```
Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi
```

Here, you start in the top-left corner; your goal is near the middle. You could start by moving down or right, but eventually you'll need to head toward the e at the bottom. From there, you can spiral around to the goal:

```
v..v<<<<
>v.vv<<^
.>vv>E^^
..v>>>^^
..>>>>>^
```

In the above diagram, the symbols indicate whether the path exits each square moving up (^), down (v), left (<), or right (>). The location that should get the best signal is still E, and . marks unvisited squares.

This path reaches the goal in 31 steps, the fewest possible.

**What is the fewest steps required to move from your current position to the location that should get the best signal?**

# Solutions

* At first I went with a DFS (Depth First Search) and dealing with the response back as the stack un-pilled but it was very messy and I couldn't figure out how to find the shortest... only able to find a path
* The proper approach was to use a BFS (Breath First Search) and keeping track of all the distance in a seperate matrix holding the values the distance to reach a given point in the matrix from a starting location

# Solution Part2

* The first solution was to calculate the shortest path from all "a" location which worked but was not very efficient -- It was still finishing quickly but I realized later there is a much better way...
* The 2nd approach I took on part-2 was to reverse the BFS algorithm to start from the END until we reach an "a" and then returning the shortest distance of all the "a" found along the way... reducing the O(time) complexity significantly since we'll only run the BFS once instead of n*times.
  * I haven't try it but might be a way to exit the BFS as soon as we find the first "a" instead of traversing the rest of the matrix until all "a" were computed.
