<script lang="ts">
    import Button, {Label} from '@smui/button';
    import {Login} from "$lib/requests/login"
    import {getData, Result} from "$lib/rest/rest"

    let another = 0;
    let clicked = 0;
    $: clicked2 = 2 * clicked;
    $: clicked3 = 2 * clicked2;

    $: clicked, updateAnother();

    function updateAnother() {
        another = 3 * clicked3;
    }

    function logLogin(l: Login) {
        console.log(l);
    }

    async function fetchIt() {
        console.log("fetching")
        const result = await getData<Login>(Login, "/example");
        console.log(result);
        logLogin(result.data);
    }

</script>

<div>
    <Button on:click={() => clicked++} variant="raised">
        <Label>Raised</Label>
    </Button>
    <Button on:click={fetchIt} variant="raised">
        <Label>Fetch</Label>
    </Button>
    <div>Clicked {clicked}</div>
    <div>Clicked {clicked2}</div>
    <div>Clicked {clicked3}</div>
    <div>Another {another}</div>
</div>