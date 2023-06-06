<script lang="ts">
    import TextField from "@smui/textfield"
    import Button, {Label} from "@smui/button"
    import {getData, postData} from "$lib/rest/rest";
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

    let status: 'welcome' | 'setting-up' | 'succeeded' | 'failed' | 'error' | 'done' = 'welcome';

    let getLeadership = async (): Promise<undefined> => {
        const r = await getData<LeaderStatus>(LeaderStatus, "/api/setup/leader-status");
        if (r.status != 200) {
            error("leadership requested ended in error", r.data.error!);
            status = 'error';
            return;
        }
        status = r.data.isLeader ? 'succeeded' : 'failed';
    }

    let finalize = async (e: any): Promise<undefined> => {
        status = 'setting-up';
        e.preventDefault();
        const data = Setup.toJson({
            adminPassword: adminPassword,
            admin: {
                id: "",
                username: adminUsername,
                firstname: adminFirstname,
                lastname: adminLastname,
            },
            family: {
                id: "",
                name: defaultFamilyName,
                description: defaultFamilyDescription,
            }
        })
        const r2 = await postData<AppErrorRs>(AppErrorRs, '/api/setup/finalize', data)
        if (r2.status == 204) {
            status = 'done';
            setTimeout(() => {
                goto('/app/login');
            }, 5000);
            return;
        }

        // otherwise handle error
    }

</script>
<div>
    <h1>Setup</h1>

    {#if status === 'welcome'}
        <p>Welcome to Secret Santa! Click next to get started.</p>
        <Button data-testid="next-btn" on:click={() => getLeadership()}>Next</Button>
    {:else if status === 'failed'}
        <p>The site is already being setup by another computer or browser.</p>
        <p>If you are unable to complete the setup, please restart the server and refresh the page.</p>
    {:else if status === 'succeeded' || status === 'setting-up'}
        <form>

            <fieldset>
                <legend class="h2">Setup Admin User</legend>

                <TextField
                        bind:value={adminUsername}
                        label="Admin Username"
                        data-testid="admin-username"
                />
                <br/>

                <TextField
                        bind:value={adminFirstname}
                        label="Admin First Name"
                        data-testid="admin-firstname"
                />
                <br/>

                <TextField
                        bind:value={adminLastname}
                        label="Admin Last Name"
                        data-testid="admin-lastname"
                />
                <br/>

                <TextField
                        type="password"
                        bind:value={adminPassword}
                        label="Admin Password"
                        data-testid="admin-password"
                />

            </fieldset>

            <fieldset>

                <legend class="h2">Default Family</legend>

                <TextField
                        bind:value={defaultFamilyName}
                        label="Family Name"
                        data-testid="default-family-name"
                />

                <TextField
                        bind:value={defaultFamilyName}
                        label="Family Description"
                />

            </fieldset>

            <fieldset>

                <Button data-testid="setup-submit" class="submit" variant="raised" on:click={finalize}>
                    <Label>Setup</Label>
                </Button>

            </fieldset>

        </form>
    {:else if status === 'error'}
        <p data-testid="status">Error occurred while setting up, please try refreshing the page.</p>
    {:else if status === 'done'}
        <p data-testid="status">Setup completed, will redirect to app in a moment...</p>
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

