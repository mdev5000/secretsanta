<script lang="ts">
    import Button from "@smui/button"
    import TextField from "@smui/textfield"
    import Autocomplete from '@smui-extra/autocomplete';
    import Chip, {Set, TrailingAction, Text} from '@smui/chips';
    import DataTable, {Head, Body, Row, Cell} from '@smui/data-table';
    import Checkbox from '@smui/checkbox';
    import FormField from '@smui/form-field';
    import Tab, { Label } from '@smui/tab';
    import TabBar from '@smui/tab-bar';

    type Family = { id: string; name: string; selected: boolean };

    let name = "";
    let date = "";

    $: slug = computeSlug(name)

    function computeSlug(name: string) {
        return name.replaceAll(" ", "-").toLowerCase();
    }

    let families: Family[] = [
        {id: "default", name: "Default", selected: true},
        {id: "another", name: "Another", selected: false},
        {id: "third", name: "Third", selected: false},
    ];
    let currentFamily: Family | undefined = undefined;
    $: unselectedFamilies = families.filter((f) => !f.selected);
    $: selectedFamilies = families.filter((f) => f.selected);

    let users = [
        {
            id: "john", name: "John Person",
            families: [{name: "Default"}, {name: "Another"}],
            excludeByDefault: false,
            ssExclusions: [],
        },
        {id: "matt", name: "Matt Somebody", families: [{name: "Default"}],
            excludeByDefault: true,
            ssExclusions: [{id: "john", name: "John"}],
        },
    ];

    let selectedUsers = users.filter((u) => !u.excludeByDefault);

    let secretSantas = [
        {
            id: "christmas-harrison---2019",
            name: "Christmas (Harrison) - 2019",
            slug: "christmas-harrison---2019",
            pairings: {
                "matt": "john",
                "john": "bob",
            }
        },
        {
            id: "christmas-dohmer---2019",
            name: "Christmas (Dohmer) - 2019",
            slug: "christmas-dohmer---2019",
            pairings: {
                "matt": "john",
                "john": "bob",
            }
        }
    ];

    let selectedSecretSantas: string[] = [];

    let steps = [
        "Step 1: Families",
        "Step 2: Users",
        "Step 3: Pairing exclusions",
        "Step 4: Finalize",
    ]
    let activeStep = steps[0];

    function addFamily() {
        if (currentFamily == undefined) {
            return;
        }
        families = familySelected(currentFamily.id, true);
        currentFamily = undefined;
    }

    function familyRemoved(e: any) {
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

<h1>Generate a new Secret Santa</h1>

<form>

    <div class="tab-bar">
        <TabBar tabs={steps} let:tab bind:active={activeStep}>
            <Tab tab={tab}>
                <Label>{tab}</Label>
            </Tab>
        </TabBar>
    </div>

    {#if activeStep === steps[0]}
        <div class="step-select-families">

            <div class="explanation">
                Select the families you wish to participate in this Secret Santa.
            </div>

            <fieldset>

                <Autocomplete options={unselectedFamilies}
                              getOptionLabel={(f) => f ? f.name : ''}
                              bind:value={currentFamily}
                              showMenuWithNoInput={false}
                              label="Family / Group"/>

                <Button on:click={addFamily}>Add</Button>

                <div class="selected">
                    <div>Selected:</div>
                    <Set chips={selectedFamilies} let:chip key={(chip) => chip.id}>
                        <Chip chip="{chip.id}" on:MDCChip:removal={familyRemoved}>
                            <Text>{chip.name}</Text>
                            <TrailingAction icon$class="material-icons">cancel</TrailingAction>
                        </Chip>
                    </Set>
                </div>

            </fieldset>

            <fieldset class="fs-submit">
                <Button variant="raised" on:click={() => activeStep = steps[1]}>Next</Button>
            </fieldset>

        </div>

    {:else if activeStep === steps[1]}

        {#if users}

            <div class="explanation">
                Select the users you wish to participate in this Secret Santa.
            </div>

            <fieldset>
                <DataTable style="width: 100%;">
                    <Head>
                        <Row>
                            <Cell checkbox>
                                <Checkbox/>
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

            <fieldset class="fs-submit">
                <Button variant="raised" on:click={() => activeStep = steps[2]}>Next</Button>
            </fieldset>
        {/if}


    {:else if activeStep === steps[2]}

        <fieldset class="skip-button">
            <Button variant="raised" on:click={() => activeStep = steps[3]}>Skip</Button>
        </fieldset>

        <div class="ss-pairing-exclusions">

            <div class="explanation">
                <p>
                    Prevent certain users from being paired this Secret Santa.
                </p>
                <p>
                    For example sometimes it's desirable to not pair users this year, if paired together in previous
                    years or to prevent spouses from being paired.
                </p>
            </div>

            <fieldset class="fs-secret-santa">

                <legend>Secret Santas</legend>

                <div class="explanation">Exclude pairings from previous Secret Santas.</div>

                <div class="secret-santas">
                    {#each secretSantas as ss}
                        <div>
                            <FormField>
                                <Checkbox
                                        bind:group={selectedSecretSantas}
                                        value={ss}
                                        valueKey={ss.id}
                                />
                                <span slot="label">{ss.name}</span>
                            </FormField>
                        </div>
                    {/each}
                </div>

                <Button variant="raised">All Secret Santas</Button>

            </fieldset>

            <Button variant="raised">Customize</Button>

            <fieldset class="fs-submit">
                <Button variant="raised" on:click={() => {activeStep = steps[3]; scroll(0, 0)}}>Next</Button>
            </fieldset>
        </div>

    {:else if activeStep === steps[3]}

        <div class="step-finalize">

            <fieldset class="fs-name-date">

                <TextField label="Name" bind:value={name}/>
                <br/>
                <TextField label="Slug" bind:value={slug} disabled/>
                <br/>

                <TextField label="Date" type="date" bind:value={date}/>
                <br/>

            </fieldset>

            <div class="families">
                <span>Families:</span>
                <span>{selectedFamilies.map((f) => f.name).join(",")}</span>
            </div>

            {#if users}
                <fieldset>
                    <legend class="h4">Users</legend>

                    <DataTable style="width: 100%;">
                        <Head>
                            <Row>
                                <Cell>Name</Cell>
                                <Cell>Will not pair with</Cell>
                            </Row>
                        </Head>
                        <Body>
                        {#each users as user (user.id)}
                            <Row>
                                <Cell>{user.name}</Cell>
                                <Cell>Example Person, Example Person 2</Cell>
                            </Row>
                        {/each}
                        </Body>
                    </DataTable>
                </fieldset>
            {/if}

            <fieldset class="fs-submit">
                <Button variant="raised">Generate</Button>
            </fieldset>

        </div>

    {/if}

</form>

<style lang="scss">

  .tab-bar {
    margin-bottom: 40px;
  }

  fieldset {
    margin-bottom: 60px;
    border: none;
    padding: 0;

    &.skip-button {
      border: none;
      margin-bottom: 20px;
      text-align: right;
    }

    &.fs-submit {
      border: none;
      margin-top: 60px;
      text-align: right;
    }

    &.fs-name-date {
      border: none;
      padding: 0;
    }
  }

  .explanation {
    margin: 20px 0 20px 10px;
    color: #3c3c3c;
  }

  .ss-pairing-exclusions {
    .explanation {
      color: #5a5a5a;
    }

    fieldset {
      .explanation {
        margin-left: 20px;
      }
    }

    .secret-santas {
      margin-bottom: 20px;
    }
  }

  .step-select-families {
    .selected {
      margin-top: 40px;
    }
  }

  .step-finalize {
    .families {
      margin-bottom: 40px;
    }
  }
</style>