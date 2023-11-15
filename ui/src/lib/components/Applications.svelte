<script>
 import { activeSources } from '$lib/stores/sources';
</script>

<ul class="list-group list-group-flush">
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
</ul>
