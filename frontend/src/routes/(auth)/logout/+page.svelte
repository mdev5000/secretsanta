<script lang="ts">
    import {postData} from "$lib/rest/rest";
    import {AppErrorRs} from "$lib/requests/core/error";
    import {goto} from "$app/navigation";
    import {onMount} from "svelte";

    async function logout(): Promise<undefined> {
        const logoutRs = await postData<AppErrorRs>(AppErrorRs, '/api/logout', null)
        if (logoutRs.status !== 200) {
            console.log('failed', logoutRs);
        }
        setTimeout(() => {
            goto('/app/login');
        }, 1000);
        return;
    }

    onMount(async () => {
        await logout();
    });
</script>

<div class="logout-wrapper">
    logging out...
</div>

<style lang="scss">
  .logout-wrapper {
    margin: auto;
  }
</style>