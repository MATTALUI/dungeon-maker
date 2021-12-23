# I know this is hacky, but I honestly wasn't planning on a full client-server
# relationship when building this, so I didn't really organize this project very
# well. This was the easiest way to keep the game classes in sync for the client
# and the server.
#
#  ¯\_(ツ)_/¯
#

echo "Clearing links"
rm *.go
echo "Making links"
ln -s ../../client/game/game.go ./game.go
ln -s ../../client/game/dungeon.go ./dungeon.go
ln -s ../../client/game/room.go ./room.go
ln -s ../../client/game/dimension.go ./dimension.go
ln -s ../../client/game/coordinates.go ./coordinates.go
ln -s ../../client/game/hero.go ./hero.go
ln -s ../../client/game/utils.go ./utils.go
ln -s ../../client/game/animatedSprite.go ./animatedSprite.go
ln -s ../../client/game/connhandlers.go ./connhandlers.go
ln -s ../../client/game/graphicalRange.go ./graphicalRange.go
ln -s ../../client/game/connectedPlayer.go ./connectedPlayer.go
ln -s ../../client/game/gameState.go ./gameState.go
ln -s ../../client/game/adventureGameState.go ./adventureGameState.go
ln -s ../../client/game/dialogState.go ./dialogState.go
ln -s ../../client/game/menuState.go ./menuState.go
ln -s ../../client/game/exitState.go ./exitState.go
ln -s ../../client/game/mapState.go ./mapState.go
echo "The deed is done"