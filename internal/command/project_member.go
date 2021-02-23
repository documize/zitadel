package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	addedMember := NewProjectMemberWriteModel(member.AggregateID, member.UserID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := r.addProjectMember(ctx, projectAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (r *CommandSide) addProjectMember(ctx context.Context, projectAgg *eventstore.Aggregate, addedMember *ProjectMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-W8m4l", "Errors.Project.Member.Invalid")
	}

	err := r.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, err
	}
	err = r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "PROJECT-PtXi1", "Errors.Project.Member.AlreadyExists")
	}

	return project.NewProjectMemberAddedEvent(ctx, projectAgg, member.UserID, member.Roles...), nil
}

//ChangeProjectMember updates an existing member
func (r *CommandSide) ChangeProjectMember(ctx context.Context, member *domain.Member, resourceOwner string) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.Member.Invalid")
	}

	existingMember, err := r.projectMemberWriteModelByID(ctx, member.AggregateID, member.UserID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "PROJECT-LiaZi", "Errors.Project.Member.RolesNotChanged")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, project.NewProjectMemberChangedEvent(ctx, projectAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveProjectMember(ctx context.Context, projectID, userID, resourceOwner string) error {
	m, err := r.projectMemberWriteModelByID(ctx, projectID, userID, resourceOwner)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	projectAgg := ProjectAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, project.NewProjectMemberRemovedEvent(ctx, projectAgg, userID))
	return err
}

func (r *CommandSide) projectMemberWriteModelByID(ctx context.Context, projectID, userID, resourceOwner string) (member *ProjectMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewProjectMemberWriteModel(projectID, userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "PROJECT-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}