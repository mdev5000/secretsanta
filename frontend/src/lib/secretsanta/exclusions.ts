
export enum ExclusionOperation {
    Default = "default",
    Exclude = "exclude",
    Include = "include"
}

export type ExclusionOpList = { [prop: string]: ExclusionOperation }
export type ExclusionItem = { [prop: string]: ExclusionOpList };

export type Reversals = {
    [prop: string]: {
        [prop: string]: boolean
    }
};

export type ExclusionList = { [prop: string]: boolean }
export type Exclusions = { [prop: string]: ExclusionList };

export function compileExclusions(exclusions: ExclusionItem[], reversals: Reversals): Exclusions {
    const result = exclusions.reduce((result, e) => {
        for (const user in e) {
            result[user] = compileUsersExclusions(fetchExclusion(result, user), e[user]);
        }
        return result;
    }, {} as Exclusions);

    for (const fromUser in reversals) {
        const userReversals = reversals[fromUser];
        for (const toUser in userReversals) {
            const exclusions = fetchExclusion(result, fromUser);
            if (userReversals[toUser]) {
                exclusions[toUser] = reverseResult(exclusions[toUser]);
            }
        }
    }

    return result;
}

export function compileUsersExclusions(result: ExclusionList, ops: ExclusionOpList): ExclusionList {
    for (const user in ops) {
        result[user] = operationToResult(result[user], ops[user]);
    }
    return result;
}

function operationToResult(current: boolean, op: ExclusionOperation): boolean {
    switch (op) {
        case ExclusionOperation.Default:
            return handleMissing(current);
        case ExclusionOperation.Exclude:
            return true;
        case ExclusionOperation.Include:
            return false;
    }
}

function reverseResult(r: boolean): boolean {
    return !handleMissing(r);
}

function handleMissing(r: boolean): boolean {
    return r === undefined ? false : r;
}

function fetchExclusion(result: Exclusions, user: string): ExclusionList {
    if (result[user] === undefined) {
        result[user] = {};
    }
    return result[user];
}
