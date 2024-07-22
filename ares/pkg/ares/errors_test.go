package ares_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"bitbucket.org/unchain/ares/gen/dto"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"github.com/bxcodec/faker"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func e() *apperr.Error {
	return apperr.New()
}

func TestErrorWraps(t *testing.T) {
	err := apperr.Conflict.Wrap(fmt.Errorf("1 : "))
	err = apperr.Internal.Wrap(err)
	//err = ErrInternal2.Wrap(err)

	require.True(t, errors.Is(err, apperr.Internal))
	require.True(t, errors.Is(err, apperr.Conflict))
	spew.Config.DisableMethods = true
	//spew.Dump(ErrInternal2)

	require.True(t, e().Is(e()))

	err2 := sql.ErrNoRows
	err2 = apperr.NotFound.Wrap(err2).WithMessagef("asdf")
	require.True(t, errors.Is(err2, apperr.NotFound))
	require.True(t, errors.Is(err2, sql.ErrNoRows))
}

func TestDTOErrors(t *testing.T) {
	dtoErr := &dto.ErrorResponse{}
	appErr := &apperr.Error{}

	for i := 0; i < 100; i++ {
		eq(t, dtoErr, appErr)
		eq(t, appErr, dtoErr)
	}
}

func eq(t *testing.T, dst, src interface{}) {
	dst = reflect.New(reflect.TypeOf(dst)).Interface()
	src = reflect.New(reflect.TypeOf(src)).Interface()

	err := faker.FakeData(src)
	require.NoError(t, err)

	data1, err := json.Marshal(src)
	require.NoError(t, err)
	require.NotEmpty(t, data1)

	err = json.Unmarshal(data1, dst)
	require.NoError(t, err)

	data2, err := json.Marshal(dst)
	require.NoError(t, err)
	require.NotEmpty(t, data2)

	//spew.Dump("data1: ", string(data1), "data2: ", string(data2))
	require.Equal(t, string(data1), string(data2))
}
