<script>
 import '../hathoris.scss'
 import "bootstrap-icons/font/bootstrap-icons.css";

 import { sources } from '$lib/stores/sources';
 sources.refresh();
 setInterval(sources.refresh, 5000);

 import { inputs } from '$lib/stores/inputs';
 inputs.refresh();
 setInterval(inputs.refresh, 4500);

 const version = fetch('api/version', {headers: {'Accept': 'application/json'}}).then((res) => res.json())
</script>

<svelte:head>
    <title>Hathoris</title>
</svelte:head>

<div class="flex-fill d-flex flex-column">
    <div class="container-fluid flex-fill d-flex flex-column justify-content-center">
        <slot></slot>
    </div>
</div>
