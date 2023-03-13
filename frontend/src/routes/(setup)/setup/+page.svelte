<script lang="ts">
    import TextField from "@smui/textfield"
    import Button, {Label} from "@smui/button"
    import {getData} from "../../../lib/rest/rest";
    import {LeaderStatus} from "../../../lib/requests/setup/leader";

    let adminUsername = "admin";
    let adminFirstname = "Admiral"
    let adminLastname = "Min"
    let adminPassword = "";
    let defaultFamilyName = "Default";
    let defaultFamilyDescription = "Default Family";

    let leadershipStatus = 'welcome';

    let getLeadership = async (): Promise<null> => {
        const r = await getData<LeaderStatus>(LeaderStatus, "/api/setup/leader-status");
        if (r.status != 200) {
            leadershipStatus = 'error';
            return;
        }
        if (r.data.isLeader) {
            leadershipStatus = 'succeeded';
        } else {
            leadershipStatus = 'failed';
        }
    }

</script>
<div>
    <h1>Setup</h1>

    {#if leadershipStatus === 'welcome' || leadershipStatus === 'loading'}
        <p>Welcome to Secret Santa! Click next to get started.</p>
        <Button on:click={() => getLeadership()}>Next</Button>
    {:else if leadershipStatus === 'failed'}
        <p>The site is already being setup by another computer or browser.</p>
        <p>If you are unable to complete the setup, please restart the server and refresh the page.</p>
    {:else if leadershipStatus === 'succeeded'}
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

                <Button class="submit" variant="raised">
                    <Label>Setup</Label>
                </Button>

            </fieldset>

        </form>
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

