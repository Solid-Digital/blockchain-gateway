package component_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"github.com/unchainio/pkg/errors"

	"bitbucket.org/unchain/ares/pkg/component"

	"github.com/davecgh/go-spew/spew"

	stderr "errors"

	pkgerr "github.com/pkg/errors"
	"github.com/unchainio/pkg/xlogger"

	"bitbucket.org/unchain/ares/pkg/ares"

	"bitbucket.org/unchain/ares/gen/dto"
	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/testhelper"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestService_CreateAction() {
	cases := map[string]struct {
		Request      *dto.CreateComponentRequest
		Organization *orm.Organization
		User         *dto.User
		Success      bool
	}{
		"create action": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(true),
			s.factory.DTOUser(true),
			true,
		},
		"organization does not exist": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(false),
			s.factory.DTOUser(true),
			false,
		},
		"user does not exist": {
			&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
			s.factory.Organization(true),
			s.factory.DTOUser(false),
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.service.CreateAction(tc.Request, tc.Organization.Name, tc.User)

			if tc.Success {
				xrequire.NoError(t, err)
				require.Equal(t, tc.Request.Name, *response.Name)
			} else {
				xrequire.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestDuplicateAction() {
	s.T().Skip()
	tc := struct {
		Request      *dto.CreateComponentRequest
		Organization *orm.Organization
		User         *dto.User
		Success      bool
	}{
		&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
		s.factory.Organization(true),
		s.factory.DTOUser(true),
		true,
	}

	//sql.Error = stderr.New("??")
	response, err := s.service.CreateAction(tc.Request, tc.Organization.Name, tc.User)
	_ = response
	//sql.Error = nil
	//s.Suite.Require().Nil(err)
	//s.Suite.Require().Equal(tc.Request.Name, *response.Name)
	//response, err = s.service.CreateAction(tc.Request, tc.Organization.Name, tc.User)
	err2 := ares.ParsePQErr(err)

	spew.Config.DisableMethods = true
	spew.Dump(err2)
	fmt.Println()

	log := xlogger.NewSimpleLogger()
	apperr.HandleError(err, "req1", log)

	b, _ := json.Marshal(err)

	fmt.Printf("%s\n", string(b))
}

func (s *TestSuite) TestDuplicateAction4() {
	tc := struct {
		Request      *dto.CreateComponentRequest
		Organization *orm.Organization
		User         *dto.User
		Success      bool
	}{
		&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
		s.factory.Organization(true),
		s.factory.DTOUser(true),
		true,
	}

	//sql.Error = stderr.New("??")
	response, err := s.service.CreateAction(tc.Request, "zz", tc.User)
	_ = response
	//sql.Error = nil
	err2 := ares.ParsePQErr(err)

	spew.Config.DisableMethods = true
	spew.Dump(err2)
	fmt.Println()

	log := xlogger.NewSimpleLogger()
	apperr.HandleError(err, "req1", log)

	b, _ := json.Marshal(err)

	fmt.Printf("%s\n", string(b))
}

func (s *TestSuite) TestDuplicateAction2() {
	tc := struct {
		Request      *dto.CreateComponentRequest
		Organization *orm.Organization
		User         *dto.User
		Success      bool
	}{
		&dto.CreateComponentRequest{Name: testhelper.Randumb(randomdata.SillyName())},
		s.factory.Organization(true),
		s.factory.DTOUser(true),
		true,
	}

	//var tx, tx2 *sql.Tx

	go func() {
		spew.Dump("begin1")
		ctx, tx, err := s.ares.DB.Begin()
		spew.Dump("begin1 - done")

		if err != nil {
			spew.Dump("err")
			spew.Dump(err)
			return
		}

		spew.Dump("exec1")
		_, _, err = component.CreateActionTx(ctx, tx, tc.Request, tc.Organization.Name, tc.User)
		spew.Dump("exec1 - done")

		if err != nil {
			spew.Dump("err")
			spew.Dump(err)
			return
		}

		time.Sleep(5 * time.Second)
		spew.Dump("commit1")
		err = tx.Commit()
		spew.Dump("commit1 - done")

		if err != nil {
			spew.Dump("commit1 err")
			spew.Dump(err)
			return
		}
	}()

	go func() {
		spew.Dump("begin2")
		ctx2, tx2, err2 := s.ares.DB.Begin()
		spew.Dump("begin2 - done")

		if err2 != nil {
			spew.Dump("err2")
			spew.Dump(err2)
			return
		}

		spew.Dump("exec2")
		time.Sleep(2 * time.Second)
		_, _, err2 = component.CreateActionTx(ctx2, tx2, tc.Request, tc.Organization.Name, tc.User)
		spew.Dump("exec2 - done")

		//if err2 != nil {
		//	spew.Dump("err2")
		//	spew.Dump(err2)
		//	return
		//}

		spew.Dump("commit2")
		err2 = tx2.Commit()
		spew.Dump("commit2 - done")

		if err2 != nil {
			spew.Dump("commit2 err2")
			spew.Dump(err2)
			return
		}
	}()

	time.Sleep(10 * time.Second)
	//ares.DebugPQErr(err)
	//fmt.Println()
	//
	//spew.Dump(err)
	//ares.HandleError(err, xlogger.NewSimpleLogger())
}

func TestName(t *testing.T) {
	var err error

	err = pkgerr.New("test err")
	err = errors.New("test err")
	err = stderr.New("test err")

	//err = errors.WithMessage(err, "more context")
	err = errors.Wrap(err, "wrap")
	err = errors.WithMessage(err, "more context 2")
	//err := errors.New("test err")
	fmt.Printf("%+v", err)
	//fmt.Printf("%+v")

	//printError(os.Stdout, err)
}
