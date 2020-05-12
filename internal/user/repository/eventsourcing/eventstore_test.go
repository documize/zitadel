package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	repo_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		user *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "user from events, ok",
			args: args{
				es:   GetMockUserByIDOK(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "UserName"}}),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &model.Profile{UserName: "UserName"}},
			},
		},
		{
			name: "no events found",
			args: args{
				es:   GetMockUserByIDNoEvents(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no id",
			args: args{
				es:   GetMockUserByIDNoEvents(ctrl),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserByID(nil, tt.args.user.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID != tt.res.user.AggregateID {
				t.Errorf("got wrong result aggregateID: expected: %v, actual: %v ", tt.res.user.AggregateID, result.AggregateID)
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result userName: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es   *UserEventstore
		ctx  context.Context
		user *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "init mail because no pw",
			args: args{
				es:  GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
		},
		{
			name: "email as username",
			args: args{
				es:  GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "EmailAddress", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
		},
		{
			name: "with verified phone number",
			args: args{
				es:  GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "EmailAddress", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true}, Phone: &repo_model.Phone{PhoneNumber: "PhoneNumber", IsPhoneVerified: true}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
						UserName:  "UserName",
					},
					Email: &model.Email{
						EmailAddress:    "UserName",
						IsEmailVerified: true,
					},
					Phone: &model.Phone{
						PhoneNumber:     "UserName",
						IsPhoneVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
					Phone: &model.Phone{
						PhoneNumber:     "UserName",
						IsPhoneVerified: true,
					},
				},
			},
		},
		{
			name: "with password",
			args: args{
				es:  GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress", IsEmailVerified: true}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
						UserName:  "UserName",
					},
					Password: &model.Password{SecretString: "Password"},
					Email: &model.Email{
						EmailAddress:    "UserName",
						IsEmailVerified: true,
					},
				},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress:    "EmailAddress",
						IsEmailVerified: true,
					},
				},
			},
		},
		{
			name: "create user invalid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateUser(tt.args.ctx, tt.args.user)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result username: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc == nil && tt.res.user.Email != nil {
				if result.IsEmailVerified != tt.res.user.IsEmailVerified {
					t.Errorf("got wrong result IsEmailVerified: expected: %v, actual: %v ", tt.res.user.IsEmailVerified, result.IsEmailVerified)
				}
			}
			if tt.res.errFunc == nil && tt.res.user.Phone != nil {
				if result.IsPhoneVerified != tt.res.user.IsPhoneVerified {
					t.Errorf("got wrong result IsPhoneVerified: expected: %v, actual: %v ", tt.res.user.IsPhoneVerified, result.IsPhoneVerified)
				}
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *UserEventstore
		ctx           context.Context
		user          *model.User
		resourceOwner string
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "register user, ok",
			args: args{
				es:  GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress"}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{
					ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
					Password: &model.Password{
						SecretString: "Password",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "email as username",
			args: args{
				es:  GetMockManipulateUserWithPasswordAndEmailCodeGen(ctrl, repo_model.User{Profile: &repo_model.Profile{UserName: "EmailAddress", FirstName: "FirstName", LastName: "LastName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress"}}),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
					Password: &model.Password{
						SecretString: "Password",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
		},
		{
			name: "invalid user",
			args: args{
				es:            GetMockManipulateUser(ctrl),
				ctx:           auth.NewMockContext("orgID", "userID"),
				user:          &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1}},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "register user no password",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
				resourceOwner: "ResourceOwner",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no resourceowner",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
				user: &model.User{ObjectRoot: es_models.ObjectRoot{Sequence: 1},
					Profile: &model.Profile{
						UserName:  "EmailAddress",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &model.Email{
						EmailAddress: "EmailAddress",
					},
				},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.RegisterUser(tt.args.ctx, tt.args.user, tt.args.resourceOwner)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.UserName, result.UserName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestDeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate user, ok",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.USERSTATE_INACTIVE},
			},
		},
		{
			name: "deactivate user with inactive state",
			args: args{
				es:       GetMockManipulateInactiveUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateUser(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestReactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate user, ok",
			args: args{
				es:       GetMockManipulateInactiveUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.USERSTATE_ACTIVE},
			},
		},
		{
			name: "reactivate user with inital state",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateUser(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "lock user, ok",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.USERSTATE_LOCKED},
			},
		},
		{
			name: "lock user with locked state",
			args: args{
				es:       GetMockManipulateLockedUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.LockUser(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestUnlockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		user    *model.User
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "unlock user, ok",
			args: args{
				es:       GetMockManipulateLockedUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				user: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, State: model.USERSTATE_ACTIVE},
			},
		},
		{
			name: "unlock user not locked state",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UnlockUser(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.State != tt.res.user.State {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.user.State, result.State)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestGetInitCodeByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		code    *model.InitUserCode
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserWithInitCode(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				code: &model.InitUserCode{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Expiry: time.Hour * 30},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.InitializeUserCodeByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Expiry != tt.res.code.Expiry {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.code.Expiry, result.Expiry)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateInitCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		code    *model.InitUserCode
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create init code",
			args: args{
				es:       GetMockManipulateUserWithInitCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}}),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				code: &model.InitUserCode{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Expiry: time.Hour * 1},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateInitializeUserCodeByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Expiry != tt.res.code.Expiry {
				t.Errorf("got wrong result expiry: expected: %v, actual: %v ", tt.res.code.Expiry, result.Expiry)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSkipMfaInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "skip mfa init",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.SkipMfaInit(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("rshould not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPasswordID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserFull(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: true},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "existing pw not found",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.UserPasswordByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.ChangeRequired != tt.res.password.ChangeRequired {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.password.ChangeRequired, result.ChangeRequired)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSetOneTimePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		password *model.Password
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create one time pw",
			args: args{
				es:       GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}}),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: true},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: ""}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetOneTimePassword(tt.args.ctx, tt.args.password)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.ChangeRequired != true {
				t.Errorf("should be one time")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestSetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		password *model.Password
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create pw",
			args: args{
				es:       GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}}),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: ""}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}, SecretString: "Password"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.SetPassword(tt.args.ctx, tt.args.password)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.ChangeRequired != false {
				t.Errorf("should not be one time")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRequestSetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es         *UserEventstore
		ctx        context.Context
		userID     string
		notifyType model.NotificationType
	}
	type res struct {
		password *model.Password
		errFunc  func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create pw",
			args: args{
				es:         GetMockManipulateUserWithPasswordCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID"}}),
				ctx:        auth.NewMockContext("orgID", "userID"),
				userID:     "AggregateID",
				notifyType: model.NOTIFICATIONTYPE_EMAIL,
			},
			res: res{
				password: &model.Password{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, ChangeRequired: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:         GetMockManipulateUser(ctrl),
				notifyType: model.NOTIFICATIONTYPE_EMAIL,
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:         GetMockManipulateUserNoEvents(ctrl),
				ctx:        auth.NewMockContext("orgID", "userID"),
				userID:     "AggregateID",
				notifyType: model.NOTIFICATIONTYPE_EMAIL,
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RequestSetPassword(tt.args.ctx, tt.args.userID, tt.args.notifyType)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestProfileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		profile *model.Profile
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserFull(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, UserName: "UserName"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProfileByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.UserName != tt.res.profile.UserName {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.profile.UserName, result.UserName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *UserEventstore
		ctx     context.Context
		profile *model.Profile
	}
	type res struct {
		profile *model.Profile
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:      GetMockManipulateUserFull(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
			res: res{
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
		},
		{
			name: "invalid profile",
			args: args{
				es:      GetMockManipulateUser(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:      GetMockManipulateUserNoEvents(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				profile: &model.Profile{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, FirstName: "FirstName Changed", LastName: "LastName Changed"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeProfile(tt.args.ctx, tt.args.profile)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.FirstName != tt.res.profile.FirstName {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.profile.FirstName, result.FirstName)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestEmailByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		email   *model.Email
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserFull(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddress"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.EmailByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.EmailAddress != tt.res.email.EmailAddress {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.EmailAddress, result.EmailAddress)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserEventstore
		ctx   context.Context
		email *model.Email
	}
	type res struct {
		email   *model.Email
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change email address, verified",
			args: args{
				es:    GetMockManipulateUserFull(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
		},
		{
			name: "change email not verified, getting code",
			args: args{
				es:    GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &repo_model.Profile{UserName: "UserName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress"}}),
				ctx:   auth.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: false},
			},
			res: res{
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:    GetMockManipulateUser(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:    GetMockManipulateUserNoEvents(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				email: &model.Email{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, EmailAddress: "EmailAddressChanged", IsEmailVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeEmail(tt.args.ctx, tt.args.email)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.EmailAddress != tt.res.email.EmailAddress {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.EmailAddress, result.EmailAddress)
			}
			if tt.res.errFunc == nil && result.IsEmailVerified != tt.res.email.IsEmailVerified {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.email.IsEmailVerified, result.IsEmailVerified)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "verify email code ok",
			args: args{
				es:     GetMockManipulateUserWithEmailCode(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				code: "Code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "Code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.VerifyEmail(tt.args.ctx, tt.args.userID, tt.args.code)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateEmailVerificationCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create email verification code ok",
			args: args{
				es:     GetMockManipulateUserWithEmailCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &repo_model.Profile{UserName: "UserName"}, Email: &repo_model.Email{EmailAddress: "EmailAddress"}}),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no email found",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "already verified",
			args: args{
				es:     GetMockManipulateUserVerifiedEmail(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CreateEmailVerificationCode(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not ger err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestPhoneByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		phone   *model.Phone
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserFull(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumber"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.PhoneByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.PhoneNumber != tt.res.phone.PhoneNumber {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.PhoneNumber, result.PhoneNumber)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangePhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es    *UserEventstore
		ctx   context.Context
		phone *model.Phone
	}
	type res struct {
		phone   *model.Phone
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change phone, verified",
			args: args{
				es:    GetMockManipulateUserFull(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: true},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: true},
			},
		},
		{
			name: "change phone not verified, getting code",
			args: args{
				es:    GetMockManipulateUserWithPhoneCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &repo_model.Profile{UserName: "UserName"}, Phone: &repo_model.Phone{PhoneNumber: "PhoneNumber"}}),
				ctx:   auth.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: false},
			},
			res: res{
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: false},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:    GetMockManipulateUser(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:    GetMockManipulateUserNoEvents(ctrl),
				ctx:   auth.NewMockContext("orgID", "userID"),
				phone: &model.Phone{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, PhoneNumber: "PhoneNumberChanged", IsPhoneVerified: true},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePhone(tt.args.ctx, tt.args.phone)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.PhoneNumber != tt.res.phone.PhoneNumber {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.PhoneNumber, result.PhoneNumber)
			}
			if tt.res.errFunc == nil && result.IsPhoneVerified != tt.res.phone.IsPhoneVerified {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.phone.IsPhoneVerified, result.IsPhoneVerified)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestVerifyPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "verify code ok",
			args: args{
				es:     GetMockManipulateUserWithPhoneCode(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "code",
			},
			res: res{},
		},
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				code: "Code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
				code:   "Code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.VerifyPhone(tt.args.ctx, tt.args.userID, tt.args.code)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not get err %v", err)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreatePhoneVerificationCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "create phone verification code okk",
			args: args{
				es:     GetMockManipulateUserWithPhoneCodeGen(ctrl, repo_model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Profile: &repo_model.Profile{UserName: "UserName"}, Phone: &repo_model.Phone{PhoneNumber: "PhoneNumber"}}),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{},
		},
		{
			name: "no userid",
			args: args{
				es:  GetMockManipulateUser(ctrl),
				ctx: auth.NewMockContext("orgID", "userID"),
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "no phone found",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "already verified",
			args: args{
				es:     GetMockManipulateUserVerifiedPhone(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CreatePhoneVerificationCode(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("should not ger err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddressByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		address *model.Address
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserFull(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "Country"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddressByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Country != tt.res.address.Country {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.address.Country, result.Country)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *UserEventstore
		ctx     context.Context
		address *model.Address
	}
	type res struct {
		address *model.Address
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change address ok",
			args: args{
				es:      GetMockManipulateUserFull(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryChanged"},
			},
			res: res{
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryChanged"},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:      GetMockManipulateUser(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}, Country: "CountryChanged"},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:      GetMockManipulateUserNoEvents(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				address: &model.Address{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}, Country: "CountryCountry"},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeAddress(tt.args.ctx, tt.args.address)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Country != tt.res.address.Country {
				t.Errorf("got wrong result change required: expected: %v, actual: %v ", tt.res.address.Country, result.Country)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestOTPByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "get by id, ok",
			args: args{
				es:       GetMockManipulateUserWithOTP(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.OTPByID(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Errorf("got wrong result change required: actual: %v ", result)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add ok",
			args: args{
				es:     GetMockManipulateUserWithOTPGen(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
		},
		{
			name: "empty userid",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "AggregateID",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddOTP(tt.args.ctx, tt.args.userID)

			if tt.res.errFunc == nil && result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if tt.res.errFunc == nil && result.Url == "" {
				t.Errorf("result has no url")
			}
			if tt.res.errFunc == nil && result.SecretString == "" {
				t.Errorf("result has no url")
			}
			if tt.res.errFunc == nil && result == nil {
				t.Errorf("got wrong result change required: actual: %v ", result)
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestRemoveOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *UserEventstore
		ctx      context.Context
		existing *model.User
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove ok",
			args: args{
				es:       GetMockManipulateUserWithOTP(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
		},
		{
			name: "empty userid",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:       GetMockManipulateUserNoEvents(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user has no otp",
			args: args{
				es:       GetMockManipulateUser(ctrl),
				ctx:      auth.NewMockContext("orgID", "userID"),
				existing: &model.User{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveOTP(tt.args.ctx, tt.args.existing.AggregateID)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCheckOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *UserEventstore
		ctx    context.Context
		userID string
		code   string
	}
	type res struct {
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "empty userid",
			args: args{
				es:   GetMockManipulateUser(ctrl),
				ctx:  auth.NewMockContext("orgID", "userID"),
				code: "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "empty code",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing user not found",
			args: args{
				es:     GetMockManipulateUserNoEvents(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsNotFound,
			},
		},
		{
			name: "user has no otp",
			args: args{
				es:     GetMockManipulateUser(ctrl),
				ctx:    auth.NewMockContext("orgID", "userID"),
				userID: "userID",
				code:   "code",
			},
			res: res{
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.CheckMfaOTP(tt.args.ctx, tt.args.userID, tt.args.code)

			if tt.res.errFunc == nil && err != nil {
				t.Errorf("result should not get err")
			}
			if tt.res.errFunc != nil && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}