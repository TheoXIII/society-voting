package httpcore

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/guildScraper"
	"github.com/CSSUoB/society-voting/internal/httpcore/htmlutil"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"strings"
)

func (endpoints) authLoginPage(ctx *fiber.Ctx) error {
	page := htmlutil.SkeletonPage(config.Get().Platform.SocietyName+" voting",
		//g.If(requestProblem != "",
		//	html.P(g.Textf("Error: %s", requestProblem), g.Attr("style", "color: red")),
		//),
		html.Script(g.Attr("src", "https://unpkg.com/htmx.org@1.9.10"), g.Attr("defer")),
		html.H1(g.Attr("class", "h3 mb-3 fw-normal"), g.Text("Welcome!")),
		html.Div(
			g.Attr("hx-get", loginActionEndpoint),
			g.Attr("hx-trigger", "load"),
			g.Attr("hx-swap", "outerHTML"),
			//g.Text("Loading..."),
		),
		html.P(g.Attr("class", "htmx-indicator"), g.Attr("id", "indicator"), g.Text("Working...")),
	)

	return htmlutil.SendPage(ctx, page)
}

func (endpoints) authLogin(ctx *fiber.Ctx) error {
	var requestProblem string

	var requestData = struct {
		StudentID            string `json:"studentid,omitempty"`
	}{
		StudentID:            strings.TrimSpace(ctx.FormValue("studentid")),
	}

	if ctx.Method() == fiber.MethodGet {
		goto askStudentID
	}

	{
		// Has SID?
		if requestData.StudentID == "" {
			// No: show form
			goto askStudentID
		}

		// SID already registered?
		user, err := database.GetUser(requestData.StudentID)
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return fmt.Errorf("authLogin get user: %w", err)
		}

		if user != nil {
			goto success
		}

		guildMember, err := guildScraper.GetMember(requestData.StudentID)
		if err != nil {
			return fmt.Errorf("authLogin get guild member: %w", err)
		}
		if guildMember == nil {
			requestProblem = "Invalid student ID - it doesn't look like that student ID corresponds to a " + config.Get().Platform.SocietyName + " member."
			goto reset
		}

		user = &database.User{
			StudentID:    requestData.StudentID,
			Name:         guildMember.FirstName + " " + guildMember.LastName,
			IsAdmin:      subtle.ConstantTimeCompare([]byte(requestData.StudentID), []byte(config.Get().Platform.AdminSID)) == 1, //Checks if user is admin
		}

		if err := user.Insert(); err != nil {
			return fmt.Errorf("authLogin insert new user: %w", err)
		}
	}

success:
	ctx.Cookie(newSessionTokenCookie(signData(requestData.StudentID)))
	ctx.Set("HX-Redirect", "/")
	ctx.Status(fiber.StatusNoContent)
	return nil

reset:
askStudentID:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.If(requestProblem != "", html.P(
			g.Attr("style", "color: red;"),
			html.Em(g.Text(requestProblem+" If you're having trouble logging in, please speak to a member of the committee.")),
		)),
		htmlutil.SmallTitle("What's your student ID?"),
		htmlutil.FormInput("text", "studentid", "Your student ID", "Student ID"),
		htmlutil.FormSubmitButton(),
	))

}

func (endpoints) authLogout(ctx *fiber.Ctx) error {
	ctx.Cookie(newSessionTokenDeletionCookie())
	titleLine := config.Get().Platform.SocietyName + " voting"

	return htmlutil.SendPage(ctx, htmlutil.SkeletonPage(
		titleLine,
		html.H1(g.Text("You're all signed out!")),
		html.A(g.Attr("href", "/auth/login"), g.Text("Click here to login again")),
	))
}
