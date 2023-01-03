import {redirect} from "@sveltejs/kit";
import {base} from "$app/paths";

export const load: import('./$types').PageLoad = ({params}) => {
    if (params.name !== "cool") {
        throw redirect(302, base + "/");
    }
}
