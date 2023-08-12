package svelte_check_wrapper

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

var exampleOutputTestdataAttr = `
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

var exampleOutputNodeModulesError = `
1691854143648 START "/home/matt/godev/secretsanta/frontend"
1691854162077 ERROR "node_modules/@smui/menu/src/SelectionGroupIcon.ts" 8:3 "Types of construct signatures are incompatible.\n  Type 'new (options: ComponentConstructorOptions<{ getElement?: (() => HTMLSpanElement) | undefined; use?: ActionArray | undefined; class?: string | undefined; id?: string | null | undefined; property?: string | ... 1 more ... | undefined; ... 198 more ...; 'on:fullscreenerror'?: EventHandler<...> | ... 1 more ... | undefined; }>) => Graphic__SvelteComponent_' is not assignable to type 'new <Props extends Record<string, any> = any, Events extends Record<string, any> = any, Slots extends Record<string, any> = any>(options: ComponentConstructorOptions<Props>) => SvelteComponent<...>'.\n    Construct signature return types 'Graphic__SvelteComponent_' and 'SvelteComponent<Props, Events, Slots>' are incompatible.\n      The types of '$$prop_def' are incompatible between these types.\n        Type '{ getElement?: (() => HTMLSpanElement) | undefined; use?: ActionArray | undefined; class?: string | undefined; id?: string | null | undefined; property?: string | ... 1 more ... | undefined; ... 198 more ...; 'on:fullscreenerror'?: EventHandler<...> | ... 1 more ... | undefined; }' is not assignable to type 'Props'.\n          '{ getElement?: (() => HTMLSpanElement) | undefined; use?: ActionArray | undefined; class?: string | undefined; id?: string | null | undefined; property?: string | ... 1 more ... | undefined; ... 198 more ...; 'on:fullscreenerror'?: EventHandler<...> | ... 1 more ... | undefined; }' is assignable to the constraint of type 'Props', but 'Props' could be instantiated with a different subtype of constraint 'Record<string, any>'."
`

func TestIgnoreDataTestIdMessage(t *testing.T) {
	b := bytes.NewBufferString(exampleOutputTestdataAttr)
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

func TestIgnoreNodeModules(t *testing.T) {
	b := bytes.NewBufferString(exampleOutputNodeModulesError)
	out := ParseLines(b, IgnoreNodeModules)
	require.Equal(t, []CheckMessage{
		{
			Num:      1691854143648,
			Type:     Start,
			Filename: "/home/matt/godev/secretsanta/frontend",
		},
	}, out)
}

func TestAnd(t *testing.T) {
	cases := []struct {
		name     string
		msg      CheckMessage
		f        Filter
		expected bool
	}{
		{
			name: "returns true when all pass",
			msg:  CheckMessage{},
			f: And(
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
			),
			expected: true,
		},
		{
			name: "returns false when one filter fails",
			msg:  CheckMessage{},
			f: And(
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return false },
				func(message CheckMessage) bool { return true },
			),
			expected: false,
		},
		{
			name: "returns false when one filter fails (final filter)",
			msg:  CheckMessage{},
			f: And(
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return true },
				func(message CheckMessage) bool { return false },
			),
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.f(tc.msg))
		})
	}

	t.Run("passes message to sub funcs", func(t *testing.T) {
		msg := CheckMessage{
			Num:            50,
			Type:           "some-type",
			ErrorMessage:   "some message",
			Filename:       "some/filename",
			FileLineAndCol: "80:80",
		}
		var called int
		numFilters := 10
		filters := make([]Filter, numFilters)
		for i := 0; i < numFilters; i++ {
			filters[i] = func(message CheckMessage) bool {
				called++
				require.Equal(t, msg, message)
				return true
			}
		}

		And(filters...)(msg)
		require.Equal(t, numFilters, called, "called all filters")
	})
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
