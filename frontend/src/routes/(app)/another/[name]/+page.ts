import {redirect} from "@sveltejs/kit";
import {base} from "$app/paths";
import type {PageLoad} from './$types';

export const load = (({params}) => {
    if (params.name !== "cool") {
        throw redirect(302, base + "/");
    }
}) satisfies PageLoad;
