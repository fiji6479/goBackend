package handler

import (
	"fmt"
	"goBackend/api/proto"
	"goBackend/internal/model"
	"goBackend/internal/service"
)

type GRPCServer struct {
	proto.UnimplementedCalculatorServer
	eval *service.Evaluator
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		eval: service.NewEvaluator(),
	}
}

func (s *GRPCServer) Calculate(stream proto.Calculator_CalculateServer) error {
	var instructions []model.Instruction

	for {
		in, err := stream.Recv()
		if err != nil {
			break
		}

		instr := model.Instruction{
			Type:  model.InstructionType(in.Type),
			Op:    in.Op,
			Var:   in.Var,
			Left:  parseOperand(in.Left),
			Right: parseOperand(in.Right),
		}
		instructions = append(instructions, instr)
	}

	result, err := s.eval.EvalInstructions(instructions)
	if err != nil {
		return stream.SendAndClose(&proto.Output{
			Items: []*proto.OutputItem{},
		})
	}

	var out []*proto.OutputItem
	for _, item := range result {
		out = append(out, &proto.OutputItem{
			Var:   item.Var,
			Value: item.Value,
		})
	}

	return stream.SendAndClose(&proto.Output{Items: out})
}

func parseOperand(value string) any {
	var v int64
	_, err := fmt.Sscanf(value, "%d", &v)
	if err == nil {
		return v
	}
	return value
}
