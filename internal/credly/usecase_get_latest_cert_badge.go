package credly

type GetLatestCertBadgeUseCase struct {
	input GetLatestCertBadgeUseCaseInput
}

func NewGetLatestCertBadgeUseCase(input GetLatestCertBadgeUseCaseInput) GetLatestCertBadgeUseCase {
	return GetLatestCertBadgeUseCase{input: input}
}

//func New(r Repository) *UseCase{
//	return &UseCase{
//		repo: r,
//	}
//}
//
type GetLatestCertBadgeUseCaseInput struct {
	Client   CredlyClient
	Username string
}

type BadgeSimple struct {
	ImageURL string
}

func (uc *GetLatestCertBadgeUseCase) Do() (BadgeSimple, error) {
	badge, err := uc.input.Client.GetMostRecentBadgeByUsername(uc.input.Username)
	if err != nil {
		return BadgeSimple{}, err
	}

	return BadgeSimple{
		ImageURL: badge.BadgeTemplate.ImageURL,
	}, nil
}
