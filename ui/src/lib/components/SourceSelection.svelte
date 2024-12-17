<script>
 import { sources } from '$lib/stores/sources';

 let activating_source = null;
 async function clickSource(src) {
     activating_source = src.id;
     try {
         if (src.enabled) {
             await src.deactivate();
             await sources.refresh();
         } else {
             await src.activate();
             await sources.refresh();
         }
     } catch {
     }
     activating_source = null;
 }
</script>

{#if !$sources}
    <div class="d-flex flex-fill justify-content-center">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
{:else}
    <div class="row row-cols-2 row-cols-sm-3 row-cols-md-4 row-cols-lg-6 row-cols-xl-auto justify-content-center g-2">
        {#each Object.keys($sources) as source}
            <div class="col d-flex flex-column justify-content-center align-items-center">
                <button
                    class="btn btn-lg"
                    class:btn-primary={$sources[source].enabled}
                    class:btn-secondary={!$sources[source].enabled}
                    disabled={activating_source !== null}
                    on:click={() => clickSource($sources[source])}
                >
                    {#if activating_source && activating_source === source}
                        <div class="spinner-grow spinner-grow-sm" role="status">
                            <span class="visually-hidden">Loading...</span>
                        </div>
                    {/if}
                    {$sources[source].name}
                </button>
            </div>
        {/each}
    </div>
{/if}
