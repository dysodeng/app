package errors

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Pipeline 使用 errgroup 的链式管道多错误处理
type Pipeline struct {
	ctx context.Context
	g   *errgroup.Group
	fns []func() error
}

func NewPipeline() *Pipeline {
	return NewPipelineWithContext(context.Background())
}

func NewPipelineWithContext(ctx context.Context) *Pipeline {
	g, ctx := errgroup.WithContext(ctx)
	return &Pipeline{
		ctx: ctx,
		g:   g,
		fns: make([]func() error, 0),
	}
}

// Then 添加一个函数到管道中
func (p *Pipeline) Then(fn ...func() error) *Pipeline {
	p.fns = append(p.fns, fn...)
	return p
}

// Execute 顺序执行所有函数（遇到错误立即停止）
func (p *Pipeline) Execute() error {
	for _, fn := range p.fns {
		select {
		case <-p.ctx.Done():
			return p.ctx.Err()
		default:
		}

		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteParallel 并发执行所有函数
func (p *Pipeline) ExecuteParallel() error {
	for _, fn := range p.fns {
		fn := fn // 避免闭包问题
		p.g.Go(func() error {
			return fn()
		})
	}
	return p.g.Wait()
}

// ExecuteParallelWithLimit 限制并发数量执行
func (p *Pipeline) ExecuteParallelWithLimit(limit int) error {
	g := new(errgroup.Group)
	g.SetLimit(limit)

	for _, fn := range p.fns {
		fn := fn // 避免闭包问题
		g.Go(func() error {
			return fn()
		})
	}
	return g.Wait()
}

// Context 获取上下文
func (p *Pipeline) Context() context.Context {
	return p.ctx
}
