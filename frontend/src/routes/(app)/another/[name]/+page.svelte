<script lang="ts">
    import {onMount} from "svelte";

    export let message = "";

    onMount(() => {
        message = "something"
    });

</script>

<div>
    {#if message}
        Hello this is my page {message}
    {:else}
        Loading...
    {/if}
</div>