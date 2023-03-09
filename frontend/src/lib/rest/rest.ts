import axios from "axios";
import type {AxiosResponse} from "axios";
import type {PartialMessage} from "@protobuf-ts/runtime";

export type Response = AxiosResponse;

export interface Retrievable<T extends object> {
    create(d?: PartialMessage<T>): T
    fromJsonString(json: string): T
}

export class Result<T extends object> {
    status: number
    data: T | undefined

    constructor(code: number, data: T | undefined = undefined) {
        this.status = code;
        this.data = data;
    }
}

export async function getData<T extends object>(r: Retrievable<T>, uri: string): Promise<Result<T>> {
    // @todo switch to parsing from string
    // console.log(Login.fromJsonString(`{"duck": true}`));
    const response = await get(uri);
    const result = r.create(response.data)
    return new Result(response.status, result);
}

export async function get(uri: string): Promise<Response> {
    const url = "http://localhost:3000"
    return await axios.get(`${url}${uri}`);
}