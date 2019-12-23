package edit

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/model"
	"github.com/jansorg/tom/go-tom/money"
	"github.com/jansorg/tom/go-tom/util"
	"github.com/jansorg/tom/go-tom/util/tristate"
)

func newEditProjectCommand(ctx *context.TomContext, parent *cobra.Command) *cobra.Command {
	var name string
	var parentNameOrID string
	var nameDelimiter string
	var hourlyRate string
	var noteRequired string

	var cmd = &cobra.Command{
		Use:   "project fullName | ID",
		Short: "edit properties of a project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("name").Changed && len(name) == 0 {
				util.Fatal("unable to use empty project name")
			} else if !cmd.Flag("name").Changed && !cmd.Flag("parent").Changed && !cmd.Flag("hourly-rate").Changed && !cmd.Flag("note-required").Changed {
				util.Fatalf("no modification defined, use --name, --parent, or --hourly-rate to update project data")
			}

			var parent *string
			if cmd.Flag("parent").Changed {
				parent = &(parentNameOrID)
			}

			var hourlyRateValue *string
			if cmd.Flag("hourly-rate").Changed {
				hourlyRateValue = &hourlyRate
			}

			var noteRequiredValue *tristate.Tristate
			if cmd.Flag("note-required").Changed {
				value, err := tristate.FromString(noteRequired)
				if err != nil {
					util.Fatalf("unable to parse value %s for note-required. Valid values: true, false, empty value", noteRequired)
				}
				noteRequiredValue = &value
			}

			if err := doEditProjectCommand(name, parent, nameDelimiter, hourlyRateValue, noteRequiredValue, args, ctx); err != nil {
				util.Fatal(err)
			} else {
				println("Successfully updated project data")
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "update the project name")
	cmd.Flags().StringVarP(&parentNameOrID, "parent", "p", "", "update the parent. Use an empty ID to make it a top-level project. A project keeps all frames and subprojects when it's assigned to a new parent project.")
	cmd.Flags().StringVarP(&nameDelimiter, "name-delimiter", "", "/", "Delimiter used in full project names")
	cmd.Flags().StringVarP(&hourlyRate, "hourly-rate", "", "", "Optional hourly rate which applies to this project and all subproject without hourly rate values")
	cmd.Flags().StringVarP(&noteRequired, "note-required", "", "", "An optional flag to enforce a note for time entries of this project and all subprojects, where this setting is not turned off.")

	parent.AddCommand(cmd)
	return cmd
}

func doEditProjectCommand(newName string, parentNameOrID *string, nameDelimiter string, hourlyRate *string, noteRequired *tristate.Tristate, projectIDsOrNames []string, ctx *context.TomContext) error {
	var err error
	var parentProjectID string

	// a non-nil, but empty parentNameOrID points to the top-level
	if parentNameOrID != nil && *parentNameOrID != "" {
		if parent, err := ctx.Query.ProjectByFullNameOrID(*parentNameOrID, nameDelimiter); err != nil {
			return fmt.Errorf("parent project %s not found", *parentNameOrID)
		} else {
			parentProjectID = parent.ID
		}
	}

	var parsedHourlyRate *money.Money
	if hourlyRate != nil {
		if *hourlyRate == "" {
			// remove the current value
			parsedHourlyRate = nil
		} else {
			if parsedHourlyRate, err = money.Parse(*hourlyRate); err != nil {
				return nil
			}
		}
	}

	// batch mode to handle many projects at once
	ctx.Store.StartBatch()
	defer ctx.Store.StopBatch()

	var projects []*model.Project
	for _, idOrName := range projectIDsOrNames {
		var project *model.Project
		if project, err = ctx.Query.ProjectByID(idOrName); err != nil {
			if project, err = ctx.Query.ProjectByFullName(strings.Split(idOrName, nameDelimiter)); err != nil {
				return err
			}
		}
		projects = append(projects, project)
	}

	for _, p := range projects {
		if len(newName) > 0 {
			p.Name = newName
		}

		if hourlyRate != nil {
			p.SetHourlyRate(parsedHourlyRate)
		}

		if noteRequired != nil {
			p.SetNoteRequired(noteRequired.ToBool())
		}

		if parentNameOrID != nil {
			if p, err = ctx.StoreHelper.MoveProject(p, parentProjectID); err != nil {
				return err
			}
		}

		if _, err = ctx.Store.UpdateProject(*p); err != nil {
			return err
		}
	}

	return nil
}
