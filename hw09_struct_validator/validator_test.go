package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{Version: "123"},
			ValidationErrors{ValidationError{Err: ErrValidateStringLength, Field: "Version"}},
		},
		{
			App{Version: "12345"},
			nil,
		},
		{
			User{
				ID:     "123",
				Name:   "test",
				Age:    15,
				Email:  "test.com",
				Role:   "stuff",
				Phones: []string{"12312312312"},
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{Err: ErrValidateStringLength, Field: "ID"},
				ValidationError{Err: ErrValidateMinimumValue, Field: "Age"},
				ValidationError{Err: ErrValidateRegexp, Field: "Email"},
			},
		},
		{
			Token{Header: nil, Payload: nil, Signature: nil},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.Nil(t, err)
				return
			}
			require.IsType(t, ValidationErrors{}, err)

			expectedErr, ok := tt.expectedErr.(ValidationErrors) //nolint: errorlint
			require.True(t, ok)

			originalErr, ok := err.(ValidationErrors) //nolint: errorlint
			require.True(t, ok)
			require.Len(t, originalErr, len(expectedErr))

			for i := 0; i < len(expectedErr); i++ {
				exp := expectedErr[i]
				orig := originalErr[i]
				require.Equal(t, exp.Field, orig.Field)
				require.ErrorIs(t, orig.Err, exp.Err)
			}
		})
	}
}

func TestMakeValidationRulesFromRawData(t *testing.T) {
	t.Run("make MIN MAX rules ", func(t *testing.T) {
		rawData := "min:26|max:65"
		result, _ := makeValidationRulesFromRawData(rawData)
		expected := ValidationRules{
			ValidationRule{
				Operator: Min,
				Value:    26,
			},
			ValidationRule{
				Operator: Max,
				Value:    65,
			},
		}
		require.Equal(t, expected, result)
	})

	t.Run("make REGEXP LEN rules ", func(t *testing.T) {
		rawData := "regexp:\\\\d+|len:11"
		result, _ := makeValidationRulesFromRawData(rawData)
		re := regexp.MustCompile(`\\d+`)
		expected := ValidationRules{
			ValidationRule{
				Operator: Regexp,
				Value:    re,
			},
			ValidationRule{
				Operator: Len,
				Value:    11,
			},
		}
		require.Equal(t, expected, result)
	})

	t.Run("make IN rules ", func(t *testing.T) {
		rawData := "in:1,2,3"
		result, _ := makeValidationRulesFromRawData(rawData)
		expected := ValidationRules{
			ValidationRule{
				Operator: In,
				Value:    []string{"1", "2", "3"},
			},
		}
		require.Equal(t, expected, result)
	})
}

func TestValidateValueByRegexp(t *testing.T) {
	t.Run("positive validate regexp ", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		result := validateValueByRegexp("123", re)

		require.True(t, result)
	})

	t.Run("negative validate regexp ", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		result := validateValueByRegexp("abc", re)

		require.False(t, result)
	})
}

func TestValidateValueContainsInRule(t *testing.T) {
	t.Run("positive validate contains string", func(t *testing.T) {
		slice := []string{"1", "2", "3"}
		result := validateValueContainsInRule("3", slice)

		require.True(t, result)
	})

	t.Run("positive validate contains int", func(t *testing.T) {
		slice := []string{"1", "2", "3"}
		result := validateValueContainsInRule(2, slice)

		require.True(t, result)
	})

	t.Run("negative validate contains ", func(t *testing.T) {
		slice := []string{"1", "2", "3"}
		result := validateValueContainsInRule("6", slice)

		require.False(t, result)
	})
}

func TestValidateStringFieldByRules(t *testing.T) {
	t.Run("positive validate by regexp and len", func(t *testing.T) {
		re := regexp.MustCompile("test")
		rules := ValidationRules{
			ValidationRule{
				Operator: Len,
				Value:    11,
			},
			ValidationRule{
				Operator: Regexp,
				Value:    re,
			},
		}
		err := validateStringFieldByRules("test string", rules)
		require.NoError(t, err)
	})

	t.Run("positive validate by contains", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: In,
				Value:    []string{"test"},
			},
		}
		err := validateStringFieldByRules("test", rules)
		require.NoError(t, err)
	})

	t.Run("negative validate by len", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: Len,
				Value:    12,
			},
		}
		err := validateStringFieldByRules("test string", rules)
		require.ErrorIs(t, err, ErrValidateStringLength)
	})

	t.Run("negative validate by regexp", func(t *testing.T) {
		re := regexp.MustCompile("test")
		rules := ValidationRules{
			ValidationRule{
				Operator: Regexp,
				Value:    re,
			},
		}
		err := validateStringFieldByRules("tes string", rules)
		require.ErrorIs(t, err, ErrValidateRegexp)
	})

	t.Run("negative validate by contains", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: In,
				Value:    []string{"test"},
			},
		}
		err := validateStringFieldByRules("tes", rules)
		require.ErrorIs(t, err, ErrValidateContain)
	})
}

func TestValidateIntFieldByRules(t *testing.T) {
	t.Run("positive validate int by min", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: Min,
				Value:    10,
			},
		}
		err := validateIntFieldByRules(12, rules)
		require.NoError(t, err)
	})

	t.Run("positive validate int by max", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: Max,
				Value:    10,
			},
		}
		err := validateIntFieldByRules(5, rules)
		require.NoError(t, err)
	})

	t.Run("positive validate int by in", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: In,
				Value:    []string{"1", "55", "32"},
			},
		}
		err := validateIntFieldByRules(32, rules)
		require.NoError(t, err)
	})

	t.Run("negative validate int by min", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: Min,
				Value:    10,
			},
		}
		err := validateIntFieldByRules(5, rules)
		require.ErrorIs(t, err, ErrValidateMinimumValue)
	})

	t.Run("negative validate int by max", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: Max,
				Value:    10,
			},
		}
		err := validateIntFieldByRules(15, rules)
		require.ErrorIs(t, err, ErrValidateMaximumValue)
	})

	t.Run("negative validate int by in", func(t *testing.T) {
		rules := ValidationRules{
			ValidationRule{
				Operator: In,
				Value:    []string{"1", "55", "32"},
			},
		}
		err := validateIntFieldByRules(22, rules)
		require.ErrorIs(t, err, ErrValidateContain)
	})
}

func TestPrepareValidateValue(t *testing.T) {
	t.Run("positive make value MIN", func(t *testing.T) {
		result, err := prepareValidateValue(Min, "1")

		require.NoError(t, err)
		require.Equal(t, 1, result)
	})

	t.Run("positive make value MAX", func(t *testing.T) {
		result, err := prepareValidateValue(Max, "5")

		require.NoError(t, err)
		require.Equal(t, 5, result)
	})

	t.Run("positive make value LEN", func(t *testing.T) {
		result, err := prepareValidateValue(Len, "11")

		require.NoError(t, err)
		require.Equal(t, 11, result)
	})

	t.Run("positive make value IN", func(t *testing.T) {
		result, err := prepareValidateValue(In, "11,12,55")

		require.NoError(t, err)
		require.Equal(t, []string{"11", "12", "55"}, result)
	})

	t.Run("positive make value REGEXP", func(t *testing.T) {
		re := regexp.MustCompile(`\d+`)
		result, err := prepareValidateValue(Regexp, `\d+`)

		require.NoError(t, err)
		require.Equal(t, re, result)
	})

	t.Run("negative make value REGEXP", func(t *testing.T) {
		_, err := prepareValidateValue(Regexp, "\\d+([")
		require.Error(t, err)
	})

	t.Run("negative make value MIN", func(t *testing.T) {
		_, err := prepareValidateValue(Min, "\\d+([")
		require.Error(t, err)
	})
}

func TestDetectValidationOperator(t *testing.T) {
	t.Run("positive detect MIN", func(t *testing.T) {
		result, err := detectValidationOperator("min")

		require.NoError(t, err)
		require.Equal(t, Min, result)
	})

	t.Run("positive detect MAX", func(t *testing.T) {
		result, err := detectValidationOperator("min")

		require.NoError(t, err)
		require.Equal(t, Min, result)
	})

	t.Run("positive detect MAX", func(t *testing.T) {
		result, err := detectValidationOperator("max")

		require.NoError(t, err)
		require.Equal(t, Max, result)
	})

	t.Run("positive detect MAX", func(t *testing.T) {
		result, err := detectValidationOperator("len")

		require.NoError(t, err)
		require.Equal(t, Len, result)
	})

	t.Run("positive detect REGEXP", func(t *testing.T) {
		result, err := detectValidationOperator("regexp")

		require.NoError(t, err)
		require.Equal(t, Regexp, result)
	})

	t.Run("negative detect operator", func(t *testing.T) {
		_, err := detectValidationOperator("non_impl")
		require.ErrorIs(t, err, ErrValidateTagNotImplemented)
	})
}
