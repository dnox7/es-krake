package main

func main() {
	// cfg := config.NewConfig()
	//
	// ctx := context.Background()
	// log.Initialize(os.Stdout, cfg, []string{"request-id", "recurringID"})
	//
	// vault, err := vaultcli.NewVaultAppRoleClient(ctx)
	//
	// pg := rdb.NewOrGetSingleton(cfg)
	// defer pg.Close()
	//
	// poolLogger := pg.StartLoggingPoolSize()
	// defer poolLogger()
	//
	// err := pg.Ping(ctx)
	// if err != nil {
	// 	slog.Error("database ping failed", "detail", err)
	// 	return
	// }
	//
	// migrateType := flag.String("type", "up", "Migration type: up, down, step (required)")
	// step := flag.Int("st", 0, "Number of steps for 'step' action")
	// module := flag.String("module", "", "Name of the module to migrate (required)")
	//
	// switch *migrateType {
	// case "step":
	// 	err = migration.MigrateStep(cfg, pg.Conn(), *module, *step)
	// case "down":
	// 	err = migration.MigrateDown(cfg, pg.Conn())
	// default:
	// 	err = migration.MigrateUp(cfg, pg.Conn())
	// }
	//
	// if err != nil {
	// 	slog.Error("cannot migrate database", "detail", err)
	// 	os.Exit(1)
	// }
}
