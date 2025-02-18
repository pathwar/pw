package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
	"pathwar.land/pathwar/v2/go/pkg/pwdb"
)

func (svc *service) AdminListActivities(ctx context.Context, in *AdminListActivities_Input) (*AdminListActivities_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}
	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var activities []*pwdb.Activity
	req := svc.db.
		Preload("Author").
		Preload("Team").
		Preload("User").
		Preload("Agent").
		Preload("Organization").
		Preload("Season").
		Preload("Challenge").
		Preload("ChallengeFlavor").
		Preload("ChallengeInstance").
		Preload("Coupon").
		Preload("SeasonChallenge").
		Preload("TeamMember").
		Preload("ChallengeSubscription").
		Order("created_at DESC")
	if in.Limit > 0 {
		req = req.Limit(in.Limit)
	}
	if in.Since != nil {
		req = req.Where("created_at > ?", *in.Since)
	}
	if in.GetTo() != nil && !in.GetTo().IsZero() {
		req = req.Where("created_at < ?", *in.To)
	}
	switch in.FilteringPreset {
	case "default", "":
	// noop
	case "registers":
		req = req.Where(&pwdb.Activity{Kind: pwdb.Activity_UserRegister})
	case "validations":
		req = req.Where(&pwdb.Activity{Kind: pwdb.Activity_ChallengeSubscriptionValidate})
	default:
		return nil, errcode.TODO
	}

	if err := req.Find(&activities).Error; err != nil {
		return nil, errcode.ErrListActivities.Wrap(err)
	}

	for i, j := 0, len(activities)-1; i < j; i, j = i+1, j-1 {
		activities[i], activities[j] = activities[j], activities[i]
	}

	out := AdminListActivities_Output{Activities: activities}
	return &out, nil
}
