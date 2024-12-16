<script>
 import '../hathoris.scss';
 import "bootstrap-icons/font/bootstrap-icons.css";

 import { activeSources, sources } from '$lib/stores/sources';
 sources.refresh();
 setInterval(sources.refresh, 5000);

 import { activeInputs, inputs } from '$lib/stores/inputs';
 inputs.refresh();
 setInterval(inputs.refresh, 4500);

 import SourceSelection from '$lib/components/SourceSelection.svelte';

 const version = fetch('api/version', {headers: {'Accept': 'application/json'}}).then((res) => res.json())
</script>

<svelte:head>
    <title>Hathoris</title>
</svelte:head>

<div class="flex-fill d-flex flex-column">
    <div class="container-fluid flex-fill d-flex flex-column justify-content-start">
        <div class="my-3">
            <SourceSelection />
        </div>

        {#if $activeSources.length === 0 && $activeInputs.length === 0}
            <div class="text-muted text-center mt-1 mb-1">
                Aucune source active pour l'instant.
            </div>
        {:else}
            <marquee>
                {#each $activeSources as source}
                    <div class="d-inline-block me-3">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                {#if source.currentTitle}
                                    <strong>{source.currentTitle}</strong> <span class="text-muted">@ {source.name}</span>
                                {:else}
                                    {source.name} activée
                                {/if}
                            </div>
                        </div>
                    </div>
                {/each}
                {#each $activeInputs as input}
                    <div class="d-inline-block me-3">
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                {#if input.streams && input.streams.length}
                                    {#each Object.keys(input.streams) as idstream}
                                        {@const title = input.streams[idstream]}
                                        <strong>{title}</strong>
                                    {/each}
                                    <span class="text-muted">@ {input.name}</span>
                {:else}
                                    {input.name} activée
                                {/if}
                            </div>
                        </div>
                    </div>
                {/each}
            </marquee>
        {/if}

        <slot></slot>
    </div>
</div>
