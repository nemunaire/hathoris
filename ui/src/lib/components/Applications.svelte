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
            {#if source.controlable || source.hasplaylist}
                <div class="d-flex gap-1">
                    {#if source.hasplaylist}
                        <div class="btn-group" role="group">
                            <button
                                class="btn btn-sm btn-primary"
                                on:click={() => source.prevtrack()}
                            >
                                <i class="bi bi-skip-backward-fill"></i>
                            </button>
                            <button
                                class="btn btn-sm btn-primary"
                                on:click={() => source.nexttrack()}
                            >
                                <i class="bi bi-skip-forward-fill"></i>
                            </button>
                        </div>
                    {/if}
                    {#if source.controlable}
                        <button
                            class="btn btn-sm btn-primary"
                            on:click={() => source.playpause()}
                        >
                            <i class="bi bi-pause"></i>
                        </button>
                    {/if}
                </div>
            {/if}
        </li>
    {/each}
</ul>
