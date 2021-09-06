package gou

import (
	"context"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/maps"
)

// Exec 运行flow
func (flow *Flow) Exec(args ...interface{}) interface{} {

	vm := &FlowVM{
		Otto: otto.New(),
	}
	res := map[string]interface{}{} // 结果集
	ctx, cancel := context.WithCancel(context.Background())

	flowCtx := &FlowContext{
		Context: &ctx,
		Cancel:  cancel,
		Res:     res,
		In:      args,
	}

	// 运行工作流节点，并处理数据
	flowProcess := "flows." + flow.Name
	for i, node := range flow.Nodes {

		// 死循环检查
		if strings.HasPrefix(node.Process, flowProcess) {
			exception.New("不能调用自身工作流(%s)", 400, node.Process)
		}

		// 执行解析器
		flow.ExecNode(&node, flowCtx, vm, i-1)
	}

	// 结果集输出处理
	return flow.FormatResult(flowCtx)
}

// FormatResult 结果集格式化输出
func (flow *Flow) FormatResult(ctx *FlowContext) interface{} {
	if flow.Output == nil {
		return ctx.Res
	}

	var data = maps.Map{"$in": ctx.In, "$res": ctx.Res}.Dot()
	return Bind(flow.Output, data)
}

// ExecNode 运行节点
func (flow *Flow) ExecNode(node *FlowNode, ctx *FlowContext, vm *FlowVM, prev int) []interface{} {
	var data = maps.Map{"$in": ctx.In, "$res": ctx.Res}.Dot()
	resp, outs := flow.RunProcess(node, ctx, data)
	_, outs = flow.RunScript(vm, node, ctx, data, resp, outs)
	return outs
}

// RunProcess 运行处理器
func (flow *Flow) RunProcess(node *FlowNode, ctx *FlowContext, data maps.Map) (interface{}, []interface{}) {
	var outs = []interface{}{}
	var resp, res interface{}

	for i := range node.Args {
		node.Args[i] = Bind(node.Args[i], data)
	}

	if node.Process != "" {
		process := NewProcess(node.Process, node.Args...)
		resp = process.Run()
	}

	if node.Outs == nil || len(node.Outs) == 0 {
		res = resp
	} else {
		data["$out"] = resp
		data = data.Dot()
		for _, value := range node.Outs {
			outs = append(outs, Bind(value, data))
		}
		res = outs
	}

	if node.Name != "" {
		ctx.Res[node.Name] = res
	}
	return resp, outs
}

// RunScript 运行数据处理脚本
func (flow *Flow) RunScript(vm *FlowVM, node *FlowNode, ctx *FlowContext, data maps.Map, processResp interface{}, processOuts []interface{}) (interface{}, []interface{}) {
	var resp, res interface{}
	if node.Script == "" {
		return processResp, processOuts
	}

	source, has := flow.Scripts[node.Script]
	if !has {
		return processResp, processOuts
	}

	filename := flow.Name + "." + node.Script + ".js"
	vm.Set("args", ctx.In)
	vm.Set("res", ctx.Res)
	vm.Set("out", processResp)
	script, err := vm.Compile(filename, source+"\nmain(args, out, res);")
	if err != nil {
		exception.Err(err, 500).Throw()
	}

	value, err := vm.Run(script)
	if err != nil {
		exception.Err(err, 500).Throw()
	}

	resp, err = value.Export()
	if err != nil {
		exception.Err(err, 500).Throw()
	}

	if node.Outs == nil || len(node.Outs) == 0 {
		res = resp
	} else {
		var outs = []interface{}{}
		data["$out"] = resp
		data = data.Dot()
		for _, value := range node.Outs {
			outs = append(outs, Bind(value, data))
		}
		res = outs
		processOuts = outs
	}

	if node.Name != "" {
		ctx.Res[node.Name] = res
	}

	return resp, processOuts
}