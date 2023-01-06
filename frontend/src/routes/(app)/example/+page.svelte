
<script>
    import Button, {Label} from '@smui/button';

    let another = 0;
    let clicked = 0;
    $: clicked2 = 2 * clicked;
    $: clicked3 = 2 * clicked2;

    $: clicked, updateAnother();

    function updateAnother() {
        another = 3 * clicked3;
    }

</script>

<div>
    <Button on:click={() => clicked++} variant="raised">
        <Label>Raised</Label>
    </Button>
    <div>Clicked {clicked}</div>
    <div>Clicked {clicked2}</div>
    <div>Clicked {clicked3}</div>
    <div>Another {another}</div>
</div>