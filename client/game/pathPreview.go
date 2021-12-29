package game

type PathPreview struct {
	CurrentRoom *Room;
	NextRoom *Room;
	PreviousRoom *Room;
	IsTarget bool;
}