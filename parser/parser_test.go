package parser

import (
	"testing"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/lexer"
	"github.com/LissaGreense/GO4SQL/token"
)

func TestParserCreateCommand(t *testing.T) {
	tests := []struct {
		input               string
		expectedTableName   string
		expectedColumnNames []string
		expectedColumTypes  []token.Token
	}{
		{"CREATE TABLE 	TBL( ONE TEXT );", "TBL", []string{"ONE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}}},
		{"CREATE TABLE 	TBL( ONE TEXT,  TWO TEXT, THREE INT);", "TBL", []string{"ONE", "TWO", "THREE"}, []token.Token{{Type: token.TEXT, Literal: "TEXT"}, {Type: token.TEXT, Literal: "TEXT"}, {Type: token.INT, Literal: "INT"}}},
		{"CREATE TABLE 	TBL(  );", "TBL", []string{}, []token.Token{}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testCreateStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedColumnNames, tt.expectedColumTypes) {
			return
		}
	}
}

func testCreateStatement(t *testing.T, command ast.Command, expectedTableName string, expectedColumnNames []string, expectedColumTypes []token.Token) bool {
	if command.TokenLiteral() != "CREATE" {
		t.Errorf("command.TokenLiteral() not 'CREATE'. got=%q", command.TokenLiteral())
		return false
	}

	actualCreateCommand, ok := command.(*ast.CreateCommand)
	if !ok {
		t.Errorf("actualCreateCommand is not %T. got=%T", &ast.CreateCommand{}, command)
		return false
	}

	if actualCreateCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualCreateCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !stringArrayEquals(actualCreateCommand.ColumnNames, expectedColumnNames) {
		t.Errorf("")
		return false
	}

	if !tokenArrayEquals(actualCreateCommand.ColumnTypes, expectedColumTypes) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseInsertCommand(t *testing.T) {
	tests := []struct {
		input                string
		expectedTableName    string
		expectedValuesTokens []token.Token
	}{
		{"INSERT INTO TBL VALUES();", "TBL", []token.Token{}},
		{"INSERT INTO TBL VALUES( 'HELLO' );", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}}},
		{"INSERT INTO TBL VALUES( 'HELLO',	 10 , 'LOL');", "TBL", []token.Token{{Type: token.IDENT, Literal: "HELLO"}, {Type: token.LITERAL, Literal: "10"}, {Type: token.IDENT, Literal: "LOL"}}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testInsertStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedValuesTokens) {
			return
		}
	}
}

func testInsertStatement(t *testing.T, command ast.Command, expectedTableName string, expectedValuesTokens []token.Token) bool {
	if command.TokenLiteral() != "INSERT" {
		t.Errorf("command.TokenLiteral() not 'INSERT'. got=%q", command.TokenLiteral())
		return false
	}

	actualInsertCommand, ok := command.(*ast.InsertCommand)
	if !ok {
		t.Errorf("actualInsertCommand is not %T. got=%T", &ast.InsertCommand{}, command)
		return false
	}

	if actualInsertCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualInsertCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !tokenArrayEquals(actualInsertCommand.Values, expectedValuesTokens) {
		t.Errorf("")
		return false
	}

	return true
}

func TestParseSelectCommand(t *testing.T) {
	tests := []struct {
		input             string
		expectedTableName string
		expectedColumns   []token.Token
	}{
		{"SELECT * FROM TBL;", "TBL", []token.Token{{Type: token.ASTERISK, Literal: "*"}}},
		{"SELECT ONE, TWO, THREE FROM TBL;", "TBL", []token.Token{{Type: token.IDENT, Literal: "ONE"}, {Type: token.IDENT, Literal: "TWO"}, {Type: token.IDENT, Literal: "THREE"}}},
		{"SELECT FROM TBL;", "TBL", []token.Token{}},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 1 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testSelectStatement(t, sequences.Commands[0], tt.expectedTableName, tt.expectedColumns) {
			return
		}
	}
}

func TestParseWhereCommand(t *testing.T) {
	firstExpression := ast.ConditionExpresion{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "fda"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	secondExpression := ast.ConditionExpresion{
		Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
		Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
		Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"},
	}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda';",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 EQUAL 6462389;",
			expectedExpression: secondExpression,
		},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 2 {
			t.Fatalf("sequences does not contain 1 statements. got=%d", len(sequences.Commands))
		}

		if !testWhereStatement(t, sequences.Commands[1], tt.expectedExpression) {
			return
		}
	}
}

func TestParseLogicOperatorsInCommand(t *testing.T) {

	firstExpression := ast.OperationExpression{
		Left: ast.ConditionExpresion{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "fda"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Right: ast.ConditionExpresion{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "123"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Operation: token.Token{Type: token.AND, Literal: "AND"},
	}

	secondExpression := ast.OperationExpression{
		Left: ast.ConditionExpresion{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName2"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.LITERAL, Literal: "6462389"}},
			Condition: token.Token{Type: token.NOT, Literal: "NOT"}},
		Right: ast.ConditionExpresion{
			Left:      ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "colName1"}},
			Right:     ast.Anonymitifier{Token: token.Token{Type: token.IDENT, Literal: "qwe"}},
			Condition: token.Token{Type: token.EQUAL, Literal: "EQUAL"}},
		Operation: token.Token{Type: token.OR, Literal: "OR"},
	}

	thirdExpression := ast.BooleanExpresion{
		Boolean: token.Token{Type: token.TRUE, Literal: "TRUE"},
	}

	tests := []struct {
		input              string
		expectedExpression ast.Expression
	}{
		{
			input:              "SELECT * FROM TBL WHERE colName1 EQUAL 'fda' AND colName2 NOT 123;",
			expectedExpression: firstExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE colName2 NOT 6462389 OR colName1 EQUAL 'qwe';",
			expectedExpression: secondExpression,
		},
		{
			input:              "SELECT * FROM TBL WHERE TRUE;",
			expectedExpression: thirdExpression,
		},
	}

	for _, tt := range tests {
		lexer := lexer.RunLexer(tt.input)
		parserInstance := New(lexer)
		sequences := parserInstance.ParseSequence()

		if len(sequences.Commands) != 2 {
			t.Fatalf("sequences does not contain 2 statements. got=%d", len(sequences.Commands))
		}

		if !testWhereStatement(t, sequences.Commands[1], tt.expectedExpression) {
			return
		}
	}
}

func testSelectStatement(t *testing.T, command ast.Command, expectedTableName string, expectedColumnsTokens []token.Token) bool {
	if command.TokenLiteral() != "SELECT" {
		t.Errorf("command.TokenLiteral() not 'SELECT'. got=%q", command.TokenLiteral())
		return false
	}

	actualSelectCommand, ok := command.(*ast.SelectCommand)
	if !ok {
		t.Errorf("actualSelectCommand is not %T. got=%T", &ast.SelectCommand{}, command)
		return false
	}

	if actualSelectCommand.Name.Token.Literal != expectedTableName {
		t.Errorf("%s != %s", actualSelectCommand.TokenLiteral(), expectedTableName)
		return false
	}

	if !tokenArrayEquals(actualSelectCommand.Space, expectedColumnsTokens) {
		t.Errorf("")
		return false
	}

	return true
}

func testWhereStatement(t *testing.T, command ast.Command, expectedExpression ast.Expression) bool {
	if command.TokenLiteral() != "WHERE" {
		t.Errorf("command.TokenLiteral() not 'WHERE'. got=%q", command.TokenLiteral())
		return false
	}

	actualWhereCommand, ok := command.(*ast.WhereCommand)
	if !ok {
		t.Errorf("actualWhereCommand is not %T. got=%T", &ast.WhereCommand{}, command)
		return false
	}

	if expressionsAreEqual(*actualWhereCommand.Expression, expectedExpression) {
		t.Errorf("Actual expression is not equal to expected one")
		return false
	}

	return true
}

func stringArrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func tokenArrayEquals(a []token.Token, b []token.Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Literal != b[i].Literal {
			return false
		}
	}
	return true
}

func expressionsAreEqual(first ast.Expression, second ast.Expression) bool {

	booleanExpresion, booleanExpresionIsValid := first.(ast.BooleanExpresion)
	if booleanExpresionIsValid {
		return validateBooleanExpressions(second, booleanExpresion)

	}

	conditionExpresion, conditionExpresionIsValid := first.(*ast.ConditionExpresion)
	if conditionExpresionIsValid {
		return validateConditionExpresion(second, conditionExpresion)
	}

	operationExpression, operationExpressionIsValid := first.(*ast.OperationExpression)
	if operationExpressionIsValid {
		return validateOperationExpression(second, operationExpression)
	}

	return false
}

func validateOperationExpression(second ast.Expression, operationExpression *ast.OperationExpression) bool {
	secondOperationExpression, secondOperationExpressionIsValid := second.(*ast.OperationExpression)

	if !secondOperationExpressionIsValid {
		return false
	}

	if operationExpression.Operation.Literal != secondOperationExpression.Operation.Literal {
		return false
	}

	return expressionsAreEqual(operationExpression, secondOperationExpression)
}

func validateConditionExpresion(second ast.Expression, conditionExpresion *ast.ConditionExpresion) bool {
	secondConditionExpresion, secondConditionExpresionIsValid := second.(*ast.ConditionExpresion)

	if !secondConditionExpresionIsValid {
		return false
	}

	if conditionExpresion.Left.GetToken().Literal != secondConditionExpresion.Left.GetToken().Literal &&
		conditionExpresion.Left.IsIdentifier() == secondConditionExpresion.Left.IsIdentifier() {
		return false
	}

	if conditionExpresion.Right.GetToken().Literal != secondConditionExpresion.Right.GetToken().Literal &&
		conditionExpresion.Right.IsIdentifier() == secondConditionExpresion.Right.IsIdentifier() {
		return false
	}

	if conditionExpresion.Condition.Literal != secondConditionExpresion.Condition.Literal {
		return false
	}
	return false
}

func validateBooleanExpressions(second ast.Expression, booleanExpresion ast.BooleanExpresion) bool {
	secondBooleanExpresion, secondBooleanExpresionIsValid := second.(ast.BooleanExpresion)

	if !secondBooleanExpresionIsValid {
		return false
	}

	if booleanExpresion.Boolean.Literal != secondBooleanExpresion.Boolean.Literal {
		return false
	}

	return true
}
