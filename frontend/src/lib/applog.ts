import type {AppError} from "./requests/core/error";

export function error(msg: string, err: AppError) {
    console.log("error: " + msg, err)
}