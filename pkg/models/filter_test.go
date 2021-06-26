package models

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLimitOffset(t *testing.T) {
	cases := []struct {
		f      ListFilter
		limit  int
		offset int
	}{
		{ListFilter{Page: 1, PageRows: 5}, 5, 0},
		{ListFilter{Page: 2, PageRows: 5}, 5, 5},
		{ListFilter{Page: 10, PageRows: 5}, 5, 45},
	}

	for _, c := range cases {
		l, o := c.f.GetLimitOffset()
		assert.Equal(t, c.limit, l)
		assert.Equal(t, c.offset, o)
	}
}

func TestParse(t *testing.T) {
	cf := CondFilter{
		QueryStr: "((aa lt 100) AND (bb eq zz))",
	}
	st, err := cf.parse()
	if assert.NoError(t, err) {
		assert.Len(t, st, 1)
	}

	cf = CondFilter{
		QueryStr: "(date eq '2016-05-01 IST') AND ((distance gt 20) OR (distance lt 10))",
	}
	st, err = cf.parse()
	if assert.NoError(t, err) {
		assert.Len(t, st, 1)
	}
}

func TestCheckParens(t *testing.T) {
	assert.NoError(t, checkParens("()"))
	assert.NoError(t, checkParens("((()))"))
	assert.NoError(t, checkParens("(((()))()())"))

	assert.Error(t, checkParens("("))
	assert.Error(t, checkParens(")"))
	assert.Error(t, checkParens("())"))
	assert.Error(t, checkParens("(()"))
	assert.Error(t, checkParens("(((()))()()"))
}

func TestValidate(t *testing.T) {
	// happy paths
	t.Run("basic symbol and str", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(name EQ 'abcd')",
			EntityTableName: "users",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		assert.NoError(t, cf.validate(es))
	})

	t.Run("basic symbol and num", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(seconds EQ 10)",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		assert.NoError(t, cf.validate(es))
	})

	t.Run("basic symbol and num", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(loc EQ 'abcd') AND (seconds GT 20)",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		assert.NoError(t, cf.validate(es))
	})

	// sad paths
	t.Run("incorrect rhs", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(loc EQ abcd)",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		err = cf.validate(es)
		assert.Equal(t, ErrInvorrectElement, errors.Cause(err))
	})
	t.Run("incorrect lhs", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(locx EQ 10)",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		err = cf.validate(es)
		assert.Equal(t, ErrInvorrectElement, errors.Cause(err))
	})
	t.Run("incorrect op", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(loc qqq 10)",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		err = cf.validate(es)
		assert.Equal(t, ErrUnknownOperator, errors.Cause(err))
	})
}

func TestGenerateQuery(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(ts EQ '2016-05-01 IST')",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		require.NoError(t, cf.validate(es))

		q, err := cf.generateQuery(es[0])
		if assert.NoError(t, err) {
			assert.Equal(t, "(ts = '2016-05-01 IST')", q)
		}
	})

	t.Run("complex case", func(t *testing.T) {
		cf := CondFilter{
			QueryStr:        "(ts EQ '2016-05-01 IST') AND ((distance GT 20) OR (distance LT 10))",
			EntityTableName: "activities",
		}
		es, err := cf.parse()
		require.NoError(t, err)
		require.NoError(t, cf.validate(es))

		q, err := cf.generateQuery(es[0])
		if assert.NoError(t, err) {
			assert.Equal(t, "((ts = '2016-05-01 IST') AND ((distance > 20) OR (distance < 10)))", q)
		}
	})
}
