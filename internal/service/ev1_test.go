package service

import (
	"goBackend/internal/model"
	"reflect"
	"testing"
)

func TestEvaluator_Example1(t *testing.T) {
	e := NewEvaluator()

	instr := []model.Instruction{
		{Type: model.CalcType, Op: "+", Var: "x", Left: 1, Right: 2},
		{Type: model.PrintType, Var: "x"},
	}

	expected := []model.OutputItem{{Var: "x", Value: 3}}
	assertResult(t, e, instr, expected)
}

func TestEvaluator_Example2(t *testing.T) {
	e := NewEvaluator()

	instr := []model.Instruction{
		{Type: model.CalcType, Op: "+", Var: "x", Left: 10, Right: 2},
		{Type: model.PrintType, Var: "x"},
		{Type: model.CalcType, Op: "-", Var: "y", Left: "x", Right: 3},
		{Type: model.CalcType, Op: "*", Var: "z", Left: "x", Right: "y"},
		{Type: model.PrintType, Var: "w"},
		{Type: model.CalcType, Op: "*", Var: "w", Left: "z", Right: 0},
	}

	expected := []model.OutputItem{
		{Var: "x", Value: 12},
		{Var: "w", Value: 0},
	}

	assertResult(t, e, instr, expected)
}

func TestEvaluator_Example3(t *testing.T) {
	e := NewEvaluator()

	instr := []model.Instruction{
		{Type: model.CalcType, Op: "+", Var: "x", Left: 10, Right: 2},
		{Type: model.CalcType, Op: "*", Var: "y", Left: "x", Right: 5},
		{Type: model.CalcType, Op: "-", Var: "q", Left: "y", Right: 20},
		{Type: model.CalcType, Op: "+", Var: "unusedA", Left: "y", Right: 100},
		{Type: model.CalcType, Op: "*", Var: "unusedB", Left: "unusedA", Right: 2},
		{Type: model.PrintType, Var: "q"},
		{Type: model.CalcType, Op: "-", Var: "z", Left: "x", Right: 15},
		{Type: model.PrintType, Var: "z"},
		{Type: model.CalcType, Op: "+", Var: "ignoreC", Left: "z", Right: "y"},
		{Type: model.PrintType, Var: "x"},
	}

	expected := []model.OutputItem{
		{Var: "q", Value: 40},
		{Var: "z", Value: -3},
		{Var: "x", Value: 12},
	}

	assertResult(t, e, instr, expected)
}

// Вспомогательная функция сравнения результатов
func assertResult(t *testing.T, e *Evaluator, input []model.Instruction, expected []model.OutputItem) {
	t.Helper()

	result, err := e.EvalInstructions(input)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось: %v, получено: %v", expected, result)
	}
}
