<script lang="ts">
    import TextField from "@smui/textfield"
    import Button, {Label} from "@smui/button"

    let adminName = "admin";
    let adminPassword = "";
    let defaultFamilyName = "Default";
</script>

<div>
    <h1>Setup</h1>

    <form>

        <fieldset>
            <legend class="h2">Setup Admin User</legend>

        <TextField
                bind:value={adminName}
                label="Admin Username"
        />
        <br/>

        <TextField
                type="password"
                bind:value={adminPassword}
                label="Admin Password"
        />

        </fieldset>

        <fieldset>

        <legend class="h2">Default Family Name</legend>

        <TextField
                bind:value={defaultFamilyName}
                label="Default Family Name"
        />

        </fieldset>

        <fieldset>

            <Button class="submit" variant="raised">
                <Label>Setup</Label>
            </Button>

        </fieldset>

    </form>
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

