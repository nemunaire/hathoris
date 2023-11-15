<script>
 import Applications from '$lib/components/Applications.svelte';
 import Inputs from '$lib/components/Inputs.svelte';
 import Mixer from '$lib/components/Mixer.svelte';
 import SourceSelection from '$lib/components/SourceSelection.svelte';
 import { activeSources } from '$lib/stores/sources';
 import { activeInputs } from '$lib/stores/inputs';

 let mixerAdvanced = false;
</script>

<div class="my-2">
    <SourceSelection />
</div>

<div class="container">
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
                            {#await source.currently()}
                                <div class="spinner-border spinner-border-sm" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div> <span class="text-muted">@ {source.name}</span>
                            {:then title}
                                <strong>{title}</strong> <span class="text-muted">@ {source.name}</span>
                            {:catch error}
                                {source.name} activée
                            {/await}
                        </div>
                    </div>
                </div>
            {/each}
            {#each $activeInputs as input}
                <div class="d-inline-block me-3">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            {#await input.streams()}
                                <div class="spinner-border spinner-border-sm" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div> <span class="text-muted">@ {input.name}</span>
                            {:then streams}
                                {#each Object.keys(streams) as idstream}
                                    {@const title = streams[idstream]}
                                    <strong>{title}</strong>
                                {/each}
                                <span class="text-muted">@ {input.name}</span>
                            {:catch error}
                                {input.name} activée
                            {/await}
                        </div>
                    </div>
                </div>
            {/each}
        </marquee>
    {/if}

    <div class="row">
        <div class="col-md">
            <div class="card my-2">
                <h4 class="card-header">
                    <div class="d-flex justify-content-between">
                        <div>
                            <i class="bi bi-sliders"></i>
                            Mixer
                        </div>
                        <button
                            class="btn btn-sm"
                            class:btn-info={mixerAdvanced}
                            class:btn-secondary={!mixerAdvanced}
                            on:click={() => { mixerAdvanced = !mixerAdvanced; }}
                        >
                            <i class="bi bi-alt"></i>
                        </button>
                    </div>
                </h4>
                <Mixer advanced={mixerAdvanced} />
            </div>
        </div>

        <div class="col-md">
            {#if $activeSources.length > 0}
                <div class="card my-2">
                    <h4 class="card-header">
                        <i class="bi bi-window-stack"></i>
                        Applications
                    </h4>
                    <Applications />
                </div>
            {/if}
            <div class="card my-2">
                <h4 class="card-header">
                    <i class="bi bi-speaker"></i>
                    Sources
                </h4>
                <Inputs />
            </div>
        </div>
    </div>
</div>
