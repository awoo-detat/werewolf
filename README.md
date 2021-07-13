# werewolf

## TODO
### Player messaging
The move from `awoo` to `werewolf` stripped out all of the messaging, because it doesn't necessarily make sense. But regardless of a message displayed via websocket or a Discord message, we need to be able to message users.

I'm thinking that we provide a channel for a Message (recreating everything I stripped out), and either provide an interface or leave the handling of the channel up to the code that uses this library. `Broadcast()` could message the game channel in a Discord, while `Message()` sends it to the player-specific mod channel?
