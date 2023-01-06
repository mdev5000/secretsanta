import {expect, test} from 'vitest'
import {ExclusionOperation as O, compileExclusions} from './exclusions'
import type {ExclusionItem} from './exclusions'

function withOp(original: O, next: O, expected: boolean) {
    const exclusions = [
        {
            'user1': {
                'user2': original
            }
        },
        {
            'user1': {
                'user2': next
            }
        },
    ];
    const reverseResults = {};
    const result = compileExclusions(exclusions, reverseResults);
    expect(result).toEqual({
        'user1': {
            'user2': expected,
        }
    });
}

test("can override excludes", () => withOp(O.Exclude, O.Include, false));

test("can override includes", () => withOp(O.Include, O.Exclude, true));

test("defaults to include the item", () => withOp(O.Default, O.Default, false));

test("default acts as pass-through", () => {
    withOp(O.Include, O.Default, false);
    withOp(O.Exclude, O.Default, true);
});

test("merges exclusion lists", () => {
    const exclusions: ExclusionItem[] = [
        {
            'user1': {
                'user2': O.Exclude,
                'user3': O.Default,
                'user4': O.Include,
                'user5': O.Include,
                'user6': O.Include,
                'user7': O.Exclude,
                'user8': O.Include,
            }
        },
        {
            'user1': {
                'user2': O.Default,
                'user4': O.Exclude,
                'user5': O.Default,
            }
        },
        {
            'user1': {
                'user5': O.Exclude,
            }
        },
    ];
    const reverseResults = {
        'user1': {
            'user6': false,
            'user7': true,
            'user8': true,
            'user9': true,
        }
    };
    const result = compileExclusions(exclusions, reverseResults);
    expect(result).toStrictEqual({
        'user1': {
            'user2': true,
            'user3': false,
            'user4': true,
            'user5': true,
            'user6': false,
            'user7': false,
            'user8': true,
            'user9': true,
        }
    })
});
