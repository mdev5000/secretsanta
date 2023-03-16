import axios from "axios";
import type {AxiosResponse} from "axios";
import type {JsonValue, PartialMessage} from "@protobuf-ts/runtime";
import { dev } from '$app/environment';

export type Response = AxiosResponse;

export interface Retrievable<T extends object> {
    create(d?: PartialMessage<T>): T
    fromJsonString(json: string): T
}

export interface Request<T extends object> {
    toJson(json: string): T
}

export class Result<T extends object> {
    status: number
    data: T

    constructor(code: number, data: T) {
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

async function get(uri: string): Promise<Response> {
    let url = "";
    if (dev) {
        url = "http://localhost:3000"
    }
    return await axios.get(`${url}${uri}`,{withCredentials: true});
}

export async function postData<T extends object>(
    response: Retrievable<T>, uri: string, requestJson: JsonValue): Promise<Result<T>> {
    const rs = await axios.post(uri, requestJson);
    const result = response.create(rs.data)
    return new Result(rs.status, result);
}

export async function postTmp(uri: string): Promise<Response> {
    let url = "";
    if (dev) {
        url = "http://localhost:3000"
    }
    return await axios.post(`${url}${uri}`,{
    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    });
}
