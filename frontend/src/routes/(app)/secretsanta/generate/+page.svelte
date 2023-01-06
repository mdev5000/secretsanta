<script lang="ts">
    import Button from "@smui/button"
    import TextField from "@smui/textfield"
    import Autocomplete from '@smui-extra/autocomplete';
    import Chip, {Set, TrailingAction, Text} from '@smui/chips';
    import Accordion, {Panel, Header, Content} from "@smui-extra/accordion";
    import DataTable, { Head, Body, Row, Cell } from '@smui/data-table';
    import Checkbox from '@smui/checkbox';

    type Family = { id: string; name: string; selected: boolean };

    let name = "";
    let date = "";

    $: slug = computeSlug(name)

    function computeSlug(name) {
        return name.replaceAll(" ", "-").toLowerCase();
    }


    let families: Family[] = [
        {id: "default", name: "Default", selected: true},
        {id: "another", name: "Another", selected: false},
        {id: "third", name: "Third", selected: false},
    ];
    let currentFamily = undefined;
    $: unselectedFamilies = families.filter((f) => !f.selected);
    $: selectedFamilies = families.filter((f) => f.selected);

    let users = [
        {id: "john", name: "John Person", enabled: true, families: [{name: "Default"}, {name: "Another"}]},
        {id: "matt", name: "Matt Somebody", enabled: true, families: [{name: "Default"}]},
    ];

    let selectedUsers = [users[0], users[1]];

    function addFamily() {
        if (currentFamily == undefined) {
            return;
        }
        families = familySelected(currentFamily.id, true);
        currentFamily = undefined;
    }

    function familyRemoved(e) {
        families = familySelected(e.detail.chipId, false);
    }

    function familySelected(id: string, selected: boolean): Family[] {
        return families.map((f) => {
            if (f.id == id) {
                f.selected = selected;
            }
            return f;
        });
    }

</script>

<form>

    <fieldset>

        <TextField label="Name" bind:value={name}/>
        <br/>
        <TextField label="Slug" bind:value={slug} disabled/>
        <br/>

        <TextField label="Date" type="date" bind:value={date}/>
        <br/>

    </fieldset>

    <fieldset>

        <legend class="h4">Families</legend>

        <Autocomplete options={unselectedFamilies}
                      getOptionLabel={(f) => f ? f.name : ''}
                      bind:value={currentFamily}
                      showMenuWithNoInput={false}
                      label="Family / Group"/>

        <Button on:click={addFamily}>Add</Button>

        <Set chips={selectedFamilies} let:chip key={(chip) => chip.id}>
            <Chip chip="{chip.id}" on:MDCChip:removal={familyRemoved}>
                <Text>{chip.name}</Text>
                <TrailingAction icon$class="material-icons">cancel</TrailingAction>
            </Chip>
        </Set>

    </fieldset>

    {#if users}
        <fieldset>
            <legend class="h4">Users</legend>

            <DataTable style="max-width: 100%;">
                <Head>
                    <Row>
                        <Cell checkbox>
                            <Checkbox />
                        </Cell>
                        <Cell>Name</Cell>
                        <Cell>Families</Cell>
                    </Row>
                </Head>
                <Body>
                {#each users as user (user.id)}
                    <Row>
                        <Cell checkbox>
                            <Checkbox
                                    bind:group={selectedUsers}
                                    value={user}
                                    valueKey={user.id}
                            />
                        </Cell>
                        <Cell>{user.name}</Cell>
                        <Cell>{user.families.map((f) => f.name).join(", ")}</Cell>
                    </Row>
                {/each}
                </Body>
            </DataTable>
        </fieldset>
    {/if}

    <Accordion>
        <Panel disabled={!selectedFamilies.length}>
            <Header>Exclusions</Header>
            <Content>
                <fieldset>
                    this is content
                </fieldset>
            </Content>
        </Panel>
    </Accordion>

</form>

<style lang="scss">
    fieldset {
      margin-bottom: 20px;
    }
</style>