package parser

import (
	"go/ast"

	"github.com/minions/bello/pkg/lexer"
)

// Node represents all Bello AST nodes with source location support.
type Node interface {
	Pos() lexer.Position
}

// Decl represents declarations.
type Decl interface {
	Node
	declNode()
}

// Expr represents expressions.
type Expr interface {
	Node
	exprNode()
}

// Stmt represents statements.
type Stmt interface {
	Node
	stmtNode()
}

// File is the top-level parsed source unit.
type File struct {
	GoFile    *ast.File
	Package   *Ident
	Imports   []*ImportSpec
	Decls     []Decl
	Posn      lexer.Position
	Filename  string
	Translated string
}

func (f *File) Pos() lexer.Position { return f.Posn }

// Base node with coordinates.
type BaseNode struct {
	Posn lexer.Position
}

func (b *BaseNode) Pos() lexer.Position { return b.Posn }

// ----------------------------------------------------------------------------
// Declarations

type ImportSpec struct {
	BaseNode
	Name *Ident
	Path string
}

func (*ImportSpec) declNode() {}

type ValueSpec struct {
	BaseNode
	Name *Ident
}

func (*ValueSpec) declNode() {}

type FuncDecl struct {
	BaseNode
	Recv   *FieldList
	Name   *Ident
	Type   *FuncType
	Body   *BlockStmt
	TypeParams *FieldList
}

func (*FuncDecl) declNode() {}

type VarDecl struct {
	BaseNode
	Names []*Ident
	Type   Expr
	Values []Expr
}

func (*VarDecl) declNode() {}

type ConstDecl struct {
	BaseNode
	Names []*Ident
	Type  Expr
	Values []Expr
}

func (*ConstDecl) declNode() {}

type TypeDecl struct {
	BaseNode
	Name *Ident
	Type Expr
}

func (*TypeDecl) declNode() {}

// ----------------------------------------------------------------------------
// Statements

type BlockStmt struct {
	BaseNode
	List []Stmt
}

func (*BlockStmt) stmtNode() {}

type ReturnStmt struct {
	BaseNode
	Results []Expr
}

func (*ReturnStmt) stmtNode() {}

type IfStmt struct {
	BaseNode
	Init Stmt
	Cond Expr
	Body *BlockStmt
	Else Stmt
}

func (*IfStmt) stmtNode() {}

type ForStmt struct {
	BaseNode
	Init   Stmt
	Cond   Expr
	Post   Stmt
	Range  *RangeClause
	Body   *BlockStmt
}

func (*ForStmt) stmtNode() {}

type RangeClause struct {
	BaseNode
	Key   Expr
	Value Expr
	X     Expr
}

type SwitchStmt struct {
	BaseNode
	Init Stmt
	Tag  Expr
	Body []Stmt
}

func (*SwitchStmt) stmtNode() {}

type TypeSwitchStmt struct {
	BaseNode
	Init   Stmt
	Assign Stmt
	Body   []Stmt
}

func (*TypeSwitchStmt) stmtNode() {}

type SelectStmt struct {
	BaseNode
	Body []Stmt
}

func (*SelectStmt) stmtNode() {}

type CommClause struct {
	BaseNode
	Comm Stmt
	Body []Stmt
}

func (*CommClause) stmtNode() {}

type DeclStmt struct {
	BaseNode
	Decl Decl
}

func (*DeclStmt) stmtNode() {}

type CaseClause struct {
	BaseNode
	Values  []Expr
	Body    []Stmt
	Default bool
}

func (*CaseClause) stmtNode() {}

type GoStmt struct {
	BaseNode
	Call *CallExpr
}

func (*GoStmt) stmtNode() {}

type DeferStmt struct {
	BaseNode
	Call *CallExpr
}

func (*DeferStmt) stmtNode() {}

type BranchStmt struct {
	BaseNode
	Tok   string
	Label *Ident
}

func (*BranchStmt) stmtNode() {}

type LabeledStmt struct {
	BaseNode
	Label *Ident
	Stmt  Stmt
}

func (*LabeledStmt) stmtNode() {}

type AssignStmt struct {
	BaseNode
	Lhs []Expr
	Op  string
	Rhs []Expr
}

func (*AssignStmt) stmtNode() {}

type SendStmt struct {
	BaseNode
	Chan  Expr
	Value Expr
}

func (*SendStmt) stmtNode() {}

type IncDecStmt struct {
	BaseNode
	X   Expr
	Tok string
}

func (*IncDecStmt) stmtNode() {}

type ExprStmt struct {
	BaseNode
	X Expr
}

func (*ExprStmt) stmtNode() {}

// ----------------------------------------------------------------------------
// Expressions

type BinaryExpr struct {
	BaseNode
	X  Expr
	Op lexer.TokenType
	Y  Expr
}

func (*BinaryExpr) exprNode() {}

type UnaryExpr struct {
	BaseNode
	Op lexer.TokenType
	X  Expr
}

func (*UnaryExpr) exprNode() {}

type CallExpr struct {
	BaseNode
	Fun      Expr
	Args     []Expr
	Ellipsis bool
}

func (*CallExpr) exprNode() {}

type IndexExpr struct {
	BaseNode
	X, Index Expr
}

func (*IndexExpr) exprNode() {}

type SliceExpr struct {
	BaseNode
	X, Low, High, Max Expr
	Full             bool
}

func (*SliceExpr) exprNode() {}

type SelectorExpr struct {
	BaseNode
	X, Sel Expr
}

func (*SelectorExpr) exprNode() {}

type TypeAssertExpr struct {
	BaseNode
	X    Expr
	Type Expr
}

func (*TypeAssertExpr) exprNode() {}

type Ident struct {
	BaseNode
	Name string
}

func (*Ident) exprNode() {}

type BasicLit struct {
	BaseNode
	Kind  lexer.TokenType
	Value string
}

func (*BasicLit) exprNode() {}

type CompositeLit struct {
	BaseNode
	Type Expr
	Elts []Expr
}

func (*CompositeLit) exprNode() {}

type KeyValueExpr struct {
	BaseNode
	Key   Expr
	Value Expr
}

func (*KeyValueExpr) exprNode() {}

type FuncLit struct {
	BaseNode
	Type *FuncType
	Body *BlockStmt
}

func (*FuncLit) exprNode() {}

// ----------------------------------------------------------------------------
// Type nodes

type FieldList struct {
	BaseNode
	List []*Field
}

type Field struct {
	BaseNode
	Names    []*Ident
	Type     Expr
	Tag      *BasicLit
	Ellipsis bool
}

type ArrayType struct {
	BaseNode
	Len Expr
	Elt Expr
}

type SliceType struct {
	BaseNode
	Elt Expr
}

type MapType struct {
	BaseNode
	Key Expr
	Value Expr
}

type ChanType struct {
	BaseNode
	Dir string
	Value Expr
}

type PointerType struct {
	BaseNode
	Elt Expr
}

type FuncType struct {
	BaseNode
	TypeParams *FieldList
	Params  *FieldList
	Results *FieldList
}

type StructType struct {
	BaseNode
	Fields *FieldList
}

type InterfaceType struct {
	BaseNode
	Methods *FieldList
}

// ----------------------------------------------------------------------------
// Helpers

func (n *FieldList) declNode() {}
func (n *Field) exprNode() {}
func (n *ArrayType) exprNode() {}
func (n *SliceType) exprNode() {}
func (n *MapType) exprNode() {}
func (n *ChanType) exprNode() {}
func (n *PointerType) exprNode() {}
func (n *FuncType) exprNode() {}
func (n *StructType) exprNode() {}
func (n *InterfaceType) exprNode() {}
