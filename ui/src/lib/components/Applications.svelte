<script>
 import { activeSources } from '$lib/stores/sources';
</script>

<ul class="list-group list-group-flush">
    {#each $activeSources as source}
        <li class="list-group-item py-3 d-flex justify-content-between align-items-center">
            <div>
                <strong>{source.name}</strong>
                {#if source.currentTitle}
                    <span class="text-muted">{source.currentTitle}</span>
                {/if}
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
