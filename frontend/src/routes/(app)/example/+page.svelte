
<script>
    import Button, {Label} from '@smui/button';

    let clicked = 0;
</script>

<div>
    <Button on:click={() => clicked++} variant="raised">
        <Label>Raised</Label>
    </Button>
    <span>Clicked {clicked}</span>
</div>