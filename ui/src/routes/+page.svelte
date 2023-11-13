<script>
 import Mixer from '$lib/components/Mixer.svelte';
 import SourceSelection from '$lib/components/SourceSelection.svelte';
 import { activeSources } from '$lib/stores/sources';

 let mixerAdvanced = false;
</script>

<div class="my-2">
    <SourceSelection />
</div>

<div class="container">
    {#if $activeSources.length === 0}
        <div class="text-muted text-center mt-1 mb-1">
            Aucune source active pour l'instant.
        </div>
    {:else}
        <marquee>
            {#each $activeSources as source}
                <div class="d-inline-block me-3">
                    <div class="d-flex justify-content-between align-items-center">
                        <div>
                            <strong>{source.name}&nbsp;:</strong>
                            {#await source.currently()}
                                <div class="spinner-border spinner-border-sm" role="status">
                                    <span class="visually-hidden">Loading...</span>
                                </div>
                            {:then title}
                                {title}
                            {:catch error}
                                activée
                            {/await}
                        </div>
                    </div>
                </div>
            {/each}
        </marquee>
    {/if}

    <div class="row">
        <div class="col">
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

        <div class="col">
            <div class="card my-2">
                <h4 class="card-header">
                    <i class="bi bi-speaker"></i>
                    Sources
                </h4>
                <div class="card-body">
                </div>
            </div>
        </div>
    </div>
</div>
