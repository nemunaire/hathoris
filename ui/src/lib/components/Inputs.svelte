<script>
 import { activeInputs, inputsList } from '$lib/stores/inputs';
 import { activeSources } from '$lib/stores/sources';

 export let showInactives = false;
</script>

<ul class="list-group list-group-flush">
    {#if $activeSources.length === 0 && ((showInactives && $inputsList.length === 0) || (!showInactives && $activeInputs.length === 0))}
        <li class="list-group-item py-3">
            <span class="text-muted">
                Aucune source active.
            </span>
        </li>
    {/if}
    {#each $activeSources as source}
        <li class="list-group-item py-3 d-flex justify-content-between">
            <div>
                <strong>{source.name}</strong>
                {#await source.currently()}
                    <div class="spinner-border spinner-border-sm" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                {:then title}
                    <span class="text-muted">{title}</span>
                {/await}
            </div>
            {#if source.controlable}
                <div>
                    <button
                        class="btn btn-sm btn-primary"
                        on:click={() => source.playpause()}
                    >
                        <i class="bi bi-pause"></i>
                    </button>
                </div>
            {/if}
        </li>
    {/each}
    {#each $inputsList as input}
        {#if showInactives || input.active}
            <li class="list-group-item py-3 d-flex justify-content-between">
                <div>
                    <strong>{input.name}</strong>
                    {#await input.currently()}
                        <div class="spinner-border spinner-border-sm" role="status">
                            <span class="visually-hidden">Loading...</span>
                        </div>
                    {:then title}
                        <span class="text-muted">{title}</span>
                    {/await}
                </div>
                {#if input.controlable}
                    <div>
                        <button
                            class="btn btn-sm btn-primary"
                            on:click={() => input.playpause()}
                        >
                            <i class="bi bi-pause"></i>
                        </button>
                    </div>
                {/if}
            </li>
        {/if}
    {/each}
</ul>
