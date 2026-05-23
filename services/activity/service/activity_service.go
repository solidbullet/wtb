package service
import ("errors"; "github.com/wtb-ordering/services/activity/model"; "github.com/wtb-ordering/services/activity/repository")
type ActivityService struct {
	annRepo *repository.AnnouncementRepo; actRepo *repository.ActivityRepo; regRepo *repository.RegistrationRepo }
func NewActivityService(annRepo *repository.AnnouncementRepo, actRepo *repository.ActivityRepo, regRepo *repository.RegistrationRepo) *ActivityService {
	return &ActivityService{annRepo: annRepo, actRepo: actRepo, regRepo: regRepo} }
func (s *ActivityService) ListAnnouncements() ([]model.Announcement, error) { return s.annRepo.ListActive() }
func (s *ActivityService) ListActivities() ([]model.Activity, error) { return s.actRepo.ListPublished() }
func (s *ActivityService) Register(userID, activityID uint, name, phone, remark string) (*model.ActivityRegistration, error) {
	act, err := s.actRepo.FindByID(activityID); if err != nil { return nil, errors.New("活动不存在") }
	if act.MaxParticipants > 0 && act.CurrentParticipants >= act.MaxParticipants { return nil, errors.New("活动名额已满") }
	if _, err := s.regRepo.FindByUserAndActivity(userID, activityID); err == nil { return nil, errors.New("已报名") }
	reg := &model.ActivityRegistration{ActivityID: activityID, UserID: userID, Name: name, Phone: phone, Remark: remark}
	if err := s.regRepo.Create(reg); err != nil { return nil, err }; s.actRepo.IncrementParticipants(activityID); return reg, nil }
func (s *ActivityService) MyRegistrations(userID uint) ([]model.ActivityRegistration, error) { return s.regRepo.ListByUser(userID) }
func (s *ActivityService) CancelRegistration(userID, activityID uint) error {
	reg, err := s.regRepo.FindByUserAndActivity(userID, activityID); if err != nil { return errors.New("未报名") }; return s.regRepo.Cancel(reg.ID) }
func (s *ActivityService) CreateAnnouncement(a *model.Announcement) error { return s.annRepo.Create(a) }
func (s *ActivityService) CreateActivity(a *model.Activity) error { return s.actRepo.Create(a) }
