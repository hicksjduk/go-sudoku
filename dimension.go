package sudoku


type DimensionChange struct {
	source   Dimension
	definite ValueSet
	possible ValueSet
}

type Dimension interface {
	addListener(listener Dimension)
	makeDefinite(values ValueSet)
	makeImpossible(values ValueSet)
	possibleValues() ValueSet
	impossibleValues() ValueSet
}

type Square struct {
	possibleValues ValueSet
}
