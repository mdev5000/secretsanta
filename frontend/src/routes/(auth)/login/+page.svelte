<script lang="ts">
    import TextField from "@smui/textfield";
    import Button, {Label} from "@smui/button";
    import {postData} from "$lib/rest/rest";
    import {AppErrorRs} from "$lib/requests/core/error";
    import {Login} from "$lib/requests/core/login";
    import {goto} from "$app/navigation";

    let username = '';
    let password = '';
    let status = 'init';

    async function login(e: any): Promise<undefined> {
        e.preventDefault();
        status = 'logging_in'

        const data = Login.toJson({
            username: username,
            password: password
        })
        const loginRs = await postData<AppErrorRs>(AppErrorRs, '/api/login', data)
        if (loginRs.status !== 200) {
            // @todo actually do something useful here
            console.log('failed');
            status = 'failed';
            return;
        }
        status = 'success';
        setTimeout(() => {
            goto('/app');
        }, 1000);
        return;
    }

</script>

<div class="logout-wrapper">
    <h1>Login</h1>

    <form>
        <TextField
                bind:value={username}
                label="Username"
        />
        <br/>
        <TextField
                type="password"
                bind:value={password}
                label="Password"
        />
        <div class="submit-wrapper">
            {#if status === 'init' }
                <Button class="submit" variant="raised" on:click={login}>
                    <Label>Setup</Label>
                </Button>
            {:else if status === 'logging_in' }
                <Button class="submit" variant="raised">
                    <Label>Logging in...</Label>
                </Button>
            {:else if status === 'success' }
                <Button class="submit" variant="raised">
                    <Label>Success</Label>
                </Button>
            {/if}
        </div>
    </form>

</div>

<style lang="scss">
  .logout-wrapper {
    margin: auto;
  }

  .submit-wrapper {
    margin-top: 40px;
  }
</style>