// Package reviews implements endpoints used
// to start a review, edit it and accept it
package reviews

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/mailer"
	"github.com/benoitkugler/maths-online/server/src/pass"
	edAPI "github.com/benoitkugler/maths-online/server/src/prof/editor"
	hoAPI "github.com/benoitkugler/maths-online/server/src/prof/homework"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	trAPI "github.com/benoitkugler/maths-online/server/src/prof/trivial"

	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/reviews"
	re "github.com/benoitkugler/maths-online/server/src/sql/reviews"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

var errAccesForbidden = errors.New("internal error: access forbidden")

type Controller struct {
	db    *sql.DB
	admin teacher.Teacher
	smtp  pass.SMTP
}

func NewController(db *sql.DB, admin teacher.Teacher, smtp pass.SMTP) *Controller {
	return &Controller{db: db, admin: admin, smtp: smtp}
}

type ReviewHeader struct {
	Id         re.IdReview
	Title      string // of the target resource
	Kind       re.ReviewKind
	OwnerMail  string
	NbComments int
}

// ReviewsList returns the list of all the reviews.
func (ct *Controller) ReviewsList(c echo.Context) error {
	out, err := ct.listReviews()
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) listReviews() (out []ReviewHeader, err error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	reviews, err := re.SelectAllReviews(tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	comments, err := re.SelectAllReviewParticipations(tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}
	commentsByReview := comments.ByIdReview()

	targets, err := re.LoadTargets(tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	teachers, err := teacher.SelectAllTeachers(tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	for idReview, ta := range targets {
		taHeader, err := ta.Load(tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, utils.SQLError(err)
		}
		tc := teachers[taHeader.Owner]

		nbComments := 0
		for _, teacher := range commentsByReview[idReview] {
			nbComments += len(teacher.Comments)
		}

		out = append(out, ReviewHeader{
			Id:         idReview,
			Title:      taHeader.Title,
			OwnerMail:  tc.Mail,
			Kind:       reviews[idReview].Kind,
			NbComments: nbComments,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })

	return out, nil
}

type ReviewComment struct {
	Comment    re.Comment
	AuthorMail string
	IsOwned    bool
}

type ReviewExt struct {
	Approvals    [3]int          // number per Approval values
	Comments     []ReviewComment // sorted by time (earlier first)
	UserApproval re.Approval     // the approval of the user doing the request, or zero
	IsDeletable  bool            // does the user have the right to delete the review ?
	IsAcceptable bool            // does the user have the right to accept the review ?
}

// ReviewLoad returns the full content of the given review.
func (ct *Controller) ReviewLoad(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.loadReview(re.IdReview(id), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) loadReview(id re.IdReview, userID uID) (out ReviewExt, err error) {
	parts, err := re.SelectReviewParticipationsByIdReviews(ct.db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tcs, err := teacher.SelectAllTeachers(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	for _, part := range parts {
		teacher := tcs[part.IdTeacher]
		out.Approvals[part.Approval] += 1
		for _, comment := range part.Comments {
			out.Comments = append(out.Comments, ReviewComment{
				Comment:    comment,
				AuthorMail: teacher.Mail,
				IsOwned:    part.IdTeacher == userID,
			})
		}

		if part.IdTeacher == userID {
			out.UserApproval = part.Approval
		}
	}

	// only admin may accept review
	out.IsAcceptable = userID == ct.admin.Id

	// admin and owner may delete a review
	targetLink, err := re.LoadTarget(ct.db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	target, err := targetLink.Load(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.IsDeletable = userID == ct.admin.Id || userID == target.Owner

	sort.Slice(out.Comments, func(i, j int) bool {
		ti, tj := out.Comments[i].Comment.Time, out.Comments[j].Comment.Time
		return ti.Before(tj)
	})

	return out, nil
}

type ReviewUpdateCommentsIn struct {
	IdReview re.IdReview
	Comments re.Comments
}

// ReviewUpdateCommnents update all the comments for one review and one teacher.
func (ct *Controller) ReviewUpdateCommnents(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ReviewUpdateCommentsIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateReview(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateReview(args ReviewUpdateCommentsIn, userID uID) error {
	part, has, err := re.SelectReviewParticipationByIdReviewAndIdTeacher(ct.db, args.IdReview, userID)
	if err != nil {
		return utils.SQLError(err)
	}

	if !has { // create the participation on first update
		part = re.ReviewParticipation{IdReview: args.IdReview, IdTeacher: userID}
	}

	part.Comments = args.Comments
	err = re.UpdateParticipation(ct.db, part)

	return err
}

type ReviewUpdateApprovalIn struct {
	IdReview re.IdReview
	Approval re.Approval
}

func (ct *Controller) ReviewUpdateApproval(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ReviewUpdateApprovalIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateApproval(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateApproval(args ReviewUpdateApprovalIn, userID uID) error {
	part, has, err := re.SelectReviewParticipationByIdReviewAndIdTeacher(ct.db, args.IdReview, userID)
	if err != nil {
		return utils.SQLError(err)
	}

	if !has { // create the participation on first update
		part = re.ReviewParticipation{IdReview: args.IdReview, IdTeacher: userID}
	}

	part.Approval = args.Approval
	err = re.UpdateParticipation(ct.db, part)

	return err
}

type ReviewCreateIn struct {
	Kind re.ReviewKind
	Id   int64 // either IdQuestion, IdExercice, IdTrivial, IdSheet
}

// ReviewCreate is trigger by a teacher who wants to publish one of his
// resource.
func (ct *Controller) ReviewCreate(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ReviewCreateIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.createReview(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createReview(args ReviewCreateIn, userID uID) (re.Review, error) {
	var target re.Target
	switch args.Kind {
	case re.KQuestion:
		target = re.ReviewQuestion{IdQuestion: editor.IdQuestiongroup(args.Id), Kind: args.Kind}
	case re.KExercice:
		target = re.ReviewExercice{IdExercice: editor.IdExercicegroup(args.Id), Kind: args.Kind}
	case re.KTrivial:
		target = re.ReviewTrivial{IdTrivial: trivial.IdTrivial(args.Id), Kind: args.Kind}
	case re.KSheet:
		target = re.ReviewSheet{IdSheet: homework.IdSheet(args.Id), Kind: args.Kind}
	default:
		return re.Review{}, fmt.Errorf("internal error: unknown target kind %d", args.Kind)
	}

	header, err := target.Load(ct.db)
	if err != nil {
		return re.Review{}, utils.SQLError(err)
	}
	// check for the owner
	if header.Owner != userID {
		return re.Review{}, errAccesForbidden
	}

	// create the review
	tx, err := ct.db.Begin()
	if err != nil {
		return re.Review{}, utils.SQLError(err)
	}
	review, err := re.Review{Kind: args.Kind}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return re.Review{}, utils.SQLError(err)
	}
	target = target.WithIdReview(review.Id) // use the newly created review
	err = target.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return re.Review{}, utils.SQLError(err)
	}
	err = tx.Commit()
	if err != nil {
		return re.Review{}, utils.SQLError(err)
	}

	return review, nil
}

// ReviewDelete completly delete the review and its messages
func (ct *Controller) ReviewDelete(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteReview(re.IdReview(id), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteReview(id re.IdReview, userID uID) error {
	target, err := reviews.LoadTarget(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	header, err := target.Load(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// delete action is granted for admin and the owner of the review (target)
	if userID != header.Owner && userID != ct.admin.Id {
		return errAccesForbidden
	}

	// all related items cascade
	_, err = reviews.DeleteReviewById(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// ReviewAccept accepts the review, changing the owner of target to admin,
// notifying the creator and deleting the review.
// Only admin account may perform this operation.
func (ct *Controller) ReviewAccept(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.acceptReview(re.IdReview(id), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) acceptReview(id re.IdReview, userID uID) error {
	if userID != ct.admin.Id {
		return errAccesForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	review, err := re.SelectReview(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	target, err := re.LoadTarget(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	header, err := target.Load(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	tc, err := teacher.SelectTeacher(tx, header.Owner)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// step 1 : move to admin
	err = target.MoveToAdmin(tx, ct.admin.Id)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// step 2 : delete the review
	_, err = re.DeleteReviewById(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// step 3 : notify the origin owner
	body := fmt.Sprintf(`Bonjour, <br/>
	
	Votre demande de partage pour la ressource %s (%s) a été acceptée ! <br/><br/>

	Merci infinement pour votre contribution au développement de la plateforme. <br/><br/>

	L'équipe Isyro
	`, header.Title, review.Kind)
	err = mailer.SendMail(ct.smtp, []string{tc.Mail}, "[Isyro] - Partage accepté", body)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// ReviewLoadTarget return the actual content of the review's target,
// so that teachers may try it.
func (ct *Controller) ReviewLoadTarget(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	idReview, err := utils.QueryParamInt64(c, "id-review")
	if err != nil {
		return err
	}

	out, err := ct.loadTargetContent(re.IdReview(idReview), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, LoadTargetOut{out})
}

func (ct *Controller) loadTargetContent(id re.IdReview, userID uID) (TargetContent, error) {
	target, err := re.LoadTarget(ct.db, id)
	if err != nil {
		return nil, err
	}

	switch target := target.(type) {
	case re.ReviewQuestion:
		return ct.loadQuestion(target, userID)
	case re.ReviewExercice:
		return ct.loadExercice(target, userID)
	case re.ReviewTrivial:
		return ct.loadTrivial(target, userID)
	case re.ReviewSheet:
		return ct.loadSheet(target, userID)
	default:
		panic("exhaustive switch")
	}
}

func (ct *Controller) loadQuestion(target re.ReviewQuestion, userID uID) (TargetQuestion, error) {
	group, err := editor.SelectQuestiongroup(ct.db, target.IdQuestion)
	if err != nil {
		return TargetQuestion{}, utils.SQLError(err)
	}

	variants, err := edAPI.LoadQuestionVariants(ct.db, target.IdQuestion)
	if err != nil {
		return TargetQuestion{}, err
	}

	questionTags, err := editor.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, target.IdQuestion)
	if err != nil {
		return TargetQuestion{}, utils.SQLError(err)
	}

	questiongroup := edAPI.NewQuestiongroupExt(group, variants, questionTags.Tags(),
		tcAPI.OptionalIdReview{InReview: true, Id: target.IdReview}, userID, ct.admin.Id)

	allTags, err := edAPI.LoadTags(ct.db, userID)
	if err != nil {
		return TargetQuestion{}, err
	}

	return TargetQuestion{
		Group:    questiongroup,
		Variants: variants,
		AllTags:  allTags,
	}, nil
}

func (ct *Controller) loadExercice(target re.ReviewExercice, userID uID) (TargetExercice, error) {
	group, err := editor.SelectExercicegroup(ct.db, target.IdExercice)
	if err != nil {
		return TargetExercice{}, utils.SQLError(err)
	}

	m, err := editor.SelectExercicesByIdGroups(ct.db, target.IdExercice)
	if err != nil {
		return TargetExercice{}, utils.SQLError(err)
	}
	variants := m.ByGroup()[target.IdExercice]

	tags, err := editor.SelectExercicegroupTagsByIdExercicegroups(ct.db, target.IdExercice)
	if err != nil {
		return TargetExercice{}, utils.SQLError(err)
	}

	exercicegroup := edAPI.NewExercicegroupExt(group, variants, tags.Tags(), tcAPI.OptionalIdReview{InReview: true, Id: target.IdReview}, userID, ct.admin.Id)

	allTags, err := edAPI.LoadTags(ct.db, userID)
	if err != nil {
		return TargetExercice{}, err
	}

	return TargetExercice{
		Group:   exercicegroup,
		AllTags: allTags,
	}, nil
}

func (ct *Controller) loadTrivial(target re.ReviewTrivial, userID uID) (TargetTrivial, error) {
	triv, err := trivial.SelectTrivial(ct.db, target.IdTrivial)
	if err != nil {
		return TargetTrivial{}, utils.SQLError(err)
	}

	// we use the admin ID since it will be the one used after validation
	nbs, err := trAPI.LoadQuestionNumbers(ct.db, triv, ct.admin.Id)
	if err != nil {
		return TargetTrivial{}, err
	}

	return TargetTrivial{Config: triv, NbQuestionsByCategories: nbs}, nil
}

func (ct *Controller) loadSheet(target re.ReviewSheet, userID uID) (TargetSheet, error) {
	sheet, err := hoAPI.LoadSheet(ct.db, target.IdSheet, userID, ct.admin.Id)
	if err != nil {
		return TargetSheet{}, err
	}
	return TargetSheet{Sheet: sheet}, nil
}
