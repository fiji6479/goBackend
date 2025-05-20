package model

type InstructionType string

const (
	CalcType  InstructionType = "calc"
	PrintType InstructionType = "print"
)

type Instruction struct {
	Type  InstructionType
	Op    string
	Var   string
	Left  any
	Right any
}

type OutputItem struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}

type Output struct {
	Items []OutputItem `json:"items"`
}
