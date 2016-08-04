package app

import (
	"encoding/json"
	"fmt"

	"git.wid.la/co-net/auth-server/infrastructure"
	"git.wid.la/co-net/auth-server/middlewares"
	"git.wid.la/co-net/auth-server/models"
	"git.wid.la/co-net/auth-server/utils"
	"github.com/boltdb/bolt"
	"github.com/go-zoo/bone"
	"github.com/solher/zest"

	_ "git.wid.la/co-net/auth-server/interactors"
	_ "git.wid.la/co-net/auth-server/validators"
)

func Run(overrideConst zest.SeqFunc) {
	if overrideConst == nil {
		overrideConst = func(z *zest.Zest) error { return nil }
	}

	appli := zest.New()

	SetCli(appli)

	appli.RegisterSequence = []zest.SeqFunc{
		Register,
	}

	appli.InitSequence = []zest.SeqFunc{
		PopulateConstants,
		overrideConst,
		InitServer,
		SetRoutes,
		ConnectDatabase,
		MigrateDatabase,
		SeedDatabase,
		LaunchGarbageCollector,
	}

	appli.ExitSequence = []zest.SeqFunc{
		CloseDatabase,
	}

	if err := appli.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

func Register(z *zest.Zest) error {
	z.Injector.Register(
		// Router
		bone.New(),
		// Used by the controllers
		zest.NewRender(),
		// Used by the controllers
		infrastructure.NewParamsGetter(),
		// App constants
		NewConstants(),
		// The database
		&bolt.DB{},
		// The garbage collector, used to backup the expired sessions
		NewGarbageCollector,
		// The config importer, used to import config files in DB
		NewConfigImporter,
	)

	return nil
}

func PopulateConstants(z *zest.Zest) error {
	d := &struct{ Const *Constants }{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	d.Const.App.Port = z.Context.GlobalInt("port")
	d.Const.App.ExitTimeout = z.Context.GlobalDuration("exitTimeout")
	d.Const.App.Config = z.Context.GlobalString("config")

	d.Const.Auth.RedirectURL = z.Context.GlobalString("redirectUrl")
	d.Const.Auth.GrantAll = z.Context.GlobalBool("grantAll")

	d.Const.GC.Location = z.Context.GlobalString("gcLocation")
	d.Const.GC.Freq = z.Context.GlobalDuration("gcFreq")

	d.Const.DB.Location = z.Context.GlobalString("dbLocation")
	d.Const.DB.Timeout = z.Context.GlobalDuration("dbTimeout")

	d.Const.Session.Validity = z.Context.GlobalDuration("sessionValidity")
	d.Const.Session.TokenLength = z.Context.GlobalInt("sessionTokenLength")

	return nil
}

func ConnectDatabase(z *zest.Zest) error {
	d := &struct {
		DB    *bolt.DB
		Const *Constants
	}{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	db, err := bolt.Open(d.Const.DB.Location, 0600, &bolt.Options{Timeout: d.Const.DB.Timeout})
	if err != nil {
		return err
	}

	*d.DB = *db

	return nil
}

func InitServer(z *zest.Zest) error {
	d := &struct {
		Router *bone.Mux
		Const  *Constants
	}{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	z.Server.Use(zest.NewRecovery())
	z.Server.Use(zest.NewLogger())
	z.Server.Use(middlewares.NewSwagger())

	z.Server.UseHandler(d.Router)

	z.Port = d.Const.App.Port
	z.ExitTimeout = d.Const.App.ExitTimeout

	return nil
}

func MigrateDatabase(z *zest.Zest) error {
	d := &struct{ DB *bolt.DB }{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	err := d.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("policies"))
		if err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte("resources")); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte("sessions")); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func SeedDatabase(z *zest.Zest) error {
	d := &struct {
		DB       *bolt.DB
		Const    *Constants
		Importer *ConfigImporter
	}{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	if len(d.Const.App.Config) != 0 {
		err := d.DB.Update(func(tx *bolt.Tx) error {
			pc := tx.Bucket([]byte("policies")).Cursor()
			rc := tx.Bucket([]byte("resources")).Cursor()

			for k, _ := pc.First(); k != nil; k, _ = pc.Next() {
				if err := pc.Delete(); err != nil {
					return err
				}
			}

			for k, _ := rc.First(); k != nil; k, _ = rc.Next() {
				if err := rc.Delete(); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		if err := d.Importer.Import(d.Const.App.Config); err != nil {
			return err
		}
	}

	err := d.DB.Update(func(tx *bolt.Tx) error {
		policies := tx.Bucket([]byte("policies"))

		if len(policies.Get([]byte("guest"))) == 0 {
			guestPolicy := &models.Policy{
				Name:        utils.StrCpy("guest"),
				Permissions: []models.Permission{},
			}

			m, _ := json.Marshal(guestPolicy)

			policies.Put([]byte("guest"), m)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func LaunchGarbageCollector(z *zest.Zest) error {
	d := &struct {
		GC    *GarbageCollector
		Const *Constants
	}{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	return d.GC.Run(d.Const.GC.Location, d.Const.GC.Freq)
}

func CloseDatabase(z *zest.Zest) error {
	d := &struct{ DB *bolt.DB }{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	d.DB.Close()

	return nil
}
