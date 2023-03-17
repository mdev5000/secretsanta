<script lang="ts">
    import TextField from "@smui/textfield"
    import Button, {Label} from "@smui/button"
    import {getData, postData, postTmp} from "$lib/rest/rest";
    import {LeaderStatus} from "$lib/requests/setup/leader";
    import {goto} from "$app/navigation";
    import {error} from "$lib/applog";
    import {Setup} from "$lib/requests/setup/setup";
    import {AppErrorRs} from "../../../lib/requests/core/error";

    let adminUsername = "admin";
    let adminFirstname = "Admiral"
    let adminLastname = "Min"
    let adminPassword = "";
    let defaultFamilyName = "Default";
    let defaultFamilyDescription = "Default Family";

    let status: 'welcome' | 'succeeded' | 'failed' | 'error' | 'done' = 'welcome';

    let getLeadership = async (): Promise<null> => {
        const r = await getData<LeaderStatus>(LeaderStatus, "/api/setup/leader-status");
        if (r.status != 200) {
            error("leadership requested ended in error", r.data.error);
            status = 'error';
            return;
        }
        status = r.data.isLeader ? 'succeeded' : 'failed';
    }

    let finalize = async (e): Promise<null> => {
        e.preventDefault();
        const data = Setup.toJson({
            adminPassword: adminPassword,
            admin: {
                username: adminUsername,
                firstname: adminFirstname,
                lastname: adminLastname,
            },
            family: {
                name: defaultFamilyName,
                description: defaultFamilyDescription,
            }
        })
        const r2 = await postData<AppErrorRs>(AppErrorRs, '/api/setup/finalize', data)
        if (r2.status == 200) {
            status = 'done';
            setTimeout(() => {
                goto('/app');
            }, 5000);
        }

        const r = await postTmp("/api/setup/finalize-quick");
        if (r.status == 200) {
            status = 'done';
            setTimeout(() => {
                goto('/app');
            }, 5000);
        }
    }

</script>
<div>
    <h1>Setup</h1>

    {#if status === 'welcome' || status === 'loading'}
        <p>Welcome to Secret Santa! Click next to get started.</p>
        <Button on:click={() => getLeadership()}>Next</Button>
    {:else if status === 'failed'}
        <p>The site is already being setup by another computer or browser.</p>
        <p>If you are unable to complete the setup, please restart the server and refresh the page.</p>
    {:else if status === 'succeeded'}
        <form>

            <fieldset>
                <legend class="h2">Setup Admin User</legend>

                <TextField
                        bind:value={adminUsername}
                        label="Admin Username"
                />
                <br/>

                <TextField
                        bind:value={adminFirstname}
                        label="Admin First Name"
                />
                <br/>

                <TextField
                        bind:value={adminLastname}
                        label="Admin Last Name"
                />
                <br/>

                <TextField
                        type="password"
                        bind:value={adminPassword}
                        label="Admin Password"
                />

            </fieldset>

            <fieldset>

                <legend class="h2">Default Family</legend>

                <TextField
                        bind:value={defaultFamilyName}
                        label="Family Name"
                />

                <TextField
                        bind:value={defaultFamilyName}
                        label="Family Description"
                />

            </fieldset>

            <fieldset>

                <Button class="submit" variant="raised" on:click={finalize}>
                    <Label>Setup</Label>
                </Button>

            </fieldset>

        </form>
    {:else if status === 'error'}
        Error occurred while setting up, please try refreshing the page.
    {:else if status === 'done'}
        Setup completed, will redirect to app in a moment...
    {/if}


</div>

<style lang="scss">
  form {
    fieldset {
      border: none;
      margin-bottom: 60px;

      legend {
        margin-bottom: 10px;
      }
    }

    :global(button.submit) {
      margin-top: 20px;
    }
  }
</style>

