<script lang="ts">
    import Button, {Label} from '@smui/button';
    import {Login} from "$lib/requests/core/login"
    import {getData} from "$lib/rest/rest"

    let another = 0;
    let clicked = 0;
    $: clicked2 = 2 * clicked;
    $: clicked3 = 2 * clicked2;

    let login: Login = {username: '', password: ''};

    $: clicked, updateAnother();

    function updateAnother() {
        another = 3 * clicked3;
    }

    function logLogin(l: Login) {
        console.log(l);
        login = l;
    }

    async function fetchIt() {
        console.log("fetching")
        const result = await getData<Login>(Login, "/api/example");
        console.log(result);
        logLogin(result.data);
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

    {#if login.username !== '' }
    <div data-testid="login-data">
        <div data-testid="username"><span>Username</span>: {login.username}</div>
        <div><span>Password</span>: {login.password}</div>
    </div>
    {/if}
</div>