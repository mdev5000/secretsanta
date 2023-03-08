<script lang="ts">
    import TextField from "@smui/textfield"
    import Button, {Label} from "@smui/button"

    let adminUsername = "admin";
    let adminFirstname = "Admiral"
    let adminLastname = "Min"
    let adminPassword = "";
    let defaultFamilyName = "Default";
    let defaultFamilyDescription = "Default Family";
</script>

<div>
    <h1>Setup</h1>

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

