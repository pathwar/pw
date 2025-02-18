package pwapi

import (
	"context"

	"pathwar.land/pathwar/v2/go/pkg/pwdb"

	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

func (svc *service) AdminListSeasonChallenges(ctx context.Context, in *AdminListSeasonChallenges_Input) (*AdminListSeasonChallenges_Output, error) {
	if !isAdminContext(ctx) {
		return nil, errcode.ErrRestrictedArea
	}

	if in == nil {
		return nil, errcode.ErrMissingInput
	}

	var seasonChallenges []*pwdb.SeasonChallenge
	if len(in.Id) == 0 {
		svc.db.Find(&seasonChallenges)
	} else {
		svc.db.Find(&seasonChallenges, in.Id)
	}

	return &AdminListSeasonChallenges_Output{SeasonChallenge: seasonChallenges}, nil
}
