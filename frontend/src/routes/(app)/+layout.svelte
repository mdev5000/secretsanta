<script lang="ts">
    import TopAppBar, {Row, Section, Title} from '@smui/top-app-bar';
    import IconButton from '@smui/icon-button';
    import {base} from "$app/paths";
    import {dev, browser} from '$app/environment';
    import {goto} from '$app/navigation';
    import {getData} from "$lib/rest/rest";
    import {Status} from "$lib/requests/setup/status";

    export const ssr = false;

    async function checkStatus() {
        const status = await getData<Status>(Status, "/api/setup/status");
        if (status.data.status !== "setup") {
            goto("/app/setup");
        }
    }

    if (browser && dev) {
        checkStatus();
    }
</script>

<div class="app">

    <div class="top-bar">
        <TopAppBar
                class="top-bar"
                variant="static"
        >
            <Row>
                <Section>
                    <IconButton class="material-icons">menu</IconButton>
                    <Title><a href="{base}/">Secret Santa</a></Title>
                </Section>
                <Section>
                    <a href="{base}/example">Example</a>
                    <a href="{base}/about">About</a>
                </Section>
                <Section align="end" toolbar>
                    <IconButton class="material-icons" aria-label="Download">file_download</IconButton>
                    <IconButton class="material-icons" aria-label="Print this page">print</IconButton>
                    <IconButton class="material-icons" aria-label="Bookmark this page">bookmark</IconButton>
                </Section>
            </Row>
        </TopAppBar>
    </div>

    <main>
        <slot/>
    </main>

    <footer>
        <p>visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to learn SvelteKit</p>
    </footer>
</div>

<style lang="scss">
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  .top-bar {
    a {
      padding-left: 5px;
      padding-right: 5px;
      color: white;
      text-decoration: none;
    }
  }

  main {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 1rem;
    width: 100%;
    max-width: 64rem;
    margin: 0 auto;
    box-sizing: border-box;
  }

  footer {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 12px;
  }

  footer a {
    font-weight: bold;
  }

  @media (min-width: 480px) {
    footer {
      padding: 12px 0;
    }
  }
</style>
