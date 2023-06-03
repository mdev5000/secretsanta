<script lang="ts">
    import Button, {Label} from '@smui/button';
    import {Status} from "$lib/requests/setup/status";
    import {getData} from "$lib/rest/rest"

    let another = 0;
    let clicked = 0;
    $: clicked2 = 2 * clicked;
    $: clicked3 = 2 * clicked2;
    $: clicked, updateAnother();

    function updateAnother() {
        another = 3 * clicked3;
    }

    let status: Status = {status: ""};

    async function fetchIt() {
        console.log("fetching")
        const result = await getData<Status>(Status, "/api/example");
        status = result.data;
    }

</script>

<div>
    <Button on:click={() => clicked++} variant="raised">
        <Label>Raised</Label>
    </Button>
    <div data-testid="fetcher">
        <Button on:click={fetchIt} variant="raised">
            <Label>Fetch</Label>
        </Button>
    </div>
    <div>Clicked {clicked}</div>
    <div>Clicked {clicked2}</div>
    <div>Clicked {clicked3}</div>
    <div>Another {another}</div>

    {#if status.status !== '' }
    <div data-testid="status-data">
        <div data-testid="status"><span>Status</span>: {status.status}</div>
    </div>
    {/if}
</div>