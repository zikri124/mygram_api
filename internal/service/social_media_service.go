package service

import (
	"context"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
)

type SocialMediaService interface {
	PostSocial(ctx context.Context, userId uint32, social model.CreateSocialMedia) (*model.CreateSocialMediaRes, error)
	GetAllComments(ctx context.Context) ([]model.SocialMediaView, error)
}

type socialMediaServiceImpl struct {
	repo repository.SocialMediaRepository
}

func NewSocialMediaService(repo repository.SocialMediaRepository) SocialMediaService {
	return &socialMediaServiceImpl{repo: repo}
}

func (s *socialMediaServiceImpl) PostSocial(ctx context.Context, userId uint32, newSocial model.CreateSocialMedia) (*model.CreateSocialMediaRes, error) {
	social := model.SocialMedia{}
	social.Name = newSocial.Name
	social.SocialMediaUrl = newSocial.SocialMediaUrl
	social.UserId = userId

	err := s.repo.CreateSocial(ctx, &social)
	if err != nil {
		return nil, err
	}

	socialMediaRes := model.CreateSocialMediaRes{}
	socialMediaRes.ID = social.ID
	socialMediaRes.Name = social.Name
	socialMediaRes.UserId = social.UserId
	socialMediaRes.SocialMediaUrl = social.SocialMediaUrl
	socialMediaRes.CreatedAt = social.CreatedAt

	return &socialMediaRes, nil
}

func (s *socialMediaServiceImpl) GetAllComments(ctx context.Context) ([]model.SocialMediaView, error) {
	socials, err := s.repo.GetAllSocials(ctx)
	if err != nil {
		return nil, err
	}

	return socials, nil
}
