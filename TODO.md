# TODO

## High Level

- finish the grid actor scheduler 
  - if it's generic enough then move this package down stream to `lytics/grid`
- add index actors - Add a required set of actors to manage each shard of an index. 
    1. as a proof of concept using a single actor per index
    2. add a query path for the actor.
    3. shard out the index over N actors, 
       - add packages for shard management. 

## Low Level

- Refactor actor creation and management 
  - Create a wrapper around the grid actor interface that adds some DEFAULT features to all GUI actors. Features:
    - HealthCheck messages (ping/pong): the wrapper should add a method to call actores that allow them to be pinged to see if they are active and running.
- The current state machine logic in PeerState is pretty inflexible.  Refactor the PeerState logic to create a DFA state machine simlar to the way `lytics/dfa` is used to manage the actor's own life cycle, this new dfa libary should be useful for the scheduler to manage the life cycle.
- Move the relocation tests from the actor-pool package to the relocation package.
