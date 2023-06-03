<script lang="ts">
    import {dev, browser} from '$app/environment';
    import { page } from '$app/stores';
    import {goto} from '$app/navigation';
    import {getData} from "$lib/rest/rest";
    import {Status} from "$lib/requests/setup/status";

    import TopMenuBar from "$lib/components/topMenuBar.svelte"
    import {isLoggedIn} from "../../lib/auth/auth";

    async function checkStatus() {
        const status = await getData<Status>(Status, "/api/setup/status");
        if (status.data.status !== "setup") {
            goto("/app/setup");
        }
    }
    if (browser && dev) {
        checkStatus();
    }

    if (!isLoggedIn()) {
        goto("/app/login");
    }
</script>

<div class="app">

    <div class="top-bar">
        <TopMenuBar/>
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
