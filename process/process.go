package process

import (
	"time"

	"go.uber.org/cadence/workflow"
)

type Status string

const (
	Running  Status = "running"
	Failed   Status = "failed"
	Finished Status = "finished"
)

func NewProcessLogActivity(service Service) ProcessLogActivity {
	return ProcessLogActivity{service: service}
}

type ProcessLogActivityInput struct {
	ID         string
	ParentID   string
	OrgID      int32
	Type       string
	Log        string
	ResourceID string
	Status     Status
	StartedAt  time.Time
	FinishedAt *time.Time
}

type ProcessEventActivityInput struct {
	ProcessID string
	Type      string
	Log       string
	Status    Status
	Timestamp time.Time
}

type Status interface {
	End(error)
}

type processLog struct {
	ctx           workflow.Context
	activityInput ProcessLogActivityInput
}

func (p *processLog) End(err error) {
	finishedAt := workflow.Now(p.ctx)
	p.activityInput.FinishedAt = &finishedAt
	if err != nil {
		p.activityInput.Status = Failed
		p.activityInput.Log = err.Error()
	} else {
		p.activityInput.Status = Finished
	}

	err = workflow.ExecuteActivity(p.ctx, ProcessLogActivityName, p.activityInput).Get(p.ctx, nil)
	if err != nil {
		workflow.GetLogger(p.ctx).Sugar().Warnf("failed to log process end: %s", err)
	}
}

func NewProcessLog(ctx workflow.Context, orgID uint, resourceID string) Status {
	winfo := workflow.GetInfo(ctx)
	parentID := ""
	if winfo.ParentWorkflowExecution != nil {
		parentID = winfo.ParentWorkflowExecution.ID
	}
	activityInput := ProcessLogActivityInput{
		ID:         winfo.WorkflowExecution.ID,
		ParentID:   parentID,
		Type:       winfo.WorkflowType.Name,
		StartedAt:  workflow.Now(ctx),
		Status:     Running,
		OrgID:      int32(orgID),
		ResourceID: resourceID,
	}
	err := workflow.ExecuteActivity(ctx, ProcessLogActivityName, activityInput).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Sugar().Warnf("failed to log process: %s", err)
	}

	return &processLog{ctx: ctx, activityInput: activityInput}
}

type processEvent struct {
	ctx           workflow.Context
	activityInput ProcessEventActivityInput
}

func (p *processEvent) End(err error) {
	p.activityInput.Timestamp = workflow.Now(p.ctx)
	if err != nil {
		p.activityInput.Status = Failed
		p.activityInput.Log = err.Error()
	} else {
		p.activityInput.Status = Finished
	}

	err = workflow.ExecuteActivity(p.ctx, ProcessEventActivityName, p.activityInput).Get(p.ctx, nil)
	if err != nil {
		workflow.GetLogger(p.ctx).Sugar().Warnf("failed to log process event end: %s", err.Error())
	}
}

func NewProcessEvent(ctx workflow.Context, activityName string) Status {
	winfo := workflow.GetInfo(ctx)

	activityInput := ProcessEventActivityInput{
		ProcessID: winfo.WorkflowExecution.ID,
		Type:      activityName,
		Timestamp: workflow.Now(ctx),
		Status:    Running,
	}

	err := workflow.ExecuteActivity(ctx, ProcessEventActivityName, activityInput).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Sugar().Warnf("failed to log process event: %s", err)
	}

	return &processEvent{ctx: ctx, activityInput: activityInput}
}
