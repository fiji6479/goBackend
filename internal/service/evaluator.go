package service

import (
	"fmt"
	"goBackend/internal/model"
	"strconv"
)

type Evaluator struct {
	results   map[string]int64
	processed map[string]bool
	deps      map[string][]string
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		results:   make(map[string]int64),
		processed: make(map[string]bool),
		deps:      make(map[string][]string),
	}
}

func (e *Evaluator) EvalInstructions(instructions []model.Instruction) ([]model.OutputItem, error) {
	// Шаг 1: собрать зависимости для calc-инструкций
	for _, instr := range instructions {
		if instr.Type == model.CalcType {
			var deps []string
			if s, ok := instr.Left.(string); ok {
				deps = append(deps, s)
			}
			if s, ok := instr.Right.(string); ok {
				deps = append(deps, s)
			}
			e.deps[instr.Var] = deps
		}
	}

	// Шаг 2: вычислить calc-инструкции
	for _, instr := range instructions {
		if instr.Type == model.CalcType {
			if e.processed[instr.Var] {
				return nil, fmt.Errorf("переменная '%s' уже определена", instr.Var)
			}
			err := e.evalCalc(instr, instructions)
			if err != nil {
				return nil, err
			}
			e.processed[instr.Var] = true
		}
	}

	// Шаг 3: собрать результаты print-инструкций
	var output []model.OutputItem
	for _, instr := range instructions {
		if instr.Type == model.PrintType {
			val, ok := e.results[instr.Var]
			if !ok {
				val = 0
			}
			output = append(output, model.OutputItem{
				Var:   instr.Var,
				Value: val,
			})
		}
	}

	return output, nil
}

func (e *Evaluator) evalCalc(instr model.Instruction, all []model.Instruction) error {
	left, err := e.resolve(instr.Left, all)
	if err != nil {
		return err
	}
	right, err := e.resolve(instr.Right, all)
	if err != nil {
		return err
	}

	var result int64
	switch instr.Op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	default:
		return fmt.Errorf("неподдерживаемая операция: %s", instr.Op)
	}

	e.results[instr.Var] = result
	return nil
}

func (e *Evaluator) resolve(val any, all []model.Instruction) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		if res, ok := e.results[v]; ok {
			return res, nil
		}
		for _, instr := range all {
			if instr.Type == model.CalcType && instr.Var == v {
				err := e.evalCalc(instr, all)
				if err != nil {
					return 0, err
				}
				return e.results[v], nil
			}
		}
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return num, nil
		}
		return 0, fmt.Errorf("неизвестная переменная: %s", v)
	default:
		return 0, fmt.Errorf("некорректный тип операнда: %T", val)
	}
}
