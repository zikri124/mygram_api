package service

import (
	"context"

	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/repository"
)

type SocialMediaService interface {
	PostSocial(ctx context.Context, userId uint32, social model.NewSocialMedia) (*model.CreateSocialMediaRes, error)
	GetAllSocialMediasByUserId(ctx context.Context, userId uint32) ([]model.SocialMediaView, error)
	GetSocialById(ctx context.Context, socialId uint32) (*model.SocialMedia, error)
	UpdateSocial(ctx context.Context, social model.SocialMedia) (*model.UpdateSocialMediaRes, error)
	DeleteSocial(ctx context.Context, socialId uint32) error
}

type socialMediaServiceImpl struct {
	repo repository.SocialMediaRepository
}

func NewSocialMediaService(repo repository.SocialMediaRepository) SocialMediaService {
	return &socialMediaServiceImpl{repo: repo}
}

func (s *socialMediaServiceImpl) PostSocial(ctx context.Context, userId uint32, newSocial model.NewSocialMedia) (*model.CreateSocialMediaRes, error) {
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

func (s *socialMediaServiceImpl) GetAllSocialMediasByUserId(ctx context.Context, userId uint32) ([]model.SocialMediaView, error) {
	socials, err := s.repo.GetAllSocialMediasByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return socials, nil
}

func (s *socialMediaServiceImpl) GetSocialById(ctx context.Context, socialId uint32) (*model.SocialMedia, error) {
	social, err := s.repo.GetSocialById(ctx, socialId)
	if err != nil {
		return nil, err
	}

	return social, nil
}

func (s *socialMediaServiceImpl) UpdateSocial(ctx context.Context, social model.SocialMedia) (*model.UpdateSocialMediaRes, error) {
	err := s.repo.UpdateSocial(ctx, &social)

	if err != nil {
		return nil, err
	}

	socialRes := model.UpdateSocialMediaRes{}
	socialRes.ID = social.ID
	socialRes.UserId = social.UserId
	socialRes.Name = social.Name
	socialRes.SocialMediaUrl = social.SocialMediaUrl
	socialRes.UpdatedAt = social.UpdatedAt

	return &socialRes, nil
}

func (s *socialMediaServiceImpl) DeleteSocial(ctx context.Context, socialId uint32) error {
	return s.repo.DeleteSocial(ctx, socialId)
}
