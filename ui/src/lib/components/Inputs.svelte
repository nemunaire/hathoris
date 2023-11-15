<script>
 import { activeInputs, inputsList } from '$lib/stores/inputs';

 export let showInactives = false;
</script>

<ul class="list-group list-group-flush">
    {#if (showInactives && $inputsList.length === 0) || (!showInactives && $activeInputs.length === 0)}
        <li class="list-group-item py-3">
            <span class="text-muted">
                Aucune source active.
            </span>
        </li>
    {/if}
    {#each $inputsList as input}
        {#if showInactives || input.active}
            <li class="list-group-item py-3 d-flex flex-column">
                <strong>{input.name}</strong>
                {#await input.streams()}
                    <div class="spinner-border spinner-border-sm" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                {:then streams}
                    {#each Object.keys(streams) as idstream}
                        {@const title = streams[idstream]}
                        <div class="d-flex justify-content-between">
                            <div>
                                <span class="text-muted">{title}</span>
                            </div>
                            {#if input.controlable}
                                <div>
                                    <button
                                        class="btn btn-sm btn-primary"
                                        on:click={() => input.playpause(idstream)}
                                    >
                                        <i class="bi bi-pause"></i>
                                    </button>
                                </div>
                            {/if}
                        </div>
                    {/each}
                {/await}
            </li>
        {/if}
    {/each}
</ul>
