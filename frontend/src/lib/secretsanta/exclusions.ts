export enum ExclusionOperation {
    Default = "default",
    Exclude = "exclude",
    Include = "include"
}

export type UserID = string;

export type ExclusionOpList = { [toUser: UserID]: ExclusionOperation };
export type ExclusionItem = { [fromUser: UserID]: ExclusionOpList };

export type Reversals = {
    [fromUser: UserID]: {
        [toUser: UserID]: boolean
    }
};

export type ExclusionList = { [toUser: UserID]: boolean }
export type ExclusionResults = { [fromUser: UserID]: ExclusionList };

export type Exclusions = { [fromUser: UserID]: UserID[] };

export function toExclusions(result: ExclusionResults): Exclusions {
    const out: Exclusions = {};

    for (const fromUser in result) {
        const userExclusions = result[fromUser];
        const userOut = [];
        for (const toUser in userExclusions) {
            if (userExclusions[toUser]) {
                userOut.push(toUser);
            }
        }

        if (userOut.length) {
            out[fromUser] = userOut;
        }
    }

    return out;
}


export function compileExclusions(exclusions: ExclusionItem[], reversals: Reversals): ExclusionResults {
    const result = exclusions.reduce((r, e) => {
        for (const user in e) {
            r[user] = compileUsersExclusions(fetchExclusion(r, user), e[user]);
        }
        return r;
    }, {} as ExclusionResults);

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

function fetchExclusion(result: ExclusionResults, user: UserID): ExclusionList {
    if (result[user] === undefined) {
        result[user] = {};
    }
    return result[user];
}
