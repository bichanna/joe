package compiler

type ASTType uint

const (
	ASTClassDecl ASTType = iota
	ASTImportDecl
	ASTModuleDecl
	ASTMethodDecl
	ASTConstructDecl
	ASTLabelDecl
	ASTOperatorDecl
	ASTVarDecl
	ASTValue
	ASTValueList
	ASTUTypeArgList
	ASTUtypeArgListOpt
	ASTVectorArray
	ASTUTypeArg
	ASTUTypeArgOpt
	ASTExpr
	ASTArrayExpr
	ASTPrimaryExpr
	ASTDotCallExpr
	ASTUType
	ASTBlock
	ASTFinallyBlock
	ASTAssemblyBlock
	ASTCatchClause
	ASTMethodReturnType
	ASTReturnStmt
	ASTStmt
	ASTIfStmt
	ASTElseIfStmt
	ASTElseStmt
	ASTTryCatchStmt
	ASTThrowStmt
	ASTContinueStmt
	ASTBreakStmt
	ASTGoToStmt
	ASTWhileStmt
	ASTDoWhileStmt
	ASTAssemblyStmt
	ASTForStmt
	ASTForExprCond
	ASTForExprIter
	ASTForEachStmt
	ASTTypeID
	ASTRefPtr
	ASTModuleName
	ASTLiteral

	ASTLiteral_E
	ASTUTypeClass_E
	ASTDotNot_E
	ASTSelf_E
	ASTBase_E
	ASTNull_E
	ASTNev_E
	ASTNot_E
	ASTPostInc_E
	ASTArray_E
	ASTDotFn_E
	ASTCast_E
	ASTPreInc_E
	ASTParen_E
	ASTVect_E
	ASTAdd_E
	ASTMult_E
	ASTShift_E
	ASTLess_E
	ASTEqual_E
	ASTAnd_E
	ASTQuestion_E
	ASTAssign_E
	ASTSizeOf_E

	ASTNone
)

type AST struct {
	Line uint
	Col  uint

	atype    ASTType
	parent   *AST
	subAsts  []*AST
	entities []*TokenEntity
}

// NewAST creates a new AST node with given information.
func NewAST(parent *AST, atype ASTType, line, col uint) *AST {
	return &AST{
		atype:  atype,
		parent: parent,
		Line:   line,
		Col:    col,
	}
}

// NewEmptyAST creates a new empty AST node.
func NewEmptyAST() *AST {
	return NewAST(nil, 0, 0, 0)
}

func (a *AST) Encapsulate(t ASTType) {

}

// GetType gets the AST type of the AST node.
func (a *AST) GetType() ASTType {
	return a.atype
}

// GetParent gets the parent node of the AST node.
func (a *AST) GetParent() *AST {
	return a.parent
}

// SubASTCount returns the number of sub ASTs the current AST has.
func (a *AST) SubASTCount() int {
	return len(a.subAsts)
}

// GetSubAST gets the sub AST of AST type t.
func (a *AST) GetSubAST(t ASTType) *AST {
	return nil
}

// GetSUbASTAt gets the sub AST at `at`.
func (a *AST) GetSUbASTAt(at int) *AST {
	return a.subAsts[at]
}

// GetLastSubAST gets the last sub AST in the list.
func (a *AST) GetLastSubAST() *AST {
	return a.subAsts[len(a.subAsts)-1]
}

// HasSubAST checks if the AST has a sub AST of the provided type.
func (a *AST) HasSubAST(t ASTType) bool {
	return false
}

// HasEntity checks if the AST has a token entity of the provided type.
func (a *AST) HasEntity(t TokenType) bool {
	return false
}

// GetEntity gets the token entity of the provided type.
func (a *AST) GetEntity(t TokenType) *TokenEntity {
	return nil
}

// GetEntityAt gets the token entity at `at`.
func (a *AST) GetEntityAt(at int) *TokenEntity {
	return a.entities[at]
}

// GetEntityCount gets the number of token entities the current AST has.
func (a *AST) GetEntityCount() int {
	return len(a.entities)
}

// AddEntity adds the given entity to the list.
func (a *AST) AddEntity(entity *TokenEntity) {
	a.entities = append(a.entities, entity)
}

// AddAST adds the given AST to the list.
func (a *AST) AddAST(ast *AST) {
	a.subAsts = append(a.subAsts, ast)
}

// SetASTType sets the AST type to the provided type.
func (a *AST) SetASTType(t ASTType) {
	a.atype = t
}
