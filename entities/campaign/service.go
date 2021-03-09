package campaign

type Service interface {
	GetAllCampaigns() ([]Campaign, error)
	GetCampaigsByUserID(userID int) ([]Campaign, error)
	GetCampaignByID(id int) (Campaign, error)
	// CreateCampaign(campaign Campaign) (Campaign, error)
	// UpdateCampaign(campaign Campaign) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetAllCampaigns() ([]Campaign, error) {
	campaigns, err := s.repository.All()

	if err != nil {
		return campaigns, err
	}

	return campaigns, err
}

func (s *service) GetCampaigsByUserID(userID int) ([]Campaign, error) {
	campaigns, err := s.repository.AllByUserID(userID)

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(id int) (Campaign, error) {
	campaign, err := s.repository.Get(id)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
