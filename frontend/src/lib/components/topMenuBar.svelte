<script lang="ts">
    import TopAppBar, {Row, Section, Title} from "@smui/top-app-bar";
    import IconButton from '@smui/icon-button';
    import Button from "@smui/button";
    import {base} from "$app/paths";
    import List, {Item} from "@smui/list";
    import MenuSurface from "@smui/menu-surface";
    import {media} from "$lib/stores/media"

    let open = false;
    let surface: MenuSurface;

    function toggleOpen() {
        open = !open;
        surface.setOpen(open);
    }
</script>

<TopAppBar variant="static">
    <Row>
        <Section align="start">
            <Button on:click={toggleOpen}>
                <IconButton class="material-icons">menu</IconButton>
            </Button>
            <MenuSurface bind:this={surface} anchorCorner="BOTTOM_LEFT">
                <List>
                    <Item>Menu Item 1</Item>
                    <Item>Menu Item 2</Item>
                    <Item>Menu Item 3</Item>
                    {#if $media.small}
                        <Item>
                            <IconButton class="material-icons" aria-label="Download">file_download</IconButton>
                            Download
                        </Item>
                    {/if}
                </List>
            </MenuSurface>
            <Title><a href="{base}/">Secret Santa</a></Title>
        </Section>
        <Section>
            <a href="{base}/example">Example</a>
            <a href="{base}/about">About</a>
        </Section>
        <Section align="end" toolbar>
            {#if !$media.small}
                <IconButton class="material-icons" aria-label="Download">file_download</IconButton>
                <IconButton class="material-icons" aria-label="Print this page">print</IconButton>
                <IconButton class="material-icons" aria-label="Bookmark this page">bookmark</IconButton>
            {/if}
        </Section>
    </Row>
</TopAppBar>

<style>
    a {
      padding-left: 5px;
      padding-right: 5px;
      color: white;
      text-decoration: none;
    }
</style>
