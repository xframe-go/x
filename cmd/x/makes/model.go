package makes

import (
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/utils/generator"
	"github.com/xframe-go/x/xdb"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func createMakeModelCmd(app contracts.Application) *cobra.Command {
	var (
		connection string
		tables     string
	)
	conf := generator.DefaultConfig
	cfg, ok := app.Config().Get(generator.ConfigSymbol)
	if ok {
		conf = cfg.(*generator.Config)
	}

	cmd := &cobra.Command{
		Use:     "make:model",
		GroupID: "make",
		Run: func(c *cobra.Command, args []string) {
			g := gen.NewGenerator(gen.Config{
				OutPath:           conf.DaoPath,
				OutFile:           conf.OutFile,
				ModelPkgPath:      conf.ModelPath,
				WithUnitTest:      conf.WithUnitTest,
				FieldNullable:     conf.FieldNullable,
				FieldCoverable:    conf.FieldCoverable,
				FieldSignable:     conf.FieldSignable,
				FieldWithIndexTag: conf.FieldWithIndexTag,
				FieldWithTypeTag:  conf.FieldWithTypeTag,
				Mode:              conf.GenerateMode,
			})

			tx := app.DB().Connection(connection)

			g.UseDB(tx)

			g.WithDataTypeMap(map[string]func(columnType gorm.ColumnType) (dataType string){
				"datetime": xdb.FieldDatetime,
			})

			g.WithImportPkgPath("github.com/dromara/carbon/v2")

			models := g.GenerateAllTable()
			g.ApplyBasic(models...)

			if conf.Config != nil {
				conf.Config(g)
			}

			g.Execute()
		},
	}

	cmd.Flags().StringVarP(&connection, "connection", "c", "", "db connection")
	cmd.Flags().StringVarP(&tables, "table", "t", "", "table name,multiple tables can be concatenated using commas")

	return cmd
}
