package translation

import (
	"errors"
	"go.temporal.io/sdk/temporal"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

func TestSuccessfulTranslateActivityHelloGerman(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(TranslateTerm)

	input := TranslationActivityInput{
		Term:         "Hello",
		LanguageCode: "de",
	}

	future, err := env.ExecuteActivity(TranslateTerm, input)
	assert.NoError(t, err)

	var output TranslationActivityOutput
	assert.NoError(t, future.Get(&output))
	assert.Equal(t, "Hallo", output.Translation)
}

func TestSuccessfulTranslateActivityGoodbyeLatvian(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(TranslateTerm)

	input := TranslationActivityInput{
		Term:         "Goodbye",
		LanguageCode: "lv",
	}

	future, err := env.ExecuteActivity(TranslateTerm, input)
	assert.NoError(t, err)

	var output TranslationActivityOutput
	assert.NoError(t, future.Get(&output))
	assert.Equal(t, "Ardievu", output.Translation)
}

func TestFailedTranslateActivityBadLanguageCode(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(TranslateTerm)

	input := TranslationActivityInput{
		Term:         "Hello",
		LanguageCode: "bad",
	}
	_, err := env.ExecuteActivity(TranslateTerm, input)

	var applicationErr *temporal.ApplicationError
	assert.True(t, errors.As(err, &applicationErr))

	assert.Equal(t, "HTTP Error 400: Unknown language code 'bad'\n", applicationErr.Error())
}
