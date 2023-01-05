import {redirect} from "@sveltejs/kit";
import {base} from "$app/paths";
import type {PageLoad} from "../../../../../.svelte-kit/types/src/routes";

export const load = (({params}) => {
    if (params.name !== "cool") {
        throw redirect(302, base + "/");
    }
}) satisfies PageLoad;
