package svelte_check_wrapper

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

var exampleOutput = `
1686094808695 START "/home/matt/godev/secretsanta/frontend"

 (x2)
 (x3)
 (x4)
 (x5)
 (x6)
 (x7)
1686094818344 ERROR "src/routes/(auth)/login/+page.svelte" 44:17 "Type '{ value: string; label: string; \"data-testid\": string; }' is not assignable to type '{ focus?: (() => void) | undefined; blur?: (() => void) | undefined; layout?: (() => void) | undefined; getElement?: (() => HTMLDivElement | HTMLLabelElement) | undefined; ... 1986 more ...; input$resizable?: boolean | undefined; }'.\n  Object literal may only specify known properties, and '\"data-testid\"' does not exist in type '{ focus?: (() => void) | undefined; blur?: (() => void) | undefined; layout?: (() => void) | undefined; getElement?: (() => HTMLDivElement | HTMLLabelElement) | undefined; ... 1986 more ...; input$resizable?: boolean | undefined; }'."
1686094818344 ERROR "src/routes/(auth)/login/+page.svelte" 51:17 "Type '{ type: string; value: string; label: string; \"data-testid\": string; }' is not assignable to type '{ focus?: (() => void) | undefined; blur?: (() => void) | undefined; layout?: (() => void) | undefined; getElement?: (() => HTMLDivElement | HTMLLabelElement) | undefined; ... 1986 more ...; input$resizable?: boolean | undefined; }'.\n  Object literal may only specify known properties, and '\"data-testid\"' does not exist in type '{ focus?: (() => void) | undefined; blur?: (() => void) | undefined; layout?: (() => void) | undefined; getElement?: (() => HTMLDivElement | HTMLLabelElement) | undefined; ... 1986 more ...; input$resizable?: boolean | undefined; }'."
1686094818344 ERROR "src/routes/(auth)/login/+page.svelte" 55:25 "Type '{ \"data-testid\": string; class: string; variant: \"raised\"; }' is not assignable to type '{ getElement?: (() => HTMLElement) | undefined; use?: ActionArray | undefined; class?: string | undefined; style?: string | undefined; ripple?: boolean | undefined; ... 326 more ...; 'sveltekit:reload'?: true | ... 1 more ... | undefined; } | { ...; } | undefined'.\n  Object literal may only specify known properties, and '\"data-testid\"' does not exist in type '{ getElement?: (() => HTMLElement) | undefined; use?: ActionArray | undefined; class?: string | undefined; style?: string | undefined; ripple?: boolean | undefined; ... 326 more ...; 'sveltekit:reload'?: true | ... 1 more ... | undefined; } | { ...; }'."
1686094818345 ERROR "src/routes/(auth)/login/+page.svelte" 55:25 "This message is ok"
`

func TestIgnoreDataTestIdMessage(t *testing.T) {
	b := bytes.NewBufferString(exampleOutput)
	out := ParseLines(b, IgnoreDataTestIdMessage)
	require.Equal(t, []CheckMessage{
		{
			Num:      1686094808695,
			Type:     Start,
			Filename: "/home/matt/godev/secretsanta/frontend",
		},
		{
			Num:            1686094818345,
			Type:           Error,
			Filename:       "src/routes/(auth)/login/+page.svelte",
			FileLineAndCol: "55:25",
			ErrorMessage:   "This message is ok",
		},
	}, out)
}

func TestParseLine(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		ok       bool
		expected CheckMessage
	}{
		{
			name: "standard line",
			line: `1686094818344 ERROR "src/routes/(auth)/login/+page.svelte" 44:17 ` +
				`"Type '{ value: string; label: string; \"data-testid\": string; }' is not assignable to type"`,
			ok: true,
			expected: CheckMessage{
				Num:            1686094818344,
				Type:           Error,
				Filename:       "src/routes/(auth)/login/+page.svelte",
				FileLineAndCol: "44:17",
				ErrorMessage:   `Type '{ value: string; label: string; "data-testid": string; }' is not assignable to type`,
			},
		},
		{
			name: "unescapes new line",
			line: `1686094818344 ERROR "filename" 1:2 "something\nanother"`,
			ok:   true,
			expected: CheckMessage{
				Num:            1686094818344,
				Type:           Error,
				Filename:       "filename",
				FileLineAndCol: "1:2",
				ErrorMessage:   "something\nanother",
			},
		},
		{
			name:     "ignores weird number thing",
			line:     ` (x2) `,
			ok:       false,
			expected: CheckMessage{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			msg, ok := ParseLine(tc.line)
			require.Equal(t, tc.ok, ok)
			require.Equal(t, tc.expected, msg)
		})
	}

}
