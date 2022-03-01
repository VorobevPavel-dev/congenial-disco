package parser

// parsingStep is a current state of parsing machine.
// Depending on it function Handle will choose a way to treat first
// N tokens in input array
type parsingStep int

const (
	// CREATE TABLE branch

	stepTableKeyword = iota
	stepTableName
	stepColumnName
	stepColumnType
	stepEngineKeyword
	stepEngineName
	stepOrderKeyword
	stepByKeyword
	stepOrderColumn
	// INSERT INTO branch

	stepInsIntoKeyword
	stepInsTableName
	stepInsColsetName
	stepInsValuesKeyword
	stepInsValuesetOpenBracket
	stepInsValuesetValue
	stepInsValuesetCloseBracket
	// SELECT branch

	stepSelColName
	stepSelFromKeyword
	stepSelTableName
	stepSelWhereKeyword
	stepSelCondColName
	stepSelCondSymbol
	stepSelCondValue
	// SHOW CREATE branch

	stepShCreateKeyword
	stepShTableName

	stepEnd
)
