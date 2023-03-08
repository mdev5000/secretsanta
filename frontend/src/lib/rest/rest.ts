import axios from "axios";
import type {AxiosResponse} from "axios";

export type Response = AxiosResponse;

export async function get(uri: string): Promise<Response> {
    const url = "http://localhost:3000"
    return await axios.get(`${url}/${uri}`);
}